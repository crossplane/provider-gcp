/*
Copyright 2019 The Crossplane Authors.

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

package container

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/google/go-cmp/cmp"
	container "google.golang.org/api/container/v1beta1"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/crossplaneio/stack-gcp/apis/compute/v1beta1"
	gcp "github.com/crossplaneio/stack-gcp/pkg/clients"
)

const (
	// BootstrapNodePoolName is the name of the node pool that is used to
	// boostrap GKE cluster creation.
	BootstrapNodePoolName = "crossplane-bootstrap"

	// BNPNameFormat is the format for the fully qualified name of the bootstrap node pool.
	BNPNameFormat = "%s/nodePools/%s"

	// ParentFormat is the format for the fully qualified name of a cluster parent.
	ParentFormat = "projects/%s/locations/%s"

	// ClusterNameFormat is the format for the fully qualified name of a cluster.
	ClusterNameFormat = "projects/%s/locations/%s/clusters/%s"
)

// GenerateNodePoolForCreate inserts the default node pool into
// *container.Cluster so that it can be provisioned successfully.
func GenerateNodePoolForCreate(in *container.Cluster) {
	pool := &container.NodePool{
		Name:             BootstrapNodePoolName,
		InitialNodeCount: 0,
	}
	in.NodePools = []*container.NodePool{pool}
}

// GenerateCluster generates *container.Cluster instance from GKEClusterParameters.
func GenerateCluster(in v1beta1.GKEClusterParameters) *container.Cluster { // nolint:gocyclo
	cluster := &container.Cluster{
		ClusterIpv4Cidr:       gcp.StringValue(in.ClusterIpv4Cidr),
		Description:           gcp.StringValue(in.Description),
		EnableKubernetesAlpha: gcp.BoolValue(in.EnableKubernetesAlpha),
		EnableTpu:             gcp.BoolValue(in.EnableTpu),
		InitialClusterVersion: gcp.StringValue(in.InitialClusterVersion),
		LabelFingerprint:      gcp.StringValue(in.LabelFingerprint),
		Locations:             in.Locations,
		LoggingService:        gcp.StringValue(in.LoggingService),
		MonitoringService:     gcp.StringValue(in.MonitoringService),
		Name:                  in.Name,
		Network:               gcp.StringValue(in.Network),
		ResourceLabels:        in.ResourceLabels,
		Subnetwork:            gcp.StringValue(in.Subnetwork),
	}

	GenerateAddonsConfig(in.AddonsConfig, cluster)
	GenerateAuthenticatorGroupsConfig(in.AuthenticatorGroupsConfig, cluster)
	GenerateAutoscaling(in.Autoscaling, cluster)
	GenerateBinaryAuthorization(in.BinaryAuthorization, cluster)
	GenerateDatabaseEncryption(in.DatabaseEncryption, cluster)
	GenerateDefaultMaxPodsConstraint(in.DefaultMaxPodsConstraint, cluster)
	GenerateIPAllocationPolicy(in.IPAllocationPolicy, cluster)
	GenerateLegacyAbac(in.LegacyAbac, cluster)
	GenerateMaintenancePolicy(in.MaintenancePolicy, cluster)
	GenerateMasterAuth(in.MasterAuth, cluster)
	GenerateMasterAuthorizedNetworksConfig(in.MasterAuthorizedNetworksConfig, cluster)
	GenerateNetworkConfig(in.NetworkConfig, cluster)
	GenerateNetworkPolicy(in.NetworkPolicy, cluster)
	GeneratePodSecurityPolicyConfig(in.PodSecurityPolicyConfig, cluster)
	GeneratePrivateClusterConfig(in.PrivateClusterConfig, cluster)
	GenerateResourceUsageExportConfig(in.ResourceUsageExportConfig, cluster)
	GenerateTierSettings(in.TierSettings, cluster)
	GenerateVerticalPodAutoscaling(in.VerticalPodAutoscaling, cluster)
	GenerateWorkloadIdentityConfig(in.WorkloadIdentityConfig, cluster)

	return cluster
}

// GenerateAddonsConfig generates *container.AddonsConfig from *AddonsConfig.
func GenerateAddonsConfig(in *v1beta1.AddonsConfig, cluster *container.Cluster) {
	if in != nil {
		out := &container.AddonsConfig{}
		if in.CloudRunConfig != nil {
			out.CloudRunConfig = &container.CloudRunConfig{
				Disabled: in.CloudRunConfig.Disabled,
			}
		}
		if in.HorizontalPodAutoscaling != nil {
			out.HorizontalPodAutoscaling = &container.HorizontalPodAutoscaling{
				Disabled: in.HorizontalPodAutoscaling.Disabled,
			}
		}
		if in.HTTPLoadBalancing != nil {
			out.HttpLoadBalancing = &container.HttpLoadBalancing{
				Disabled: in.HTTPLoadBalancing.Disabled,
			}
		}
		if in.IstioConfig != nil {
			out.IstioConfig = &container.IstioConfig{
				Auth:     gcp.StringValue(in.IstioConfig.Auth),
				Disabled: gcp.BoolValue(in.IstioConfig.Disabled),
			}
		}
		if in.KubernetesDashboard != nil {
			out.KubernetesDashboard = &container.KubernetesDashboard{
				Disabled: in.KubernetesDashboard.Disabled,
			}
		}
		if in.NetworkPolicyConfig != nil {
			out.NetworkPolicyConfig = &container.NetworkPolicyConfig{
				Disabled: in.NetworkPolicyConfig.Disabled,
			}
		}

		cluster.AddonsConfig = out
	}
}

// GenerateAuthenticatorGroupsConfig generates *container.AuthenticatorGroupsConfig from *AuthenticatorGroupsConfig.
func GenerateAuthenticatorGroupsConfig(in *v1beta1.AuthenticatorGroupsConfig, cluster *container.Cluster) {
	if in != nil {
		out := &container.AuthenticatorGroupsConfig{
			Enabled:       gcp.BoolValue(in.Enabled),
			SecurityGroup: gcp.StringValue(in.SecurityGroup),
		}

		cluster.AuthenticatorGroupsConfig = out
	}
}

// GenerateAutoscaling generates *container.ClusterAutoscaling from *ClusterAutoscaling.
func GenerateAutoscaling(in *v1beta1.ClusterAutoscaling, cluster *container.Cluster) {
	if in != nil {
		out := &container.ClusterAutoscaling{
			AutoprovisioningLocations:  in.AutoprovisioningLocations,
			EnableNodeAutoprovisioning: gcp.BoolValue(in.EnableNodeAutoprovisioning),
		}

		if in.AutoprovisioningNodePoolDefaults != nil {
			out.AutoprovisioningNodePoolDefaults = &container.AutoprovisioningNodePoolDefaults{
				OauthScopes:    in.AutoprovisioningNodePoolDefaults.OauthScopes,
				ServiceAccount: gcp.StringValue(in.AutoprovisioningNodePoolDefaults.ServiceAccount),
			}
		}

		for _, limit := range in.ResourceLimits {
			if limit != nil {
				out.ResourceLimits = append(out.ResourceLimits, &container.ResourceLimit{
					Maximum:      gcp.Int64Value(limit.Maximum),
					Minimum:      gcp.Int64Value(limit.Minimum),
					ResourceType: gcp.StringValue(limit.ResourceType),
				})
			}
		}

		cluster.Autoscaling = out
	}
}

// GenerateBinaryAuthorization generates *container.BinaryAuthorization from *BinaryAuthorization.
func GenerateBinaryAuthorization(in *v1beta1.BinaryAuthorization, cluster *container.Cluster) {
	if in != nil {
		out := &container.BinaryAuthorization{
			Enabled: in.Enabled,
		}

		cluster.BinaryAuthorization = out
	}
}

// GenerateDatabaseEncryption generates *container.DatabaseEncryption from *DatabaseEncryption.
func GenerateDatabaseEncryption(in *v1beta1.DatabaseEncryption, cluster *container.Cluster) {
	if in != nil {
		out := &container.DatabaseEncryption{
			KeyName: gcp.StringValue(in.KeyName),
			State:   gcp.StringValue(in.State),
		}

		cluster.DatabaseEncryption = out
	}
}

// GenerateDefaultMaxPodsConstraint generates *container.MaxPodsConstraint from *DefaultMaxPodsConstraint.
func GenerateDefaultMaxPodsConstraint(in *v1beta1.MaxPodsConstraint, cluster *container.Cluster) {
	if in != nil {
		out := &container.MaxPodsConstraint{
			MaxPodsPerNode: in.MaxPodsPerNode,
		}

		cluster.DefaultMaxPodsConstraint = out
	}
}

// GenerateIPAllocationPolicy generates *container.MaxPodsConstraint from *IpAllocationPolicy.
func GenerateIPAllocationPolicy(in *v1beta1.IPAllocationPolicy, cluster *container.Cluster) {
	if in != nil {
		out := &container.IPAllocationPolicy{
			AllowRouteOverlap:          gcp.BoolValue(in.AllowRouteOverlap),
			ClusterIpv4CidrBlock:       gcp.StringValue(in.ClusterIpv4CidrBlock),
			ClusterSecondaryRangeName:  gcp.StringValue(in.ClusterSecondaryRangeName),
			CreateSubnetwork:           gcp.BoolValue(in.CreateSubnetwork),
			NodeIpv4CidrBlock:          gcp.StringValue(in.NodeIpv4CidrBlock),
			ServicesIpv4CidrBlock:      gcp.StringValue(in.ServicesIpv4CidrBlock),
			ServicesSecondaryRangeName: gcp.StringValue(in.DeepCopy().ServicesSecondaryRangeName),
			SubnetworkName:             gcp.StringValue(in.SubnetworkName),
			TpuIpv4CidrBlock:           gcp.StringValue(in.TpuIpv4CidrBlock),
			UseIpAliases:               gcp.BoolValue(in.UseIPAliases),
		}

		cluster.IpAllocationPolicy = out
	}
}

// GenerateLegacyAbac generates *container.LegacyAbac from *LegacyAbac.
func GenerateLegacyAbac(in *v1beta1.LegacyAbac, cluster *container.Cluster) {
	if in != nil {
		out := &container.LegacyAbac{
			Enabled: in.Enabled,
		}

		cluster.LegacyAbac = out
	}
}

// GenerateMaintenancePolicy generates *container.MaintenancePolicy from *MaintenancePolicy.
func GenerateMaintenancePolicy(in *v1beta1.MaintenancePolicySpec, cluster *container.Cluster) {
	if in != nil {
		out := &container.MaintenancePolicy{
			Window: &container.MaintenanceWindow{
				DailyMaintenanceWindow: &container.DailyMaintenanceWindow{
					StartTime: in.Window.DailyMaintenanceWindow.StartTime,
				},
			},
		}

		cluster.MaintenancePolicy = out
	}
}

// GenerateMasterAuth generates *container.MasterAuth from *MasterAuth.
func GenerateMasterAuth(in *v1beta1.MasterAuth, cluster *container.Cluster) {
	if in != nil {
		out := &container.MasterAuth{
			Password: gcp.StringValue(in.Password),
			Username: gcp.StringValue(in.Username),
		}

		if in.ClientCertificateConfig != nil {
			out.ClientCertificateConfig = &container.ClientCertificateConfig{
				IssueClientCertificate: in.ClientCertificateConfig.IssueClientCertificate,
			}
		}

		cluster.MasterAuth = out
	}
}

// GenerateMasterAuthorizedNetworksConfig generates *container.MasterAuthorizedNetworksConfig from *MasterAuthorizedNetworksConfig.
func GenerateMasterAuthorizedNetworksConfig(in *v1beta1.MasterAuthorizedNetworksConfig, cluster *container.Cluster) {
	if in != nil {
		out := &container.MasterAuthorizedNetworksConfig{
			Enabled: gcp.BoolValue(in.Enabled),
		}

		for _, cidr := range in.CidrBlocks {
			if cidr != nil {
				out.CidrBlocks = append(out.CidrBlocks, &container.CidrBlock{
					CidrBlock:   cidr.CidrBlock,
					DisplayName: gcp.StringValue(cidr.DisplayName),
				})
			}
		}

		cluster.MasterAuthorizedNetworksConfig = out
	}
}

// GenerateNetworkConfig generates *container.NetworkConfig from *NetworkConfig.
func GenerateNetworkConfig(in *v1beta1.NetworkConfigSpec, cluster *container.Cluster) {
	if in != nil {
		out := &container.NetworkConfig{
			EnableIntraNodeVisibility: in.EnableIntraNodeVisibility,
		}

		cluster.NetworkConfig = out
	}
}

// GenerateNetworkPolicy generates *container.NetworkPolicy from *NetworkPolicy.
func GenerateNetworkPolicy(in *v1beta1.NetworkPolicy, cluster *container.Cluster) {
	if in != nil {
		out := &container.NetworkPolicy{
			Enabled:  gcp.BoolValue(in.Enabled),
			Provider: gcp.StringValue(in.Provider),
		}

		cluster.NetworkPolicy = out
	}
}

// GeneratePodSecurityPolicyConfig generates *container.PodSecurityPolicyConfig from *PodSecurityPolicyConfig.
func GeneratePodSecurityPolicyConfig(in *v1beta1.PodSecurityPolicyConfig, cluster *container.Cluster) {
	if in != nil {
		out := &container.PodSecurityPolicyConfig{
			Enabled: in.Enabled,
		}

		cluster.PodSecurityPolicyConfig = out
	}
}

// GeneratePrivateClusterConfig generates *container.PrivateClusterConfig from *PrivateClusterConfig.
func GeneratePrivateClusterConfig(in *v1beta1.PrivateClusterConfigSpec, cluster *container.Cluster) {
	if in != nil {
		out := &container.PrivateClusterConfig{
			EnablePeeringRouteSharing: gcp.BoolValue(in.EnablePeeringRouteSharing),
			EnablePrivateEndpoint:     gcp.BoolValue(in.EnablePrivateEndpoint),
			EnablePrivateNodes:        gcp.BoolValue(in.EnablePrivateNodes),
			MasterIpv4CidrBlock:       gcp.StringValue(in.MasterIpv4CidrBlock),
		}

		cluster.PrivateClusterConfig = out
	}
}

// GenerateResourceUsageExportConfig generates *container.ResourceUsageExportConfig from *ResourceUsageExportConfig.
func GenerateResourceUsageExportConfig(in *v1beta1.ResourceUsageExportConfig, cluster *container.Cluster) {
	if in != nil {
		out := &container.ResourceUsageExportConfig{
			EnableNetworkEgressMetering: gcp.BoolValue(in.EnableNetworkEgressMetering),
		}

		if in.BigqueryDestination != nil {
			out.BigqueryDestination = &container.BigQueryDestination{
				DatasetId: in.BigqueryDestination.DatasetID,
			}
		}

		if in.ConsumptionMeteringConfig != nil {
			out.ConsumptionMeteringConfig = &container.ConsumptionMeteringConfig{
				Enabled: in.ConsumptionMeteringConfig.Enabled,
			}
		}

		cluster.ResourceUsageExportConfig = out
	}
}

// GenerateTierSettings generates *container.TierSettings from *TierSettings.
func GenerateTierSettings(in *v1beta1.TierSettings, cluster *container.Cluster) {
	if in != nil {
		out := &container.TierSettings{
			Tier: in.Tier,
		}

		cluster.TierSettings = out
	}
}

// GenerateVerticalPodAutoscaling generates *container.VerticalPodAutoscaling from *VerticalPodAutoscaling.
func GenerateVerticalPodAutoscaling(in *v1beta1.VerticalPodAutoscaling, cluster *container.Cluster) {
	if in != nil {
		out := &container.VerticalPodAutoscaling{
			Enabled: in.Enabled,
		}

		cluster.VerticalPodAutoscaling = out
	}
}

// GenerateWorkloadIdentityConfig generates *container.WorkloadIdentityConfig from *WorkloadIdentityConfig.
func GenerateWorkloadIdentityConfig(in *v1beta1.WorkloadIdentityConfig, cluster *container.Cluster) {
	if in != nil {
		out := &container.WorkloadIdentityConfig{
			IdentityNamespace: in.IdentityNamespace,
		}

		cluster.WorkloadIdentityConfig = out
	}
}

// GenerateObservation produces GKEClusterObservation object from *sqladmin.DatabaseInstance object.
func GenerateObservation(in container.Cluster) v1beta1.GKEClusterObservation { // nolint:gocyclo
	o := v1beta1.GKEClusterObservation{
		CreateTime:           in.CreateTime,
		CurrentMasterVersion: in.CurrentMasterVersion,
		CurrentNodeCount:     in.CurrentNodeCount,
		CurrentNodeVersion:   in.CurrentNodeVersion,
		Endpoint:             in.Endpoint,
		ExpireTime:           in.ExpireTime,
		Location:             in.Location,
		NodeIpv4CidrSize:     in.NodeIpv4CidrSize,
		SelfLink:             in.SelfLink,
		ServicesIpv4Cidr:     in.ServicesIpv4Cidr,
		Status:               in.Status,
		StatusMessage:        in.StatusMessage,
		TpuIpv4CidrBlock:     in.TpuIpv4CidrBlock,
		Zone:                 in.Zone,
	}

	if in.MaintenancePolicy != nil {
		if in.MaintenancePolicy.Window != nil {
			if in.MaintenancePolicy.Window.DailyMaintenanceWindow != nil {
				o.MaintenancePolicy = &v1beta1.MaintenancePolicyStatus{
					Window: v1beta1.MaintenanceWindowStatus{
						DailyMaintenanceWindow: v1beta1.DailyMaintenanceWindowStatus{
							Duration: in.MaintenancePolicy.Window.DailyMaintenanceWindow.Duration,
						},
					},
				}
			}
		}
	}

	if in.NetworkConfig != nil {
		o.NetworkConfig = &v1beta1.NetworkConfigStatus{
			Network:    in.NetworkConfig.Network,
			Subnetwork: in.NetworkConfig.Subnetwork,
		}
	}

	if in.PrivateClusterConfig != nil {
		o.PrivateClusterConfig = &v1beta1.PrivateClusterConfigStatus{
			PrivateEndpoint: in.PrivateClusterConfig.PrivateEndpoint,
			PublicEndpoint:  in.PrivateClusterConfig.PublicEndpoint,
		}
	}

	for _, condition := range in.Conditions {
		if condition != nil {
			o.Conditions = append(o.Conditions, &v1beta1.StatusCondition{
				Code:    condition.Code,
				Message: condition.Message,
			})
		}
	}

	for _, nodePool := range in.NodePools {
		if nodePool != nil {
			conditions := []*v1beta1.StatusCondition{}
			for _, condition := range nodePool.Conditions {
				if condition != nil {
					conditions = append(conditions, &v1beta1.StatusCondition{
						Code:    condition.Code,
						Message: condition.Message,
					})
				}
			}
			o.NodePools = append(o.NodePools, &v1beta1.NodePoolClusterStatus{
				Conditions:        conditions,
				InstanceGroupUrls: nodePool.InstanceGroupUrls,
				Name:              nodePool.Name,
				PodIpv4CidrSize:   nodePool.PodIpv4CidrSize,
				SelfLink:          nodePool.SelfLink,
				Status:            nodePool.Status,
				StatusMessage:     nodePool.StatusMessage,
				Version:           nodePool.Version,
			})
		}
	}

	return o
}

// LateInitializeSpec fills unassigned fields with the values in container.Cluster object.
func LateInitializeSpec(spec *v1beta1.GKEClusterParameters, in container.Cluster) { // nolint:gocyclo
	if in.AddonsConfig != nil {
		if spec.AddonsConfig == nil {
			spec.AddonsConfig = &v1beta1.AddonsConfig{}
		}
		if spec.AddonsConfig.CloudRunConfig == nil && in.AddonsConfig.CloudRunConfig != nil {
			spec.AddonsConfig.CloudRunConfig = &v1beta1.CloudRunConfig{
				Disabled: in.AddonsConfig.CloudRunConfig.Disabled,
			}
		}
		if spec.AddonsConfig.HorizontalPodAutoscaling == nil && in.AddonsConfig.HorizontalPodAutoscaling != nil {
			spec.AddonsConfig.HorizontalPodAutoscaling = &v1beta1.HorizontalPodAutoscaling{
				Disabled: in.AddonsConfig.HorizontalPodAutoscaling.Disabled,
			}
		}
		if spec.AddonsConfig.HTTPLoadBalancing == nil && in.AddonsConfig.HttpLoadBalancing != nil {
			spec.AddonsConfig.HTTPLoadBalancing = &v1beta1.HTTPLoadBalancing{
				Disabled: in.AddonsConfig.HttpLoadBalancing.Disabled,
			}
		}
		if in.AddonsConfig.IstioConfig != nil {
			if spec.AddonsConfig.IstioConfig == nil {
				spec.AddonsConfig.IstioConfig = &v1beta1.IstioConfig{}
			}
			spec.AddonsConfig.IstioConfig.Auth = gcp.LateInitializeString(spec.AddonsConfig.IstioConfig.Auth, in.AddonsConfig.IstioConfig.Auth)
			spec.AddonsConfig.IstioConfig.Disabled = gcp.LateInitializeBool(spec.AddonsConfig.IstioConfig.Disabled, in.AddonsConfig.IstioConfig.Disabled)
		}
		if spec.AddonsConfig.KubernetesDashboard == nil && in.AddonsConfig.KubernetesDashboard != nil {
			spec.AddonsConfig.KubernetesDashboard = &v1beta1.KubernetesDashboard{
				Disabled: in.AddonsConfig.KubernetesDashboard.Disabled,
			}
		}
		if spec.AddonsConfig.NetworkPolicyConfig == nil && in.AddonsConfig.NetworkPolicyConfig != nil {
			spec.AddonsConfig.NetworkPolicyConfig = &v1beta1.NetworkPolicyConfig{
				Disabled: in.AddonsConfig.NetworkPolicyConfig.Disabled,
			}
		}
	}

	if in.AuthenticatorGroupsConfig != nil {
		if spec.AuthenticatorGroupsConfig == nil {
			spec.AuthenticatorGroupsConfig = &v1beta1.AuthenticatorGroupsConfig{}
		}
		spec.AuthenticatorGroupsConfig.Enabled = gcp.LateInitializeBool(spec.AuthenticatorGroupsConfig.Enabled, in.AuthenticatorGroupsConfig.Enabled)
		spec.AuthenticatorGroupsConfig.SecurityGroup = gcp.LateInitializeString(spec.AuthenticatorGroupsConfig.SecurityGroup, in.AuthenticatorGroupsConfig.SecurityGroup)
	}

	if in.Autoscaling != nil {
		if spec.Autoscaling == nil {
			spec.Autoscaling = &v1beta1.ClusterAutoscaling{}
		}
		spec.Autoscaling.AutoprovisioningLocations = gcp.LateInitializeStringSlice(spec.Autoscaling.AutoprovisioningLocations, in.Autoscaling.AutoprovisioningLocations)
		if in.Autoscaling.AutoprovisioningNodePoolDefaults != nil {
			if spec.Autoscaling.AutoprovisioningNodePoolDefaults == nil {
				spec.Autoscaling.AutoprovisioningNodePoolDefaults = &v1beta1.AutoprovisioningNodePoolDefaults{}
			}
			spec.Autoscaling.AutoprovisioningNodePoolDefaults.OauthScopes = gcp.LateInitializeStringSlice(spec.Autoscaling.AutoprovisioningNodePoolDefaults.OauthScopes, in.Autoscaling.AutoprovisioningNodePoolDefaults.OauthScopes)
			spec.Autoscaling.AutoprovisioningNodePoolDefaults.ServiceAccount = gcp.LateInitializeString(spec.Autoscaling.AutoprovisioningNodePoolDefaults.ServiceAccount, in.Autoscaling.AutoprovisioningNodePoolDefaults.ServiceAccount)
		}
		spec.Autoscaling.EnableNodeAutoprovisioning = gcp.LateInitializeBool(spec.Autoscaling.EnableNodeAutoprovisioning, in.Autoscaling.EnableNodeAutoprovisioning)
		if len(in.Autoscaling.ResourceLimits) != 0 && len(spec.Autoscaling.ResourceLimits) == 0 {
			spec.Autoscaling.ResourceLimits = make([]*v1beta1.ResourceLimit, len(in.Autoscaling.ResourceLimits))
			for i, limit := range in.Autoscaling.ResourceLimits {
				spec.Autoscaling.ResourceLimits[i] = &v1beta1.ResourceLimit{
					Maximum:      &limit.Maximum,
					Minimum:      &limit.Minimum,
					ResourceType: &limit.ResourceType,
				}
			}
		}
	}

	if spec.BinaryAuthorization == nil && in.BinaryAuthorization != nil {
		spec.BinaryAuthorization = &v1beta1.BinaryAuthorization{
			Enabled: in.BinaryAuthorization.Enabled,
		}
	}

	spec.ClusterIpv4Cidr = gcp.LateInitializeString(spec.ClusterIpv4Cidr, in.ClusterIpv4Cidr)

	if in.DatabaseEncryption != nil {
		if spec.DatabaseEncryption == nil {
			spec.DatabaseEncryption = &v1beta1.DatabaseEncryption{}
		}
		spec.DatabaseEncryption.KeyName = gcp.LateInitializeString(spec.DatabaseEncryption.KeyName, in.DatabaseEncryption.KeyName)
		spec.DatabaseEncryption.State = gcp.LateInitializeString(spec.DatabaseEncryption.State, in.DatabaseEncryption.State)
	}

	if spec.DefaultMaxPodsConstraint == nil && in.DefaultMaxPodsConstraint != nil {
		spec.DefaultMaxPodsConstraint = &v1beta1.MaxPodsConstraint{
			MaxPodsPerNode: in.DefaultMaxPodsConstraint.MaxPodsPerNode,
		}
	}

	if spec.Description == nil {
		spec.Description = &in.Description
	}

	spec.EnableKubernetesAlpha = gcp.LateInitializeBool(spec.EnableKubernetesAlpha, in.EnableKubernetesAlpha)
	spec.EnableTpu = gcp.LateInitializeBool(spec.EnableTpu, in.EnableTpu)
	spec.InitialClusterVersion = gcp.LateInitializeString(spec.InitialClusterVersion, in.InitialClusterVersion)

	if in.IpAllocationPolicy != nil {
		if spec.IPAllocationPolicy == nil {
			spec.IPAllocationPolicy = &v1beta1.IPAllocationPolicy{}
		}
		spec.IPAllocationPolicy.AllowRouteOverlap = gcp.LateInitializeBool(spec.IPAllocationPolicy.AllowRouteOverlap, in.IpAllocationPolicy.AllowRouteOverlap)
		spec.IPAllocationPolicy.ClusterIpv4CidrBlock = gcp.LateInitializeString(spec.IPAllocationPolicy.ClusterIpv4CidrBlock, in.IpAllocationPolicy.ClusterIpv4CidrBlock)
		spec.IPAllocationPolicy.ClusterSecondaryRangeName = gcp.LateInitializeString(spec.IPAllocationPolicy.ClusterSecondaryRangeName, in.IpAllocationPolicy.ClusterSecondaryRangeName)
		spec.IPAllocationPolicy.CreateSubnetwork = gcp.LateInitializeBool(spec.IPAllocationPolicy.CreateSubnetwork, in.IpAllocationPolicy.CreateSubnetwork)
		spec.IPAllocationPolicy.NodeIpv4CidrBlock = gcp.LateInitializeString(spec.IPAllocationPolicy.NodeIpv4CidrBlock, in.IpAllocationPolicy.NodeIpv4CidrBlock)
		spec.IPAllocationPolicy.ServicesIpv4CidrBlock = gcp.LateInitializeString(spec.IPAllocationPolicy.ServicesIpv4CidrBlock, in.IpAllocationPolicy.ServicesIpv4CidrBlock)
		spec.IPAllocationPolicy.SubnetworkName = gcp.LateInitializeString(spec.IPAllocationPolicy.SubnetworkName, in.IpAllocationPolicy.SubnetworkName)
		spec.IPAllocationPolicy.TpuIpv4CidrBlock = gcp.LateInitializeString(spec.IPAllocationPolicy.TpuIpv4CidrBlock, in.IpAllocationPolicy.TpuIpv4CidrBlock)
		spec.IPAllocationPolicy.UseIPAliases = gcp.LateInitializeBool(spec.IPAllocationPolicy.UseIPAliases, in.IpAllocationPolicy.UseIpAliases)
	}

	spec.LabelFingerprint = gcp.LateInitializeString(spec.LabelFingerprint, in.LabelFingerprint)

	if spec.LegacyAbac == nil && in.LegacyAbac != nil {
		spec.LegacyAbac = &v1beta1.LegacyAbac{
			Enabled: in.LegacyAbac.Enabled,
		}
	}

	spec.Locations = gcp.LateInitializeStringSlice(spec.Locations, in.Locations)
	spec.LoggingService = gcp.LateInitializeString(spec.LoggingService, in.LoggingService)

	if spec.MaintenancePolicy == nil && in.MaintenancePolicy != nil {
		if in.MaintenancePolicy.Window != nil {
			if in.MaintenancePolicy.Window.DailyMaintenanceWindow != nil {
				spec.MaintenancePolicy = &v1beta1.MaintenancePolicySpec{
					Window: v1beta1.MaintenanceWindowSpec{
						DailyMaintenanceWindow: v1beta1.DailyMaintenanceWindowSpec{
							StartTime: in.MaintenancePolicy.Window.DailyMaintenanceWindow.StartTime,
						},
					},
				}
			}
		}
	}

	if in.MasterAuth != nil {
		if spec.MasterAuth == nil {
			spec.MasterAuth = &v1beta1.MasterAuth{}
		}
		if spec.MasterAuth.ClientCertificateConfig == nil && in.MasterAuth.ClientCertificateConfig != nil {
			spec.MasterAuth.ClientCertificateConfig = &v1beta1.ClientCertificateConfig{
				IssueClientCertificate: in.MasterAuth.ClientCertificateConfig.IssueClientCertificate,
			}
		}
		spec.MasterAuth.Password = gcp.LateInitializeString(spec.MasterAuth.Password, in.MasterAuth.Password)
		spec.MasterAuth.Username = gcp.LateInitializeString(spec.MasterAuth.Username, in.MasterAuth.Username)
	}

	if in.MasterAuthorizedNetworksConfig != nil {
		if spec.MasterAuthorizedNetworksConfig == nil {
			spec.MasterAuthorizedNetworksConfig = &v1beta1.MasterAuthorizedNetworksConfig{}
		}
		if len(in.MasterAuthorizedNetworksConfig.CidrBlocks) != 0 && len(spec.MasterAuthorizedNetworksConfig.CidrBlocks) == 0 {
			spec.MasterAuthorizedNetworksConfig.CidrBlocks = make([]*v1beta1.CidrBlock, len(in.MasterAuthorizedNetworksConfig.CidrBlocks))
			for i, block := range in.MasterAuthorizedNetworksConfig.CidrBlocks {
				spec.MasterAuthorizedNetworksConfig.CidrBlocks[i] = &v1beta1.CidrBlock{
					CidrBlock:   block.CidrBlock,
					DisplayName: &block.DisplayName,
				}
			}
		}
		spec.MasterAuthorizedNetworksConfig.Enabled = gcp.LateInitializeBool(spec.MasterAuthorizedNetworksConfig.Enabled, in.MasterAuthorizedNetworksConfig.Enabled)
	}

	spec.MonitoringService = gcp.LateInitializeString(spec.MonitoringService, in.MonitoringService)
	spec.Network = gcp.LateInitializeString(spec.Network, in.Network)

	if spec.NetworkConfig == nil && in.NetworkConfig != nil {
		spec.NetworkConfig = &v1beta1.NetworkConfigSpec{
			EnableIntraNodeVisibility: in.NetworkConfig.EnableIntraNodeVisibility,
		}
	}

	if in.NetworkPolicy != nil {
		if spec.NetworkPolicy == nil {
			spec.NetworkPolicy = &v1beta1.NetworkPolicy{}
		}
		spec.NetworkPolicy.Enabled = gcp.LateInitializeBool(spec.NetworkPolicy.Enabled, in.NetworkPolicy.Enabled)
		spec.NetworkPolicy.Provider = gcp.LateInitializeString(spec.NetworkPolicy.Provider, in.NetworkPolicy.Provider)
	}

	if spec.PodSecurityPolicyConfig == nil && in.PodSecurityPolicyConfig != nil {
		spec.PodSecurityPolicyConfig = &v1beta1.PodSecurityPolicyConfig{
			Enabled: in.PodSecurityPolicyConfig.Enabled,
		}
	}

	if in.PrivateClusterConfig != nil {
		if spec.PrivateClusterConfig == nil {
			spec.PrivateClusterConfig = &v1beta1.PrivateClusterConfigSpec{}
		}
		spec.PrivateClusterConfig.EnablePeeringRouteSharing = gcp.LateInitializeBool(spec.PrivateClusterConfig.EnablePeeringRouteSharing, in.PrivateClusterConfig.EnablePeeringRouteSharing)
		spec.PrivateClusterConfig.EnablePrivateEndpoint = gcp.LateInitializeBool(spec.PrivateClusterConfig.EnablePrivateEndpoint, in.PrivateClusterConfig.EnablePrivateEndpoint)
		spec.PrivateClusterConfig.EnablePrivateNodes = gcp.LateInitializeBool(spec.PrivateClusterConfig.EnablePrivateNodes, in.PrivateClusterConfig.EnablePrivateNodes)
		spec.PrivateClusterConfig.MasterIpv4CidrBlock = gcp.LateInitializeString(spec.PrivateClusterConfig.MasterIpv4CidrBlock, in.PrivateClusterConfig.MasterIpv4CidrBlock)
	}

	spec.ResourceLabels = gcp.LateInitializeStringMap(spec.ResourceLabels, in.ResourceLabels)

	if in.ResourceUsageExportConfig != nil {
		if spec.ResourceUsageExportConfig == nil {
			spec.ResourceUsageExportConfig = &v1beta1.ResourceUsageExportConfig{}
		}
		if spec.ResourceUsageExportConfig.BigqueryDestination == nil && in.ResourceUsageExportConfig.BigqueryDestination != nil {
			spec.ResourceUsageExportConfig.BigqueryDestination = &v1beta1.BigQueryDestination{
				DatasetID: in.ResourceUsageExportConfig.BigqueryDestination.DatasetId,
			}
		}
		if spec.ResourceUsageExportConfig.ConsumptionMeteringConfig == nil && in.ResourceUsageExportConfig.ConsumptionMeteringConfig != nil {
			spec.ResourceUsageExportConfig.ConsumptionMeteringConfig = &v1beta1.ConsumptionMeteringConfig{
				Enabled: in.ResourceUsageExportConfig.ConsumptionMeteringConfig.Enabled,
			}
		}
		spec.ResourceUsageExportConfig.EnableNetworkEgressMetering = gcp.LateInitializeBool(spec.ResourceUsageExportConfig.EnableNetworkEgressMetering, in.ResourceUsageExportConfig.EnableNetworkEgressMetering)
	}

	spec.Subnetwork = gcp.LateInitializeString(spec.Subnetwork, in.Subnetwork)

	if spec.TierSettings == nil && in.TierSettings != nil {
		spec.TierSettings = &v1beta1.TierSettings{
			Tier: in.TierSettings.Tier,
		}
	}

	if spec.VerticalPodAutoscaling == nil && in.VerticalPodAutoscaling != nil {
		spec.VerticalPodAutoscaling = &v1beta1.VerticalPodAutoscaling{
			Enabled: in.VerticalPodAutoscaling.Enabled,
		}
	}

	if spec.WorkloadIdentityConfig == nil && in.WorkloadIdentityConfig != nil {
		spec.WorkloadIdentityConfig = &v1beta1.WorkloadIdentityConfig{
			IdentityNamespace: in.WorkloadIdentityConfig.IdentityNamespace,
		}
	}

}

// UpdateFn returns a function that updates a cluster.
type UpdateFn func(container.Service, context.Context, string) (*container.Operation, error)

// NewAddonsConfigUpdate returns a function that updates the AddonsConfig of a cluster.
func NewAddonsConfigUpdate(in *v1beta1.AddonsConfig) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		out := &container.Cluster{}
		GenerateAddonsConfig(in, out)
		update := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredAddonsConfig: out.AddonsConfig,
			},
		}
		return s.Projects.Locations.Clusters.Update(name, update).Context(ctx).Do()
	}
}

// NewAutoscalingUpdate returns a function that updates the Autoscaling of a cluster.
func NewAutoscalingUpdate(in *v1beta1.ClusterAutoscaling) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		out := &container.Cluster{}
		GenerateAutoscaling(in, out)
		update := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredClusterAutoscaling: out.Autoscaling,
			},
		}
		return s.Projects.Locations.Clusters.Update(name, update).Context(ctx).Do()
	}
}

// NewBinaryAuthorizationUpdate returns a function that updates the BinaryAuthorization of a cluster.
func NewBinaryAuthorizationUpdate(in *v1beta1.BinaryAuthorization) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		out := &container.Cluster{}
		GenerateBinaryAuthorization(in, out)
		update := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredBinaryAuthorization: out.BinaryAuthorization,
			},
		}
		return s.Projects.Locations.Clusters.Update(name, update).Context(ctx).Do()
	}
}

// NewDatabaseEncryptionUpdate returns a function that updates the DatabaseEncryption of a cluster.
func NewDatabaseEncryptionUpdate(in *v1beta1.DatabaseEncryption) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		out := &container.Cluster{}
		GenerateDatabaseEncryption(in, out)
		update := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredDatabaseEncryption: out.DatabaseEncryption,
			},
		}
		return s.Projects.Locations.Clusters.Update(name, update).Context(ctx).Do()
	}
}

// NewLegacyAbacUpdate returns a function that updates the LegacyAbac of a cluster.
func NewLegacyAbacUpdate(in *v1beta1.LegacyAbac) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		out := &container.Cluster{}
		GenerateLegacyAbac(in, out)
		update := &container.SetLegacyAbacRequest{
			Enabled: out.LegacyAbac.Enabled,
		}
		return s.Projects.Locations.Clusters.SetLegacyAbac(name, update).Context(ctx).Do()
	}
}

// NewLocationsUpdate returns a function that updates the Locations of a cluster.
func NewLocationsUpdate(in []string) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		update := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredLocations: in,
			},
		}
		return s.Projects.Locations.Clusters.Update(name, update).Context(ctx).Do()
	}
}

// NewLoggingServiceUpdate returns a function that updates the LoggingService of a cluster.
func NewLoggingServiceUpdate(in *string) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		update := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredLoggingService: gcp.StringValue(in),
			},
		}
		return s.Projects.Locations.Clusters.Update(name, update).Context(ctx).Do()
	}
}

// NewMaintenancePolicyUpdate returns a function that updates the MaintenancePolicy of a cluster.
func NewMaintenancePolicyUpdate(in *v1beta1.MaintenancePolicySpec) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		out := &container.Cluster{}
		GenerateMaintenancePolicy(in, out)
		update := &container.SetMaintenancePolicyRequest{
			MaintenancePolicy: out.MaintenancePolicy,
		}
		return s.Projects.Locations.Clusters.SetMaintenancePolicy(name, update).Context(ctx).Do()
	}
}

// NewMasterAuthUpdate returns a function that updates the MasterAuth of a cluster.
func NewMasterAuthUpdate(in *v1beta1.MasterAuth) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		out := &container.Cluster{}
		GenerateMasterAuth(in, out)
		update := &container.SetMasterAuthRequest{
			// TODO(hasheddan): need to set Action here?
			Update: out.MasterAuth,
		}
		return s.Projects.Locations.Clusters.SetMasterAuth(name, update).Context(ctx).Do()
	}
}

// NewMasterAuthorizedNetworksConfigUpdate returns a function that updates the MasterAuthorizedNetworksConfig of a cluster.
func NewMasterAuthorizedNetworksConfigUpdate(in *v1beta1.MasterAuthorizedNetworksConfig) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		out := &container.Cluster{}
		update := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredMasterAuthorizedNetworksConfig: out.MasterAuthorizedNetworksConfig,
			},
		}
		return s.Projects.Locations.Clusters.Update(name, update).Context(ctx).Do()
	}
}

// NewMonitoringServiceUpdate returns a function that updates the MonitoringService of a cluster.
func NewMonitoringServiceUpdate(in *string) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		update := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredMonitoringService: gcp.StringValue(in),
			},
		}
		return s.Projects.Locations.Clusters.Update(name, update).Context(ctx).Do()
	}
}

// NewNetworkConfigUpdate returns a function that updates the NetworkConfig of a cluster.
func NewNetworkConfigUpdate(in *v1beta1.NetworkConfigSpec) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		out := &container.Cluster{}
		GenerateNetworkConfig(in, out)
		update := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredIntraNodeVisibilityConfig: &container.IntraNodeVisibilityConfig{
					Enabled: out.NetworkConfig.EnableIntraNodeVisibility,
				},
			},
		}
		return s.Projects.Locations.Clusters.Update(name, update).Context(ctx).Do()
	}
}

// NewNetworkPolicyUpdate returns a function that updates the NetworkPolicy of a cluster.
func NewNetworkPolicyUpdate(in *v1beta1.NetworkPolicy) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		out := &container.Cluster{}
		GenerateNetworkPolicy(in, out)
		update := &container.SetNetworkPolicyRequest{
			NetworkPolicy: out.NetworkPolicy,
		}
		return s.Projects.Locations.Clusters.SetNetworkPolicy(name, update).Context(ctx).Do()
	}
}

// NewPodSecurityPolicyConfigUpdate returns a function that updates the PodSecurityPolicyConfig of a cluster.
func NewPodSecurityPolicyConfigUpdate(in *v1beta1.PodSecurityPolicyConfig) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		out := &container.Cluster{}
		GeneratePodSecurityPolicyConfig(in, out)
		update := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredPodSecurityPolicyConfig: out.PodSecurityPolicyConfig,
			},
		}
		return s.Projects.Locations.Clusters.Update(name, update).Context(ctx).Do()
	}
}

// NewPrivateClusterConfigUpdate returns a function that updates the PrivateClusterConfig of a cluster.
func NewPrivateClusterConfigUpdate(in *v1beta1.PrivateClusterConfigSpec) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		out := &container.Cluster{}
		GeneratePrivateClusterConfig(in, out)
		update := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredPrivateClusterConfig: out.PrivateClusterConfig,
			},
		}
		return s.Projects.Locations.Clusters.Update(name, update).Context(ctx).Do()
	}
}

// NewResourceLabelsUpdate returns a function that updates the ResourceLabels of a cluster.
func NewResourceLabelsUpdate(in map[string]string) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		update := &container.SetLabelsRequest{
			ResourceLabels: in,
		}
		return s.Projects.Locations.Clusters.SetResourceLabels(name, update).Context(ctx).Do()
	}
}

// NewResourceUsageExportConfigUpdate returns a function that updates the ResourceUsageExportConfig of a cluster.
func NewResourceUsageExportConfigUpdate(in *v1beta1.ResourceUsageExportConfig) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		out := &container.Cluster{}
		GenerateResourceUsageExportConfig(in, out)
		update := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredResourceUsageExportConfig: out.ResourceUsageExportConfig,
			},
		}
		return s.Projects.Locations.Clusters.Update(name, update).Context(ctx).Do()
	}
}

// NewVerticalPodAutoscalingUpdate returns a function that updates the VerticalPodAutoscaling of a cluster.
func NewVerticalPodAutoscalingUpdate(in *v1beta1.VerticalPodAutoscaling) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		out := &container.Cluster{}
		GenerateVerticalPodAutoscaling(in, out)
		update := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredVerticalPodAutoscaling: out.VerticalPodAutoscaling,
			},
		}
		return s.Projects.Locations.Clusters.Update(name, update).Context(ctx).Do()
	}
}

// NewWorkloadIdentityConfigUpdate returns a function that updates the WorkloadIdentityConfig of a cluster.
func NewWorkloadIdentityConfigUpdate(in *v1beta1.WorkloadIdentityConfig) UpdateFn {
	return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
		out := &container.Cluster{}
		GenerateWorkloadIdentityConfig(in, out)
		update := &container.UpdateClusterRequest{
			Update: &container.ClusterUpdate{
				DesiredWorkloadIdentityConfig: out.WorkloadIdentityConfig,
			},
		}
		return s.Projects.Locations.Clusters.Update(name, update).Context(ctx).Do()
	}
}

// CheckForBootstrapNodePool checks if the bootstrap node pool exists for the
// cluster and returns a function to delete it if so.
func CheckForBootstrapNodePool(c container.Cluster) UpdateFn {
	for _, pool := range c.NodePools {
		if pool == nil || pool.Name != BootstrapNodePoolName {
			continue
		}
		return func(s container.Service, ctx context.Context, name string) (*container.Operation, error) {
			return s.Projects.Locations.Clusters.NodePools.Delete(GetFullyQualifiedBNP(name)).Context(ctx).Do()
		}
	}
	return nil
}

// IsUpToDate checks whether current state is up-to-date compared to the given
// set of parameters.
// NOTE(hasheddan): This function is significantly above our cyclomatic
// complexity limit, but is necessary due to the fact that the GKE API only
// allows for update of one field at a time.
func IsUpToDate(in *v1beta1.GKEClusterParameters, currentState container.Cluster) (bool, UpdateFn) { // nolint:gocyclo
	currentParams := &v1beta1.GKEClusterParameters{}
	LateInitializeSpec(currentParams, currentState)
	if fn := CheckForBootstrapNodePool(currentState); fn != nil {
		return false, fn
	}
	if !cmp.Equal(in.AddonsConfig, currentParams.AddonsConfig) {
		return false, NewAddonsConfigUpdate(in.AddonsConfig)
	}
	if !cmp.Equal(in.Autoscaling, currentParams.Autoscaling) {
		return false, NewAutoscalingUpdate(in.Autoscaling)
	}
	if !cmp.Equal(in.BinaryAuthorization, currentParams.BinaryAuthorization) {
		return false, NewBinaryAuthorizationUpdate(in.BinaryAuthorization)
	}
	if !cmp.Equal(in.DatabaseEncryption, currentParams.DatabaseEncryption) {
		return false, NewDatabaseEncryptionUpdate(in.DatabaseEncryption)
	}
	if !cmp.Equal(in.LegacyAbac, currentParams.LegacyAbac) {
		return false, NewLegacyAbacUpdate(in.LegacyAbac)
	}
	if !cmp.Equal(in.Locations, currentParams.Locations) {
		return false, NewLocationsUpdate(in.Locations)
	}
	if !cmp.Equal(in.LoggingService, currentParams.LoggingService) {
		return false, NewLoggingServiceUpdate(in.LoggingService)
	}
	if !cmp.Equal(in.MaintenancePolicy, currentParams.MaintenancePolicy) {
		return false, NewMaintenancePolicyUpdate(in.MaintenancePolicy)
	}
	if !cmp.Equal(in.MasterAuth, currentParams.MasterAuth) {
		return false, NewMasterAuthUpdate(in.MasterAuth)
	}
	if !cmp.Equal(in.MasterAuthorizedNetworksConfig, currentParams.MasterAuthorizedNetworksConfig) {
		return false, NewMasterAuthorizedNetworksConfigUpdate(in.MasterAuthorizedNetworksConfig)
	}
	if !cmp.Equal(in.MonitoringService, currentParams.MonitoringService) {
		return false, NewMonitoringServiceUpdate(in.MonitoringService)
	}
	if !cmp.Equal(in.NetworkConfig, currentParams.NetworkConfig) {
		return false, NewNetworkConfigUpdate(in.NetworkConfig)
	}
	if !cmp.Equal(in.NetworkPolicy, currentParams.NetworkPolicy) {
		return false, NewNetworkPolicyUpdate(in.NetworkPolicy)
	}
	if !cmp.Equal(in.PodSecurityPolicyConfig, currentParams.PodSecurityPolicyConfig) {
		return false, NewPodSecurityPolicyConfigUpdate(in.PodSecurityPolicyConfig)
	}
	if !cmp.Equal(in.PrivateClusterConfig, currentParams.PrivateClusterConfig) {
		return false, NewPrivateClusterConfigUpdate(in.PrivateClusterConfig)
	}
	if !cmp.Equal(in.ResourceLabels, currentParams.ResourceLabels) {
		return false, NewResourceLabelsUpdate(in.ResourceLabels)
	}
	if !cmp.Equal(in.ResourceUsageExportConfig, currentParams.ResourceUsageExportConfig) {
		return false, NewResourceUsageExportConfigUpdate(in.ResourceUsageExportConfig)
	}
	if !cmp.Equal(in.VerticalPodAutoscaling, currentParams.VerticalPodAutoscaling) {
		return false, NewVerticalPodAutoscalingUpdate(in.VerticalPodAutoscaling)
	}
	if !cmp.Equal(in.WorkloadIdentityConfig, currentParams.WorkloadIdentityConfig) {
		return false, NewWorkloadIdentityConfigUpdate(in.WorkloadIdentityConfig)
	}
	return true, nil
}

// GetFullyQualifiedParent builds the fully qualified name of the cluster
// parent.
func GetFullyQualifiedParent(project string, p v1beta1.GKEClusterParameters) string {
	return fmt.Sprintf(ParentFormat, project, p.Location)
}

// GetFullyQualifiedName builds the fully qualified name of the cluster.
func GetFullyQualifiedName(project string, p v1beta1.GKEClusterParameters) string {
	return fmt.Sprintf(ClusterNameFormat, project, p.Location, p.Name)
}

// GetFullyQualifiedBNP build the fully qualified name of the bootstrap node
// pool.
func GetFullyQualifiedBNP(clusterName string) string {
	return fmt.Sprintf(BNPNameFormat, clusterName, BootstrapNodePoolName)
}

// GenerateClientConfig generates a clientcmdapi.Config that can be used by any
// kubernetes client.
func GenerateClientConfig(cluster *container.Cluster) (clientcmdapi.Config, error) {
	c := clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			cluster.Name: {
				Server: fmt.Sprintf("https://%s", cluster.Endpoint),
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			cluster.Name: {
				Cluster:  cluster.Name,
				AuthInfo: cluster.Name,
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			cluster.Name: {
				Username: cluster.MasterAuth.Username,
				Password: cluster.MasterAuth.Password,
			},
		},
		CurrentContext: cluster.Name,
	}

	val, err := base64.StdEncoding.DecodeString(cluster.MasterAuth.ClusterCaCertificate)
	if err != nil {
		return clientcmdapi.Config{}, err
	}
	c.Clusters[cluster.Name].CertificateAuthorityData = val

	val, err = base64.StdEncoding.DecodeString(cluster.MasterAuth.ClientCertificate)
	if err != nil {
		return clientcmdapi.Config{}, err
	}
	c.AuthInfos[cluster.Name].ClientCertificateData = val

	val, err = base64.StdEncoding.DecodeString(cluster.MasterAuth.ClientKey)
	if err != nil {
		return clientcmdapi.Config{}, err
	}
	c.AuthInfos[cluster.Name].ClientKeyData = val

	return c, nil
}
