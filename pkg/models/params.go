package models

type ListBuildsParams struct {
	Namespace     string `json:"namespace"`
	Prefix        string `json:"prefix,omitempty"`
	LabelSelector string `json:"label-selector,omitempty"`
}

type GetBuildParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

type CreateBuildParams struct {
	Name         string            `json:"name"`
	Namespace    string            `json:"namespace,omitempty"`
	SourceType   string            `json:"source-type"`
	SourceURL    string            `json:"source-url"`
	ContextDir   string            `json:"context-dir,omitempty"`
	Revision     string            `json:"revision,omitempty"`
	Strategy     string            `json:"strategy"`
	StrategyKind string            `json:"strategy-kind,omitempty"`
	OutputImage  string            `json:"output-image"`
	Parameters   map[string]string `json:"parameters,omitempty"`
	Timeout      string            `json:"timeout,omitempty"`
}

type ListBuildRunsParams struct {
	Namespace     string `json:"namespace"`
	Prefix        string `json:"prefix,omitempty"`
	LabelSelector string `json:"label-selector,omitempty"`
}

type GetBuildRunParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

type CreateBuildRunParams struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	BuildName string `json:"build-name,omitempty"`

	SourceType     string            `json:"source-type,omitempty"`
	SourceURL      string            `json:"source-url,omitempty"`
	ContextDir     string            `json:"context-dir,omitempty"`
	Revision       string            `json:"revision,omitempty"`
	Strategy       string            `json:"strategy,omitempty"`
	StrategyKind   string            `json:"strategy-kind,omitempty"`
	OutputImage    string            `json:"output-image,omitempty"`
	Parameters     map[string]string `json:"parameters,omitempty"`
	Timeout        string            `json:"timeout,omitempty"`
	ServiceAccount string            `json:"service-account,omitempty"`
}

type RestartBuildRunParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

type DeleteBuildParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

type DeleteBuildRunParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

type ListBuildStrategiesParams struct {
	Namespace     string `json:"namespace"`
	Prefix        string `json:"prefix,omitempty"`
	LabelSelector string `json:"label-selector,omitempty"`
}

type ListClusterBuildStrategiesParams struct {
	Prefix        string `json:"prefix,omitempty"`
	LabelSelector string `json:"label-selector,omitempty"`
}
