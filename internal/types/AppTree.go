package types

type NodeLabel struct {
	Label                string `json:"label"`
	ResourceType         string `json:"resourceType"`
	ResourceGroup        string `json:"resourceGroup"`
	ResourceData         any    `json:"resourceData"`
	ResourceOrigin       string `json:"resourceOrigin"`
	ControlledByResource string `json:"controlledByResource"`
	ControlledByKind     string `json:"controlledByKind"`
	ServiceSelector      string `json:"serviceSelector"`
}

type AppTreeNode struct {
	ID       string              `json:"id"`
	Data     NodeLabel           `json:"data"`
	Position AppTreeNodePosition `json:"position"`
}

type AppTreeNodePosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type AppTreeEdge struct {
	ID       string `json:"id"`
	Source   string `json:"source"`
	Target   string `json:"target"`
	EdgeType string `json:"edgeType"`
	Animated bool   `json:"animated"`
}

type AppTree struct {
	Nodes []AppTreeNode `json:"nodes"`
	Edges []AppTreeEdge `json:"edges"`
}
