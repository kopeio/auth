/*
Copyright 2015 The Kubernetes Authors.

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

import kruntime "k8s.io/apimachinery/pkg/runtime"

//const (
//	defaultRootDir = "/var/lib/kubelet"
//
//	// When these values are updated, also update test/e2e/framework/util.go
//	defaultPodInfraContainerImageName    = "gcr.io/google_containers/pause"
//	defaultPodInfraContainerImageVersion = "3.0"
//	defaultPodInfraContainerImage        = defaultPodInfraContainerImageName +
//		"-" + runtime.GOARCH + ":" +
//		defaultPodInfraContainerImageVersion
//
//	// From pkg/kubelet/rkt/rkt.go to avoid circular import
//	defaultRktAPIServiceEndpoint = "localhost:15441"
//
//	AutoDetectCloudProvider = "auto-detect"
//
//	defaultIPTablesMasqueradeBit = 14
//	defaultIPTablesDropBit       = 15
//)

//var (
//	zeroDuration = metav1.Duration{}
//	// Refer to [Node Allocatable](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/node-allocatable.md) doc for more information.
//	defaultNodeAllocatableEnforcement = []string{"pods"}
//)

func addDefaultingFuncs(scheme *kruntime.Scheme) error {
	RegisterDefaults(scheme)
	return scheme.AddDefaultingFuncs(
		SetDefaults_AuthConfiguration,
	)
}

func SetDefaults_AuthConfiguration(obj *AuthConfiguration) {
	//if obj.BindAddress == "" {
	//	obj.BindAddress = "0.0.0.0"
	//}
	//if obj.HealthzPort == 0 {
	//	obj.HealthzPort = 10249
	//}
	//if obj.HealthzBindAddress == "" {
	//	obj.HealthzBindAddress = "127.0.0.1"
	//}
	//if obj.OOMScoreAdj == nil {
	//	temp := int32(qos.KubeProxyOOMScoreAdj)
	//	obj.OOMScoreAdj = &temp
	//}
	//if obj.ResourceContainer == "" {
	//	obj.ResourceContainer = "/kube-proxy"
	//}
	//if obj.IPTablesSyncPeriod.Duration == 0 {
	//	obj.IPTablesSyncPeriod = metav1.Duration{Duration: 30 * time.Second}
	//}
	//zero := metav1.Duration{}
	//if obj.UDPIdleTimeout == zero {
	//	obj.UDPIdleTimeout = metav1.Duration{Duration: 250 * time.Millisecond}
	//}
	//// If ConntrackMax is set, respect it.
	//if obj.ConntrackMax == 0 {
	//	// If ConntrackMax is *not* set, use per-core scaling.
	//	if obj.ConntrackMaxPerCore == 0 {
	//		obj.ConntrackMaxPerCore = 32 * 1024
	//	}
	//	if obj.ConntrackMin == 0 {
	//		obj.ConntrackMin = 128 * 1024
	//	}
	//}
	//if obj.IPTablesMasqueradeBit == nil {
	//	temp := int32(14)
	//	obj.IPTablesMasqueradeBit = &temp
	//}
	//if obj.ConntrackTCPEstablishedTimeout == zero {
	//	obj.ConntrackTCPEstablishedTimeout = metav1.Duration{Duration: 24 * time.Hour} // 1 day (1/5 default)
	//}
	//if obj.ConntrackTCPCloseWaitTimeout == zero {
	//	// See https://github.com/kubernetes/kubernetes/issues/32551.
	//	//
	//	// CLOSE_WAIT conntrack state occurs when the the Linux kernel
	//	// sees a FIN from the remote server. Note: this is a half-close
	//	// condition that persists as long as the local side keeps the
	//	// socket open. The condition is rare as it is typical in most
	//	// protocols for both sides to issue a close; this typically
	//	// occurs when the local socket is lazily garbage collected.
	//	//
	//	// If the CLOSE_WAIT conntrack entry expires, then FINs from the
	//	// local socket will not be properly SNAT'd and will not reach the
	//	// remote server (if the connection was subject to SNAT). If the
	//	// remote timeouts for FIN_WAIT* states exceed the CLOSE_WAIT
	//	// timeout, then there will be an inconsistency in the state of
	//	// the connection and a new connection reusing the SNAT (src,
	//	// port) pair may be rejected by the remote side with RST. This
	//	// can cause new calls to connect(2) to return with ECONNREFUSED.
	//	//
	//	// We set CLOSE_WAIT to one hour by default to better match
	//	// typical server timeouts.
	//	obj.ConntrackTCPCloseWaitTimeout = metav1.Duration{Duration: 1 * time.Hour}
	//}
}

func boolVar(b bool) *bool {
	return &b
}

//var (
//	defaultCfg = KubeletConfiguration{}
//)
