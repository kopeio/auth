/*
Copyright 2017 The Kubernetes Authors.

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

package fake

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "kope.io/auth/pkg/apis/componentconfig/v1alpha1"
)

// FakeAuthProviders implements AuthProviderInterface
type FakeAuthProviders struct {
	Fake *FakeConfigV1alpha1
}

var authprovidersResource = schema.GroupVersionResource{Group: "config.auth.kope.io", Version: "v1alpha1", Resource: "authproviders"}

var authprovidersKind = schema.GroupVersionKind{Group: "config.auth.kope.io", Version: "v1alpha1", Kind: "AuthProvider"}

// Get takes name of the authProvider, and returns the corresponding authProvider object, and an error if there is any.
func (c *FakeAuthProviders) Get(name string, options v1.GetOptions) (result *v1alpha1.AuthProvider, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(authprovidersResource, name), &v1alpha1.AuthProvider{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AuthProvider), err
}

// List takes label and field selectors, and returns the list of AuthProviders that match those selectors.
func (c *FakeAuthProviders) List(opts v1.ListOptions) (result *v1alpha1.AuthProviderList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(authprovidersResource, authprovidersKind, opts), &v1alpha1.AuthProviderList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.AuthProviderList{}
	for _, item := range obj.(*v1alpha1.AuthProviderList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested authProviders.
func (c *FakeAuthProviders) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(authprovidersResource, opts))
}

// Create takes the representation of a authProvider and creates it.  Returns the server's representation of the authProvider, and an error, if there is any.
func (c *FakeAuthProviders) Create(authProvider *v1alpha1.AuthProvider) (result *v1alpha1.AuthProvider, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(authprovidersResource, authProvider), &v1alpha1.AuthProvider{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AuthProvider), err
}

// Update takes the representation of a authProvider and updates it. Returns the server's representation of the authProvider, and an error, if there is any.
func (c *FakeAuthProviders) Update(authProvider *v1alpha1.AuthProvider) (result *v1alpha1.AuthProvider, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(authprovidersResource, authProvider), &v1alpha1.AuthProvider{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AuthProvider), err
}

// Delete takes name of the authProvider and deletes it. Returns an error if one occurs.
func (c *FakeAuthProviders) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(authprovidersResource, name), &v1alpha1.AuthProvider{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAuthProviders) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(authprovidersResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.AuthProviderList{})
	return err
}

// Patch applies the patch and returns the patched authProvider.
func (c *FakeAuthProviders) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.AuthProvider, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(authprovidersResource, name, data, subresources...), &v1alpha1.AuthProvider{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AuthProvider), err
}
