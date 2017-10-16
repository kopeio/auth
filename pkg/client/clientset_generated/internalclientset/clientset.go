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

package internalclientset

import (
	glog "github.com/golang/glog"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
	authinternalversion "kope.io/auth/pkg/client/clientset_generated/internalclientset/typed/auth/internalversion"
	authv1alpha1 "kope.io/auth/pkg/client/clientset_generated/internalclientset/typed/auth/v1alpha1"
	configinternalversion "kope.io/auth/pkg/client/clientset_generated/internalclientset/typed/config/internalversion"
	configv1alpha1 "kope.io/auth/pkg/client/clientset_generated/internalclientset/typed/config/v1alpha1"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	Auth() authinternalversion.AuthInterface
	AuthV1alpha1() authv1alpha1.AuthV1alpha1Interface
	Config() configinternalversion.ConfigInterface
	ConfigV1alpha1() configv1alpha1.ConfigV1alpha1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	auth           *authinternalversion.AuthClient
	authV1alpha1   *authv1alpha1.AuthV1alpha1Client
	config         *configinternalversion.ConfigClient
	configV1alpha1 *configv1alpha1.ConfigV1alpha1Client
}

// Auth retrieves the AuthClient
func (c *Clientset) Auth() authinternalversion.AuthInterface {
	return c.auth
}

// AuthV1alpha1 retrieves the AuthV1alpha1Client
func (c *Clientset) AuthV1alpha1() authv1alpha1.AuthV1alpha1Interface {
	return c.authV1alpha1
}

// Config retrieves the ConfigClient
func (c *Clientset) Config() configinternalversion.ConfigInterface {
	return c.config
}

// ConfigV1alpha1 retrieves the ConfigV1alpha1Client
func (c *Clientset) ConfigV1alpha1() configv1alpha1.ConfigV1alpha1Interface {
	return c.configV1alpha1
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
	cs.auth, err = authinternalversion.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.authV1alpha1, err = authv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.config, err = configinternalversion.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.configV1alpha1, err = configv1alpha1.NewForConfig(&configShallowCopy)
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
	cs.auth = authinternalversion.NewForConfigOrDie(c)
	cs.authV1alpha1 = authv1alpha1.NewForConfigOrDie(c)
	cs.config = configinternalversion.NewForConfigOrDie(c)
	cs.configV1alpha1 = configv1alpha1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.auth = authinternalversion.New(c)
	cs.authV1alpha1 = authv1alpha1.New(c)
	cs.config = configinternalversion.New(c)
	cs.configV1alpha1 = configv1alpha1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
