package providers

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sullivtr/k8s_platform/internal/modules"
)

// StartDataSink starts the data collection process for various Kubernetes resources. It periodically collects data for
// pods, deployments, daemonsets, replicasets, statefulsets, jobs, cronjobs, services, ingresses, nodes,cluster events, and
// the resource tree map. The collection interval is specified by the 'interval' parameter.
//
// The method uses the 'Poll' function from the 'modules' package to perform the collection. The 'Poll' function
// takes a context, an interval, and a function to call at each interval. The function to be called are one of
// the 'collect' methods of the 'ModuleProviders' struct.
//
// After starting the collection for all resources, the method waits for 10 seconds before starting the collection
// for the resource tree map. This is to ensure that the latest data for all resources is available when building the resource tree map.
//
// Note: This method runs indefinitely until the provided context is cancelled. It should typically be run in a separate goroutine.
func (p *ModuleProviders) StartDataSink(ctx context.Context, intervaSecond int) {
	// Todo: Use a bounded pool for goroutines for more control over the aysnc processes.
	fastInterval := time.Duration(intervaSecond) * time.Second
	slowInterval := time.Duration(intervaSecond+20) * time.Second

	modules.Poll(ctx, fastInterval, p.collectPods)
	modules.Poll(ctx, fastInterval, p.collectDeployments)
	modules.Poll(ctx, fastInterval, p.collectDaemonsets)
	modules.Poll(ctx, fastInterval, p.collectReplicasets)
	modules.Poll(ctx, fastInterval, p.collectStatefulsets)
	modules.Poll(ctx, fastInterval, p.collectJobs)
	modules.Poll(ctx, fastInterval, p.collectCronJobs)
	modules.Poll(ctx, fastInterval, p.collectServices)
	modules.Poll(ctx, fastInterval, p.collectIngresses)
	modules.Poll(ctx, fastInterval, p.collectConfigMaps)
	modules.Poll(ctx, slowInterval, p.collectNodes)
	modules.Poll(ctx, fastInterval, p.collectClusterEvents)
	time.Sleep(10 * time.Second)
}

func (p *ModuleProviders) collectPods() {
	p.collectK8sResource(p.K8sProvider.GetPods, "pods")
}

func (p *ModuleProviders) collectDeployments() {
	p.collectK8sResource(p.K8sProvider.GetDeployments, "deployments")
}

func (p *ModuleProviders) collectDaemonsets() {
	p.collectK8sResource(p.K8sProvider.GetDaemonsets, "daemonsets")
}

func (p *ModuleProviders) collectReplicasets() {
	p.collectK8sResource(p.K8sProvider.GetReplicasets, "replicasets")
}

func (p *ModuleProviders) collectStatefulsets() {
	p.collectK8sResource(p.K8sProvider.GetStatefulsets, "statefulsets")
}

func (p *ModuleProviders) collectJobs() {
	p.collectK8sResource(p.K8sProvider.GetJobs, "jobs")
}

func (p *ModuleProviders) collectCronJobs() {
	p.collectK8sResource(p.K8sProvider.GetCronJobs, "cronjobs")
}

func (p *ModuleProviders) collectServices() {
	p.collectK8sResource(p.K8sProvider.GetServices, "services")
}

func (p *ModuleProviders) collectIngresses() {
	p.collectK8sResource(p.K8sProvider.GetIngresses, "ingresses")
}

func (p *ModuleProviders) collectConfigMaps() {
	p.collectK8sResource(p.K8sProvider.GetConfigMaps, "configmaps")
}

func (p *ModuleProviders) collectNodes() {
	p.collectGlobalK8sResource(p.K8sProvider.GetNodes, "nodes")
}

func (p *ModuleProviders) collectClusterEvents() {
	p.collectGlobalK8sResource(p.K8sProvider.GetClusterEvents, "clusterevents")
}

// collectK8sResource collects namespaced kubernetes resources like pods, deployments, daemonsets, replicasets, statefulsets,
// jobs, cronjobs, services, ingresses, and configmaps.
func (p *ModuleProviders) collectK8sResource(getFunc func(namespaces []string) (any, error), resource string) {
	dac, err := p.StorageProvider.GetDynamicAppConfig()
	if err != nil {
		log.Error().Msgf("unable to get dynamic app config: %s", err.Error())
	}
	log.Info().Msgf("collecting %s data", resource)
	res, err := getFunc(dac.Data.K8sClusterNamespaces)
	if err != nil {
		log.Error().Msgf("unable to get %s: %s", resource, err.Error())
	}
	key := fmt.Sprintf("%s_%s", dac.Data.K8sClusterName, resource)
	if err := p.CacheProvider.Put(key, res); err != nil {
		log.Error().Msgf("unable to collect %s data: %s", resource, err.Error())
	}
}

// collectGlobalK8sResource collects global kubernetes resources like nodes and cluster events.
// the difference is that this method does not require namespaces to be specified.
func (p *ModuleProviders) collectGlobalK8sResource(getFunc func() (any, error), resource string) {
	dac, err := p.StorageProvider.GetDynamicAppConfig()
	if err != nil {
		log.Error().Msgf("unable to get dynamic app config: %s", err.Error())
	}
	log.Info().Msgf("collecting %s data", resource)
	res, err := getFunc()
	if err != nil {
		log.Error().Msgf("unable to get %s: %s", resource, err.Error())
	}
	key := fmt.Sprintf("%s_%s", dac.Data.K8sClusterName, resource)
	if err := p.CacheProvider.Put(key, res); err != nil {
		log.Error().Msgf("unable to collect %s data: %s", resource, err.Error())
	}
}
