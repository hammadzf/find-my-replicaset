package webhook

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	appsv1 "k8s.io/api/apps/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type DepValidator struct{}

func (v *DepValidator) validate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	dep, ok := obj.(*appsv1.Deployment)
	if !ok {
		return nil, fmt.Errorf("expected a Deployment object but got %T", obj)
	}

	log := logf.FromContext(ctx)
	log.Info("validating deployment", "deployment name", dep.ObjectMeta.Name)

	key := "owned-replica-set"
	key, found := dep.Labels[key]
	if !found {
		log.Info("did not find the required label")
		return nil, fmt.Errorf("deployment does not contain the required label %s", key)
	}

	return nil, nil
}

func (v *DepValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return v.validate(ctx, obj)
}

func (v *DepValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	return v.validate(ctx, newObj)
}

func (v *DepValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return v.validate(ctx, obj)
}
