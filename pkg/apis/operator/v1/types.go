// Copyright (c) 2019-2020 Tigera, Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// InstallationSpec defines configuration for a Calico or Calico Enterprise installation.
// +k8s:openapi-gen=true
type InstallationSpec struct {
	// Variant is the product to install - one of Calico or TigeraSecureEnterprise
	// Default: Calico
	// +optional
	// +kubebuilder:validation:Enum=Calico,TigeraSecureEnterprise
	Variant ProductVariant `json:"variant,omitempty"`

	// Registry is the default Docker registry used for component Docker images. If specified,
	// all images will be pulled from this registry. If not specified then the default registries
	// will be used.
	//
	// Image format:
	//    `<registry>/<imagePath>/<imageName>:<image-tag>`
	//
	// This option allows configuring the `<registry>` portion of the above format.
	// +optional
	Registry string `json:"registry,omitempty"`

	// ImagePath allows for the path part of an image to be specified. If specified
	// then the specified value will be used as the image path for each image. If not specified
	// or empty, the default for each image will be used.
	//
	// Image format:
	//    `<registry>/<imagePath>/<imageName>:<image-tag>`
	//
	// This option allows configuring the `<imagePath>` portion of the above format.
	// +optional
	ImagePath string `json:"imagePath,omitempty"`

	// ImagePullSecrets is an array of references to container registry pull secrets to use. These are
	// applied to all images to be pulled.
	// +optional
	ImagePullSecrets []v1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// KubernetesProvider specifies a particular provider of the Kubernetes platform and enables provider-specific configuration.
	// If the specified value is empty, the Operator will attempt to automatically determine the current provider.
	// If the specified value is not empty, the Operator will still attempt auto-detection, but
	// will additionally compare the auto-detected value to the specified value to confirm they match.
	// +optional
	// +kubebuilder:validation:Enum=,EKS,GKE,AKS,OpenShift,DockerEnterprise
	KubernetesProvider Provider `json:"kubernetesProvider,omitempty"`

	// CalicoNetwork specifies configuration options for Calico provided pod networking.
	// +optional
	CalicoNetwork *CalicoNetworkSpec `json:"calicoNetwork,omitempty"`

	// How the cluster is managed. Valid values for this field are: Standalone, Management, Managed.
	// Standalone clusters are fully self-contained installations of Calico Enterprise. Management clusters provide
	// a single view to manage any number of Managed clusters, which are a lighter weight installation.
	// This option is applicable only when variant is TigeraSecureEnterprise.
	// Default: Standalone
	// +optional
	// +kubebuilder:validation:Enum=Standalone,Management,Managed
	ClusterManagementType ClusterManagementType `json:"clusterManagementType,omitempty"`

	// ControlPlaneNodeSelector is used to select control plane nodes on which to run specific Calico
	// components. This currently only applies to kube-controllers and the apiserver.
	// +optional
	ControlPlaneNodeSelector map[string]string `json:"controlPlaneNodeSelector,omitempty"`

	// NodeMetricsPort specifies which port calico/node serves prometheus metrics on. By default, metrics are not enabled.
	// If specified, this overrides any FelixConfiguration resources which may exist. If omitted, then
	// prometheus metrics may still be configured through FelixConfiguration.
	// +optional
	NodeMetricsPort *int32 `json:"nodeMetricsPort,omitempty"`

	// FlexVolumePath optionally specifies a custom path for FlexVolume. If not specified, FlexVolume will be
	// enabled by default. If set to 'None', FlexVolume will be disabled. The default is based on the
	// kubernetesProvider.
	// +optional
	FlexVolumePath string `json:"flexVolumePath,omitempty"`

	// NodeUpdateStrategy can be used to customize the desired update strategy, such as the MaxUnavailable
	// field.
	// +optional
	NodeUpdateStrategy appsv1.DaemonSetUpdateStrategy `json:"nodeUpdateStrategy,omitempty"`
}

// Provider represents a particular provider or flavor of Kubernetes. Valid options
// are: EKS, GKE, AKS, OpenShift, DockerEnterprise.
type Provider string

var (
	ProviderNone      Provider = ""
	ProviderEKS       Provider = "EKS"
	ProviderGKE       Provider = "GKE"
	ProviderAKS       Provider = "AKS"
	ProviderOpenShift Provider = "OpenShift"
	ProviderDockerEE  Provider = "DockerEnterprise"
)

// ProductVariant represents the variant of the product. Valid options are: Calico, TigeraSecureEnterprise.
type ProductVariant string

var (
	Calico                 ProductVariant = "Calico"
	TigeraSecureEnterprise ProductVariant = "TigeraSecureEnterprise"
)

// ClusterManagementType represents the type of multicluster management to use. Valid options for this field are: "Standalone",
// "Management", "Managed".
type ClusterManagementType string

const (
	ClusterManagementTypeStandalone ClusterManagementType = "Standalone"
	ClusterManagementTypeManagement ClusterManagementType = "Management"
	ClusterManagementTypeManaged    ClusterManagementType = "Managed"
)

// HostPortsType specifies if the HostPorts plugin enabled status.
type HostPortsType string

const (
	HostPortsEnabled  HostPortsType = "Enabled"
	HostPortsDisabled HostPortsType = "Disabled"
)

var HostPortsTypes []HostPortsType = []HostPortsType{
	HostPortsEnabled,
	HostPortsDisabled,
}
var HostPortsTypesString []string = []string{
	HostPortsEnabled.String(),
	HostPortsDisabled.String(),
}

type MultiInterfaceMode string

func (m MultiInterfaceMode) Value() string {
	return strings.ToLower(string(m))
}

const (
	MultiInterfaceModeNone   MultiInterfaceMode = "None"
	MultiInterfaceModeMultus MultiInterfaceMode = "Multus"
)

func (nt HostPortsType) String() string {
	return string(nt)
}

// CalicoNetworkSpec specifies configuration options for Calico provided pod networking.
type CalicoNetworkSpec struct {
	// IPPools contains a list of IP pools to use for allocating pod IP addresses. At most one IP pool
	// may be specified. If omitted, a single pool will be configured when needed.
	// +optional
	IPPools []IPPool `json:"ipPools,omitempty"`

	// MTU specifies the maximum transmission unit to use for pods on the Calico network.
	// Default: 1410
	// +optional
	MTU *int32 `json:"mtu,omitempty"`

	// NodeAddressAutodetectionV4 specifies an approach to automatically detect node IPv4 addresses. If not specified,
	// will use default auto-detection settings to acquire an IPv4 address for each node.
	// +optional
	NodeAddressAutodetectionV4 *NodeAddressAutodetection `json:"nodeAddressAutodetectionV4,omitempty"`

	// NodeAddressAutodetectionV6 specifies an approach to automatically detect node IPv4 addresses. If not specified,
	// IPv6 addresses will not be auto-detected.
	// +optional
	NodeAddressAutodetectionV6 *NodeAddressAutodetection `json:"nodeAddressAutodetectionV6,omitempty"`

	// HostPorts configures whether or not Calico will support Kubernetes HostPorts.
	// Default: Enabled
	// +optional
	// +kubebuilder:validation:Enum=Enabled,Disabled
	HostPorts *HostPortsType `json:"hostPorts,omitempty"`

	// MultiInterfaceMode configures what will configure multiple interface per pod. Only valid for Calico Enterprise installations.
	// Default: None
	// +optional
	// +kubebuilder:validation:Enum=None,Multus
	MultiInterfaceMode *MultiInterfaceMode `json:"multiInterfaceMode,omitempty"`
}

// NodeAddressAutodetection provides configuration options for auto-detecting node addresses. At most one option
// can be used. If no detection option is specified, then IP auto detection will be disabled for this address family and IPs
// must be specified directly on the Node resource.
type NodeAddressAutodetection struct {
	// FirstFound uses default interface matching parameters to select an interface, performing best-effort
	// filtering based on well-known interface names.
	// +optional
	FirstFound *bool `json:"firstFound,omitempty"`

	// Interface enables IP auto-detection based on interfaces that match the given regex.
	// +optional
	Interface string `json:"interface,omitempty"`

	// SkipInterface enables IP auto-detection based on interfaces that do not match
	// the given regex.
	// +optional
	SkipInterface string `json:"skipInterface,omitempty"`

	// CanReach enables IP auto-detection based on which source address on the node is used to reach the
	// specified IP or domain.
	// +optional
	CanReach string `json:"canReach,omitempty"`
}

// EncapsulationType is the type of encapsulation to use on an IP pool. Valid
// options are: IPIP, VXLAN, IPIPCrossSubnet, VXLANCrossSubnet, None.
type EncapsulationType string

const (
	EncapsulationIPIPCrossSubnet  EncapsulationType = "IPIPCrossSubnet"
	EncapsulationIPIP             EncapsulationType = "IPIP"
	EncapsulationVXLAN            EncapsulationType = "VXLAN"
	EncapsulationVXLANCrossSubnet EncapsulationType = "VXLANCrossSubnet"
	EncapsulationNone             EncapsulationType = "None"
	EncapsulationDefault          EncapsulationType = EncapsulationIPIP
)

var EncapsulationTypes []EncapsulationType = []EncapsulationType{
	EncapsulationIPIPCrossSubnet,
	EncapsulationIPIP,
	EncapsulationVXLAN,
	EncapsulationVXLANCrossSubnet,
	EncapsulationNone,
}
var EncapsulationTypesString []string = []string{
	EncapsulationIPIPCrossSubnet.String(),
	EncapsulationIPIP.String(),
	EncapsulationVXLAN.String(),
	EncapsulationVXLANCrossSubnet.String(),
	EncapsulationNone.String(),
}

func (et EncapsulationType) String() string {
	return string(et)
}

// NATOutgoingType describe the type of outgoing NAT to use.
type NATOutgoingType string

const (
	NATOutgoingEnabled  NATOutgoingType = "Enabled"
	NATOutgoingDisabled NATOutgoingType = "Disabled"
)

var NATOutgoingTypes []NATOutgoingType = []NATOutgoingType{
	NATOutgoingEnabled,
	NATOutgoingDisabled,
}
var NATOutgoingTypesString []string = []string{
	NATOutgoingEnabled.String(),
	NATOutgoingDisabled.String(),
}

func (nt NATOutgoingType) String() string {
	return string(nt)
}

const NodeSelectorDefault string = "all()"

type IPPool struct {
	// CIDR contains the address range for the IP Pool in classless inter-domain routing format.
	CIDR string `json:"cidr"`

	// Encapsulation specifies the encapsulation type that will be used with
	// the IP Pool.
	// Default: IPIP
	// +optional
	// +kubebuilder:validation:Enum=IPIPCrossSubnet,IPIP,VXLAN,VXLANCrossSubnet,None
	Encapsulation EncapsulationType `json:"encapsulation,omitempty"`

	// NATOutgoing specifies if NAT will be enabled or disabled for outgoing traffic.
	// Default: Enabled
	// +optional
	// +kubebuilder:validation:Enum=Enabled,Disabled
	NATOutgoing NATOutgoingType `json:"natOutgoing,omitempty"`

	// NodeSelector specifies the node selector that will be set for the IP Pool.
	// Default: 'all()'
	// +optional
	NodeSelector string `json:"nodeSelector,omitempty"`

	// BlockSize specifies the CIDR prefex length to use when allocating per-node IP blocks from
	// the main IP pool CIDR.
	// Default: 26 (IPv4), 122 (IPv6)
	// +optional
	BlockSize *int32 `json:"blockSize,omitempty"`
}

// InstallationStatus defines the observed state of the Calico or Calico Enterprise installation.
// +k8s:openapi-gen=true
type InstallationStatus struct {
	// Variant is the most recently observed installed variant - one of Calico or TigeraSecureEnterprise
	// +kubebuilder:validation:Enum=Calico,TigeraSecureEnterprise
	Variant ProductVariant `json:"variant,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +genclient
// +genclient:nonNamespaced

// Installation configures an installation of Calico or Calico Enterprise. At most one instance
// of this resource is supported. It must be named "default". The Installation API installs core networking
// and network policy components, and provides general install-time configuration.
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Installation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Specification of the desired state for the Calico or Calico Enterprise installation.
	Spec InstallationSpec `json:"spec,omitempty"`

	// Most recently observed state for the Calico or Calico Enterprise installation.
	Status InstallationStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InstallationList contains a list of Installation
type InstallationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Installation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Installation{}, &InstallationList{})
}
