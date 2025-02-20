package modules

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sullivtr/k8s_platform/internal/types"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

func (sdk *K8sSDK) ExecutePodExecPlugin(namespace, podName string, plugin types.K8sPodExecPlugin) (string, string, error) {
	cmd := strings.Split(plugin.Command, " ")
	execOptions := &v1.PodExecOptions{
		Command:   cmd,
		Container: plugin.Container,
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       true,
	}

	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	execReq := sdk.client.CoreV1().RESTClient().
		Post().
		Namespace(namespace).
		Resource("pods").
		Name(podName).
		SubResource("exec").
		VersionedParams(execOptions, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(sdk.restClientConfig, "POST", execReq.URL())
	if err != nil {
		return "", "", fmt.Errorf("%w Failed initialize remote executor for (%s: %s) on %v/%v", err, plugin.Container, plugin.Command, namespace, podName)
	}

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = exec.StreamWithContext(ctxWithTimeout, remotecommand.StreamOptions{
		Stdout: buf,
		Stderr: errBuf,
	})
	if err != nil {
		cancel()
		return "", "", fmt.Errorf("%w Failed executing command (%s) on %v/%v", err, plugin.Command, namespace, podName)
	}
	return buf.String(), errBuf.String(), nil
}
