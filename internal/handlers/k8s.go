package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"github.com/sullivtr/k8s_platform/internal/providers"
	"github.com/sullivtr/k8s_platform/internal/types"
)

type K8sSessionHandler struct {
	provider *providers.ModuleProviders
}

// Get godoc
// @Summary Get K8s Cluster Name Context
// @Description get cluster name
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 200 {object} string
// @Router /api/k8s/name [get]
func (c K8sSessionHandler) GetClusterName(ctx echo.Context) error {
	dac, err := c.provider.StorageProvider.GetDynamicAppConfig()
	if err != nil {
		log.Error().Msgf("k8s handler: unable to get dynamic app config: %s", err.Error())
	}
	return ctx.JSON(http.StatusOK, dac.Data.K8sClusterName)
}

// GetPods godoc
// @Summary Get Pods via WebSocket
// @Description get pods' information via WebSocket
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} string "Bad Request"
// @Router /api/k8s/pods [get]
func (c K8sSessionHandler) GetPods(ctx echo.Context) error {
	return c.k8sDataHandler(ctx, "pods")
}

// DeletePod godoc
// @Summary Delets a given pod in the given namespace
// @Description Delets a given pod in the given namespace
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 204 {string} string "No Content"
// @Failure 400 {object} string "Bad Request"
// @Router /api/k8s/pods [delete]
func (c K8sSessionHandler) DeletePod(ctx echo.Context) error {
	podName := ctx.QueryParam("name")
	if podName == "" {
		return ctx.JSON(http.StatusBadRequest, "pod name must be specified")
	}
	namespace := ctx.QueryParam("namespace")
	if namespace == "" {
		return ctx.JSON(http.StatusBadRequest, "namespace must be specified")
	}
	p, err := c.provider.K8sProvider.GetPod(namespace, podName)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("pod does not exist: %v", err))
	}

	// Fetch user's permissions from session data
	sess, err := session.Get("user-permissions", ctx)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, fmt.Sprintf("forbidden. Unable to get auth info: %s", err.Error()))
	}

	userPermissions := []string{}
	if sess.Values["permissions"] != nil {
		userPermissions = sess.Values["permissions"].([]string)
	}

	if !c.hasWritePermissions(userPermissions, p.Labels) {
		return ctx.JSON(http.StatusForbidden, "forbidden. You do not have write permissions for this resource")
	}

	err = c.provider.K8sProvider.DeletePod(namespace, podName)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to delete pod: %v", err))
	}
	return ctx.NoContent(http.StatusNoContent)
}

// GetDeployments godoc
// @Summary Get deployments via WebSocket
// @Description get deployments' information via WebSocket
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} string "Bad Request"
// @Router /api/k8s/deployments [get]
func (c K8sSessionHandler) GetDeployments(ctx echo.Context) error {
	return c.k8sDataHandler(ctx, "deployments")
}

// ScaleDeployment godoc
// @Summary Scale a deployment to a given number of replicas
// @Description Scale a deployment to a given number of replicas
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Successfully scaled deployment"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "unable to scale deployment"
// @Router /api/k8s/deployments/scale [post]
func (c K8sSessionHandler) ScaleDeployment(ctx echo.Context) error {
	sess, err := session.Get("user-permissions", ctx)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, fmt.Sprintf("forbidden. Unable to get auth info: %s", err.Error()))
	}

	userPermissions := []string{}
	if sess.Values["permissions"] != nil {
		userPermissions = sess.Values["permissions"].([]string)
	}

	scaleInfoData, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("unable to read deploy scale info from json body %s", err.Error()))
	}

	scaleInfo := types.ScaleInfo{}

	err = json.Unmarshal(scaleInfoData, &scaleInfo)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("unable to unmarshal deployment scale info from json body %s", err.Error()))
	}

	if !c.hasWritePermissions(userPermissions, scaleInfo.ResourceLabels) {
		return ctx.JSON(http.StatusForbidden, "forbidden. You do not have write permissions for this resource")
	}

	dac, ok := ctx.Get("dynamicAppConfig").(types.DynamicAppConfig)
	if !ok {
		log.Warn().Msg("dynamic config format unknown")
	}

	scaleLimit := dac.Data.DefaultReplicaScaleLimit

	for _, v := range scaleInfo.ResourceLabels {
		if l, ok := dac.Data.ReplicaScaleLimits[v]; ok {
			scaleLimit = l
			break
		}
	}

	if scaleInfo.Replicas < 0 {
		return ctx.JSON(http.StatusBadRequest, "desired replicas must be greater than or equal to 0")
	}

	if scaleInfo.Replicas > int32(scaleLimit) {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("desired replicas must be less than or equal to %d (the limit set for this resource). If you need more replicas, increase the limit for this resource.", scaleLimit))
	}

	err = c.provider.K8sProvider.ScaleDeployment(scaleInfo.Namespace, scaleInfo.Name, scaleInfo.Replicas)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to scale deployment: %v", err))
	}

	return ctx.JSON(http.StatusOK, fmt.Sprintf("successfully scaled deployment %s to %d replicas", scaleInfo.Name, scaleInfo.Replicas))
}

// GetReplicasets godoc
// @Summary Get daemonsets via WebSocket
// @Description get daemonsets' information via WebSocket
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} string "Bad Request"
// @Router /api/k8s/replicasets [get]
func (c K8sSessionHandler) GetReplicasets(ctx echo.Context) error {
	return c.k8sDataHandler(ctx, "replicasets")
}

// GetDaemonsets godoc
// @Summary Get daemonsets via WebSocket
// @Description get daemonsets' information via WebSocket
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} string "Bad Request"
// @Router /api/k8s/daemonsets [get]
func (c K8sSessionHandler) GetDaemonsets(ctx echo.Context) error {
	return c.k8sDataHandler(ctx, "daemonsets")
}

// GetStatefulsets godoc
// @Summary Get statefulsets via WebSocket
// @Description get statefulsets' information via WebSocket
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} string "Bad Request"
// @Router /api/k8s/statefulsets [get]
func (c K8sSessionHandler) GetStatefulsets(ctx echo.Context) error {
	return c.k8sDataHandler(ctx, "statefulsets")
}

// GetJobs godoc
// @Summary Get jobs via WebSocket
// @Description get jobs' information via WebSocket
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} string "Bad Request"
// @Router /api/k8s/jobs [get]
func (c K8sSessionHandler) GetJobs(ctx echo.Context) error {
	return c.k8sDataHandler(ctx, "jobs")
}

// GetCronJobs godoc
// @Summary Get cronjobs via WebSocket
// @Description get cronjobs' information via WebSocket
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} string "Bad Request"
// @Router /api/k8s/cronjobs [get]
func (c K8sSessionHandler) GetCronJobs(ctx echo.Context) error {
	return c.k8sDataHandler(ctx, "cronjobs")
}

// GetServices godoc
// @Summary Get services via WebSocket
// @Description get services information via WebSocket
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} string "Bad Request"
// @Router /api/k8s/services [get]
func (c K8sSessionHandler) GetServices(ctx echo.Context) error {
	return c.k8sDataHandler(ctx, "services")
}

// GetIngresses godoc
// @Summary Get ingresses via WebSocket
// @Description get ingresses information via WebSocket
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} string "Bad Request"
// @Router /api/k8s/ingresses [get]
func (c K8sSessionHandler) GetIngresses(ctx echo.Context) error {
	return c.k8sDataHandler(ctx, "ingresses")
}

// GetConfigMaps godoc
// @Summary Get configmaps via WebSocket
// @Description get configmaps information via WebSocket
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} string "Bad Request"
// @Router /api/k8s/configmaps [get]
func (c K8sSessionHandler) GetConfigMaps(ctx echo.Context) error {
	return c.k8sDataHandler(ctx, "configmaps")
}

// GetNodes godoc
// @Summary Get nodes via WebSocket
// @Description get nodes information via WebSocket
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} string "Bad Request"
// @Router /api/k8s/nodes [get]
func (c K8sSessionHandler) GetNodes(ctx echo.Context) error {
	return c.k8sDataHandler(ctx, "nodes")
}

// GetClusterEvents godoc
// @Summary Get cluster events via WebSocket
// @Description get cluster events information via WebSocket
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 101 {string} string "Switching Protocols"
// @Failure 400 {object} string "Bad Request"
// @Router /api/k8s/clusterevents [get]
func (c K8sSessionHandler) GetClusterEvents(ctx echo.Context) error {
	return c.k8sDataHandler(ctx, "clusterevents")
}

// RolloutRestart godoc
// @Summary Initiate a restart for a deployment, daemonset, or statefulset
// @Description initiate a restart for a deployment, daemonset, or statefulset
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 400 {object} string "Bad Request"
// @Router /api/k8s/rolloutrestart [post]
func (c K8sSessionHandler) RolloutRestart(ctx echo.Context) error {
	sess, err := session.Get("user-permissions", ctx)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, fmt.Sprintf("forbidden. Unable to get auth info: %s", err.Error()))
	}

	userPermissions := []string{}
	if sess.Values["permissions"] != nil {
		userPermissions = sess.Values["permissions"].([]string)
	}

	resourceInfoData, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("unable to read resource info from json body %s", err.Error()))
	}

	resourceInfo := types.ResourceInfo{}

	err = json.Unmarshal(resourceInfoData, &resourceInfo)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("unable to unmarshal resource info from json body %s", err.Error()))
	}

	if !c.hasWritePermissions(userPermissions, resourceInfo.Labels) {
		return ctx.JSON(http.StatusForbidden, "forbidden. You do not have write permissions for this resource")
	}

	if resourceInfo.Kind == "deployment" {
		err = c.provider.K8sProvider.RolloutRestartDeployment(resourceInfo.Namespace, resourceInfo.Name)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to restart deployment: %v", err))
		}
	} else if resourceInfo.Kind == "daemonset" {
		err = c.provider.K8sProvider.RolloutRestartDaemonSet(resourceInfo.Namespace, resourceInfo.Name)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to restart daemonset: %v", err))
		}
	} else if resourceInfo.Kind == "statefulset" {
		err = c.provider.K8sProvider.RolloutRestartStatefulSet(resourceInfo.Namespace, resourceInfo.Name)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to restart statefulset: %v", err))
		}
	} else {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("unknown or invalid resource type for rollout restart: %s", resourceInfo.Kind))
	}

	return ctx.JSON(http.StatusOK, fmt.Sprintf("successfully initiated restart for %s: %s", resourceInfo.Kind, resourceInfo.Name))
}

// RunPodExecPlugin godoc
// @Summary Executes a pod exec plugin on a specific pod
// @Description Executes a pod exec plugin on a specific pod
// @Tags K8s
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Successfully initiated pod exec plugin"
// @Failure 400 {object} string "Bad Request"
// @Failure 403 {object} string "Forbidden"
// @Failure 500 {object} string "unable to run pod exec plugin"
// @Router /api/k8s/execplugin [post]
func (c K8sSessionHandler) RunPodExecPlugin(ctx echo.Context) error {
	resourceInfoData, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("unable to read resource info from json body for pod exec plugin request %s", err.Error()))
	}
	resourceInfo := types.ResourceInfo{}

	err = json.Unmarshal(resourceInfoData, &resourceInfo)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("unable to unmarshal resource info from json body for pod exec plugin request %s", err.Error()))
	}

	sess, err := session.Get("user-permissions", ctx)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, fmt.Sprintf("forbidden. Unable to get auth info: %s", err.Error()))
	}

	userPermissions := []string{}
	if sess.Values["permissions"] != nil {
		userPermissions = sess.Values["permissions"].([]string)
	}

	if !c.hasWritePermissions(userPermissions, resourceInfo.Labels) {
		return ctx.JSON(http.StatusForbidden, "forbidden. You do not have write permissions for this resource")
	}

	// Query the exec plugin to ensure no malicious commands are being injected from the client
	appConfig, err := c.provider.StorageProvider.Session.SDK.GetDynamicAppConfig()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to load exec plugin configurations: %v", err))
	}

	plugin := types.K8sPodExecPlugin{}
	pluginExists := false
	for _, p := range appConfig.Data.K8sPodExecPlugins {
		if p.Name == resourceInfo.Plugin.Name {
			plugin = p
			pluginExists = true
			break
		}
	}

	if !pluginExists {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("plugin with name, %s, does not exist", resourceInfo.Plugin.Name))
	}

	if plugin.Command != resourceInfo.Plugin.Command {
		return ctx.JSON(http.StatusBadRequest, fmt.Sprintf("plugin command does not match the expected command for %s", resourceInfo.Plugin.Name))
	}

	stdout, stderr, err := c.provider.K8sProvider.ExecutePodExecPlugin(resourceInfo.Namespace, resourceInfo.Name, plugin)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("unable to run pod exec plugin: %v", err))
	}
	outStr := fmt.Sprintf("stdout: %s\nstderr: %s", stdout, stderr)
	return ctx.JSON(http.StatusOK, outStr)
}

func (c K8sSessionHandler) k8sDataHandler(ctx echo.Context, resource string) error {

	sess, err := session.Get("user-permissions", ctx)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, fmt.Sprintf("forbidden. Unable to get auth info: %s", err.Error()))
	}

	userPermissions := []string{}
	if sess.Values["permissions"] != nil {
		userPermissions = sess.Values["permissions"].([]string)
	}

	upgradeHeader := ctx.Request().Header.Get("Upgrade")
	if upgradeHeader != "websocket" {
		data, err := c.GetK8sData(ctx, userPermissions, resource)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
		return ctx.JSON(http.StatusOK, data)
	} else {
		ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		for {
			// gracefully handle client closure
			_, msg, err := ws.ReadMessage()
			if ce, ok := err.(*websocket.CloseError); ok {
				switch ce.Code {
				case websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived:
					log.Debug().Msgf("client closed connection: %s", ce.Text)
					ws.Close()
				}
			}
			log.Debug().Msgf("message: %s", string(msg))

			data, err := c.GetK8sData(ctx, userPermissions, resource)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}
			// Write
			err = ws.WriteJSON(data)
			if err != nil {
				return err
			}
		}
	}
}

func (c K8sSessionHandler) GetK8sData(ctx echo.Context, userPermissions []string, resource string) ([]types.K8sResourceWrapper, error) {
	dac, ok := ctx.Get("dynamicAppConfig").(types.DynamicAppConfig)
	if !ok {
		log.Warn().Msg("dynamic config format unknown")
	}
	if dac.Data.K8sClusterName == "" {
		dac.Data.K8sClusterName = "default"
	}

	data, err := c.provider.CacheProvider.Get(fmt.Sprintf("%s_%s", dac.Data.K8sClusterName, resource))
	if err != nil {
		return nil, fmt.Errorf("unable to get %s: %v", resource, err)
	}

	rd, ok := data.([]any)
	if !ok {
		return nil, fmt.Errorf("unexpected data type for %s: %T", resource, data)
	}

	resp := []types.K8sResourceWrapper{}

	for _, p := range userPermissions {
		if p == "*" || resource == "clusterevents" {
			for _, r := range rd {
				resp = append(resp, types.K8sResourceWrapper{
					Data:  r,
					Write: true,
				})
			}
			return resp, nil
		}
	}

	if resource != "clusterevents" {
		permissionMap := make(map[string]bool)
		for _, p := range userPermissions {
			permissionMap[p] = true
		}

		for _, r := range rd {
			d, ok := r.(map[string]interface{})
			if !ok {
				continue
			}

			md, ok := d["metadata"].(map[string]interface{})
			if !ok {
				continue
			}

			labelMap := map[string]interface{}{}
			if mdl, ok := md["labels"]; ok {
				labelMap, ok = mdl.(map[string]interface{})
				if !ok {
					continue
				}
			}

			writeAdded := false
			readAdded := false
			for _, v := range labelMap {
				writePermission := fmt.Sprintf("%s_write", v)
				readPermission := fmt.Sprintf("%s_read", v)

				if !writeAdded && permissionMap[writePermission] {
					resp = append(resp, types.K8sResourceWrapper{
						Data:  r,
						Write: true,
					})
					writeAdded = true
				} else if !readAdded && (permissionMap["global_read_only"] || permissionMap[readPermission]) {
					resp = append(resp, types.K8sResourceWrapper{
						Data:  r,
						Write: false,
					})
					readAdded = true
				}

				if writeAdded && readAdded {
					break
				}
			}
		}
	}

	return resp, nil
}

func (c K8sSessionHandler) hasWritePermissions(userPermissions []string, labels map[string]string) bool {
	for _, p := range userPermissions {
		if p == "*" {
			return true
		}
		for _, v := range labels {
			if p == fmt.Sprintf("%s_write", v) {
				return true
			}
		}
	}
	return false
}
