// Copyright Aeraki Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lazyxds

import (
	"context"
	"fmt"

	istioclient "istio.io/client-go/pkg/clientset/versioned"

	"k8s.io/apimachinery/pkg/api/errors"

	v1 "k8s.io/api/core/v1"
	controllerclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var (
	namespacePredicates = predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return true
		},
	}
)

// namespaceController creates bootstrap configMap for sidecar proxies
type namespaceController struct {
	controllerclient.Client
	istioClient *istioclient.Clientset
}

// Reconcile watch namespace change and create bootstrap configmap for sidecar proxies
func (c *namespaceController) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	controllerLog.Infof("reconcile namespace: %s", request.Name)

	ns := &v1.Namespace{}
	err := c.Get(ctx, request.NamespacedName, ns)
	if errors.IsNotFound(err) {
		return reconcile.Result{}, nil
	}

	if err != nil {
		return reconcile.Result{}, fmt.Errorf("could not fetch Namespace: %+v", err)
	}

	if c.shouldHandle(ns) {
		c.createDefaultSidecar(ctx, ns.Name)
	}
	return reconcile.Result{}, nil
}

// AddNamespaceController adds namespaceController
func AddNamespaceController(mgr manager.Manager, istioClient *istioclient.Clientset) error {
	namespaceCtrl := &namespaceController{
		Client:      mgr.GetClient(),
		istioClient: istioClient,
	}
	c, err := controller.New("lazyxds-namespace-controller", mgr,
		controller.Options{Reconciler: namespaceCtrl})
	if err != nil {
		return err
	}
	// Watch for changes on Namespace CRD
	err = c.Watch(&source.Kind{Type: &v1.Namespace{}}, &handler.EnqueueRequestForObject{},
		namespacePredicates)
	if err != nil {
		return err
	}

	controllerLog.Infof("NamespaceController (used to for lazyxds) registered")
	return nil
}

func (c *namespaceController) shouldHandle(ns *v1.Namespace) bool {
	if ns.Annotations["lazy-xds"] == "true" {
		return true
	}
	return false
}
