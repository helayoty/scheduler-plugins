/*
Copyright 2021 The Kubernetes Authors.

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

// +k8s:defaulter-gen=true

package v1beta3

import (
	"strconv"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	schedulerconfigv1beta3 "k8s.io/kube-scheduler/config/v1beta3"
	k8sschedulerconfigv1beta3 "k8s.io/kubernetes/pkg/scheduler/apis/config/v1beta3"
)

var (
	defaultPermitWaitingTimeSeconds      int64 = 60
	defaultDeniedPGExpirationTimeSeconds int64 = 20

	defaultNodeResourcesAllocatableMode = Least

	// defaultResourcesToWeightMap is used to set the default resourceToWeight map for CPU and memory
	// used by the NodeResourcesAllocatable scoring plugin.
	// The base unit for CPU is millicore, while the base using for memory is a byte.
	// The default CPU weight is 1<<20 and default memory weight is 1. That means a millicore
	// has a weighted score equivalent to 1 MiB.
	defaultNodeResourcesAllocatableResourcesToWeightMap = []schedulerconfigv1beta3.ResourceSpec{
		{Name: "cpu", Weight: 1 << 20}, {Name: "memory", Weight: 1},
	}

	// Defaults for TargetLoadPacking plugin

	// Default 1 core CPU usage for containers without requests and limits i.e. Best Effort QoS.
	DefaultRequestsMilliCores int64 = 1000
	// DefaultRequestsMultiplier for containers without limits predicted as 1.5*requests i.e. Burstable QoS class
	DefaultRequestsMultiplier = "1.5"
	// DefaultTargetUtilizationPercent Recommended to keep -10 than desired limit.
	DefaultTargetUtilizationPercent int64 = 40

	// Defaults for LoadVariationRiskBalancing plugin

	// Risk is usually calculated as average (aka. mu) plus standard deviation (aka. sigma).
	// In order to allow customization in the calculation of risk, two parameters are provided:
	// Margin and Sensitivity. Margin is a multiplier of sigma, and Sensitivity is a root power of sigma.
	// For example, Margin=3 and Sensitivity=2 leads to a risk evaluated as: mu + 3 sqrt(sigma).
	// The default value for both parameters is 1, leading to: mu + sigma.
	// DefaultSafeVarianceMargin is one
	DefaultSafeVarianceMargin = 1.0
	// DefaultSafeVarianceSensitivity is one
	DefaultSafeVarianceSensitivity = 1.0

	// Defaults for MetricProviderSpec
	// DefaultMetricProviderType is the Kubernetes metrics server
	DefaultMetricProviderType = KubernetesMetricsServer
	// DefaultInsecureSkipVerify is whether to skip the certificate verification
	DefaultInsecureSkipVerify = true

	defaultResourceSpec = []schedulerconfigv1beta3.ResourceSpec{
		{Name: string(v1.ResourceCPU), Weight: 1},
		{Name: string(v1.ResourceMemory), Weight: 1},
	}
)

// SetDefaultsCoschedulingArgs sets the default parameters for Coscheduling plugin.
func SetDefaultsCoschedulingArgs(obj *CoschedulingArgs) {
	if obj.PermitWaitingTimeSeconds == nil {
		obj.PermitWaitingTimeSeconds = &defaultPermitWaitingTimeSeconds
	}
	if obj.DeniedPGExpirationTimeSeconds == nil {
		obj.DeniedPGExpirationTimeSeconds = &defaultDeniedPGExpirationTimeSeconds
	}
}

// SetDefaultsNodeResourcesAllocatableArgs sets the defaults parameters for NodeResourceAllocatable.
func SetDefaultsNodeResourcesAllocatableArgs(obj *NodeResourcesAllocatableArgs) {
	if len(obj.Resources) == 0 {
		obj.Resources = defaultNodeResourcesAllocatableResourcesToWeightMap
	}

	if obj.Mode == "" {
		obj.Mode = defaultNodeResourcesAllocatableMode
	}
}

// SetDefaultTargetLoadPackingArgs sets the default parameters for TargetLoadPacking plugin
func SetDefaultTargetLoadPackingArgs(args *TargetLoadPackingArgs) {
	if args.DefaultRequests == nil {
		args.DefaultRequests = v1.ResourceList{v1.ResourceCPU: resource.MustParse(
			strconv.FormatInt(DefaultRequestsMilliCores, 10) + "m")}
	}
	if args.DefaultRequestsMultiplier == nil {
		args.DefaultRequestsMultiplier = &DefaultRequestsMultiplier
	}
	if args.TargetUtilization == nil || *args.TargetUtilization <= 0 {
		args.TargetUtilization = &DefaultTargetUtilizationPercent
	}
	if args.WatcherAddress == nil && args.MetricProvider.Type == "" {
		args.MetricProvider.Type = MetricProviderType(DefaultMetricProviderType)
	}
	if args.MetricProvider.Type == Prometheus && args.MetricProvider.InsecureSkipVerify == nil {
		args.MetricProvider.InsecureSkipVerify = &DefaultInsecureSkipVerify
	}
}

// SetDefaultLoadVariationRiskBalancingArgs sets the default parameters for LoadVariationRiskBalancing plugin
func SetDefaultLoadVariationRiskBalancingArgs(args *LoadVariationRiskBalancingArgs) {
	if args.WatcherAddress == nil && args.MetricProvider.Type == "" {
		args.MetricProvider.Type = MetricProviderType(DefaultMetricProviderType)
	}
	if args.SafeVarianceMargin == nil || *args.SafeVarianceMargin < 0 {
		args.SafeVarianceMargin = &DefaultSafeVarianceMargin
	}
	if args.SafeVarianceSensitivity == nil || *args.SafeVarianceSensitivity < 0 {
		args.SafeVarianceSensitivity = &DefaultSafeVarianceSensitivity
	}
	if args.MetricProvider.Type == Prometheus && args.MetricProvider.InsecureSkipVerify == nil {
		args.MetricProvider.InsecureSkipVerify = &DefaultInsecureSkipVerify
	}
}

// SetDefaultsNodeResourceTopologyMatchArgs sets the default parameters for NodeResourceTopologyMatch plugin.
func SetDefaultsNodeResourceTopologyMatchArgs(obj *NodeResourceTopologyMatchArgs) {
	if obj.ScoringStrategy == nil {
		obj.ScoringStrategy = &ScoringStrategy{
			Type:      LeastAllocated,
			Resources: defaultResourceSpec,
		}
	}

	if len(obj.ScoringStrategy.Resources) == 0 {
		// If no resources specified, use the default set.
		obj.ScoringStrategy.Resources = append(obj.ScoringStrategy.Resources, defaultResourceSpec...)
	}

	for i := range obj.ScoringStrategy.Resources {
		if obj.ScoringStrategy.Resources[i].Weight == 0 {
			obj.ScoringStrategy.Resources[i].Weight = 1
		}
	}
}

// PreemptionTolerationArgs reuses SetDefaults_DefaultPreemptionArgs
func SetDefaultsPreemptionTolerationArgs(obj *PreemptionTolerationArgs) {
	k8sschedulerconfigv1beta3.SetDefaults_DefaultPreemptionArgs((*schedulerconfigv1beta3.DefaultPreemptionArgs)(obj))
}