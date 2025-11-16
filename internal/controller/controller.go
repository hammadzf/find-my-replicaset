package controller

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type ReconcileDeployment struct {
	client.Client
}

var _ reconcile.Reconciler = &ReconcileDeployment{}

func (r *ReconcileDeployment) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := log.FromContext(ctx)

	// get the deployment reosurce
	dep := &appsv1.Deployment{}
	err := r.Get(ctx, request.NamespacedName, dep)

	if apierrors.IsNotFound(err) {
		log.Error(nil, "could not find the deployment")
		return reconcile.Result{}, nil
	}
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch deployment: %+w", err)
	}

	log.Info("reconciling deployment", "deployment name", dep.ObjectMeta.Name)

	// get the list of ReplicaSets and filter them out using labels
	rs := &appsv1.ReplicaSetList{}
	err = r.List(ctx, rs, client.InNamespace(request.Namespace), client.MatchingLabels(dep.Spec.Selector.MatchLabels))
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not list ReplicaSets: %+w", err)
	}
	if len(rs.Items) == 0 {
		log.Error(nil, "no ReplicaSets found")
		return reconcile.Result{}, nil
	}

	for _, item := range rs.Items {
		// confirm the ownership
		if metav1.IsControlledBy(&item.ObjectMeta, &dep.ObjectMeta) {
			log.Info("found the ReplicaSet owned by deployment", "deployment name", dep.ObjectMeta.Name)
			// write the name of the ReplicaSet as a label in the owner resource (deployment)
			dep.Labels["owned-replica-set"] = item.ObjectMeta.Name
			err := r.Update(ctx, dep)
			if err != nil {
				return reconcile.Result{}, fmt.Errorf("could not update the deployment: %+w", err)
			}
			return reconcile.Result{}, nil
		}
	}

	return reconcile.Result{}, nil
}
