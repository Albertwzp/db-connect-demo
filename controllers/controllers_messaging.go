/*
Copyright 2026.

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

package controllers

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	dbv1 "db-connect-demo/api/v1"
	"db-connect-demo/lib"
)

// KafkaConnectionReconciler reconciles a KafkaConnection object
type KafkaConnectionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=db.connect.local,resources=kafkaconnections,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=db.connect.local,resources=kafkaconnections/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=db.connect.local,resources=kafkaconnections/finalizers,verbs=update

func (r *KafkaConnectionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	conn := &dbv1.KafkaConnection{}
	if err := r.Get(ctx, req.NamespacedName, conn); err != nil {
		log.Error(err, "unable to fetch KafkaConnection")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if conn.DeletionTimestamp != nil {
		err := lib.CloseBackend(conn.Name)
		if err != nil {
			log.Error(err, "unable to close backend")
		}
		conn.Status.Connected = false
		meta.SetStatusCondition(&conn.Status.Conditions, metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			ObservedGeneration: conn.GetGeneration(),
			Reason:             "Deleted",
			Message:            "Connection closed",
			LastTransitionTime: metav1.Now(),
		})
		if err := r.Status().Update(ctx, conn); err != nil {
			log.Error(err, "unable to update connection status")
		}
		return ctrl.Result{}, nil
	}

	brokerList := fmt.Sprintf("[%s]", conn.Spec.Brokers)
	if err := lib.RegisterBackend(conn.Name, "kafka", brokerList); err != nil {
		log.Error(err, "unable to register backend", "name", conn.Name)
		conn.Status.Connected = false
		conn.Status.Error = err.Error()
		lib.MarkBackendFailed(conn.Name, err)
		meta.SetStatusCondition(&conn.Status.Conditions, metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			ObservedGeneration: conn.GetGeneration(),
			Reason:             "ConnectionFailed",
			Message:            err.Error(),
			LastTransitionTime: metav1.Now(),
		})
	} else {
		conn.Status.Connected = true
		conn.Status.Error = ""
		meta.SetStatusCondition(&conn.Status.Conditions, metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			ObservedGeneration: conn.GetGeneration(),
			Reason:             "ConnectionSucceeded",
			Message:            "Backend registered successfully",
			LastTransitionTime: metav1.Now(),
		})
	}
	conn.Status.LastProbeTime = metav1.NewTime(time.Now())

	if err := r.Status().Update(ctx, conn); err != nil {
		log.Error(err, "unable to update connection status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func (r *KafkaConnectionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dbv1.KafkaConnection{}).
		Complete(r)
}

// SolaceConnectionReconciler reconciles a SolaceConnection object
type SolaceConnectionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=db.connect.local,resources=solaceconnections,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=db.connect.local,resources=solaceconnections/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=db.connect.local,resources=solaceconnections/finalizers,verbs=update

func (r *SolaceConnectionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	conn := &dbv1.SolaceConnection{}
	if err := r.Get(ctx, req.NamespacedName, conn); err != nil {
		log.Error(err, "unable to fetch SolaceConnection")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if conn.DeletionTimestamp != nil {
		err := lib.CloseBackend(conn.Name)
		if err != nil {
			log.Error(err, "unable to close backend")
		}
		conn.Status.Connected = false
		meta.SetStatusCondition(&conn.Status.Conditions, metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			ObservedGeneration: conn.GetGeneration(),
			Reason:             "Deleted",
			Message:            "Connection closed",
			LastTransitionTime: metav1.Now(),
		})
		if err := r.Status().Update(ctx, conn); err != nil {
			log.Error(err, "unable to update connection status")
		}
		return ctrl.Result{}, nil
	}

	dsn := fmt.Sprintf("%s|%s|%s", conn.Spec.BrokerURL, conn.Spec.Username, conn.Spec.Password)
	if err := lib.RegisterBackend(conn.Name, "solace", dsn); err != nil {
		log.Error(err, "unable to register backend", "name", conn.Name)
		conn.Status.Connected = false
		conn.Status.Error = err.Error()
		lib.MarkBackendFailed(conn.Name, err)
		meta.SetStatusCondition(&conn.Status.Conditions, metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionFalse,
			ObservedGeneration: conn.GetGeneration(),
			Reason:             "ConnectionFailed",
			Message:            err.Error(),
			LastTransitionTime: metav1.Now(),
		})
	} else {
		conn.Status.Connected = true
		conn.Status.Error = ""
		meta.SetStatusCondition(&conn.Status.Conditions, metav1.Condition{
			Type:               "Ready",
			Status:             metav1.ConditionTrue,
			ObservedGeneration: conn.GetGeneration(),
			Reason:             "ConnectionSucceeded",
			Message:            "Backend registered successfully",
			LastTransitionTime: metav1.Now(),
		})
	}
	conn.Status.LastProbeTime = metav1.NewTime(time.Now())

	if err := r.Status().Update(ctx, conn); err != nil {
		log.Error(err, "unable to update connection status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
}

func (r *SolaceConnectionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dbv1.SolaceConnection{}).
		Complete(r)
}
