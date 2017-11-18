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

package v1alpha1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1alpha1 "kope.io/auth/pkg/apis/componentconfig/v1alpha1"
	scheme "kope.io/auth/pkg/client/clientset_generated/internalclientset/scheme"
)

// AuthConfigurationsGetter has a method to return a AuthConfigurationInterface.
// A group's client should implement this interface.
type AuthConfigurationsGetter interface {
	AuthConfigurations() AuthConfigurationInterface
}

// AuthConfigurationInterface has methods to work with AuthConfiguration resources.
type AuthConfigurationInterface interface {
	Create(*v1alpha1.AuthConfiguration) (*v1alpha1.AuthConfiguration, error)
	Update(*v1alpha1.AuthConfiguration) (*v1alpha1.AuthConfiguration, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.AuthConfiguration, error)
	List(opts v1.ListOptions) (*v1alpha1.AuthConfigurationList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.AuthConfiguration, err error)
	AuthConfigurationExpansion
}

// authConfigurations implements AuthConfigurationInterface
type authConfigurations struct {
	client rest.Interface
}

// newAuthConfigurations returns a AuthConfigurations
func newAuthConfigurations(c *ConfigV1alpha1Client) *authConfigurations {
	return &authConfigurations{
		client: c.RESTClient(),
	}
}

// Get takes name of the authConfiguration, and returns the corresponding authConfiguration object, and an error if there is any.
func (c *authConfigurations) Get(name string, options v1.GetOptions) (result *v1alpha1.AuthConfiguration, err error) {
	result = &v1alpha1.AuthConfiguration{}
	err = c.client.Get().
		Resource("authconfigurations").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of AuthConfigurations that match those selectors.
func (c *authConfigurations) List(opts v1.ListOptions) (result *v1alpha1.AuthConfigurationList, err error) {
	result = &v1alpha1.AuthConfigurationList{}
	err = c.client.Get().
		Resource("authconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested authConfigurations.
func (c *authConfigurations) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Resource("authconfigurations").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a authConfiguration and creates it.  Returns the server's representation of the authConfiguration, and an error, if there is any.
func (c *authConfigurations) Create(authConfiguration *v1alpha1.AuthConfiguration) (result *v1alpha1.AuthConfiguration, err error) {
	result = &v1alpha1.AuthConfiguration{}
	err = c.client.Post().
		Resource("authconfigurations").
		Body(authConfiguration).
		Do().
		Into(result)
	return
}

// Update takes the representation of a authConfiguration and updates it. Returns the server's representation of the authConfiguration, and an error, if there is any.
func (c *authConfigurations) Update(authConfiguration *v1alpha1.AuthConfiguration) (result *v1alpha1.AuthConfiguration, err error) {
	result = &v1alpha1.AuthConfiguration{}
	err = c.client.Put().
		Resource("authconfigurations").
		Name(authConfiguration.Name).
		Body(authConfiguration).
		Do().
		Into(result)
	return
}

// Delete takes name of the authConfiguration and deletes it. Returns an error if one occurs.
func (c *authConfigurations) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("authconfigurations").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *authConfigurations) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Resource("authconfigurations").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched authConfiguration.
func (c *authConfigurations) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.AuthConfiguration, err error) {
	result = &v1alpha1.AuthConfiguration{}
	err = c.client.Patch(pt).
		Resource("authconfigurations").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
