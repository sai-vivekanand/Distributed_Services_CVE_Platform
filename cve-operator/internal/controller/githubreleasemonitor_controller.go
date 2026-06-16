package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	su24v1 "github.com/cyse7125-su24-team12/cve-operator/api/v1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	finalizerName          = "githubreleasemonitor.finalizers.su24.csye7125.t12.io"
	ghReleaseFinalizerName = "githubrelease.finalizers.su24.csye7125.t12.io"
)

// GithubReleaseMonitorReconciler reconciles a GithubReleaseMonitor object
type GithubReleaseMonitorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=su24.csye7125.t12.io,resources=githubreleasemonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=su24.csye7125.t12.io,resources=githubreleasemonitors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=su24.csye7125.t12.io,resources=githubreleasemonitors/finalizers,verbs=update

func (r *GithubReleaseMonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch the GithubReleaseMonitor instance
	var monitor su24v1.GithubReleaseMonitor
	if err := r.Get(ctx, req.NamespacedName, &monitor); err != nil {
		logger.Error(err, "unable to fetch GithubReleaseMonitor")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Handle deletion logic
	if monitor.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer, add it
		if !controllerutil.ContainsFinalizer(&monitor, finalizerName) {
			controllerutil.AddFinalizer(&monitor, finalizerName)
			if err := r.Update(ctx, &monitor); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&monitor, finalizerName) {
			// Run finalization logic for finalizerName
			if err := r.finalizeMonitor(ctx, logger, &monitor); err != nil {
				return ctrl.Result{}, err
			}

			// Remove finalizer
			controllerutil.RemoveFinalizer(&monitor, finalizerName)
			if err := r.Update(ctx, &monitor); err != nil {
				return ctrl.Result{}, err
			}
		}
		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	// Extract the URL and monitorFrom date from the spec
	url := monitor.Spec.URL
	monitorFrom := monitor.Spec.MonitorFrom

	// Convert monitorFrom to a time.Time object
	monitorFromDate, err := time.Parse(time.RFC3339, monitorFrom)
	if err != nil {
		logger.Error(err, "unable to parse monitorFrom date")
		return ctrl.Result{}, err
	}

	// Poll the GitHub releases URL
	releases, err := r.fetchGitHubReleases(url)
	if err != nil {
		logger.Error(err, "unable to fetch GitHub releases")
		return ctrl.Result{RequeueAfter: 2 * time.Minute}, err
	}

	// Filter releases based on monitorFrom date and search for "_delta_" assets
	deltaLinks := r.getDeltaLinks(releases, monitorFromDate, logger)

	// Initialize or update the status field
	statusMap := make(map[string]string)
	for _, release := range monitor.Status.Releases {
		statusMap[release.Release] = release.Status
	}

	// Check all created GithubRelease CRs for failed conditions
	if err := r.checkAndRecreateFailedGithubReleases(ctx, logger); err != nil {
		logger.Error(err, "error checking and recreating failed GithubRelease CRs")
		return ctrl.Result{}, err
	}

	// Create GithubRelease CRs for each delta link and update the status accordingly
	for _, deltaLink := range deltaLinks {
		status := "processed successfully"
		exists, err := r.githubReleaseCRExists(ctx, deltaLink, logger)
		if err != nil {
			status = "processing failed"
			logger.Error(err, "unable to create GithubRelease CR", "deltaLink", deltaLink)
		} else if exists {
			status = "already processed"
		} else if err := r.createGithubReleaseCR(ctx, deltaLink, logger); err != nil {
			status = "processing failed"
			logger.Error(err, "unable to create GithubRelease CR", "deltaLink", deltaLink)
		}
		statusMap[deltaLink] = status
	}

	// Update the status field in the monitor
	var updatedReleases []su24v1.ReleaseStatus
	for link, status := range statusMap {
		updatedReleases = append(updatedReleases, su24v1.ReleaseStatus{
			Release: link,
			Status:  status,
		})
	}
	monitor.Status.Releases = updatedReleases
	monitor.Status.RetrievedTime = time.Now().UTC().Format(time.RFC3339)
	monitor.Status.MonitorFromTimestamp = monitorFromDate.Format(time.RFC3339)

	if err := r.Status().Update(ctx, &monitor); err != nil {
		logger.Error(err, "unable to update GithubReleaseMonitor status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: 2 * time.Minute}, nil
}

func (r *GithubReleaseMonitorReconciler) finalizeMonitor(ctx context.Context, logger logr.Logger, monitor *su24v1.GithubReleaseMonitor) error {
	// Get associated GithubRelease CRs
	var releases su24v1.GithubReleaseList
	if err := r.List(ctx, &releases, client.InNamespace(monitor.Namespace)); err != nil {
		return err
	}

	for _, release := range releases.Items {
		// Delete each associated GithubRelease CR
		if err := r.Delete(ctx, &release); err != nil {
			return err
		}
		logger.Info("Deleted GithubRelease associated with GithubReleaseMonitor", "release", release.Name)
	}

	// Ensure all GithubRelease CRs are deleted before removing the finalizer
	for _, release := range releases.Items {
		// Wait for each GithubRelease CR to be deleted
		if err := r.waitForDeletion(ctx, &release); err != nil {
			return err
		}
	}

	return nil
}

func (r *GithubReleaseMonitorReconciler) waitForDeletion(ctx context.Context, release *su24v1.GithubRelease) error {
	timeout := time.Minute * 5
	interval := time.Second * 10

	endTime := time.Now().Add(timeout)
	for time.Now().Before(endTime) {
		var fetchedRelease su24v1.GithubRelease
		err := r.Get(ctx, client.ObjectKey{
			Namespace: release.Namespace,
			Name:      release.Name,
		}, &fetchedRelease)
		if errors.IsNotFound(err) {
			// Release is deleted
			return nil
		}
		if err != nil {
			return err
		}
		time.Sleep(interval)
	}

	return fmt.Errorf("timed out waiting for release %s to be deleted", release.Name)
}

func (r *GithubReleaseMonitorReconciler) fetchGitHubReleases(url string) ([]map[string]interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+ os.Getenv("GITHUB_TOKEN"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch releases, status code: %d", resp.StatusCode)
	}

	var releases []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, err
	}

	return releases, nil
}

func (r *GithubReleaseMonitorReconciler) getDeltaLinks(releases []map[string]interface{}, monitorFromDate time.Time, logger logr.Logger) []string {
	var deltaLinks []string

	for _, release := range releases {
		publishedAtStr, ok := release["published_at"].(string)
		if !ok {
			logger.Info("Missing published_at field in release")
			continue
		}

		publishedAt, err := time.Parse(time.RFC3339, publishedAtStr)
		if err != nil {
			logger.Error(err, "unable to parse published_at date")
			continue
		}

		if publishedAt.After(monitorFromDate) {
			assets, ok := release["assets"].([]interface{})
			if !ok {
				logger.Info("Missing assets field in release")
				continue
			}

			for _, asset := range assets {
				assetMap, ok := asset.(map[string]interface{})
				if !ok {
					continue
				}

				name, ok := assetMap["name"].(string)
				if !ok {
					continue
				}

				if strings.Contains(name, "_delta_") {
					browserDownloadURL, ok := assetMap["browser_download_url"].(string)
					if ok {
						deltaLinks = append(deltaLinks, browserDownloadURL)
					}
				}
			}
		}
	}

	return deltaLinks
}

func (r *GithubReleaseMonitorReconciler) checkAndRecreateFailedGithubReleases(ctx context.Context, logger logr.Logger) error {
	var githubReleases su24v1.GithubReleaseList
	if err := r.List(ctx, &githubReleases); err != nil {
		return err
	}

	for _, release := range githubReleases.Items {
		for _, condition := range release.Status.Conditions {
			if condition.Type == "Failed" && condition.Status == metav1.ConditionTrue {
				// Delete the failed GithubRelease CR
				if err := r.Delete(ctx, &release); err != nil {
					logger.Error(err, "unable to delete failed GithubRelease CR", "release", release.Name)
					return err
				}
				logger.Info("Deleted failed GithubRelease CR", "release", release.Name)

				// Recreate the GithubRelease CR
				newRelease := su24v1.GithubRelease{
					ObjectMeta: metav1.ObjectMeta{
						Name:      release.Name,
						Namespace: release.Namespace,
					},
					Spec: release.Spec,
				}
				if err := r.Create(ctx, &newRelease); err != nil {
					logger.Error(err, "unable to recreate GithubRelease CR", "release", release.Name)
					return err
				}
				logger.Info("Recreated GithubRelease CR", "release", release.Name)
			}
		}
	}

	return nil
}

func (r *GithubReleaseMonitorReconciler) githubReleaseCRExists(ctx context.Context, deltaLink string, logger logr.Logger) (bool, error) {
	var existingReleases su24v1.GithubReleaseList
	if err := r.List(ctx, &existingReleases, client.MatchingFields{"spec.Url": deltaLink}); err != nil {
		if !errors.IsNotFound(err) {
			logger.Error(err, "Failed to list GithubRelease CRs", "deltaLink", deltaLink)
			return false, err
		}
		return false, nil
	}

	for _, existingRelease := range existingReleases.Items {
		if existingRelease.Spec.Url == deltaLink {
			return true, nil
		}
	}
	return false, nil
}

func (r *GithubReleaseMonitorReconciler) createGithubReleaseCR(ctx context.Context, deltaLink string, logger logr.Logger) error {
	// Create the GithubRelease CR
	release := &su24v1.GithubRelease{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "githubrelease-",
			Namespace:    "cve-operator",
		},
		Spec: su24v1.GithubReleaseSpec{
			Url: deltaLink,
		},
	}

	// Add finalizer to the GithubRelease CR
	controllerutil.AddFinalizer(release, ghReleaseFinalizerName)

	if err := r.Create(ctx, release); err != nil {
		return err
	}

	logger.Info("Created GithubRelease CR", "deltaLink", deltaLink)
	return nil
}

func (r *GithubReleaseMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&su24v1.GithubReleaseMonitor{}).
		Complete(r)
}
