package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DeployerServiceSpec defines the desired state of DeployerService
// +k8s:openapi-gen=true
type DeployerServiceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Name is deployer name
	Name string `json:"name"`

	// Image is the container image to run as the job.
	Image string `json:"image"`

	// Gateway address (dotscience.com)
	GatewayAddress string `json:"gateway"`

	// token key (generated when creating a new deployer via https://cloud.dotscience.com/deployers)
	Token string `json:"token"`

	// ServiceAccountName for the deployer to use. Must have permissions to create resources
	ServiceAccountName string `json:"serviceAccountName"`
}

// DeployerPhase is the phase of the deployer at a given point in time.
type DeployerPhase string

// Constants for operator defaults values and different phases.
const (
	DeployerPhaseInitial     DeployerPhase = ""
	DeployerPhaseRunning     DeployerPhase = "Running"
	DeployerPhaseCreating    DeployerPhase = "Creating"
	DeployerPhasePending     DeployerPhase = "Pending"
	DeployerPhaseTerminating DeployerPhase = "Terminating"
)

// DeployerServiceStatus defines the observed state of DeployerService
type DeployerServiceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	Status DeployerPhase `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DeployerService is the Schema for the deployerservices API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=deployerservices,scope=Namespaced
type DeployerService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeployerServiceSpec   `json:"spec,omitempty"`
	Status DeployerServiceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DeployerServiceList contains a list of DeployerService
type DeployerServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeployerService `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DeployerService{}, &DeployerServiceList{})
}
