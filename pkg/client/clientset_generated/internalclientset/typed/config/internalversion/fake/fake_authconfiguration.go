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
	componentconfig "kope.io/auth/pkg/apis/componentconfig"
)

// FakeAuthConfigurations implements AuthConfigurationInterface
type FakeAuthConfigurations struct {
	Fake *FakeConfig
	ns   string
}

var authconfigurationsResource = schema.GroupVersionResource{Group: "config.auth.kope.io", Version: "", Resource: "authconfigurations"}

var authconfigurationsKind = schema.GroupVersionKind{Group: "config.auth.kope.io", Version: "", Kind: "AuthConfiguration"}

// Get takes name of the authConfiguration, and returns the corresponding authConfiguration object, and an error if there is any.
func (c *FakeAuthConfigurations) Get(name string, options v1.GetOptions) (result *componentconfig.AuthConfiguration, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(authconfigurationsResource, c.ns, name), &componentconfig.AuthConfiguration{})

	if obj == nil {
		return nil, err
	}
	return obj.(*componentconfig.AuthConfiguration), err
}

// List takes label and field selectors, and returns the list of AuthConfigurations that match those selectors.
func (c *FakeAuthConfigurations) List(opts v1.ListOptions) (result *componentconfig.AuthConfigurationList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(authconfigurationsResource, authconfigurationsKind, c.ns, opts), &componentconfig.AuthConfigurationList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &componentconfig.AuthConfigurationList{}
	for _, item := range obj.(*componentconfig.AuthConfigurationList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested authConfigurations.
func (c *FakeAuthConfigurations) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(authconfigurationsResource, c.ns, opts))

}

// Create takes the representation of a authConfiguration and creates it.  Returns the server's representation of the authConfiguration, and an error, if there is any.
func (c *FakeAuthConfigurations) Create(authConfiguration *componentconfig.AuthConfiguration) (result *componentconfig.AuthConfiguration, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(authconfigurationsResource, c.ns, authConfiguration), &componentconfig.AuthConfiguration{})

	if obj == nil {
		return nil, err
	}
	return obj.(*componentconfig.AuthConfiguration), err
}

// Update takes the representation of a authConfiguration and updates it. Returns the server's representation of the authConfiguration, and an error, if there is any.
func (c *FakeAuthConfigurations) Update(authConfiguration *componentconfig.AuthConfiguration) (result *componentconfig.AuthConfiguration, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(authconfigurationsResource, c.ns, authConfiguration), &componentconfig.AuthConfiguration{})

	if obj == nil {
		return nil, err
	}
	return obj.(*componentconfig.AuthConfiguration), err
}

// Delete takes name of the authConfiguration and deletes it. Returns an error if one occurs.
func (c *FakeAuthConfigurations) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(authconfigurationsResource, c.ns, name), &componentconfig.AuthConfiguration{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAuthConfigurations) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(authconfigurationsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &componentconfig.AuthConfigurationList{})
	return err
}

// Patch applies the patch and returns the patched authConfiguration.
func (c *FakeAuthConfigurations) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *componentconfig.AuthConfiguration, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(authconfigurationsResource, c.ns, name, data, subresources...), &componentconfig.AuthConfiguration{})

	if obj == nil {
		return nil, err
	}
	return obj.(*componentconfig.AuthConfiguration), err
}
