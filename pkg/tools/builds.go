package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	buildv1beta1 "github.com/shipwright-io/build/pkg/apis/build/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/shipwright-io/build/server/pkg/models"
)

var k8sClient client.Client

func SetClient(c client.Client) {
	k8sClient = c
}

func ListBuilds(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[models.ListBuildsParams]) (*mcp.CallToolResultFor[any], error) {
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

	var builds []buildv1beta1.Build
	for _, build := range buildList.Items {
		if params.Arguments.Prefix == "" || strings.HasPrefix(build.Name, params.Arguments.Prefix) {
			builds = append(builds, build)
		}
	}

	if len(builds) == 0 {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "No builds found"}},
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d build(s):\n\n", len(builds)))

	for _, build := range builds {
		result.WriteString(fmt.Sprintf("Name: %s\n", build.Name))
		result.WriteString(fmt.Sprintf("Namespace: %s\n", build.Namespace))
		result.WriteString(fmt.Sprintf("Strategy: %s (%s)\n", build.Spec.Strategy.Name, *build.Spec.Strategy.Kind))
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

func GetBuild(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[models.GetBuildParams]) (*mcp.CallToolResultFor[any], error) {
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
	result.WriteString(fmt.Sprintf("Strategy: %s (%s)\n", build.Spec.Strategy.Name, *build.Spec.Strategy.Kind))
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

func CreateBuild(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[models.CreateBuildParams]) (*mcp.CallToolResultFor[any], error) {
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = "default"
	}

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

	strategyKind := params.Arguments.StrategyKind
	if strategyKind == "" {
		strategyKind = "ClusterBuildStrategy"
	}
	kind := buildv1beta1.BuildStrategyKind(strategyKind)
	build.Spec.Strategy.Kind = &kind

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

func DeleteBuild(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[models.DeleteBuildParams]) (*mcp.CallToolResultFor[any], error) {
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = "default"
	}

	if params.Arguments.Name == "" {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Build name is required"}},
		}, nil
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

	if err := k8sClient.Delete(ctx, build); err != nil {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to delete build: %v", err)}},
		}, nil
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Successfully deleted Build '%s' from namespace '%s'", params.Arguments.Name, namespace)}},
	}, nil
}
