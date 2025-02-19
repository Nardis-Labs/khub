package modules

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/sullivtr/k8s_platform/internal/types"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

type K8sSDK struct {
	client           kubernetes.Interface
	metricsClient    metricsclientset.Interface
	restClientConfig *restclient.Config
}

func NewK8sSDK(client kubernetes.Interface, metricsClient metricsclientset.Interface, restClientConfig *restclient.Config) K8sSDK {
	return K8sSDK{client: client, metricsClient: metricsClient, restClientConfig: restClientConfig}
}

/*
/    PODS
*/

// GetPods returns a list of pods from the k8s cluster. If namespaces are specified in the K8sSDK,
// it returns pods from those namespaces. Otherwise, it returns pods from all namespaces.
func (sdk *K8sSDK) GetPods(namespaces []string) ([]v1.Pod, error) {
	if len(namespaces) == 0 {
		pods, err := sdk.client.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Error().Msgf("unable to get pods from all namespaces: %s", err.Error())
			return nil, err
		}
		return pods.Items, nil
	} else {
		allNSPods := []v1.Pod{}
		for _, ns := range namespaces {
			pods, err := sdk.client.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Error().Msgf("unable to get pods from %s namespace: %s", ns, err.Error())
			}
			allNSPods = append(allNSPods, pods.Items...)
		}
		return allNSPods, nil
	}
}

// GetPod will fetch a pod from the given namespace.
func (sdk *K8sSDK) GetPod(namespace, podName string) (*v1.Pod, error) {
	pod, err := sdk.client.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		log.Error().Msgf("unable to get pod, %s: %s", podName, err.Error())
	}
	return pod, err
}

// DeletePod will delete a pod from the given namespace,
func (sdk *K8sSDK) DeletePod(namespace, podName string) error {
	err := sdk.client.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		log.Error().Msgf("unable to delete pod, %s: %s", podName, err.Error())
	}
	return err
}

/*
/    DEPLOYMENTS
*/

// GetDeployments returns a list of deployments from the k8s cluster. If namespaces are specified in the K8sSDK,
// it returns deployments from those namespaces. Otherwise, it returns deployments from all namespaces.
func (sdk *K8sSDK) GetDeployments(namespaces []string) ([]appsv1.Deployment, error) {
	if len(namespaces) == 0 {
		deployments, err := sdk.client.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			log.Error().Msgf("unable to get deployments from all namespaces: %s", err.Error())
			return nil, err
		}
		return deployments.Items, nil
	} else {
		allNSDeployments := []appsv1.Deployment{}
		for _, ns := range namespaces {
			deployments, err := sdk.client.AppsV1().Deployments(ns).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Error().Msgf("unable to get deployments from %s namespace: %s", ns, err.Error())
			}
			allNSDeployments = append(allNSDeployments, deployments.Items...)
		}
		return allNSDeployments, nil
	}
}

// ScaleDeployment scales the specified deployment to the specified number of replicas.
func (sdk *K8sSDK) ScaleDeployment(namespace, deployName string, replicas int32) error {
	scale, err := sdk.client.AppsV1().Deployments(namespace).GetScale(context.TODO(), deployName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	scale.Spec.Replicas = replicas
	_, err = sdk.client.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deployName, scale, metav1.UpdateOptions{})
	return err
}

// RolloutRestartDeployment restarts the specified deployment by updating the "kubectl.kubernetes.io/restartedAt"
// annotation on the deployment's pod template. This causes Kubernetes to recreate all pods in the deployment.
func (sdk *K8sSDK) RolloutRestartDeployment(ctx context.Context, deployName, namespace string) error {
	deploy, err := sdk.client.AppsV1().Deployments(namespace).Get(ctx, deployName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if deploy.Spec.Template.ObjectMeta.Annotations == nil {
		deploy.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}
	deploy.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err = sdk.client.AppsV1().Deployments(namespace).Update(ctx, deploy, metav1.UpdateOptions{})
	return err
}

/*
/    DAEMONSETS
*/

// GetDaemonsets returns a list of daemonsets from the k8s cluster. If namespaces are specified in the K8sSDK,
// it returns daemonsets from those namespaces. Otherwise, it returns daemonsets from all namespaces.
func (sdk *K8sSDK) GetDaemonsets(namespaces []string) ([]appsv1.DaemonSet, error) {
	if len(namespaces) == 0 {
		daemonsets, err := sdk.client.AppsV1().DaemonSets("").List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			log.Error().Msgf("unable to get daemonsets from all namespaces: %s", err.Error())
			return nil, err
		}
		return daemonsets.Items, nil
	} else {
		allNSDaemonsets := []appsv1.DaemonSet{}
		for _, ns := range namespaces {
			daemonsets, err := sdk.client.AppsV1().DaemonSets(ns).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Error().Msgf("unable to get daemonsets from %s namespace: %s", ns, err.Error())
			}
			allNSDaemonsets = append(allNSDaemonsets, daemonsets.Items...)
		}
		return allNSDaemonsets, nil
	}
}

// RolloutRestartDaemonSet restarts the specified DaemonSet by updating the "kubectl.kubernetes.io/restartedAt"
// annotation on the DaemonSet's pod template. This causes Kubernetes to recreate all pods in the DaemonSet.
func (sdk *K8sSDK) RolloutRestartDaemonSet(ctx context.Context, daemonSetName, namespace string) error {
	daemonSet, err := sdk.client.AppsV1().DaemonSets(namespace).Get(ctx, daemonSetName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if daemonSet.Spec.Template.ObjectMeta.Annotations == nil {
		daemonSet.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}
	daemonSet.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err = sdk.client.AppsV1().DaemonSets(namespace).Update(ctx, daemonSet, metav1.UpdateOptions{})
	return err
}

/*
/    REPLICASETS
*/

// GetReplicasets returns a list of replicasets from the k8s cluster. If namespaces are specified in the K8sSDK,
// it returns replicasets from those namespaces. Otherwise, it returns replicasets from all namespaces.
func (sdk *K8sSDK) GetReplicasets(namespaces []string) ([]appsv1.ReplicaSet, error) {
	if len(namespaces) == 0 {
		replicasets, err := sdk.client.AppsV1().ReplicaSets("").List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			log.Error().Msgf("unable to get replicasets from all namespaces: %s", err.Error())
			return nil, err
		}
		return replicasets.Items, nil
	} else {
		allNSReplicasets := []appsv1.ReplicaSet{}
		for _, ns := range namespaces {
			replicasets, err := sdk.client.AppsV1().ReplicaSets(ns).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Error().Msgf("unable to get replicasets from %s namespace: %s", ns, err.Error())
			}
			allNSReplicasets = append(allNSReplicasets, replicasets.Items...)
		}
		return allNSReplicasets, nil
	}
}

/*
/    STATEFULSETS
*/

// GetStatefulsets returns a list of statefulsets from the k8s cluster. If namespaces are specified in the K8sSDK,
// it returns statefulsets from those namespaces. Otherwise, it returns statefulsets from all namespaces.
func (sdk *K8sSDK) GetStatefulsets(namespaces []string) ([]appsv1.StatefulSet, error) {
	if len(namespaces) == 0 {
		statefulsets, err := sdk.client.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			log.Error().Msgf("unable to get statefulsets from all namespaces: %s", err.Error())
			return nil, err
		}
		return statefulsets.Items, nil
	} else {
		allNSStatefulsets := []appsv1.StatefulSet{}
		for _, ns := range namespaces {
			statefulsets, err := sdk.client.AppsV1().StatefulSets(ns).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Error().Msgf("unable to get statefulsets from %s namespace: %s", ns, err.Error())
			}
			allNSStatefulsets = append(allNSStatefulsets, statefulsets.Items...)
		}
		return allNSStatefulsets, nil
	}
}

// RolloutRestartStatefulSet restarts the specified StatefulSet by updating the "kubectl.kubernetes.io/restartedAt"
// annotation on the StatefulSet's pod template. This causes Kubernetes to recreate all pods in the StatefulSet.
func (sdk *K8sSDK) RolloutRestartStatefulSet(ctx context.Context, statefulSetName, namespace string) error {
	statefulSet, err := sdk.client.AppsV1().StatefulSets(namespace).Get(ctx, statefulSetName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	if statefulSet.Spec.Template.ObjectMeta.Annotations == nil {
		statefulSet.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}
	statefulSet.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)

	_, err = sdk.client.AppsV1().StatefulSets(namespace).Update(ctx, statefulSet, metav1.UpdateOptions{})
	return err
}

/*
/    JOBS
*/

// GetJobs returns a list of jobs from the k8s cluster. If namespaces are specified in the K8sSDK,
// it returns jobs from those namespaces. Otherwise, it returns jobs from all namespaces.
func (sdk *K8sSDK) GetJobs(namespaces []string) ([]batchv1.Job, error) {
	if len(namespaces) == 0 {
		jobs, err := sdk.client.BatchV1().Jobs("").List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			log.Error().Msgf("unable to get jobs from all namespaces: %s", err.Error())
			return nil, err
		}
		return jobs.Items, nil
	} else {
		allNSJobs := []batchv1.Job{}
		for _, ns := range namespaces {
			jobs, err := sdk.client.BatchV1().Jobs(ns).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Error().Msgf("unable to get jobs from %s namespace: %s", ns, err.Error())
			}
			allNSJobs = append(allNSJobs, jobs.Items...)
		}
		return allNSJobs, nil
	}
}

/*
/    CRONJOBS
*/

// GetCronJobs returns a list of cronjobs from the k8s cluster. If namespaces are specified in the K8sSDK,
// it returns cronjobs from those namespaces. Otherwise, it returns cronjobs from all namespaces.
func (sdk *K8sSDK) GetCronJobs(namespaces []string) ([]batchv1.CronJob, error) {
	if len(namespaces) == 0 {
		cronjobs, err := sdk.client.BatchV1().CronJobs("").List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			log.Error().Msgf("unable to get cronjobs from all namespaces: %s", err.Error())
			return nil, err
		}
		return cronjobs.Items, nil
	} else {
		allNSCronJobs := []batchv1.CronJob{}
		for _, ns := range namespaces {
			jobs, err := sdk.client.BatchV1().CronJobs(ns).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Error().Msgf("unable to get cronjobs from %s namespace: %s", ns, err.Error())
			}
			allNSCronJobs = append(allNSCronJobs, jobs.Items...)
		}
		return allNSCronJobs, nil
	}
}

/*
/    SERVICES
*/

// GetServices returns a list of services from the k8s cluster. If namespaces are specified in the K8sSDK,
// it returns services from those namespaces. Otherwise, it returns services from all namespaces.
func (sdk *K8sSDK) GetServices(namespaces []string) ([]v1.Service, error) {
	if len(namespaces) == 0 {
		services, err := sdk.client.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			log.Error().Msgf("unable to get services from all namespaces: %s", err.Error())
			return nil, err
		}
		return services.Items, nil
	} else {
		allNSServices := []v1.Service{}
		for _, ns := range namespaces {
			services, err := sdk.client.CoreV1().Services(ns).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Error().Msgf("unable to get services from %s namespace: %s", ns, err.Error())
			}
			allNSServices = append(allNSServices, services.Items...)
		}
		return allNSServices, nil
	}
}

/*
/    INGRESS
*/

// GetIngresses returns a list of ingresses from the k8s cluster. If namespaces are specified in the K8sSDK,
// it returns ingresses from those namespaces. Otherwise, it returns ingresses from all namespaces.
func (sdk *K8sSDK) GetIngresses(namespaces []string) ([]networkingv1.Ingress, error) {
	if len(namespaces) == 0 {
		ingresses, err := sdk.client.NetworkingV1().Ingresses("").List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			log.Error().Msgf("unable to get ingresses from all namespaces: %s", err.Error())
			return nil, err
		}
		return ingresses.Items, nil
	} else {
		allNSIngresses := []networkingv1.Ingress{}
		for _, ns := range namespaces {
			ingresses, err := sdk.client.NetworkingV1().Ingresses(ns).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Error().Msgf("unable to get ingresses from %s namespace: %s", ns, err.Error())
			}
			allNSIngresses = append(allNSIngresses, ingresses.Items...)
		}
		return allNSIngresses, nil
	}
}

/*
/    CONIFIGMAPS
*/

// GetConfigMaps returns a list of configMaps from the k8s cluster. If namespaces are specified in the K8sSDK,
// it returns configMaps from those namespaces. Otherwise, it returns configMaps from all namespaces.
func (sdk *K8sSDK) GetConfigMaps(namespaces []string) ([]v1.ConfigMap, error) {
	if len(namespaces) == 0 {
		configMaps, err := sdk.client.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{})

		if err != nil {
			log.Error().Msgf("unable to get config maps from all namespaces: %s", err.Error())
			return nil, err
		}
		return configMaps.Items, nil
	} else {
		allNSConfigMaps := []v1.ConfigMap{}
		for _, ns := range namespaces {
			configMaps, err := sdk.client.CoreV1().ConfigMaps(ns).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Error().Msgf("unable to get config maps from %s namespace: %s", ns, err.Error())
			}
			allNSConfigMaps = append(allNSConfigMaps, configMaps.Items...)
		}
		return allNSConfigMaps, nil
	}
}

/*
/    NODES
*/

// GetNodes returns a list of nodes from the k8s cluster.
func (sdk *K8sSDK) GetNodes() ([]types.K8sNodeWrapper, error) {
	k8sNodes := []types.K8sNodeWrapper{}
	nodes, err := sdk.client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Error().Msgf("unable to get nodes: %s", err.Error())
		return nil, err
	}

	nodeMetrics, err := sdk.topNode()
	if err != nil || len(nodeMetrics) == 0 {
		// If the metrics API fails, just proceed and return the node data.
		for _, node := range nodes.Items {
			k8sNodes = append(k8sNodes, types.K8sNodeWrapper{Node: node})
		}
		return k8sNodes, nil
	}

	for _, node := range nodes.Items {
		for _, metrics := range nodeMetrics {
			if node.ObjectMeta.Name == metrics.Name {
				k8sNodes = append(k8sNodes, types.K8sNodeWrapper{Node: node, Metrics: metrics})
			}
		}
	}

	return k8sNodes, nil
}

/*
/    Cluster Events
*/

// GetClusterEvents returns a list of events from the k8s cluster. It wraps each event in a K8sEventWrapper,
// which includes additional information such as the interval and the object involved in the event.
func (sdk *K8sSDK) GetClusterEvents() ([]types.K8sEventWrapper, error) {
	wrappedEvents := []types.K8sEventWrapper{}
	events, err := sdk.client.CoreV1().Events("").List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		log.Error().Msgf("unable to get events cluster: %s", err.Error())
		return nil, err
	}

	for _, e := range events.Items {
		event := new(types.K8sEventWrapper)
		event.Interval = types.GetInterval(e)
		event.Object = fmt.Sprintf("%s/%s", e.InvolvedObject.Kind, e.InvolvedObject.Name)
		event.Event = e
		wrappedEvents = append(wrappedEvents, *event)
	}
	return wrappedEvents, nil
}
