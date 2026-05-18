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

// MySQLConnectionReconciler reconciles a MySQLConnection object
type MySQLConnectionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=db.connect.local,resources=mysqlconnections,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=db.connect.local,resources=mysqlconnections/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=db.connect.local,resources=mysqlconnections/finalizers,verbs=update

func (r *MySQLConnectionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	conn := &dbv1.MySQLConnection{}
	if err := r.Get(ctx, req.NamespacedName, conn); err != nil {
		log.Error(err, "unable to fetch MySQLConnection")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if conn.DeletionTimestamp != nil {
		err := lib.CloseBackend(conn.Name)
		if err != nil {
			log.Error(err, "unable to close backend", "name", conn.Name)
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

	port := int32(3306)
	if conn.Spec.Port != 0 {
		port = conn.Spec.Port
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		conn.Spec.Username, conn.Spec.Password, conn.Spec.Host, port, conn.Spec.Database)

	if err := lib.RegisterBackend(conn.Name, "mysql", dsn); err != nil {
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

func (r *MySQLConnectionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dbv1.MySQLConnection{}).
		Complete(r)
}

// SQLiteConnectionReconciler reconciles a SQLiteConnection object
type SQLiteConnectionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=db.connect.local,resources=sqliteconnections,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=db.connect.local,resources=sqliteconnections/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=db.connect.local,resources=sqliteconnections/finalizers,verbs=update

func (r *SQLiteConnectionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	conn := &dbv1.SQLiteConnection{}
	if err := r.Get(ctx, req.NamespacedName, conn); err != nil {
		log.Error(err, "unable to fetch SQLiteConnection")
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

	if err := lib.RegisterBackend(conn.Name, "sqlite", conn.Spec.FilePath); err != nil {
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

func (r *SQLiteConnectionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dbv1.SQLiteConnection{}).
		Complete(r)
}
