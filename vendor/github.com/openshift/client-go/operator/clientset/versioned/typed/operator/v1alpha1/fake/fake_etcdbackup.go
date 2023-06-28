// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	v1alpha1 "github.com/openshift/api/operator/v1alpha1"
	operatorv1alpha1 "github.com/openshift/client-go/operator/applyconfigurations/operator/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeEtcdBackups implements EtcdBackupInterface
type FakeEtcdBackups struct {
	Fake *FakeOperatorV1alpha1
}

var etcdbackupsResource = schema.GroupVersionResource{Group: "operator.openshift.io", Version: "v1alpha1", Resource: "etcdbackups"}

var etcdbackupsKind = schema.GroupVersionKind{Group: "operator.openshift.io", Version: "v1alpha1", Kind: "EtcdBackup"}

// Get takes name of the etcdBackup, and returns the corresponding etcdBackup object, and an error if there is any.
func (c *FakeEtcdBackups) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.EtcdBackup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(etcdbackupsResource, name), &v1alpha1.EtcdBackup{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EtcdBackup), err
}

// List takes label and field selectors, and returns the list of EtcdBackups that match those selectors.
func (c *FakeEtcdBackups) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.EtcdBackupList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(etcdbackupsResource, etcdbackupsKind, opts), &v1alpha1.EtcdBackupList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.EtcdBackupList{ListMeta: obj.(*v1alpha1.EtcdBackupList).ListMeta}
	for _, item := range obj.(*v1alpha1.EtcdBackupList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested etcdBackups.
func (c *FakeEtcdBackups) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(etcdbackupsResource, opts))
}

// Create takes the representation of a etcdBackup and creates it.  Returns the server's representation of the etcdBackup, and an error, if there is any.
func (c *FakeEtcdBackups) Create(ctx context.Context, etcdBackup *v1alpha1.EtcdBackup, opts v1.CreateOptions) (result *v1alpha1.EtcdBackup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(etcdbackupsResource, etcdBackup), &v1alpha1.EtcdBackup{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EtcdBackup), err
}

// Update takes the representation of a etcdBackup and updates it. Returns the server's representation of the etcdBackup, and an error, if there is any.
func (c *FakeEtcdBackups) Update(ctx context.Context, etcdBackup *v1alpha1.EtcdBackup, opts v1.UpdateOptions) (result *v1alpha1.EtcdBackup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(etcdbackupsResource, etcdBackup), &v1alpha1.EtcdBackup{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EtcdBackup), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeEtcdBackups) UpdateStatus(ctx context.Context, etcdBackup *v1alpha1.EtcdBackup, opts v1.UpdateOptions) (*v1alpha1.EtcdBackup, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(etcdbackupsResource, "status", etcdBackup), &v1alpha1.EtcdBackup{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EtcdBackup), err
}

// Delete takes name of the etcdBackup and deletes it. Returns an error if one occurs.
func (c *FakeEtcdBackups) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(etcdbackupsResource, name, opts), &v1alpha1.EtcdBackup{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeEtcdBackups) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(etcdbackupsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.EtcdBackupList{})
	return err
}

// Patch applies the patch and returns the patched etcdBackup.
func (c *FakeEtcdBackups) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.EtcdBackup, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(etcdbackupsResource, name, pt, data, subresources...), &v1alpha1.EtcdBackup{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EtcdBackup), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied etcdBackup.
func (c *FakeEtcdBackups) Apply(ctx context.Context, etcdBackup *operatorv1alpha1.EtcdBackupApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.EtcdBackup, err error) {
	if etcdBackup == nil {
		return nil, fmt.Errorf("etcdBackup provided to Apply must not be nil")
	}
	data, err := json.Marshal(etcdBackup)
	if err != nil {
		return nil, err
	}
	name := etcdBackup.Name
	if name == nil {
		return nil, fmt.Errorf("etcdBackup.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(etcdbackupsResource, *name, types.ApplyPatchType, data), &v1alpha1.EtcdBackup{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EtcdBackup), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeEtcdBackups) ApplyStatus(ctx context.Context, etcdBackup *operatorv1alpha1.EtcdBackupApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.EtcdBackup, err error) {
	if etcdBackup == nil {
		return nil, fmt.Errorf("etcdBackup provided to Apply must not be nil")
	}
	data, err := json.Marshal(etcdBackup)
	if err != nil {
		return nil, err
	}
	name := etcdBackup.Name
	if name == nil {
		return nil, fmt.Errorf("etcdBackup.Name must be provided to Apply")
	}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(etcdbackupsResource, *name, types.ApplyPatchType, data, "status"), &v1alpha1.EtcdBackup{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EtcdBackup), err
}
