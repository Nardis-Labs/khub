package modules

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

var (
	testResourceOne = metav1.ObjectMeta{
		Name:      "test-resource",
		Namespace: "default",
	}

	testResourceTwo = metav1.ObjectMeta{
		Name:      "test-resource-two",
		Namespace: "other",
	}

	getDeployReturns = []runtime.Object{
		&appsv1.Deployment{
			ObjectMeta: testResourceOne,
		},
		&appsv1.Deployment{
			ObjectMeta: testResourceTwo,
		},
	}

	getPodReturns = []runtime.Object{
		&v1.Pod{
			ObjectMeta: testResourceOne,
		},
		&v1.Pod{
			ObjectMeta: testResourceTwo,
		},
	}

	getStatefulsetReturns = []runtime.Object{
		&appsv1.StatefulSet{
			ObjectMeta: testResourceOne,
		},
		&appsv1.StatefulSet{
			ObjectMeta: testResourceTwo,
		},
	}

	getReplicasetReturns = []runtime.Object{
		&appsv1.ReplicaSet{
			ObjectMeta: testResourceOne,
		},
		&appsv1.ReplicaSet{
			ObjectMeta: testResourceTwo,
		},
	}

	getDaemonsetReturns = []runtime.Object{
		&appsv1.DaemonSet{
			ObjectMeta: testResourceOne,
		},
		&appsv1.DaemonSet{
			ObjectMeta: testResourceTwo,
		},
	}

	getJobsReturns = []runtime.Object{
		&batchv1.Job{
			ObjectMeta: testResourceOne,
		},
		&batchv1.Job{
			ObjectMeta: testResourceTwo,
		},
	}

	getCronjobsReturns = []runtime.Object{
		&batchv1.CronJob{
			ObjectMeta: testResourceOne,
		},
		&batchv1.CronJob{
			ObjectMeta: testResourceTwo,
		},
	}
)

func TestGetResources(t *testing.T) {
	tests := []struct {
		name          string
		sdk           *K8sSDK
		namespaces    []string
		mockType      string
		expectedCount int
		expectedError error
	}{
		{
			name:          "Get deployments from specific namespaces",
			expectedCount: 2,
			sdk: &K8sSDK{
				client: fake.NewSimpleClientset(getDeployReturns...),
			},
			namespaces:    []string{},
			mockType:      "deployments",
			expectedError: nil,
		},
		{
			name:          "Get deployments from specific namespaces",
			expectedCount: 1,
			sdk: &K8sSDK{
				client: fake.NewSimpleClientset(getDeployReturns...),
			},
			namespaces:    []string{"default", "kube-system"},
			mockType:      "deployments",
			expectedError: nil,
		},
		{
			name:          "Get statefulsets from all namespaces",
			expectedCount: 2,
			sdk: &K8sSDK{
				client: fake.NewSimpleClientset(getStatefulsetReturns...),
			},
			namespaces:    []string{},
			mockType:      "statefulsets",
			expectedError: nil,
		},
		{
			name:          "Get statefulsets from all namespaces",
			expectedCount: 1,
			sdk: &K8sSDK{
				client: fake.NewSimpleClientset(getStatefulsetReturns...),
			},
			namespaces:    []string{"default", "kube-system"},
			mockType:      "statefulsets",
			expectedError: nil,
		},
		{
			name:          "Get replicasets from all namespaces",
			expectedCount: 2,
			sdk: &K8sSDK{
				client: fake.NewSimpleClientset(getReplicasetReturns...),
			},
			namespaces:    []string{},
			mockType:      "replicasets",
			expectedError: nil,
		},
		{
			name:          "Get replicasets from all namespaces",
			expectedCount: 1,
			sdk: &K8sSDK{
				client: fake.NewSimpleClientset(getReplicasetReturns...),
			},
			namespaces:    []string{"default", "kube-system"},
			mockType:      "replicasets",
			expectedError: nil,
		},
		{
			name:          "Get daemonsets from specific namespaces",
			expectedCount: 2,
			sdk: &K8sSDK{
				client: fake.NewSimpleClientset(getDaemonsetReturns...),
			},
			namespaces:    []string{},
			mockType:      "daemonsets",
			expectedError: nil,
		},
		{
			name:          "Get daemonsets from specific namespaces",
			expectedCount: 1,
			sdk: &K8sSDK{
				client: fake.NewSimpleClientset(getDaemonsetReturns...),
			},
			namespaces:    []string{"default", "kube-system"},
			mockType:      "daemonsets",
			expectedError: nil,
		},
		{
			name:          "Get pods from all namespaces",
			expectedCount: 2,
			sdk: &K8sSDK{
				client: fake.NewSimpleClientset(getPodReturns...),
			},
			namespaces:    []string{},
			mockType:      "pods",
			expectedError: nil,
		},
		{
			name:          "Get pods from all namespaces",
			expectedCount: 1,
			sdk: &K8sSDK{
				client: fake.NewSimpleClientset(getPodReturns...),
			},
			namespaces:    []string{"default", "kube-system"},
			mockType:      "pods",
			expectedError: nil,
		},
		{
			name:          "Get jobs from specific namespaces",
			expectedCount: 1,
			sdk: &K8sSDK{
				client: fake.NewSimpleClientset(getJobsReturns...),
			},
			namespaces:    []string{"default", "kube-system"},
			mockType:      "jobs",
			expectedError: nil,
		},
		{
			name:          "Get jobs from specific namespaces",
			expectedCount: 2,
			sdk: &K8sSDK{
				client: fake.NewSimpleClientset(getJobsReturns...),
			},
			namespaces:    []string{},
			mockType:      "jobs",
			expectedError: nil,
		},
		{
			name:          "Get cronjobs from all namespaces",
			expectedCount: 2,
			sdk: &K8sSDK{
				client: fake.NewSimpleClientset(getCronjobsReturns...),
			},
			namespaces:    []string{},
			mockType:      "cronjobs",
			expectedError: nil,
		},
		{
			name:          "Get cronjobs from all namespaces",
			expectedCount: 1,
			sdk: &K8sSDK{
				client: fake.NewSimpleClientset(getCronjobsReturns...),
			},
			namespaces:    []string{"default", "kube-system"},
			mockType:      "cronjobs",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fake client with a pre-existing pod
			switch tt.mockType {
			case "deployments":
				resources, err := tt.sdk.GetDeployments(tt.namespaces)
				if err != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				if len(resources) != tt.expectedCount {
					t.Errorf("expected %d %s, got %v", tt.expectedCount, tt.mockType, len(resources))
				}
			case "daemonsets":
				resources, err := tt.sdk.GetDaemonsets(tt.namespaces)
				if err != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				if len(resources) != tt.expectedCount {
					t.Errorf("expected %d %s, got %v", tt.expectedCount, tt.mockType, len(resources))
				}
			case "statefulsets":
				resources, err := tt.sdk.GetStatefulsets(tt.namespaces)
				if err != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				if len(resources) != tt.expectedCount {
					t.Errorf("expected %d %s, got %v", tt.expectedCount, tt.mockType, len(resources))
				}
			case "replicasets":
				resources, err := tt.sdk.GetReplicasets(tt.namespaces)
				if err != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				if len(resources) != tt.expectedCount {
					t.Errorf("expected %d %s, got %v", tt.expectedCount, tt.mockType, len(resources))
				}
			case "pods":
				resources, err := tt.sdk.GetPods(tt.namespaces)
				if err != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				if len(resources) != tt.expectedCount {
					t.Errorf("expected %d %s, got %v", tt.expectedCount, tt.mockType, len(resources))
				}
			case "jobs":
				resources, err := tt.sdk.GetJobs(tt.namespaces)
				if err != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				if len(resources) != tt.expectedCount {
					t.Errorf("expected %d %s, got %v", tt.expectedCount, tt.mockType, len(resources))
				}
			case "cronjobs":
				resources, err := tt.sdk.GetCronJobs(tt.namespaces)
				if err != nil && err.Error() != tt.expectedError.Error() {
					t.Errorf("expected error %v, got %v", tt.expectedError, err)
				}
				if len(resources) != tt.expectedCount {
					t.Errorf("expected %d %s, got %v", tt.expectedCount, tt.mockType, len(resources))
				}
			}

		})
	}
}
