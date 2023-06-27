package backupcontroller

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	backupv1alpha1 "github.com/openshift/api/backup/v1alpha1"
	backupv1client "github.com/openshift/client-go/backup/clientset/versioned/typed/backup/v1alpha1"
	configv1client "github.com/openshift/client-go/config/clientset/versioned/typed/config/v1"
	configv1informers "github.com/openshift/client-go/config/informers/externalversions/config/v1"
	"github.com/openshift/cluster-etcd-operator/pkg/operator/etcd_assets"
	"github.com/openshift/cluster-etcd-operator/pkg/operator/health"
	"github.com/openshift/cluster-etcd-operator/pkg/operator/operatorclient"
	"github.com/openshift/library-go/pkg/controller/factory"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/storage/names"
	batchv1client "k8s.io/client-go/kubernetes/typed/batch/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

const (
	backupLabel      = "cluster-backup-job"
	recentBackupPath = "/etc/kubernetes/cluster-backup"
	backupDirEnvName = "CLUSTER_BACKUP_PATH"
)

type BackupController struct {
	operatorClient        v1helpers.OperatorClient
	clusterOperatorClient configv1client.ClusterOperatorsGetter
	backupsClient         backupv1client.EtcdBackupsGetter
	kubeClient            kubernetes.Interface
	targetImagePullSpec   string
}

func NewBackupController(
	livenessChecker *health.MultiAlivenessChecker,
	operatorClient v1helpers.OperatorClient,
	clusterOperatorClient configv1client.ClusterOperatorsGetter,
	backupsClient backupv1client.EtcdBackupsGetter,
	kubeClient kubernetes.Interface,
	kubeInformers v1helpers.KubeInformersForNamespaces,
	clusterVersionInformer configv1informers.ClusterVersionInformer,
	clusterOperatorInformer configv1informers.ClusterOperatorInformer,
	eventRecorder events.Recorder,
	targetImagePullSpec string,
) factory.Controller {
	c := &BackupController{
		operatorClient:        operatorClient,
		clusterOperatorClient: clusterOperatorClient,
		backupsClient:         backupsClient,
		kubeClient:            kubeClient,
		targetImagePullSpec:   targetImagePullSpec,
	}

	syncer := health.NewDefaultCheckingSyncWrapper(c.sync)
	livenessChecker.Add("BackupController", syncer)

	return factory.New().ResyncEvery(5*time.Minute).WithInformers(
		kubeInformers.InformersFor(operatorclient.TargetNamespace).Core().V1().Pods().Informer(),
		operatorClient.Informer(),
		clusterVersionInformer.Informer(),
		clusterOperatorInformer.Informer(),
	).WithSync(syncer.Sync).ToController("BackupController", eventRecorder.WithComponentSuffix("backup-controller"))
}

func (c *BackupController) sync(ctx context.Context, syncCtx factory.SyncContext) error {
	jobsClient := c.kubeClient.BatchV1().Jobs(operatorclient.TargetNamespace)
	currentJobs, err := jobsClient.List(ctx, v1.ListOptions{LabelSelector: "app=" + backupLabel})
	if err != nil {
		return fmt.Errorf("BackupController could not list backup jobs, error was: %w", err)
	}

	jobIndexed := indexJobsByBackupLabelName(currentJobs)
	backups, err := c.backupsClient.EtcdBackups().List(ctx, v1.ListOptions{})
	if err != nil {
		return fmt.Errorf("BackupController could not list etcdbackups CRDs, error was: %w", err)
	}

	if backups == nil {
		return nil
	}

	klog.V(4).Infof("BackupController backup: %v", backups)
	var backupsToRun []backupv1alpha1.EtcdBackup
	for _, item := range backups.Items {
		if backupJob, ok := jobIndexed[item.Name]; ok {
			klog.V(4).Infof("BackupController backup with name [%s] found, skipping", backupJob.Name)
			// TODO(thomas): reconcile its status, add another label, so we don't have to re-process it
			continue
		}

		backupsToRun = append(backupsToRun, *item.DeepCopy())
	}

	if len(backupsToRun) == 0 {
		klog.V(4).Infof("BackupController no backups to reconcile, skipping")
		return nil
	}

	// we only allow to run one at a time, if there's currently a job running then we will skip it in this reconciliation step
	runningJobs := findRunningJobs(currentJobs)
	if len(runningJobs) > 0 {
		klog.V(4).Infof("BackupController already found [%d] running jobs, skipping", len(runningJobs))
		return nil
	}

	// in case of multiple backups requested, we're trying to reconcile in order of their names
	sort.Slice(backupsToRun, func(i, j int) bool {
		return strings.Compare(backupsToRun[i].Name, backupsToRun[j].Name) < 0
	})

	for _, backup := range backupsToRun {
		err := createBackupJob(ctx, backup, c.targetImagePullSpec, jobsClient)
		if err != nil {
			return err
		}

		// only ever reconcile one item at a time
		// TODO(thomas): we should mark the remainder as "Skipped"
		return nil
	}

	return nil
}

func indexJobsByBackupLabelName(jobs *batchv1.JobList) map[string]batchv1.Job {
	m := map[string]batchv1.Job{}
	if jobs == nil {
		return m
	}

	for _, j := range jobs.Items {
		backupCrdName := j.Labels["backup-name"]
		m[backupCrdName] = *j.DeepCopy()
	}

	return m
}

func findRunningJobs(jobs *batchv1.JobList) []batchv1.Job {
	var running []batchv1.Job
	if jobs == nil {
		return running
	}

	for _, j := range jobs.Items {
		if !isJobFinished(&j) {
			running = append(running, *j.DeepCopy())
		}
	}

	return running
}

// isJobFinished checks whether the given Job has finished execution.
// It does not discriminate between successful and failed terminations.
func isJobFinished(j *batchv1.Job) bool {
	for _, c := range j.Status.Conditions {
		if (c.Type == batchv1.JobComplete || c.Type == batchv1.JobFailed) && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func createBackupJob(ctx context.Context, backup backupv1alpha1.EtcdBackup, targetImagePullSpec string, client batchv1client.JobInterface) error {
	scheme := runtime.NewScheme()
	codec := serializer.NewCodecFactory(scheme)
	err := batchv1.AddToScheme(scheme)
	if err != nil {
		return fmt.Errorf("BackupController could not add batchv1 scheme: %w", err)
	}

	obj, err := runtime.Decode(codec.UniversalDecoder(batchv1.SchemeGroupVersion), etcd_assets.MustAsset("etcd/cluster-backup-job.yaml"))
	if err != nil {
		return fmt.Errorf("BackupController could not decode batchv1 job scheme: %w", err)
	}

	backupFileName := fmt.Sprintf("backup-%s-%s", backup.Name, time.Now().Format("2006-01-02_150405"))

	job := obj.(*batchv1.Job)
	job.Name = names.SimpleNameGenerator.GenerateName(job.Name)
	job.Labels["backup-name"] = backup.Name

	job.Spec.Template.Spec.Containers[0].Image = targetImagePullSpec
	job.Spec.Template.Spec.Containers[0].Env = []corev1.EnvVar{
		{Name: backupDirEnvName, Value: fmt.Sprintf("%s/%s", recentBackupPath, backupFileName)},
	}

	injected := false
	for _, mount := range job.Spec.Template.Spec.Volumes {
		if mount.Name == "etc-kubernetes-cluster-backup" {
			mount.PersistentVolumeClaim.ClaimName = backup.Spec.PVCName
			injected = true
			break
		}
	}

	if !injected {
		return fmt.Errorf("could not inject PVC into Job template, please check the included cluster-backup-job.yaml")
	}

	klog.Infof("BackupController starts with backup [%s] as job [%s], writing to filename [%s]", backup.Name, job.Name, backupFileName)
	_, err = client.Create(ctx, job, v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("BackupController could create job: %w", err)
	}
	return nil
}
