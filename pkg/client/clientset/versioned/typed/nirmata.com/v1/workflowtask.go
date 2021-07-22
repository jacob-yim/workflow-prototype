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

package v1

import (
	"context"
	"time"

	v1 "github.com/jacob-yim/workflow-prototype/pkg/apis/nirmata.com/v1"
	scheme "github.com/jacob-yim/workflow-prototype/pkg/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// WorkflowTasksGetter has a method to return a WorkflowTaskInterface.
// A group's client should implement this interface.
type WorkflowTasksGetter interface {
	WorkflowTasks(namespace string) WorkflowTaskInterface
}

// WorkflowTaskInterface has methods to work with WorkflowTask resources.
type WorkflowTaskInterface interface {
	Create(ctx context.Context, workflowTask *v1.WorkflowTask, opts metav1.CreateOptions) (*v1.WorkflowTask, error)
	Update(ctx context.Context, workflowTask *v1.WorkflowTask, opts metav1.UpdateOptions) (*v1.WorkflowTask, error)
	UpdateStatus(ctx context.Context, workflowTask *v1.WorkflowTask, opts metav1.UpdateOptions) (*v1.WorkflowTask, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.WorkflowTask, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.WorkflowTaskList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.WorkflowTask, err error)
	WorkflowTaskExpansion
}

// workflowTasks implements WorkflowTaskInterface
type workflowTasks struct {
	client rest.Interface
	ns     string
}

// newWorkflowTasks returns a WorkflowTasks
func newWorkflowTasks(c *NirmataV1Client, namespace string) *workflowTasks {
	return &workflowTasks{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the workflowTask, and returns the corresponding workflowTask object, and an error if there is any.
func (c *workflowTasks) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.WorkflowTask, err error) {
	result = &v1.WorkflowTask{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("workflowtasks").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of WorkflowTasks that match those selectors.
func (c *workflowTasks) List(ctx context.Context, opts metav1.ListOptions) (result *v1.WorkflowTaskList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.WorkflowTaskList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("workflowtasks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested workflowTasks.
func (c *workflowTasks) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("workflowtasks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a workflowTask and creates it.  Returns the server's representation of the workflowTask, and an error, if there is any.
func (c *workflowTasks) Create(ctx context.Context, workflowTask *v1.WorkflowTask, opts metav1.CreateOptions) (result *v1.WorkflowTask, err error) {
	result = &v1.WorkflowTask{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("workflowtasks").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(workflowTask).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a workflowTask and updates it. Returns the server's representation of the workflowTask, and an error, if there is any.
func (c *workflowTasks) Update(ctx context.Context, workflowTask *v1.WorkflowTask, opts metav1.UpdateOptions) (result *v1.WorkflowTask, err error) {
	result = &v1.WorkflowTask{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("workflowtasks").
		Name(workflowTask.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(workflowTask).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *workflowTasks) UpdateStatus(ctx context.Context, workflowTask *v1.WorkflowTask, opts metav1.UpdateOptions) (result *v1.WorkflowTask, err error) {
	result = &v1.WorkflowTask{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("workflowtasks").
		Name(workflowTask.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(workflowTask).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the workflowTask and deletes it. Returns an error if one occurs.
func (c *workflowTasks) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("workflowtasks").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *workflowTasks) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("workflowtasks").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched workflowTask.
func (c *workflowTasks) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.WorkflowTask, err error) {
	result = &v1.WorkflowTask{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("workflowtasks").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
