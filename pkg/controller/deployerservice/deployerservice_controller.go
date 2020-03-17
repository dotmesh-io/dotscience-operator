package deployerservice

import (
	"context"
	"os"
	"time"

	deployerv1 "github.com/dotmesh-io/dotscience-operator/pkg/apis/deployer/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Deployer image override for offline mode
const (
	EnvRelatedImage = "RELATED_IMAGE_DEPLOYER"
)

var log = logf.Log.WithName("controller_deployerservice")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new DeployerService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileDeployerService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("deployerservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource DeployerService
	err = c.Watch(&source.Kind{Type: &deployerv1.DeployerService{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner DeployerService
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &deployerv1.DeployerService{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileDeployerService implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileDeployerService{}

// ReconcileDeployerService reconciles a DeployerService object
type ReconcileDeployerService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a DeployerService object and makes changes based on the state read
// and what is in the DeployerService.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileDeployerService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling DeployerService")

	reconcilePeriod := 10 * time.Second
	reconcileResult := reconcile.Result{RequeueAfter: reconcilePeriod}

	// Return this for a immediate retry of the reconciliation loop with the
	// same request object.
	immediateRetryResult := reconcile.Result{Requeue: true}

	// Fetch the DeployerService instance
	instance := &deployerv1.DeployerService{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	opts := &podOpts{
		image: os.Getenv(EnvRelatedImage),
	}

	// Define a new Pod object
	pod := newPodForCR(opts, instance)

	// Set DeployerService instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {

		instance.Status.Status = deployerv1.DeployerPhaseCreating
		if updateErr := r.client.Status().Update(context.Background(), instance); updateErr != nil {
			reqLogger.Info("Failed to update status", "Error", err.Error())
		}

		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		return immediateRetryResult, nil
	} else if err != nil {

		return reconcile.Result{}, err
	}

	instance.Status.Status = deployerv1.DeployerPhaseRunning

	if err := r.client.Status().Update(context.Background(), instance); err != nil {
		return immediateRetryResult, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcileResult, nil
}

type podOpts struct {
	image string
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(opts *podOpts, cr *deployerv1.DeployerService) *corev1.Pod {

	if opts.image == "" {
		opts.image = cr.Spec.Image
	}

	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: cr.Spec.ServiceAccountName,
			Containers: []corev1.Container{
				{
					Name:    cr.Spec.Name,
					Image:   opts.image,
					Command: []string{"ds-deployer", "run"},
					Env: []corev1.EnvVar{
						{
							Name:  "GATEWAY_ADDRESS",
							Value: cr.Spec.GatewayAddress,
						},
						{
							Name:  "TOKEN",
							Value: cr.Spec.Token,
						},
						{
							Name:  "HEALTH_PORT",
							Value: "9300",
						},
					},
					// adding healthcheck port
					Ports: []corev1.ContainerPort{
						corev1.ContainerPort{
							ContainerPort: 9300,
						},
					},
					LivenessProbe: &corev1.Probe{
						Handler: corev1.Handler{
							HTTPGet: &corev1.HTTPGetAction{
								Path: "/health",
								Port: intstr.FromInt(9300),
							},
						},
						InitialDelaySeconds: 30,
						TimeoutSeconds:      10,
					},
				},
			},
		},
	}
}
