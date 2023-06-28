package backupcontroller

import (
	"context"
	"strings"
	"testing"

	backupv1alpha1 "github.com/openshift/api/operator/v1alpha1"
	fake "github.com/openshift/client-go/operator/clientset/versioned/fake"
	"github.com/openshift/cluster-etcd-operator/pkg/operator/operatorclient"
	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfakeclient "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestSyncLoopHappyPath(t *testing.T) {
	backup := backupv1alpha1.EtcdBackup{ObjectMeta: v1.ObjectMeta{Name: "test-backup"}, Spec: backupv1alpha1.EtcdBackupSpec{PVCName: "backup-happy-path-pvc"}}
	operatorFake := fake.NewSimpleClientset([]runtime.Object{&backup}...)
	client := k8sfakeclient.NewSimpleClientset()

	controller := BackupController{
		backupsClient:       operatorFake.OperatorV1alpha1(),
		kubeClient:          client,
		targetImagePullSpec: "pullspec-image",
	}

	err := controller.sync(context.TODO(), nil)
	require.NoError(t, err)

}

func TestJobCreationHappyPath(t *testing.T) {
	backup := backupv1alpha1.EtcdBackup{ObjectMeta: v1.ObjectMeta{Name: "test-backup"}, Spec: backupv1alpha1.EtcdBackupSpec{PVCName: "backup-happy-path-pvc"}}
	operatorFake := fake.NewSimpleClientset([]runtime.Object{&backup}...)
	client := k8sfakeclient.NewSimpleClientset()

	err := createBackupJob(context.Background(),
		backup,
		"pullspec-image",
		client.BatchV1().Jobs(operatorclient.TargetNamespace),
		operatorFake.OperatorV1alpha1().EtcdBackups(),
	)
	require.NoError(t, err)
	actions := operatorFake.Fake.Actions()
	require.Equal(t, 1, len(actions))
	createAction := actions[0].(k8stesting.CreateActionImpl)
	require.Equal(t, operatorclient.TargetNamespace, createAction.GetNamespace())
	require.Equal(t, "create", createAction.GetVerb())
	createdJob := createAction.Object.(*batchv1.Job)

	require.Truef(t, strings.HasPrefix(createdJob.Name, "cluster-backup-job"), "expected job.name [%s] to have prefix [cluster-backup-job]", createdJob.Name)
	require.Equal(t, operatorclient.TargetNamespace, createdJob.Namespace)
	require.Equal(t, "test-backup", createdJob.Labels[backupJobLabel])
	require.Equal(t, "pullspec-image", createdJob.Spec.Template.Spec.Containers[0].Image)

	foundVolume := false
	for _, volume := range createdJob.Spec.Template.Spec.Volumes {
		if volume.Name == "etc-kubernetes-cluster-backup" {
			foundVolume = true
			require.Equal(t, "backup-happy-path-pvc", volume.PersistentVolumeClaim.ClaimName)
		}
	}

	require.Truef(t, foundVolume, "could not find injected PVC volume in %v", createdJob.Spec.Template.Spec.Volumes)

}

func TestIndexJobsByBackupLabelName(t *testing.T) {
	jobList := &batchv1.JobList{
		Items: []batchv1.Job{
			{ObjectMeta: v1.ObjectMeta{Name: "test-1", Labels: map[string]string{backupJobLabel: "test-1"}}},
			{ObjectMeta: v1.ObjectMeta{Name: "test-2", Labels: map[string]string{backupJobLabel: "test-2"}}},
			{ObjectMeta: v1.ObjectMeta{Name: "test-3", Labels: map[string]string{backupJobLabel: "test-3"}}},
			{ObjectMeta: v1.ObjectMeta{Name: "test-4", Labels: map[string]string{"some-other-label": "value"}}},
		},
	}
	expected := map[string]batchv1.Job{}
	expected["test-1"] = jobList.Items[0]
	expected["test-2"] = jobList.Items[1]
	expected["test-3"] = jobList.Items[2]

	m := indexJobsByBackupLabelName(jobList)
	require.Equal(t, expected, m)
}

func TestIsJobComplete(t *testing.T) {
	tests := map[string]struct {
		condition batchv1.JobConditionType
		complete  bool
	}{
		"no condition": {condition: "", complete: false},
		"suspended":    {condition: batchv1.JobSuspended, complete: false},
		"complete":     {condition: batchv1.JobComplete, complete: true},
		"failed":       {condition: batchv1.JobFailed, complete: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			j := &batchv1.Job{
				Status: batchv1.JobStatus{
					Conditions: []batchv1.JobCondition{
						{Type: test.condition, Status: corev1.ConditionTrue},
					},
				},
			}
			finished := isJobFinished(j)
			require.Equal(t, test.complete, finished)
		})
	}

}
