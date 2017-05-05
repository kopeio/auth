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

package clientset

import (
	glog "github.com/golang/glog"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
	authv1alpha1 "kope.io/auth/pkg/client/clientset_generated/clientset/typed/auth/v1alpha1"
	componentconfigv1alpha1 "kope.io/auth/pkg/client/clientset_generated/clientset/typed/componentconfig/v1alpha1"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	AuthV1alpha1() authv1alpha1.AuthV1alpha1Interface
	// Deprecated: please explicitly pick a version if possible.
	Auth() authv1alpha1.AuthV1alpha1Interface
	ComponentconfigV1alpha1() componentconfigv1alpha1.ComponentconfigV1alpha1Interface
	// Deprecated: please explicitly pick a version if possible.
	Componentconfig() componentconfigv1alpha1.ComponentconfigV1alpha1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	*authv1alpha1.AuthV1alpha1Client
	*componentconfigv1alpha1.ComponentconfigV1alpha1Client
}

// AuthV1alpha1 retrieves the AuthV1alpha1Client
func (c *Clientset) AuthV1alpha1() authv1alpha1.AuthV1alpha1Interface {
	if c == nil {
		return nil
	}
	return c.AuthV1alpha1Client
}

// Deprecated: Auth retrieves the default version of AuthClient.
// Please explicitly pick a version.
func (c *Clientset) Auth() authv1alpha1.AuthV1alpha1Interface {
	if c == nil {
		return nil
	}
	return c.AuthV1alpha1Client
}

// ComponentconfigV1alpha1 retrieves the ComponentconfigV1alpha1Client
func (c *Clientset) ComponentconfigV1alpha1() componentconfigv1alpha1.ComponentconfigV1alpha1Interface {
	if c == nil {
		return nil
	}
	return c.ComponentconfigV1alpha1Client
}

// Deprecated: Componentconfig retrieves the default version of ComponentconfigClient.
// Please explicitly pick a version.
func (c *Clientset) Componentconfig() componentconfigv1alpha1.ComponentconfigV1alpha1Interface {
	if c == nil {
		return nil
	}
	return c.ComponentconfigV1alpha1Client
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.AuthV1alpha1Client, err = authv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.ComponentconfigV1alpha1Client, err = componentconfigv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		glog.Errorf("failed to create the DiscoveryClient: %v", err)
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.AuthV1alpha1Client = authv1alpha1.NewForConfigOrDie(c)
	cs.ComponentconfigV1alpha1Client = componentconfigv1alpha1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.AuthV1alpha1Client = authv1alpha1.New(c)
	cs.ComponentconfigV1alpha1Client = componentconfigv1alpha1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
