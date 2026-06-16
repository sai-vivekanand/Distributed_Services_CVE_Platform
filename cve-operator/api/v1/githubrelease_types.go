/*
Copyright 2024.

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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GithubReleaseSpec defines the desired state of GithubRelease
type GithubReleaseSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// url for extraction for the job
	// +kubebuilder:validation:Required
	Url string `json:"url"`
}

// GithubReleaseStatus defines the observed state of GithubRelease
type GithubReleaseStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Active holds pointers to the active jobs of this CronJob.
	// +optional
	Active *corev1.ObjectReference `json:"active,omitempty"`

	// LastScheduleTime keeps tracks of the last time the job was successfully scheduled.
	// +optional
	LastScheduleTime *metav1.Time `json:"lastScheduleTime,omitempty"`

	// JobReference of the job that was most recently successfully scheduled.
	// +optional
	JobReference *corev1.ObjectReference `json:"jobReference,omitempty"`

	// Conditions represent the latest available observations of a Job's current state.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// GithubRelease is the Schema for the githubreleases API
type GithubRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GithubReleaseSpec   `json:"spec,omitempty"`
	Status GithubReleaseStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GithubReleaseList contains a list of GithubRelease
type GithubReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GithubRelease `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GithubRelease{}, &GithubReleaseList{})
}
