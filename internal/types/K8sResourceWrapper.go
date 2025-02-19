package types

import (
	v1 "k8s.io/api/core/v1"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
)

type K8sResourceWrapper struct {
	Data  any  `json:"data"`
	Write bool `json:"write"`
}

type K8sNodeWrapper struct {
	Node    v1.Node                `json:"node"`
	Metrics metricsapi.NodeMetrics `json:"metrics"`
}
