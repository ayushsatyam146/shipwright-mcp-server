package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	buildv1beta1 "github.com/shipwright-io/build/pkg/apis/build/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ListBuildsParams parameters for listing builds
type ListBuildsParams struct {
	Namespace     string `json:"namespace"`
	Prefix        string `json:"prefix,omitempty"`
	LabelSelector string `json:"label-selector,omitempty"`
}

// GetBuildParams parameters for getting a build
type GetBuildParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

// CreateBuildParams parameters for creating a build
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

// ListBuildRunsParams parameters for listing buildruns
type ListBuildRunsParams struct {
	Namespace     string `json:"namespace"`
	Prefix        string `json:"prefix,omitempty"`
	LabelSelector string `json:"label-selector,omitempty"`
}

// GetBuildRunParams parameters for getting a buildrun
type GetBuildRunParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

// CreateBuildRunParams parameters for creating a buildrun
type CreateBuildRunParams struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	BuildName string `json:"build-name,omitempty"`
	// Inline build spec fields
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

// RestartBuildRunParams parameters for restarting a buildrun
type RestartBuildRunParams struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

// ListBuildStrategiesParams parameters for listing build strategies
type ListBuildStrategiesParams struct {
	Namespace     string `json:"namespace"`
	Prefix        string `json:"prefix,omitempty"`
	LabelSelector string `json:"label-selector,omitempty"`
}

// ListClusterBuildStrategiesParams parameters for listing cluster build strategies
type ListClusterBuildStrategiesParams struct {
	Prefix        string `json:"prefix,omitempty"`
	LabelSelector string `json:"label-selector,omitempty"`
}

var k8sClient client.Client

func main() {
	// Set log output to stderr so it doesn't interfere with MCP protocol on stdout
	log.SetOutput(os.Stderr)
	log.Printf("🚀 Starting Shipwright Build MCP Server v1.1.0")

	// Initialize Kubernetes client
	log.Printf("🔧 Initializing Kubernetes client...")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("⚠️  Not running in cluster, trying kubeconfig...")
		// Fallback to kubeconfig
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			log.Fatalf("❌ Failed to create Kubernetes config: %v", err)
		}
	}

	scheme := runtime.NewScheme()
	if err := buildv1beta1.AddToScheme(scheme); err != nil {
		log.Fatalf("❌ Failed to add build scheme: %v", err)
	}

	k8sClient, err = client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		log.Fatalf("❌ Failed to create Kubernetes client: %v", err)
	}

	log.Printf("✅ Kubernetes client initialized successfully")

	// Create a server with Shipwright Build tools
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "shipwright-build-mcp-server",
		Version: "v1.1.0",
	}, nil)

	log.Printf("🔨 Registering MCP tools...")

	// Add tools for managing Build resources
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_builds",
		Description: "List Builds in a namespace with filtering options",
	}, listBuilds)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_build",
		Description: "Get a specific Build by name",
	}, getBuild)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_build",
		Description: "Create a new Build resource",
	}, createBuild)

	// Add tools for managing BuildRun resources
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_buildruns",
		Description: "List BuildRuns in a namespace with filtering options",
	}, listBuildRuns)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_buildrun",
		Description: "Get a specific BuildRun by name",
	}, getBuildRun)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_buildrun",
		Description: "Create a new BuildRun resource (either from existing Build or inline)",
	}, createBuildRun)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "restart_buildrun",
		Description: "Restart a BuildRun by creating a new one",
	}, restartBuildRun)

	// Add tools for managing BuildStrategy resources
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_buildstrategies",
		Description: "List BuildStrategies in a namespace with filtering options",
	}, listBuildStrategies)

	// Add tools for managing ClusterBuildStrategy resources
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_clusterbuildstrategies",
		Description: "List ClusterBuildStrategies with filtering options",
	}, listClusterBuildStrategies)

	// Run the server over stdin/stdout
	log.Printf("🎯 MCP Server ready and listening on stdin/stdout...")
	log.Printf("🛠️  Available tools: list_builds, get_build, create_build, list_buildruns, get_buildrun, create_buildrun, restart_buildrun, list_buildstrategies, list_clusterbuildstrategies")
	log.Printf("📡 Waiting for MCP client connections...")

	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Fatal(err)
	}
}

// Tool implementations

func listBuilds(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[ListBuildsParams]) (*mcp.CallToolResultFor[any], error) {
	buildList := &buildv1beta1.BuildList{}

	listOpts := []client.ListOption{
		client.InNamespace(params.Arguments.Namespace),
	}

	if params.Arguments.LabelSelector != "" {
		selector, err := metav1.ParseToLabelSelector(params.Arguments.LabelSelector)
		if err != nil {
			return &mcp.CallToolResultFor[any]{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Invalid label selector: %v", err)}},
			}, nil
		}
		selectorObj, err := metav1.LabelSelectorAsSelector(selector)
		if err != nil {
			return &mcp.CallToolResultFor[any]{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Invalid label selector: %v", err)}},
			}, nil
		}
		listOpts = append(listOpts, client.MatchingLabelsSelector{Selector: selectorObj})
	}

	if err := k8sClient.List(ctx, buildList, listOpts...); err != nil {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to list builds: %v", err)}},
		}, nil
	}

	var filteredBuilds []buildv1beta1.Build
	for _, build := range buildList.Items {
		if params.Arguments.Prefix == "" || strings.HasPrefix(build.Name, params.Arguments.Prefix) {
			filteredBuilds = append(filteredBuilds, build)
		}
	}

	if len(filteredBuilds) == 0 {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "No builds found"}},
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d build(s):\n\n", len(filteredBuilds)))

	for _, build := range filteredBuilds {
		result.WriteString(fmt.Sprintf("Name: %s\n", build.Name))
		result.WriteString(fmt.Sprintf("Namespace: %s\n", build.Namespace))
		result.WriteString(fmt.Sprintf("Strategy: %s (%s)\n", build.Spec.Strategy.Name, build.Spec.Strategy.Kind))
		if build.Spec.Source != nil {
			result.WriteString(fmt.Sprintf("Source Type: %s\n", build.Spec.Source.Type))
			if build.Spec.Source.Git != nil {
				result.WriteString(fmt.Sprintf("Git URL: %s\n", build.Spec.Source.Git.URL))
			}
		}
		result.WriteString(fmt.Sprintf("Output Image: %s\n", build.Spec.Output.Image))
		result.WriteString(fmt.Sprintf("Created: %s\n", build.CreationTimestamp.Format("2006-01-02 15:04:05")))
		result.WriteString("---\n")
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: result.String()}},
	}, nil
}

func getBuild(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[GetBuildParams]) (*mcp.CallToolResultFor[any], error) {
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = "default"
	}

	build := &buildv1beta1.Build{}
	if err := k8sClient.Get(ctx, client.ObjectKey{
		Name:      params.Arguments.Name,
		Namespace: namespace,
	}, build); err != nil {
		if errors.IsNotFound(err) {
			return &mcp.CallToolResultFor[any]{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Build '%s' not found in namespace '%s'", params.Arguments.Name, namespace)}},
			}, nil
		}
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to get build: %v", err)}},
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Build: %s\n", build.Name))
	result.WriteString(fmt.Sprintf("Namespace: %s\n", build.Namespace))
	result.WriteString(fmt.Sprintf("Strategy: %s (%s)\n", build.Spec.Strategy.Name, build.Spec.Strategy.Kind))
	if build.Spec.Source != nil {
		result.WriteString(fmt.Sprintf("Source Type: %s\n", build.Spec.Source.Type))
		if build.Spec.Source.Git != nil {
			result.WriteString(fmt.Sprintf("Git URL: %s\n", build.Spec.Source.Git.URL))
			if build.Spec.Source.Git.Revision != nil {
				result.WriteString(fmt.Sprintf("Git Revision: %s\n", *build.Spec.Source.Git.Revision))
			}
		}
		if build.Spec.Source.ContextDir != nil {
			result.WriteString(fmt.Sprintf("Context Dir: %s\n", *build.Spec.Source.ContextDir))
		}
	}
	result.WriteString(fmt.Sprintf("Output Image: %s\n", build.Spec.Output.Image))
	if build.Spec.Timeout != nil {
		result.WriteString(fmt.Sprintf("Timeout: %s\n", build.Spec.Timeout.Duration))
	}
	if len(build.Spec.ParamValues) > 0 {
		result.WriteString("Parameters:\n")
		for _, param := range build.Spec.ParamValues {
			result.WriteString(fmt.Sprintf("  %s: %v\n", param.Name, param.Value))
		}
	}
	result.WriteString(fmt.Sprintf("Created: %s\n", build.CreationTimestamp.Format("2006-01-02 15:04:05")))

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: result.String()}},
	}, nil
}

func createBuild(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[CreateBuildParams]) (*mcp.CallToolResultFor[any], error) {
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = "default"
	}

	// Validate required parameters
	if params.Arguments.Name == "" {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Build name is required"}},
		}, nil
	}
	if params.Arguments.SourceURL == "" {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Source URL is required"}},
		}, nil
	}
	if params.Arguments.Strategy == "" {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Strategy name is required"}},
		}, nil
	}
	if params.Arguments.OutputImage == "" {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Output image is required"}},
		}, nil
	}

	// Create Build object
	build := &buildv1beta1.Build{
		ObjectMeta: metav1.ObjectMeta{
			Name:      params.Arguments.Name,
			Namespace: namespace,
		},
		Spec: buildv1beta1.BuildSpec{
			Strategy: buildv1beta1.Strategy{
				Name: params.Arguments.Strategy,
			},
			Output: buildv1beta1.Image{
				Image: params.Arguments.OutputImage,
			},
		},
	}

	// Set strategy kind
	strategyKind := params.Arguments.StrategyKind
	if strategyKind == "" {
		strategyKind = "ClusterBuildStrategy"
	}
	kind := buildv1beta1.BuildStrategyKind(strategyKind)
	build.Spec.Strategy.Kind = &kind

	// Set source
	sourceType := buildv1beta1.BuildSourceType(params.Arguments.SourceType)
	build.Spec.Source = &buildv1beta1.Source{
		Type: sourceType,
	}

	if params.Arguments.ContextDir != "" {
		build.Spec.Source.ContextDir = &params.Arguments.ContextDir
	}

	switch sourceType {
	case buildv1beta1.GitType:
		build.Spec.Source.Git = &buildv1beta1.Git{
			URL: params.Arguments.SourceURL,
		}
		if params.Arguments.Revision != "" {
			build.Spec.Source.Git.Revision = &params.Arguments.Revision
		}
	case buildv1beta1.OCIArtifactType:
		build.Spec.Source.OCIArtifact = &buildv1beta1.OCIArtifact{
			Image: params.Arguments.SourceURL,
		}
	default:
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Source type must be 'Git' or 'OCI'"}},
		}, nil
	}

	// Set parameters
	if len(params.Arguments.Parameters) > 0 {
		for name, value := range params.Arguments.Parameters {
			param := buildv1beta1.ParamValue{
				Name: name,
				SingleValue: &buildv1beta1.SingleValue{
					Value: &value,
				},
			}
			build.Spec.ParamValues = append(build.Spec.ParamValues, param)
		}
	}

	// Set timeout
	if params.Arguments.Timeout != "" {
		duration, err := time.ParseDuration(params.Arguments.Timeout)
		if err != nil {
			return &mcp.CallToolResultFor[any]{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Invalid timeout duration: %v", err)}},
			}, nil
		}
		build.Spec.Timeout = &metav1.Duration{Duration: duration}
	}

	if err := k8sClient.Create(ctx, build); err != nil {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to create build: %v", err)}},
		}, nil
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Successfully created Build '%s' in namespace '%s'", params.Arguments.Name, namespace)}},
	}, nil
}

func listBuildRuns(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[ListBuildRunsParams]) (*mcp.CallToolResultFor[any], error) {
	buildRunList := &buildv1beta1.BuildRunList{}

	listOpts := []client.ListOption{
		client.InNamespace(params.Arguments.Namespace),
	}

	if params.Arguments.LabelSelector != "" {
		selector, err := metav1.ParseToLabelSelector(params.Arguments.LabelSelector)
		if err != nil {
			return &mcp.CallToolResultFor[any]{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Invalid label selector: %v", err)}},
			}, nil
		}
		selectorObj, err := metav1.LabelSelectorAsSelector(selector)
		if err != nil {
			return &mcp.CallToolResultFor[any]{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Invalid label selector: %v", err)}},
			}, nil
		}
		listOpts = append(listOpts, client.MatchingLabelsSelector{Selector: selectorObj})
	}

	if err := k8sClient.List(ctx, buildRunList, listOpts...); err != nil {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to list buildruns: %v", err)}},
		}, nil
	}

	var filteredBuildRuns []buildv1beta1.BuildRun
	for _, buildRun := range buildRunList.Items {
		if params.Arguments.Prefix == "" || strings.HasPrefix(buildRun.Name, params.Arguments.Prefix) {
			filteredBuildRuns = append(filteredBuildRuns, buildRun)
		}
	}

	if len(filteredBuildRuns) == 0 {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "No buildruns found"}},
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d buildrun(s):\n\n", len(filteredBuildRuns)))

	for _, buildRun := range filteredBuildRuns {
		result.WriteString(fmt.Sprintf("Name: %s\n", buildRun.Name))
		result.WriteString(fmt.Sprintf("Namespace: %s\n", buildRun.Namespace))

		// Get build reference
		if buildRun.Spec.Build.Name != nil {
			result.WriteString(fmt.Sprintf("Build: %s\n", *buildRun.Spec.Build.Name))
		}

		// Get status
		if len(buildRun.Status.Conditions) > 0 {
			for _, condition := range buildRun.Status.Conditions {
				if condition.Type == buildv1beta1.Succeeded {
					result.WriteString(fmt.Sprintf("Status: %s\n", condition.Status))
					result.WriteString(fmt.Sprintf("Reason: %s\n", condition.Reason))
					if condition.Message != "" {
						result.WriteString(fmt.Sprintf("Message: %s\n", condition.Message))
					}
					break
				}
			}
		}

		if buildRun.Status.StartTime != nil {
			result.WriteString(fmt.Sprintf("Started: %s\n", buildRun.Status.StartTime.Format("2006-01-02 15:04:05")))
		}
		if buildRun.Status.CompletionTime != nil {
			result.WriteString(fmt.Sprintf("Completed: %s\n", buildRun.Status.CompletionTime.Format("2006-01-02 15:04:05")))
		}

		result.WriteString(fmt.Sprintf("Created: %s\n", buildRun.CreationTimestamp.Format("2006-01-02 15:04:05")))
		result.WriteString("---\n")
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: result.String()}},
	}, nil
}

func getBuildRun(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[GetBuildRunParams]) (*mcp.CallToolResultFor[any], error) {
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = "default"
	}

	buildRun := &buildv1beta1.BuildRun{}
	if err := k8sClient.Get(ctx, client.ObjectKey{
		Name:      params.Arguments.Name,
		Namespace: namespace,
	}, buildRun); err != nil {
		if errors.IsNotFound(err) {
			return &mcp.CallToolResultFor[any]{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("BuildRun '%s' not found in namespace '%s'", params.Arguments.Name, namespace)}},
			}, nil
		}
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to get buildrun: %v", err)}},
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("BuildRun: %s\n", buildRun.Name))
	result.WriteString(fmt.Sprintf("Namespace: %s\n", buildRun.Namespace))

	// Get build reference
	if buildRun.Spec.Build.Name != nil {
		result.WriteString(fmt.Sprintf("Build: %s\n", *buildRun.Spec.Build.Name))
	}

	// Get status details
	if len(buildRun.Status.Conditions) > 0 {
		for _, condition := range buildRun.Status.Conditions {
			if condition.Type == buildv1beta1.Succeeded {
				result.WriteString(fmt.Sprintf("Status: %s\n", condition.Status))
				result.WriteString(fmt.Sprintf("Reason: %s\n", condition.Reason))
				if condition.Message != "" {
					result.WriteString(fmt.Sprintf("Message: %s\n", condition.Message))
				}
				result.WriteString(fmt.Sprintf("Last Transition: %s\n", condition.LastTransitionTime.Format("2006-01-02 15:04:05")))
				break
			}
		}
	}

	if buildRun.Status.TaskRunName != nil {
		result.WriteString(fmt.Sprintf("TaskRun: %s\n", *buildRun.Status.TaskRunName))
	}

	if buildRun.Status.StartTime != nil {
		result.WriteString(fmt.Sprintf("Started: %s\n", buildRun.Status.StartTime.Format("2006-01-02 15:04:05")))
	}
	if buildRun.Status.CompletionTime != nil {
		result.WriteString(fmt.Sprintf("Completed: %s\n", buildRun.Status.CompletionTime.Format("2006-01-02 15:04:05")))
	}

	// Output information
	if buildRun.Status.Output != nil {
		result.WriteString("Output:\n")
		if buildRun.Status.Output.Digest != "" {
			result.WriteString(fmt.Sprintf("  Digest: %s\n", buildRun.Status.Output.Digest))
		}
		if buildRun.Status.Output.Size > 0 {
			result.WriteString(fmt.Sprintf("  Size: %d bytes\n", buildRun.Status.Output.Size))
		}
	}

	// Failure details
	if buildRun.Status.FailureDetails != nil {
		result.WriteString("Failure Details:\n")
		result.WriteString(fmt.Sprintf("  Reason: %s\n", buildRun.Status.FailureDetails.Reason))
		result.WriteString(fmt.Sprintf("  Message: %s\n", buildRun.Status.FailureDetails.Message))
		if buildRun.Status.FailureDetails.Location != nil {
			result.WriteString(fmt.Sprintf("  Pod: %s\n", buildRun.Status.FailureDetails.Location.Pod))
			result.WriteString(fmt.Sprintf("  Container: %s\n", buildRun.Status.FailureDetails.Location.Container))
		}
	}

	result.WriteString(fmt.Sprintf("Created: %s\n", buildRun.CreationTimestamp.Format("2006-01-02 15:04:05")))

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: result.String()}},
	}, nil
}

func createBuildRun(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[CreateBuildRunParams]) (*mcp.CallToolResultFor[any], error) {
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = "default"
	}

	// Validate that we have either a BuildName or inline build spec
	if params.Arguments.BuildName == "" && params.Arguments.Strategy == "" {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Either build-name or inline build spec (strategy, source-url, output-image) must be provided"}},
		}, nil
	}

	// Create BuildRun object
	buildRun := &buildv1beta1.BuildRun{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
		},
		Spec: buildv1beta1.BuildRunSpec{},
	}

	// Set the name (generate if not provided)
	if params.Arguments.Name != "" {
		buildRun.Name = params.Arguments.Name
	} else {
		buildRun.GenerateName = "buildrun-"
	}

	// Set service account if provided
	if params.Arguments.ServiceAccount != "" {
		buildRun.Spec.ServiceAccount = &params.Arguments.ServiceAccount
	}

	// Set timeout if provided
	if params.Arguments.Timeout != "" {
		duration, err := time.ParseDuration(params.Arguments.Timeout)
		if err != nil {
			return &mcp.CallToolResultFor[any]{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Invalid timeout duration: %v", err)}},
			}, nil
		}
		buildRun.Spec.Timeout = &metav1.Duration{Duration: duration}
	}

	// Set parameters if provided
	if len(params.Arguments.Parameters) > 0 {
		for name, value := range params.Arguments.Parameters {
			param := buildv1beta1.ParamValue{
				Name: name,
				SingleValue: &buildv1beta1.SingleValue{
					Value: &value,
				},
			}
			buildRun.Spec.ParamValues = append(buildRun.Spec.ParamValues, param)
		}
	}

	if params.Arguments.BuildName != "" {
		// Reference existing Build
		buildRun.Spec.Build = buildv1beta1.ReferencedBuild{
			Name: &params.Arguments.BuildName,
		}
	} else {
		// Create inline build spec
		if params.Arguments.SourceURL == "" || params.Arguments.OutputImage == "" {
			return &mcp.CallToolResultFor[any]{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: "source-url and output-image are required for inline build spec"}},
			}, nil
		}

		buildSpec := &buildv1beta1.BuildSpec{
			Strategy: buildv1beta1.Strategy{
				Name: params.Arguments.Strategy,
			},
			Output: buildv1beta1.Image{
				Image: params.Arguments.OutputImage,
			},
		}

		// Set strategy kind
		strategyKind := params.Arguments.StrategyKind
		if strategyKind == "" {
			strategyKind = "ClusterBuildStrategy"
		}
		kind := buildv1beta1.BuildStrategyKind(strategyKind)
		buildSpec.Strategy.Kind = &kind

		// Set source
		sourceType := buildv1beta1.BuildSourceType(params.Arguments.SourceType)
		if sourceType == "" {
			sourceType = buildv1beta1.GitType
		}

		buildSpec.Source = &buildv1beta1.Source{
			Type: sourceType,
		}

		if params.Arguments.ContextDir != "" {
			buildSpec.Source.ContextDir = &params.Arguments.ContextDir
		}

		switch sourceType {
		case buildv1beta1.GitType:
			buildSpec.Source.Git = &buildv1beta1.Git{
				URL: params.Arguments.SourceURL,
			}
			if params.Arguments.Revision != "" {
				buildSpec.Source.Git.Revision = &params.Arguments.Revision
			}
		case buildv1beta1.OCIArtifactType:
			buildSpec.Source.OCIArtifact = &buildv1beta1.OCIArtifact{
				Image: params.Arguments.SourceURL,
			}
		default:
			return &mcp.CallToolResultFor[any]{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: "Source type must be 'Git' or 'OCI'"}},
			}, nil
		}

		buildRun.Spec.Build = buildv1beta1.ReferencedBuild{
			Spec: buildSpec,
		}
	}

	if err := k8sClient.Create(ctx, buildRun); err != nil {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to create buildrun: %v", err)}},
		}, nil
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Successfully created BuildRun '%s' in namespace '%s'", buildRun.Name, namespace)}},
	}, nil
}

func restartBuildRun(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[RestartBuildRunParams]) (*mcp.CallToolResultFor[any], error) {
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = "default"
	}

	// Get the original BuildRun
	originalBuildRun := &buildv1beta1.BuildRun{}
	if err := k8sClient.Get(ctx, client.ObjectKey{
		Name:      params.Arguments.Name,
		Namespace: namespace,
	}, originalBuildRun); err != nil {
		if errors.IsNotFound(err) {
			return &mcp.CallToolResultFor[any]{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("BuildRun '%s' not found in namespace '%s'", params.Arguments.Name, namespace)}},
			}, nil
		}
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to get buildrun: %v", err)}},
		}, nil
	}

	// Create a new BuildRun based on the original
	newBuildRun := &buildv1beta1.BuildRun{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: originalBuildRun.Name + "-restart-",
			Namespace:    namespace,
			Labels:       originalBuildRun.Labels,
			Annotations:  originalBuildRun.Annotations,
		},
		Spec: originalBuildRun.Spec,
	}

	// Remove status-related annotations that shouldn't be copied
	if newBuildRun.Annotations != nil {
		delete(newBuildRun.Annotations, "kubectl.kubernetes.io/last-applied-configuration")
	}

	if err := k8sClient.Create(ctx, newBuildRun); err != nil {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to create new buildrun: %v", err)}},
		}, nil
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Successfully restarted BuildRun '%s' as '%s' in namespace '%s'", params.Arguments.Name, newBuildRun.Name, namespace)}},
	}, nil
}

func listBuildStrategies(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[ListBuildStrategiesParams]) (*mcp.CallToolResultFor[any], error) {
	buildStrategyList := &buildv1beta1.BuildStrategyList{}

	listOpts := []client.ListOption{
		client.InNamespace(params.Arguments.Namespace),
	}

	if params.Arguments.LabelSelector != "" {
		selector, err := metav1.ParseToLabelSelector(params.Arguments.LabelSelector)
		if err != nil {
			return &mcp.CallToolResultFor[any]{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Invalid label selector: %v", err)}},
			}, nil
		}
		selectorObj, err := metav1.LabelSelectorAsSelector(selector)
		if err != nil {
			return &mcp.CallToolResultFor[any]{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Invalid label selector: %v", err)}},
			}, nil
		}
		listOpts = append(listOpts, client.MatchingLabelsSelector{Selector: selectorObj})
	}

	if err := k8sClient.List(ctx, buildStrategyList, listOpts...); err != nil {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to list buildstrategies: %v", err)}},
		}, nil
	}

	var filteredStrategies []buildv1beta1.BuildStrategy
	for _, strategy := range buildStrategyList.Items {
		if params.Arguments.Prefix == "" || strings.HasPrefix(strategy.Name, params.Arguments.Prefix) {
			filteredStrategies = append(filteredStrategies, strategy)
		}
	}

	if len(filteredStrategies) == 0 {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "No buildstrategies found"}},
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d buildstrategy(ies):\n\n", len(filteredStrategies)))

	for _, strategy := range filteredStrategies {
		result.WriteString(fmt.Sprintf("Name: %s\n", strategy.Name))
		result.WriteString(fmt.Sprintf("Namespace: %s\n", strategy.Namespace))
		result.WriteString(fmt.Sprintf("Steps: %d\n", len(strategy.Spec.Steps)))
		if len(strategy.Spec.Parameters) > 0 {
			result.WriteString(fmt.Sprintf("Parameters: %d\n", len(strategy.Spec.Parameters)))
		}
		result.WriteString(fmt.Sprintf("Created: %s\n", strategy.CreationTimestamp.Format("2006-01-02 15:04:05")))
		result.WriteString("---\n")
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: result.String()}},
	}, nil
}

func listClusterBuildStrategies(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[ListClusterBuildStrategiesParams]) (*mcp.CallToolResultFor[any], error) {
	clusterBuildStrategyList := &buildv1beta1.ClusterBuildStrategyList{}

	var listOpts []client.ListOption

	if params.Arguments.LabelSelector != "" {
		selector, err := metav1.ParseToLabelSelector(params.Arguments.LabelSelector)
		if err != nil {
			return &mcp.CallToolResultFor[any]{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Invalid label selector: %v", err)}},
			}, nil
		}
		selectorObj, err := metav1.LabelSelectorAsSelector(selector)
		if err != nil {
			return &mcp.CallToolResultFor[any]{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Invalid label selector: %v", err)}},
			}, nil
		}
		listOpts = append(listOpts, client.MatchingLabelsSelector{Selector: selectorObj})
	}

	if err := k8sClient.List(ctx, clusterBuildStrategyList, listOpts...); err != nil {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to list clusterbuildstrategies: %v", err)}},
		}, nil
	}

	var filteredStrategies []buildv1beta1.ClusterBuildStrategy
	for _, strategy := range clusterBuildStrategyList.Items {
		if params.Arguments.Prefix == "" || strings.HasPrefix(strategy.Name, params.Arguments.Prefix) {
			filteredStrategies = append(filteredStrategies, strategy)
		}
	}

	if len(filteredStrategies) == 0 {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "No clusterbuildstrategies found"}},
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d clusterbuildstrategy(ies):\n\n", len(filteredStrategies)))

	for _, strategy := range filteredStrategies {
		result.WriteString(fmt.Sprintf("Name: %s\n", strategy.Name))
		result.WriteString(fmt.Sprintf("Steps: %d\n", len(strategy.Spec.Steps)))
		if len(strategy.Spec.Parameters) > 0 {
			result.WriteString(fmt.Sprintf("Parameters: %d\n", len(strategy.Spec.Parameters)))
		}
		result.WriteString(fmt.Sprintf("Created: %s\n", strategy.CreationTimestamp.Format("2006-01-02 15:04:05")))
		result.WriteString("---\n")
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: result.String()}},
	}, nil
}
