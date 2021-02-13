/*
Copyright 2021.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SecretKeyRefSpec struct {
	Name string `json:"name,omitempty"`
	Key  string `json:"key,omitempty"`
}

// GopassRepositorySpec defines the desired state of GopassRepository
type GopassRepositorySpec struct {
	// RepositoryUrl points to the URL of the repository
	RepositoryURL string `json:"repositoryUrl,omitempty"`
	// RefreshInterval denotes how often the repository should be updated
	RefreshInterval string `json:"refreshInterval,omitempty"`
	// UserName used to authenticate authenticate with
	UserName string `json:"userName,omitempty"`
	// SecretKeyRef references the Secret to be used to authenticate
	SecretKeyRef SecretKeyRefSpec `json:"secretKeyRef,omitempty"`
	GpgKeyRef    SecretKeyRefSpec `json:"gpgKeyRef,omitempty"`
}

// GopassRepositoryStatus defines the observed state of GopassRepository
type GopassRepositoryStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// GopassRepository is the Schema for the gopassrepositories API
type GopassRepository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GopassRepositorySpec   `json:"spec,omitempty"`
	Status GopassRepositoryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GopassRepositoryList contains a list of GopassRepository
type GopassRepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GopassRepository `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GopassRepository{}, &GopassRepositoryList{})
}
