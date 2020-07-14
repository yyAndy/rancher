/*
Copyright 2020 Rancher Labs, Inc.

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

// Code generated by main. DO NOT EDIT.

package v3

import (
	"context"
	"time"

	"github.com/rancher/lasso/pkg/client"
	"github.com/rancher/lasso/pkg/controller"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

type ProjectAlertHandler func(string, *v3.ProjectAlert) (*v3.ProjectAlert, error)

type ProjectAlertController interface {
	generic.ControllerMeta
	ProjectAlertClient

	OnChange(ctx context.Context, name string, sync ProjectAlertHandler)
	OnRemove(ctx context.Context, name string, sync ProjectAlertHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ProjectAlertCache
}

type ProjectAlertClient interface {
	Create(*v3.ProjectAlert) (*v3.ProjectAlert, error)
	Update(*v3.ProjectAlert) (*v3.ProjectAlert, error)
	UpdateStatus(*v3.ProjectAlert) (*v3.ProjectAlert, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.ProjectAlert, error)
	List(namespace string, opts metav1.ListOptions) (*v3.ProjectAlertList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.ProjectAlert, err error)
}

type ProjectAlertCache interface {
	Get(namespace, name string) (*v3.ProjectAlert, error)
	List(namespace string, selector labels.Selector) ([]*v3.ProjectAlert, error)

	AddIndexer(indexName string, indexer ProjectAlertIndexer)
	GetByIndex(indexName, key string) ([]*v3.ProjectAlert, error)
}

type ProjectAlertIndexer func(obj *v3.ProjectAlert) ([]string, error)

type projectAlertController struct {
	controller    controller.SharedController
	client        *client.Client
	gvk           schema.GroupVersionKind
	groupResource schema.GroupResource
}

func NewProjectAlertController(gvk schema.GroupVersionKind, resource string, namespaced bool, controller controller.SharedControllerFactory) ProjectAlertController {
	c := controller.ForResourceKind(gvk.GroupVersion().WithResource(resource), gvk.Kind, namespaced)
	return &projectAlertController{
		controller: c,
		client:     c.Client(),
		gvk:        gvk,
		groupResource: schema.GroupResource{
			Group:    gvk.Group,
			Resource: resource,
		},
	}
}

func FromProjectAlertHandlerToHandler(sync ProjectAlertHandler) generic.Handler {
	return func(key string, obj runtime.Object) (ret runtime.Object, err error) {
		var v *v3.ProjectAlert
		if obj == nil {
			v, err = sync(key, nil)
		} else {
			v, err = sync(key, obj.(*v3.ProjectAlert))
		}
		if v == nil {
			return nil, err
		}
		return v, err
	}
}

func (c *projectAlertController) Updater() generic.Updater {
	return func(obj runtime.Object) (runtime.Object, error) {
		newObj, err := c.Update(obj.(*v3.ProjectAlert))
		if newObj == nil {
			return nil, err
		}
		return newObj, err
	}
}

func UpdateProjectAlertDeepCopyOnChange(client ProjectAlertClient, obj *v3.ProjectAlert, handler func(obj *v3.ProjectAlert) (*v3.ProjectAlert, error)) (*v3.ProjectAlert, error) {
	if obj == nil {
		return obj, nil
	}

	copyObj := obj.DeepCopy()
	newObj, err := handler(copyObj)
	if newObj != nil {
		copyObj = newObj
	}
	if obj.ResourceVersion == copyObj.ResourceVersion && !equality.Semantic.DeepEqual(obj, copyObj) {
		return client.Update(copyObj)
	}

	return copyObj, err
}

func (c *projectAlertController) AddGenericHandler(ctx context.Context, name string, handler generic.Handler) {
	c.controller.RegisterHandler(ctx, name, controller.SharedControllerHandlerFunc(handler))
}

func (c *projectAlertController) AddGenericRemoveHandler(ctx context.Context, name string, handler generic.Handler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), handler))
}

func (c *projectAlertController) OnChange(ctx context.Context, name string, sync ProjectAlertHandler) {
	c.AddGenericHandler(ctx, name, FromProjectAlertHandlerToHandler(sync))
}

func (c *projectAlertController) OnRemove(ctx context.Context, name string, sync ProjectAlertHandler) {
	c.AddGenericHandler(ctx, name, generic.NewRemoveHandler(name, c.Updater(), FromProjectAlertHandlerToHandler(sync)))
}

func (c *projectAlertController) Enqueue(namespace, name string) {
	c.controller.Enqueue(namespace, name)
}

func (c *projectAlertController) EnqueueAfter(namespace, name string, duration time.Duration) {
	c.controller.EnqueueAfter(namespace, name, duration)
}

func (c *projectAlertController) Informer() cache.SharedIndexInformer {
	return c.controller.Informer()
}

func (c *projectAlertController) GroupVersionKind() schema.GroupVersionKind {
	return c.gvk
}

func (c *projectAlertController) Cache() ProjectAlertCache {
	return &projectAlertCache{
		indexer:  c.Informer().GetIndexer(),
		resource: c.groupResource,
	}
}

func (c *projectAlertController) Create(obj *v3.ProjectAlert) (*v3.ProjectAlert, error) {
	result := &v3.ProjectAlert{}
	return result, c.client.Create(context.TODO(), obj.Namespace, obj, result, metav1.CreateOptions{})
}

func (c *projectAlertController) Update(obj *v3.ProjectAlert) (*v3.ProjectAlert, error) {
	result := &v3.ProjectAlert{}
	return result, c.client.Update(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *projectAlertController) UpdateStatus(obj *v3.ProjectAlert) (*v3.ProjectAlert, error) {
	result := &v3.ProjectAlert{}
	return result, c.client.UpdateStatus(context.TODO(), obj.Namespace, obj, result, metav1.UpdateOptions{})
}

func (c *projectAlertController) Delete(namespace, name string, options *metav1.DeleteOptions) error {
	if options == nil {
		options = &metav1.DeleteOptions{}
	}
	return c.client.Delete(context.TODO(), namespace, name, *options)
}

func (c *projectAlertController) Get(namespace, name string, options metav1.GetOptions) (*v3.ProjectAlert, error) {
	result := &v3.ProjectAlert{}
	return result, c.client.Get(context.TODO(), namespace, name, result, options)
}

func (c *projectAlertController) List(namespace string, opts metav1.ListOptions) (*v3.ProjectAlertList, error) {
	result := &v3.ProjectAlertList{}
	return result, c.client.List(context.TODO(), namespace, result, opts)
}

func (c *projectAlertController) Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
	return c.client.Watch(context.TODO(), namespace, opts)
}

func (c *projectAlertController) Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (*v3.ProjectAlert, error) {
	result := &v3.ProjectAlert{}
	return result, c.client.Patch(context.TODO(), namespace, name, pt, data, result, metav1.PatchOptions{}, subresources...)
}

type projectAlertCache struct {
	indexer  cache.Indexer
	resource schema.GroupResource
}

func (c *projectAlertCache) Get(namespace, name string) (*v3.ProjectAlert, error) {
	obj, exists, err := c.indexer.GetByKey(namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(c.resource, name)
	}
	return obj.(*v3.ProjectAlert), nil
}

func (c *projectAlertCache) List(namespace string, selector labels.Selector) (ret []*v3.ProjectAlert, err error) {

	err = cache.ListAllByNamespace(c.indexer, namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v3.ProjectAlert))
	})

	return ret, err
}

func (c *projectAlertCache) AddIndexer(indexName string, indexer ProjectAlertIndexer) {
	utilruntime.Must(c.indexer.AddIndexers(map[string]cache.IndexFunc{
		indexName: func(obj interface{}) (strings []string, e error) {
			return indexer(obj.(*v3.ProjectAlert))
		},
	}))
}

func (c *projectAlertCache) GetByIndex(indexName, key string) (result []*v3.ProjectAlert, err error) {
	objs, err := c.indexer.ByIndex(indexName, key)
	if err != nil {
		return nil, err
	}
	result = make([]*v3.ProjectAlert, 0, len(objs))
	for _, obj := range objs {
		result = append(result, obj.(*v3.ProjectAlert))
	}
	return result, nil
}

type ProjectAlertStatusHandler func(obj *v3.ProjectAlert, status v3.AlertStatus) (v3.AlertStatus, error)

type ProjectAlertGeneratingHandler func(obj *v3.ProjectAlert, status v3.AlertStatus) ([]runtime.Object, v3.AlertStatus, error)

func RegisterProjectAlertStatusHandler(ctx context.Context, controller ProjectAlertController, condition condition.Cond, name string, handler ProjectAlertStatusHandler) {
	statusHandler := &projectAlertStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromProjectAlertHandlerToHandler(statusHandler.sync))
}

func RegisterProjectAlertGeneratingHandler(ctx context.Context, controller ProjectAlertController, apply apply.Apply,
	condition condition.Cond, name string, handler ProjectAlertGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &projectAlertGeneratingHandler{
		ProjectAlertGeneratingHandler: handler,
		apply:                         apply,
		name:                          name,
		gvk:                           controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterProjectAlertStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type projectAlertStatusHandler struct {
	client    ProjectAlertClient
	condition condition.Cond
	handler   ProjectAlertStatusHandler
}

func (a *projectAlertStatusHandler) sync(key string, obj *v3.ProjectAlert) (*v3.ProjectAlert, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		var newErr error
		obj.Status = newStatus
		obj, newErr = a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
	}
	return obj, err
}

type projectAlertGeneratingHandler struct {
	ProjectAlertGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *projectAlertGeneratingHandler) Remove(key string, obj *v3.ProjectAlert) (*v3.ProjectAlert, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.ProjectAlert{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *projectAlertGeneratingHandler) Handle(obj *v3.ProjectAlert, status v3.AlertStatus) (v3.AlertStatus, error) {
	objs, newStatus, err := a.ProjectAlertGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
