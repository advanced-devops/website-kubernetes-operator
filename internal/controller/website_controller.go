package controller

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"fmt"

	demov1alpha1 "github.com/arvindpathare/website-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/tools/record"
)

type WebsiteReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=demo.homelab.io,resources=websites,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=demo.homelab.io,resources=websites/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=demo.homelab.io,resources=websites/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

func (r *WebsiteReconciler) Reconcile(
	ctx context.Context,
	req ctrl.Request,
) (ctrl.Result, error) {

	logger := log.FromContext(ctx)

	logger.Info(
		"Reconciling Website",
		"name", req.Name,
		"namespace", req.Namespace,
	)

	// ---------------------------------------------------------
	// STEP 1: Get Website
	// ---------------------------------------------------------

	website := &demov1alpha1.Website{}

	err := r.Get(ctx, req.NamespacedName, website)

	if err != nil {
		if apierrors.IsNotFound(err) {
			logger.Info("Website resource not found")
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	logger.Info(
		"Website found",
		"image", website.Spec.Image,
		"replicas", website.Spec.Replicas,
		"port", website.Spec.Port,
	)

	// ---------------------------------------------------------
	// STEP 2: Define desired Deployment
	// ---------------------------------------------------------

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      website.Name,
			Namespace: website.Namespace,
		},
	}

	// ---------------------------------------------------------
	// STEP 3: Create or update Deployment
	// ---------------------------------------------------------

	result, err := controllerutil.CreateOrUpdate(
		ctx,
		r.Client,
		deployment,
		func() error {

			labels := map[string]string{
				"app": website.Name,
			}

			deployment.Spec.Replicas = &website.Spec.Replicas

			deployment.Spec.Selector = &metav1.LabelSelector{
				MatchLabels: labels,
			}

			deployment.Spec.Template.ObjectMeta.Labels = labels

			deployment.Spec.Template.Spec.Containers = []corev1.Container{
				{
					Name:  "website",
					Image: website.Spec.Image,
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: website.Spec.Port,
						},
					},
				},
			}

			return controllerutil.SetControllerReference(
				website,
				deployment,
				r.Scheme,
			)
		},
	)

	if err != nil {
		logger.Error(err, "Failed to reconcile Deployment")
		return ctrl.Result{}, err
	}
	switch result {
	case controllerutil.OperationResultCreated:
		r.Recorder.Event(
			website,
			corev1.EventTypeNormal,
			"DeploymentCreated",
			"Created Deployment for Website",
		)

	case controllerutil.OperationResultUpdated:
		r.Recorder.Event(
			website,
			corev1.EventTypeNormal,
			"DeploymentUpdated",
			"Updated Deployment for Website",
		)
	}

	logger.Info(
		"Deployment reconciled",
		"operation", result,
	)

	// ---------------------------------------------------------
	// STEP 4: Define desired Service
	// ---------------------------------------------------------

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      website.Name,
			Namespace: website.Namespace,
		},
	}

	// ---------------------------------------------------------
	// STEP 5: Create or update Service
	// ---------------------------------------------------------

	serviceResult, err := controllerutil.CreateOrUpdate(
		ctx,
		r.Client,
		service,
		func() error {

			labels := map[string]string{
				"app": website.Name,
			}

			service.Spec.Selector = labels

			service.Spec.Ports = []corev1.ServicePort{
				{
					Name:       "http",
					Port:       website.Spec.Port,
					TargetPort: intstr.FromInt32(website.Spec.Port),
					Protocol:   corev1.ProtocolTCP,
				},
			}

			return controllerutil.SetControllerReference(
				website,
				service,
				r.Scheme,
			)
		},
	)
	switch serviceResult {
	case controllerutil.OperationResultCreated:
		r.Recorder.Event(
			website,
			corev1.EventTypeNormal,
			"ServiceCreated",
			"Created Service for Website",
		)

	case controllerutil.OperationResultUpdated:
		r.Recorder.Event(
			website,
			corev1.EventTypeNormal,
			"ServiceUpdated",
			"Updated Service for Website",
		)
	}

	if err != nil {
		logger.Error(err, "Failed to reconcile Service")
		return ctrl.Result{}, err
	}

	logger.Info(
		"Service reconciled",
		"operation", serviceResult,
	)

	// ---------------------------------------------------------
	// STEP 6: Update Website Status
	// ---------------------------------------------------------

	readyReplicas := deployment.Status.ReadyReplicas

	website.Status.ReadyReplicas = readyReplicas
	website.Status.URL = fmt.Sprintf(
		"http://%s.%s.svc.cluster.local:%d",
		service.Name,
		service.Namespace,
		website.Spec.Port,
	)

	website.Status.ObservedGeneration = website.Generation

	if readyReplicas == website.Spec.Replicas {
		meta.SetStatusCondition(
			&website.Status.Conditions,
			metav1.Condition{
				Type:               "Ready",
				Status:             metav1.ConditionTrue,
				Reason:             "WebsiteAvailable",
				Message:            "All Website replicas are available",
				ObservedGeneration: website.Generation,
			},
		)
	} else {
		meta.SetStatusCondition(
			&website.Status.Conditions,
			metav1.Condition{
				Type:   "Ready",
				Status: metav1.ConditionFalse,
				Reason: "DeploymentProgressing",
				Message: fmt.Sprintf(
					"Waiting for replicas: %d/%d ready",
					readyReplicas,
					website.Spec.Replicas,
				),
				ObservedGeneration: website.Generation,
			},
		)
	}

	if err := r.Status().Update(ctx, website); err != nil {
		logger.Error(err, "Failed to update Website status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *WebsiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demov1alpha1.Website{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Named("website").
		Complete(r)
}
