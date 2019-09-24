package health

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kyma-project/helm-broker/pkg/apis/addons/v1alpha1"
	"github.com/sirupsen/logrus"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// HandleHealth starts health endpoints
func HandleHealth(port string, client client.Client, lg *logrus.Entry) {
	rtr := mux.NewRouter()
	rtr.HandleFunc(HandleControllerLive(client, lg)).Methods("GET")
	rtr.HandleFunc(HandleReady()).Methods("GET")
	http.ListenAndServe(port, rtr)
}

// HandleBrokerLive returns handler for livness probes for broker container
func HandleBrokerLive() (string, func(w http.ResponseWriter, req *http.Request)) {
	return "/live", healthResponse()
}

// HandleControllerLive returns handler for livness probes for controller container
func HandleControllerLive(client client.Client, lg *logrus.Entry) (string, func(w http.ResponseWriter, req *http.Request)) {
	return "/live", runFullControllersCycle(client, lg)
}

// HandleReady returns handler for readiness proves
func HandleReady() (string, func(w http.ResponseWriter, req *http.Request)) {
	return "/ready", healthResponse()
}

func healthResponse() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}
}

func runFullControllersCycle(client client.Client, lg *logrus.Entry) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := runAddonsConfigurationControllerCycle(client, lg); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := runClusterAddonsConfigurationControllerCycle(client, lg); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	}
}

func runAddonsConfigurationControllerCycle(client client.Client, lg *logrus.Entry) error {
	probeName := "liveness-probe"
	probeNamespace := "default"

	addonsConfiguration := &v1alpha1.AddonsConfiguration{
		ObjectMeta: v1.ObjectMeta{
			Name:      probeName,
			Namespace: probeNamespace,
		},
		Spec: v1alpha1.AddonsConfigurationSpec{
			CommonAddonsConfigurationSpec: v1alpha1.CommonAddonsConfigurationSpec{
				Repositories: []v1alpha1.SpecRepository{{URL: ""}},
			},
		},
	}

	lg.Infof("[liveness-probe] Creating liveness probe addonsConfiguration in %q namespace", probeNamespace)
	err := client.Create(context.TODO(), addonsConfiguration)
	if err != nil {
		lg.Infof("[liveness-probe] Cannot create liveness probe addonsConfiguration: %s", err)
		return err
	}

	lg.Info("[liveness-probe] Waiting for liveness probe addonsConfiguration desirable status")
	err = wait.Poll(1*time.Second, 10*time.Second, func() (done bool, err error) {
		key := types.NamespacedName{Name: probeName, Namespace: probeNamespace}
		err = client.Get(context.TODO(), key, addonsConfiguration)
		if apierrors.IsNotFound(err) {
			lg.Info("[liveness-probe] Liveness probe addonsConfiguration not found")
			return false, nil
		}
		if err != nil {
			return false, err
		}

		if len(addonsConfiguration.Status.Repositories) != 1 {
			lg.Info("[liveness-probe] Liveness probe addonsConfiguration repositories status not set")
			return false, nil
		}

		status := addonsConfiguration.Status.Repositories[0].Status
		reason := addonsConfiguration.Status.Repositories[0].Reason
		if status == v1alpha1.RepositoryStatusFailed {
			if reason == v1alpha1.RepositoryURLFetchingError {
				lg.Info("[liveness-probe] Liveness probe addonsConfiguration has achieved the desired status")
				return true, nil
			}
		}

		lg.Infof("[liveness-probe] Liveness probe addonsConfiguration current status: %s: %s", status, reason)
		return false, nil
	})
	if err != nil {
		lg.Infof("[liveness-probe] Waiting for liveness probe addonsConfiguration failed: %s", err)
		return err
	}

	lg.Info("[liveness-probe] Removing liveness probe addonsConfiguration")
	err = client.Delete(context.TODO(), addonsConfiguration)
	if err != nil {
		lg.Infof("[liveness-probe] Cannot delete liveness probe addonsConfiguration: %s", err)
		return err
	}

	lg.Info("[liveness-probe] AddonsConfiguration controller is live")
	return nil
}

func runClusterAddonsConfigurationControllerCycle(client client.Client, lg *logrus.Entry) error {
	probeName := "liveness-probe"

	clusterAddonsConfiguration := &v1alpha1.ClusterAddonsConfiguration{
		ObjectMeta: v1.ObjectMeta{
			Name: probeName,
		},
		Spec: v1alpha1.ClusterAddonsConfigurationSpec{
			CommonAddonsConfigurationSpec: v1alpha1.CommonAddonsConfigurationSpec{
				Repositories: []v1alpha1.SpecRepository{{URL: ""}},
			},
		},
	}

	lg.Info("[liveness-probe] Creating liveness probe clusterAddonsConfiguration")
	err := client.Create(context.TODO(), clusterAddonsConfiguration)
	if err != nil {
		lg.Infof("[liveness-probe] Cannot create liveness probe clusterAddonsConfiguration: %s", err)
		return err
	}

	lg.Info("[liveness-probe] Waiting for liveness probe clusterAddonsConfiguration desirable status")
	err = wait.Poll(1*time.Second, 10*time.Second, func() (done bool, err error) {
		key := types.NamespacedName{Name: probeName, Namespace: v1.NamespaceAll}
		err = client.Get(context.TODO(), key, clusterAddonsConfiguration)
		if apierrors.IsNotFound(err) {
			lg.Info("[liveness-probe] Liveness probe clusterAddonsConfiguration not found")
			return false, nil
		}
		if err != nil {
			return false, err
		}

		if len(clusterAddonsConfiguration.Status.Repositories) != 1 {
			lg.Info("[liveness-probe] Liveness probe addonsConfiguration repositories status not set")
			return false, nil
		}

		status := clusterAddonsConfiguration.Status.Repositories[0].Status
		reason := clusterAddonsConfiguration.Status.Repositories[0].Reason
		if status == v1alpha1.RepositoryStatusFailed {
			if reason == v1alpha1.RepositoryURLFetchingError {
				lg.Info("[liveness-probe] Liveness probe clusterAddonsConfiguration has achieved the desired status")
				return true, nil
			}
		}

		lg.Infof("[liveness-probe] Liveness probe clusterAddonsConfiguration current status: %s: %s", status, reason)
		return false, nil
	})
	if err != nil {
		lg.Infof("[liveness-probe] Waiting for liveness probe clusterAddonsConfiguration failed: %s", err)
		return err
	}

	lg.Info("[liveness-probe] Removing liveness probe clusterAddonsConfiguration")
	err = client.Delete(context.TODO(), clusterAddonsConfiguration)
	if err != nil {
		lg.Infof("[liveness-probe] Cannot delete liveness probe clusterAddonsConfiguration: %s", err)
		return err
	}

	lg.Info("[liveness-probe] ClusterAddonsConfiguration controller is live")
	return nil
}
