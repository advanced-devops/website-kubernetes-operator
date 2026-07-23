package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TLSSpec struct {
	Enabled bool `json:"enabled"`
}
type MonitoringSpec struct {
	Enabled bool `json:"enabled"`
}

type WebsiteSpec struct {
	// Image is the container image used by the Website.
	//
	// +kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	// Replicas is the desired number of Website replicas.
	//
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	Replicas int32 `json:"replicas,omitempty"`

	// Port is the port exposed by the Website container.
	//
	// +kubebuilder:default=80
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port     int32  `json:"port,omitempty"`
	Hostname string `json:"hostname"`

	TLS TLSSpec `json:"tls,omitempty"`

	Monitoring MonitoringSpec `json:"monitoring,omitempty"`
}

type WebsiteStatus struct {
	// ReadyReplicas is the number of available Website replicas.
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// URL is the internal URL of the Website Service.
	URL string `json:"url,omitempty"`

	// ObservedGeneration is the latest Website generation processed
	// by the controller.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions represent the current state of the Website.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`
// +kubebuilder:printcolumn:name="Ready Replicas",type=integer,JSONPath=`.status.readyReplicas`
// +kubebuilder:printcolumn:name="URL",type=string,JSONPath=`.status.url`
// +kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

type Website struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebsiteSpec   `json:"spec"`
	Status WebsiteStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

type WebsiteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Website `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Website{}, &WebsiteList{})
}
