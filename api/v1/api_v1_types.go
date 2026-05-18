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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PostgreSQLConnectionSpec defines the desired state of PostgreSQL connection
type PostgreSQLConnectionSpec struct {
	Host     string `json:"host"`
	Port     int32  `json:"port,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"` // Consider using SecretKeyRef in production
	Database string `json:"database"`
}

// PostgreSQLConnectionStatus defines the observed state
type PostgreSQLConnectionStatus struct {
	Connected     bool               `json:"connected"`
	Error         string             `json:"error,omitempty"`
	LastProbeTime metav1.Time        `json:"lastProbeTime,omitempty"`
	Conditions    []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// PostgreSQLConnection is the Schema for the postgresqlconnections API
type PostgreSQLConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PostgreSQLConnectionSpec   `json:"spec,omitempty"`
	Status            PostgreSQLConnectionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PostgreSQLConnectionList contains a list of PostgreSQLConnection
type PostgreSQLConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PostgreSQLConnection `json:"items"`
}

// MySQLConnectionSpec defines the desired state of MySQL connection
type MySQLConnectionSpec struct {
	Host     string `json:"host"`
	Port     int32  `json:"port,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// MySQLConnectionStatus defines the observed state
type MySQLConnectionStatus struct {
	Connected     bool               `json:"connected"`
	Error         string             `json:"error,omitempty"`
	LastProbeTime metav1.Time        `json:"lastProbeTime,omitempty"`
	Conditions    []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// MySQLConnection is the Schema for the mysqlconnections API
type MySQLConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MySQLConnectionSpec   `json:"spec,omitempty"`
	Status            MySQLConnectionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MySQLConnectionList contains a list of MySQLConnection
type MySQLConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MySQLConnection `json:"items"`
}

// SQLiteConnectionSpec defines the desired state of SQLite connection
type SQLiteConnectionSpec struct {
	FilePath string `json:"filePath"`
}

// SQLiteConnectionStatus defines the observed state
type SQLiteConnectionStatus struct {
	Connected     bool               `json:"connected"`
	Error         string             `json:"error,omitempty"`
	LastProbeTime metav1.Time        `json:"lastProbeTime,omitempty"`
	Conditions    []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SQLiteConnection is the Schema for the sqliteconnections API
type SQLiteConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SQLiteConnectionSpec   `json:"spec,omitempty"`
	Status            SQLiteConnectionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SQLiteConnectionList contains a list of SQLiteConnection
type SQLiteConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SQLiteConnection `json:"items"`
}

// KafkaConnectionSpec defines the desired state of Kafka connection
type KafkaConnectionSpec struct {
	Brokers []string `json:"brokers"`
}

// KafkaConnectionStatus defines the observed state
type KafkaConnectionStatus struct {
	Connected     bool               `json:"connected"`
	Error         string             `json:"error,omitempty"`
	LastProbeTime metav1.Time        `json:"lastProbeTime,omitempty"`
	Conditions    []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KafkaConnection is the Schema for the kafkaconnections API
type KafkaConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              KafkaConnectionSpec   `json:"spec,omitempty"`
	Status            KafkaConnectionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KafkaConnectionList contains a list of KafkaConnection
type KafkaConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KafkaConnection `json:"items"`
}

// SolaceConnectionSpec defines the desired state of Solace connection
type SolaceConnectionSpec struct {
	BrokerURL string `json:"brokerURL"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
}

// SolaceConnectionStatus defines the observed state
type SolaceConnectionStatus struct {
	Connected     bool               `json:"connected"`
	Error         string             `json:"error,omitempty"`
	LastProbeTime metav1.Time        `json:"lastProbeTime,omitempty"`
	Conditions    []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SolaceConnection is the Schema for the solaceconnections API
type SolaceConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SolaceConnectionSpec   `json:"spec,omitempty"`
	Status            SolaceConnectionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SolaceConnectionList contains a list of SolaceConnection
type SolaceConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SolaceConnection `json:"items"`
}

func init() {
	SchemeBuilder.Register(
		&PostgreSQLConnection{}, &PostgreSQLConnectionList{},
		&MySQLConnection{}, &MySQLConnectionList{},
		&SQLiteConnection{}, &SQLiteConnectionList{},
		&KafkaConnection{}, &KafkaConnectionList{},
		&SolaceConnection{}, &SolaceConnectionList{},
	)
}
