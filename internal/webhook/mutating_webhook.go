package webhook

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	appsv1 "k8s.io/api/apps/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// labels deployments with "owned-replica-set" label
type DepLabeler struct{}

func (l *DepLabeler) Default(ctx context.Context, obj runtime.Object) error {
	dep, ok := obj.(*appsv1.Deployment)

	if !ok {
		return fmt.Errorf("expected a Deployment object but got %T", obj)
	}

	log := logf.FromContext(ctx)
	if dep.Labels == nil {
		dep.Labels = map[string]string{}
	}
	// leaving empty value as it will be set by reconciler
	dep.Labels["owned-replica-set"] = ""
	log.Info("labeling deployment", "deployment name", dep.ObjectMeta.Name)

	return nil
}
