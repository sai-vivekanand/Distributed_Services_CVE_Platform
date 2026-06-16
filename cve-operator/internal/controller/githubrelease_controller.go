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

package controller

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	kbatch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	// ref "k8s.io/client-go/tools/reference"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	su24v1 "github.com/cyse7125-su24-team12/cve-operator/api/v1"
	"github.com/go-logr/logr"

	// "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// "sigs.k8s.io/controller-runtime/pkg/handler"
	// "sigs.k8s.io/controller-runtime/pkg/reconcile"
	// "sigs.k8s.io/controller-runtime/pkg/source"
	// // "k8s.io/apimachinery/pkg/api/errors"
	// "sigs.k8s.io/controller-runtime/pkg/builder"
	// "sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	githubReleaseFinalizerName = "githubrelease.finalizers.su24.csye7125.t12.io"
)

// GithubReleaseReconciler reconciles a GithubRelease object
type GithubReleaseReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=su24.csye7125.t12.io,resources=githubreleases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=su24.csye7125.t12.io,resources=githubreleases/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=su24.csye7125.t12.io,resources=githubreleases/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=batch,resources=jobs,verbs=get;list;watch;create;update;patch;delete

var (
	scheduledTimeAnnotation = "su24.csye7125.t12.io/scheduled-time"
)

func getScheduledTimeForJob(job *kbatch.Job) (*time.Time, error) {
	timeRaw := job.Annotations[scheduledTimeAnnotation]
	if len(timeRaw) == 0 {
		return nil, nil
	}

	timeParsed, err := time.Parse(time.RFC3339, timeRaw)
	if err != nil {
		return nil, err
	}
	return &timeParsed, nil
}

// isJobFinished determines if a job has finished, and if so, whether it was successful or failed.
func isJobFinished(job *kbatch.Job) (bool, string) {
	for _, condition := range job.Status.Conditions {
		if (condition.Type == kbatch.JobFailed || condition.Type == kbatch.JobComplete) && condition.Status == corev1.ConditionTrue {
			return true, string(condition.Type)
		}
	}
	return false, ""
}

// updateCondition updates or adds the given condition in the conditions array
func updateCondition(conditions []metav1.Condition, newCondition metav1.Condition) []metav1.Condition {
	for i, condition := range conditions {
		if condition.Type == newCondition.Type {
			conditions[i] = newCondition
			return conditions
		}
	}
	return append(conditions, newCondition)
}

// extracting name from URL
func extractJobName(url string) string {
	// Split the URL by '/'
	parts := strings.Split(url, "/")

	// Get the last part of the URL
	lastPart := parts[len(parts)-1]

	// Remove the file extension
	lastPart = strings.TrimSuffix(lastPart, ".zip")

	// Replace underscores with hyphens
	lastPart = strings.ReplaceAll(lastPart, "_", "-")

	// Convert to lowercase
	lastPart = strings.ToLower(lastPart)

	// Remove invalid characters
	reg, _ := regexp.Compile("[^a-z0-9-.]")
	lastPart = reg.ReplaceAllString(lastPart, "")

	// Ensure it starts and ends with an alphanumeric character
	lastPart = strings.TrimFunc(lastPart, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	return lastPart
}

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GithubRelease object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile
func (r *GithubReleaseReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	// TODO(user): your logic here
	var githubRelease su24v1.GithubRelease
	if err := r.Get(ctx, req.NamespacedName, &githubRelease); err != nil {
		log.Error(err, "unable to fetch GithubRelease")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Reconciling GithubRelease")

	// Handle deletion logic
	if githubRelease.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer, add it
		if !controllerutil.ContainsFinalizer(&githubRelease, githubReleaseFinalizerName) {
			controllerutil.AddFinalizer(&githubRelease, githubReleaseFinalizerName)
			if err := r.Update(ctx, &githubRelease); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&githubRelease, githubReleaseFinalizerName) {
			// Run finalization logic for githubReleaseFinalizerName
			if err := r.finalizeRelease(ctx, log, &githubRelease); err != nil {
				return ctrl.Result{}, err
			}

			// Remove finalizer
			controllerutil.RemoveFinalizer(&githubRelease, githubReleaseFinalizerName)
			if err := r.Update(ctx, &githubRelease); err != nil {
				return ctrl.Result{}, err
			}
		}
		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	// check if jobName already exists from the job it creates if so don't create a new job
	var job kbatch.Job
	var extractedJobName = extractJobName(githubRelease.Spec.Url)
	err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: extractedJobName}, &job)
	if err == nil {
		log.V(1).Info("job already exists, skipping")
		finished, finishedType := isJobFinished(&job)
		if finished {
			// if the job is finished, update the status of the GithubRelease
			var newCondition metav1.Condition
			switch finishedType {
			case string(kbatch.JobComplete):
				newCondition = metav1.Condition{
					Type:               "Complete",
					Status:             metav1.ConditionTrue,
					LastTransitionTime: metav1.Now(),
					Reason:             "JobCompleted",
					Message:            "Job completed successfully",
				}
			case string(kbatch.JobFailed):
				newCondition = metav1.Condition{
					Type:               "Failed",
					Status:             metav1.ConditionTrue,
					LastTransitionTime: metav1.Now(),
					Reason:             "JobFailed",
					Message:            "Job failed",
				}
			}
			githubRelease.Status.Conditions = updateCondition(githubRelease.Status.Conditions, newCondition)
			githubRelease.Status.Active = nil

			// Update the status of the GithubRelease 
			if err := r.Status().Update(ctx, &githubRelease); err != nil {
				log.Error(err, "unable to update GithubRelease status")
				return ctrl.Result{}, err
			}

			log.V(1).Info("job is finished, updating status")
			// job finished don't request requeue
			return ctrl.Result{}, nil
		} else {
			// if the job is not finished, update the status of the GithubRelease
			log.V(1).Info("job is not finished, updating status")
			githubRelease.Status.Active = &corev1.ObjectReference{
				Kind:       "Job",
				Namespace:  job.Namespace,
				Name:       job.Name,
				APIVersion: "batch/v1",
			}
			// return ctrl.Result{}, nil
		}
		scheduledTimeForJob, err := getScheduledTimeForJob(&job)
		if err != nil {
			log.Error(err, "unable to parse schedule time for child job", "job", &job)
		} else {
			if scheduledTimeForJob != nil {
				githubRelease.Status.LastScheduleTime = &metav1.Time{Time: *scheduledTimeForJob}
			} else {
				githubRelease.Status.LastScheduleTime = nil
			}
		}

		// Update the status of the GithubRelease
		if err := r.Status().Update(ctx, &githubRelease); err != nil {
			log.Error(err, "unable to update GithubRelease status")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	log.V(1).Info("No existing job found, creating a new one")
	// constructJobForGithubRelease creates a new Job object
	constructJobForGithubRelease := func(githubRelease *su24v1.GithubRelease, scheduledTime time.Time) (*kbatch.Job, error) {
		name := extractJobName(githubRelease.Spec.Url)
		backoffLimit := int32(3)
		job := &kbatch.Job{
			ObjectMeta: metav1.ObjectMeta{
				Labels: map[string]string{
					"githubrelease": githubRelease.Name,
				},
				Annotations: make(map[string]string),
				Name:        name,
				Namespace:   githubRelease.Namespace,
			},
			Spec: kbatch.JobSpec{
				BackoffLimit: &backoffLimit,
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							"proxy.istio.io/config": `{ "holdApplicationUntilProxyStarts": true }`,
						},
					},
					Spec: corev1.PodSpec{
						RestartPolicy: corev1.RestartPolicyOnFailure,
						ImagePullSecrets: []corev1.LocalObjectReference{
							{
								Name: "docker-secret",
							},
						},
						Containers: []corev1.Container{
							{
								Name:            "webapp",
								Image:           "bala699/cve-webapp:latest",
								ImagePullPolicy: corev1.PullAlways,
								Env: []corev1.EnvVar{
									{
										Name: "KAFKA_HOST",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: "kafka-secret",
												},
												Key: "KAFKA_HOST",
											},
										},
									},
									{
										Name: "KAFKA_USER",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: "kafka-secret",
												},
												Key: "KAFKA_USER",
											},
										},
									},
									{
										Name: "KAFKA_PASSWORD",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: "kafka-secret",
												},
												Key: "KAFKA_PASSWORD",
											},
										},
									},
									{
										Name: "TOPIC_NAME",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: "kafka-secret",
												},
												Key: "TOPIC_NAME",
											},
										},
									},
									{
										Name:  "CVE_DATA_URL",
										Value: githubRelease.Spec.Url,
									},
								},
							},
						},
					},
				},
			},
		}
		if err := ctrl.SetControllerReference(githubRelease, job, r.Scheme); err != nil {
			return nil, err
		}
		return job, nil
	}

	// actually make the job...
	jobTemplate, err := constructJobForGithubRelease(&githubRelease, time.Now().Add(time.Second*10))
	if err != nil {
		log.Error(err, "unable to construct job from template")
		// don't bother requeuing until we get a change to the spec
		return ctrl.Result{}, nil
	}

	githubRelease.Status.JobReference = &corev1.ObjectReference{
		Kind:       "Job",
		Namespace:  githubRelease.Namespace,
		Name:       jobTemplate.Name,
		APIVersion: "batch/v1",
		UID:        jobTemplate.UID,
	}

	// Update the status of the GithubRelease
	if err := r.Status().Update(ctx, &githubRelease); err != nil {
		log.Error(err, "unable to update GithubRelease status")
		return ctrl.Result{}, err
	}
	// ...and create it on the cluster
	if err := r.Create(ctx, jobTemplate); err != nil {
		log.Error(err, "unable to create Job for GithubRelease", "job", job)
		return ctrl.Result{}, err
	}

	log.V(1).Info("created Job for GithubRelease run", "job", job)

	return ctrl.Result{}, nil
}

func (r *GithubReleaseReconciler) finalizeRelease(ctx context.Context, logger logr.Logger, release *su24v1.GithubRelease) error {
	// List all Jobs associated with the GithubRelease CR
	var jobs kbatch.JobList
	if err := r.List(ctx, &jobs, client.InNamespace(release.Namespace), client.MatchingLabels{"githubrelease": release.Name}); err != nil {
		return err
	}

	// Delete each Job
	for _, job := range jobs.Items {
		deleteOptions := metav1.DeleteOptions{
			PropagationPolicy: func() *metav1.DeletionPropagation {
				propagation := metav1.DeletePropagationForeground
				return &propagation
			}(),
		}
		if err := r.Delete(ctx, &job, &client.DeleteOptions{
			PropagationPolicy: deleteOptions.PropagationPolicy,
		}); err != nil {
			return err
		}
		logger.Info("Deleted Job associated with GithubRelease", "job", job.Name)
	}

	// Wait for each Job to be deleted
	for _, job := range jobs.Items {
		if err := r.waitForJobDeletion(ctx, &job); err != nil {
			return err
		}
	}

	return nil
}

func (r *GithubReleaseReconciler) waitForJobDeletion(ctx context.Context, job *kbatch.Job) error {
	timeout := time.Minute * 5
	interval := time.Second * 10

	endTime := time.Now().Add(timeout)
	for time.Now().Before(endTime) {
		var fetchedJob kbatch.Job
		err := r.Get(ctx, client.ObjectKey{
			Namespace: job.Namespace,
			Name:      job.Name,
		}, &fetchedJob)
		if errors.IsNotFound(err) {
			// Job is deleted
			return nil
		}
		if err != nil {
			return err
		}
		time.Sleep(interval)
	}

	return fmt.Errorf("timed out waiting for job %s to be deleted", job.Name)
}

func (r *GithubReleaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.TODO(), &su24v1.GithubRelease{}, "spec.Url", func(rawObj client.Object) []string {
		release := rawObj.(*su24v1.GithubRelease)
		return []string{release.Spec.Url}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&su24v1.GithubRelease{}).
		Owns(&kbatch.Job{}).
		Complete(r)
}
