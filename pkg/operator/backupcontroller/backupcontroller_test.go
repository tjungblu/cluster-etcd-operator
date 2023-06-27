package backupcontroller

import (
	"context"
	"strings"
	"testing"

	backupv1alpha1 "github.com/openshift/api/backup/v1alpha1"
	"github.com/openshift/cluster-etcd-operator/pkg/operator/operatorclient"
	"github.com/stretchr/testify/require"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestJobCreationHappyPath(t *testing.T) {
	client := fake.NewSimpleClientset()
	err := createBackupJob(context.Background(),
		backupv1alpha1.EtcdBackup{
			ObjectMeta: v1.ObjectMeta{
				Name:      "test-backup",
				Namespace: operatorclient.TargetNamespace,
			},
			Spec: backupv1alpha1.EtcdBackupSpec{PVCName: "backup-happy-path-pvc"}},
		"pullspec-image",
		client.BatchV1().Jobs(operatorclient.TargetNamespace),
	)
	require.NoError(t, err)
	actions := client.Fake.Actions()
	require.Equal(t, 1, len(actions))
	createAction := actions[0].(k8stesting.CreateActionImpl)
	require.Equal(t, operatorclient.TargetNamespace, createAction.GetNamespace())
	require.Equal(t, "create", createAction.GetVerb())
	createdJob := createAction.Object.(*batchv1.Job)

	require.Truef(t, strings.HasPrefix(createdJob.Name, "cluster-backup-job"), "expected job.name [%s] to have prefix [cluster-backup-job]", createdJob.Name)
	require.Equal(t, operatorclient.TargetNamespace, createdJob.Namespace)
	require.Equal(t, "test-backup", createdJob.Labels["backup-name"])
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
