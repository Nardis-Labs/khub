package modules

import (
	"context"
	"errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
	metricsV1beta1api "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

var (
	supportedMetricsAPIVersions = []string{
		"v1beta1",
	}
)

func supportedMetricsAPIVersionAvailable(discoveredAPIGroups *metav1.APIGroupList) bool {
	for _, discoveredAPIGroup := range discoveredAPIGroups.Groups {
		if discoveredAPIGroup.Name != metricsapi.GroupName {
			continue
		}
		for _, version := range discoveredAPIGroup.Versions {
			for _, supportedVersion := range supportedMetricsAPIVersions {
				if version.Version == supportedVersion {
					return true
				}
			}
		}
	}
	return false
}

func (sdk *K8sSDK) topNode() ([]metricsapi.NodeMetrics, error) {
	var err error
	selector := labels.Everything()

	apiGroups, err := sdk.client.Discovery().ServerGroups()
	if err != nil {
		return nil, err
	}

	metricsAPIAvailable := supportedMetricsAPIVersionAvailable(apiGroups)

	if !metricsAPIAvailable {
		return nil, errors.New("metrics API not available")
	}

	metrics, err := getNodeMetricsFromMetricsAPI(sdk.metricsClient, selector)
	if err != nil {
		return nil, err
	}

	if len(metrics.Items) == 0 {
		return nil, errors.New("metrics not available yet")
	}

	return metrics.Items, nil
}

func getNodeMetricsFromMetricsAPI(metricsClient metricsclientset.Interface, selector labels.Selector) (*metricsapi.NodeMetricsList, error) {
	mc := metricsClient.MetricsV1beta1()
	nm := mc.NodeMetricses()
	versionedMetrics, err := nm.List(context.TODO(), metav1.ListOptions{LabelSelector: selector.String()})
	if err != nil {
		return nil, err
	}

	metrics := &metricsapi.NodeMetricsList{}
	err = metricsV1beta1api.Convert_v1beta1_NodeMetricsList_To_metrics_NodeMetricsList(versionedMetrics, metrics, nil)
	if err != nil {
		return nil, err
	}
	return metrics, nil
}
