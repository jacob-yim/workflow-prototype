/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	workflowv1 "github.com/jacob-yim/workflow-prototype/pkg/api/workflow/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeWorkflowTasks implements WorkflowTaskInterface
type FakeWorkflowTasks struct {
	Fake *FakeWorkflowV1
	ns   string
}

var workflowtasksResource = schema.GroupVersionResource{Group: "workflow", Version: "v1", Resource: "workflowtasks"}

var workflowtasksKind = schema.GroupVersionKind{Group: "workflow", Version: "v1", Kind: "WorkflowTask"}

// Get takes name of the workflowTask, and returns the corresponding workflowTask object, and an error if there is any.
func (c *FakeWorkflowTasks) Get(ctx context.Context, name string, options v1.GetOptions) (result *workflowv1.WorkflowTask, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(workflowtasksResource, c.ns, name), &workflowv1.WorkflowTask{})

	if obj == nil {
		return nil, err
	}
	return obj.(*workflowv1.WorkflowTask), err
}

// List takes label and field selectors, and returns the list of WorkflowTasks that match those selectors.
func (c *FakeWorkflowTasks) List(ctx context.Context, opts v1.ListOptions) (result *workflowv1.WorkflowTaskList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(workflowtasksResource, workflowtasksKind, c.ns, opts), &workflowv1.WorkflowTaskList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &workflowv1.WorkflowTaskList{ListMeta: obj.(*workflowv1.WorkflowTaskList).ListMeta}
	for _, item := range obj.(*workflowv1.WorkflowTaskList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested workflowTasks.
func (c *FakeWorkflowTasks) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(workflowtasksResource, c.ns, opts))

}

// Create takes the representation of a workflowTask and creates it.  Returns the server's representation of the workflowTask, and an error, if there is any.
func (c *FakeWorkflowTasks) Create(ctx context.Context, workflowTask *workflowv1.WorkflowTask, opts v1.CreateOptions) (result *workflowv1.WorkflowTask, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(workflowtasksResource, c.ns, workflowTask), &workflowv1.WorkflowTask{})

	if obj == nil {
		return nil, err
	}
	return obj.(*workflowv1.WorkflowTask), err
}

// Update takes the representation of a workflowTask and updates it. Returns the server's representation of the workflowTask, and an error, if there is any.
func (c *FakeWorkflowTasks) Update(ctx context.Context, workflowTask *workflowv1.WorkflowTask, opts v1.UpdateOptions) (result *workflowv1.WorkflowTask, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(workflowtasksResource, c.ns, workflowTask), &workflowv1.WorkflowTask{})

	if obj == nil {
		return nil, err
	}
	return obj.(*workflowv1.WorkflowTask), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeWorkflowTasks) UpdateStatus(ctx context.Context, workflowTask *workflowv1.WorkflowTask, opts v1.UpdateOptions) (*workflowv1.WorkflowTask, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(workflowtasksResource, "status", c.ns, workflowTask), &workflowv1.WorkflowTask{})

	if obj == nil {
		return nil, err
	}
	return obj.(*workflowv1.WorkflowTask), err
}

// Delete takes name of the workflowTask and deletes it. Returns an error if one occurs.
func (c *FakeWorkflowTasks) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(workflowtasksResource, c.ns, name), &workflowv1.WorkflowTask{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeWorkflowTasks) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(workflowtasksResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &workflowv1.WorkflowTaskList{})
	return err
}

// Patch applies the patch and returns the patched workflowTask.
func (c *FakeWorkflowTasks) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *workflowv1.WorkflowTask, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(workflowtasksResource, c.ns, name, pt, data, subresources...), &workflowv1.WorkflowTask{})

	if obj == nil {
		return nil, err
	}
	return obj.(*workflowv1.WorkflowTask), err
}
