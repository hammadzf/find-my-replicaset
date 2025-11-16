package main

import (
	"os"

	appsv1 "k8s.io/api/apps/v1"

	"github.com/hammadzf/find-my-replicaset/internal/controller"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {
	ctrl.SetLogger(zap.New())
	entryLog := ctrl.Log.WithName("entrypoint")

	// create a manager
	mgr, err := manager.New(ctrl.GetConfigOrDie(), manager.Options{})
	if err != nil {
		entryLog.Error(err, "unable to set the controller manager")
		os.Exit(1)
	}

	// set up a controller
	entryLog.Info("setting up the controller")

	err = ctrl.
		NewControllerManagedBy(mgr).
		Named("toy-controller").
		For(&appsv1.Deployment{}).
		Owns(&appsv1.ReplicaSet{}).
		Complete(&controller.ReconcileDeployment{Client: mgr.GetClient()})

	if err != nil {
		entryLog.Error(err, "unable to create the controller")
		os.Exit(1)
	}

	// start manager
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		entryLog.Error(err, "unable to run the controller manager")
		os.Exit(1)
	}
}
