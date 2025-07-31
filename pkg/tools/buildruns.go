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

func ListBuildRuns(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[models.ListBuildRunsParams]) (*mcp.CallToolResultFor[any], error) {
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

	var buildRuns []buildv1beta1.BuildRun
	for _, buildRun := range buildRunList.Items {
		if params.Arguments.Prefix == "" || strings.HasPrefix(buildRun.Name, params.Arguments.Prefix) {
			buildRuns = append(buildRuns, buildRun)
		}
	}

	if len(buildRuns) == 0 {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "No buildruns found"}},
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d buildrun(s):\n\n", len(buildRuns)))

	for _, buildRun := range buildRuns {
		result.WriteString(fmt.Sprintf("Name: %s\n", buildRun.Name))
		result.WriteString(fmt.Sprintf("Namespace: %s\n", buildRun.Namespace))

		if buildRun.Spec.Build.Name != nil {
			result.WriteString(fmt.Sprintf("Build: %s\n", *buildRun.Spec.Build.Name))
		}

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

func GetBuildRun(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[models.GetBuildRunParams]) (*mcp.CallToolResultFor[any], error) {
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

	if buildRun.Spec.Build.Name != nil {
		result.WriteString(fmt.Sprintf("Build: %s\n", *buildRun.Spec.Build.Name))
	}

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

	if buildRun.Status.Output != nil {
		result.WriteString("Output:\n")
		if buildRun.Status.Output.Digest != "" {
			result.WriteString(fmt.Sprintf("  Digest: %s\n", buildRun.Status.Output.Digest))
		}
		if buildRun.Status.Output.Size > 0 {
			result.WriteString(fmt.Sprintf("  Size: %d bytes\n", buildRun.Status.Output.Size))
		}
	}

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

func CreateBuildRun(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[models.CreateBuildRunParams]) (*mcp.CallToolResultFor[any], error) {
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = "default"
	}

	if params.Arguments.BuildName == "" && params.Arguments.Strategy == "" {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Either build-name or inline build spec (strategy, source-url, output-image) must be provided"}},
		}, nil
	}

	buildRun := &buildv1beta1.BuildRun{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
		},
		Spec: buildv1beta1.BuildRunSpec{},
	}

	if params.Arguments.Name != "" {
		buildRun.Name = params.Arguments.Name
	} else {
		buildRun.GenerateName = "buildrun-"
	}

	if params.Arguments.ServiceAccount != "" {
		buildRun.Spec.ServiceAccount = &params.Arguments.ServiceAccount
	}

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
		buildRun.Spec.Build = buildv1beta1.ReferencedBuild{
			Name: &params.Arguments.BuildName,
		}
	} else {
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

		strategyKind := params.Arguments.StrategyKind
		if strategyKind == "" {
			strategyKind = "ClusterBuildStrategy"
		}
		kind := buildv1beta1.BuildStrategyKind(strategyKind)
		buildSpec.Strategy.Kind = &kind

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

func RestartBuildRun(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[models.RestartBuildRunParams]) (*mcp.CallToolResultFor[any], error) {
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = "default"
	}

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

	newBuildRun := &buildv1beta1.BuildRun{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: originalBuildRun.Name + "-restart-",
			Namespace:    namespace,
			Labels:       originalBuildRun.Labels,
			Annotations:  originalBuildRun.Annotations,
		},
		Spec: originalBuildRun.Spec,
	}

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

func DeleteBuildRun(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[models.DeleteBuildRunParams]) (*mcp.CallToolResultFor[any], error) {
	namespace := params.Arguments.Namespace
	if namespace == "" {
		namespace = "default"
	}

	if params.Arguments.Name == "" {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "BuildRun name is required"}},
		}, nil
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

	if err := k8sClient.Delete(ctx, buildRun); err != nil {
		return &mcp.CallToolResultFor[any]{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Failed to delete buildrun: %v", err)}},
		}, nil
	}

	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Successfully deleted BuildRun '%s' from namespace '%s'", params.Arguments.Name, namespace)}},
	}, nil
}
