package providers

import (
	"context"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/sullivtr/k8s_platform/internal/modules"
	"github.com/sullivtr/k8s_platform/internal/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

// K8sApiProvider is a port for the kubernetes client.
// It provides an api for interacting with a given kubernetes cluster for the needs of khub
type K8sApiProvider struct {
	Session K8sSession
}

// Compile time proof of implementation
var _ IK8sProvider = (*K8sApiProvider)(nil)

// CloudAPIAWSSession represents the AWS session, and the account it belongs to
type K8sSession struct {
	SDK modules.K8sSDK
}

// InitK8sProvider will initialize the K8sApiProvider implementation.
func (p *ModuleProviders) InitK8sProvider() {
	var config *rest.Config
	var err error
	if !p.Config.K8sInCluster {
		var kcPath string
		if home := homedir.HomeDir(); home != "" {
			kcPath = filepath.Join(home, ".kube", "config")
		} else {
			log.Fatal().Msg("unable to locate home directory for kubeconfig path")
		}

		// use the current context in kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", kcPath)
		if err != nil {
			log.Fatal().Msg("unable to build kubernetes client config from kubeconfig file")
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatal().Msg("unable to load kubernetes client config from in-cluster config")
		}
	}

	// creates the clientset from the loaded config
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal().Msg("unable to load kubernetes clientset from config")
	}

	metricsClient, err := metricsclientset.NewForConfig(config)
	if err != nil {
		log.Fatal().Msg("unable to load kubernetes metrics clientset from config")
	}

	k8sSDK := modules.NewK8sSDK(clientset, metricsClient, config)
	p.K8sProvider = &K8sApiProvider{
		Session: K8sSession{
			SDK: k8sSDK,
		},
	}
}

func (p *K8sApiProvider) GetPods(namespaces []string) (any, error) {
	return p.Session.SDK.GetPods(namespaces)
}

func (p *K8sApiProvider) GetPod(namespace, podName string) (*v1.Pod, error) {
	return p.Session.SDK.GetPod(namespace, podName)
}

func (p *K8sApiProvider) DeletePod(namespace, podName string) error {
	return p.Session.SDK.DeletePod(namespace, podName)
}

func (p *K8sApiProvider) GetDeployments(namespaces []string) (any, error) {
	return p.Session.SDK.GetDeployments(namespaces)
}

func (p *K8sApiProvider) ScaleDeployment(namespace, name string, replicas int32) error {
	return p.Session.SDK.ScaleDeployment(namespace, name, replicas)
}

func (p *K8sApiProvider) GetDaemonsets(namespaces []string) (any, error) {
	return p.Session.SDK.GetDaemonsets(namespaces)
}

func (p *K8sApiProvider) GetReplicasets(namespaces []string) (any, error) {
	return p.Session.SDK.GetReplicasets(namespaces)
}

func (p *K8sApiProvider) GetStatefulsets(namespaces []string) (any, error) {
	return p.Session.SDK.GetStatefulsets(namespaces)
}

func (p *K8sApiProvider) GetJobs(namespaces []string) (any, error) {
	return p.Session.SDK.GetJobs(namespaces)
}

func (p *K8sApiProvider) GetCronJobs(namespaces []string) (any, error) {
	return p.Session.SDK.GetCronJobs(namespaces)
}

func (p *K8sApiProvider) GetClusterEvents() (any, error) {
	return p.Session.SDK.GetClusterEvents()
}

func (p *K8sApiProvider) GetServices(namespaces []string) (any, error) {
	return p.Session.SDK.GetServices(namespaces)
}

func (p *K8sApiProvider) GetIngresses(namespaces []string) (any, error) {
	return p.Session.SDK.GetIngresses(namespaces)
}

func (p *K8sApiProvider) GetNodes() (any, error) {
	return p.Session.SDK.GetNodes()
}

func (p *K8sApiProvider) GetConfigMaps(namespaces []string) (any, error) {
	return p.Session.SDK.GetConfigMaps(namespaces)
}

func (p *K8sApiProvider) RolloutRestartDeployment(namespace, name string) error {
	return p.Session.SDK.RolloutRestartDeployment(context.Background(), name, namespace)
}

func (p *K8sApiProvider) RolloutRestartDaemonSet(namespace, name string) error {
	return p.Session.SDK.RolloutRestartDaemonSet(context.Background(), name, namespace)
}

func (p *K8sApiProvider) RolloutRestartStatefulSet(namespace, name string) error {
	return p.Session.SDK.RolloutRestartStatefulSet(context.Background(), name, namespace)
}

func (p *K8sApiProvider) ExecutePodExecPlugin(namespace, podName string, plugin types.K8sPodExecPlugin) (string, string, error) {
	return p.Session.SDK.ExecutePodExecPlugin(namespace, podName, plugin)
}
