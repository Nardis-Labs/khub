package types

type ScaleInfo struct {
	Namespace      string            `json:"namespace"`
	Name           string            `json:"name"`
	Replicas       int32             `json:"replicas"`
	ResourceLabels map[string]string `json:"resourceLabels"`
}
