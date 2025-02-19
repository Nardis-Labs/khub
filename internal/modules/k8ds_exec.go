package modules

import (
	"bytes"
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

func (sdk *K8sSDK) TakeTomcatThreadDump(namespace, podName string) (string, string, error) {
	cmd := []string{"/bin/bash", "/opt/reports/threaddump.sh"}

	eOptions := &v1.PodExecOptions{
		Command:   cmd,
		Container: "tomcat",
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
		VersionedParams(eOptions, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(sdk.restClientConfig, "POST", execReq.URL())
	if err != nil {
		return "", "", fmt.Errorf("%w Failed initialize remote executor for (tomcat thread dump) on %v/%v", err, namespace, podName)
	}

	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = exec.StreamWithContext(ctxWithTimeout, remotecommand.StreamOptions{
		Stdout: buf,
		Stderr: errBuf,
	})
	if err != nil {
		cancel()
		return "", "", fmt.Errorf("%w Failed executing command (tomcat thread dump) on %v/%v", err, namespace, podName)
	}
	return buf.String(), errBuf.String(), nil
}
