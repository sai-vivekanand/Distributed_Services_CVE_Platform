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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GithubReleaseMonitorSpec defines the desired state of GithubReleaseMonitor
type GithubReleaseMonitorSpec struct {
	// URL for the release such as https://github.com/CVEProject/cvelistV5/releases.
	URL string `json:"url"`

	// MonitorFrom specifies the date to start monitoring releases from.
	// Accepted values are in date format (YYYY-MM-DD).
	MonitorFrom string `json:"monitorFrom"`
}

// ReleaseStatus defines the status of a single release
type ReleaseStatus struct {
	Release string `json:"release"`
	Status  string `json:"status"`
}

// GithubReleaseMonitorStatus defines the observed state of GithubReleaseMonitor
type GithubReleaseMonitorStatus struct {
	// List of releases found on the GitHub repository with their processing statuses.
	Releases []ReleaseStatus `json:"releases"`

	// Time (in UTC) the value was retrieved.
	RetrievedTime string `json:"retrievedTime"`

	// Value of monitorFrom in the form of timestamp (in UTC).
	MonitorFromTimestamp string `json:"monitorFromTimestamp"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// GithubReleaseMonitor is the Schema for the githubreleasemonitors API
type GithubReleaseMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GithubReleaseMonitorSpec   `json:"spec,omitempty"`
	Status GithubReleaseMonitorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GithubReleaseMonitorList contains a list of GithubReleaseMonitor
type GithubReleaseMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GithubReleaseMonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GithubReleaseMonitor{}, &GithubReleaseMonitorList{})
}
