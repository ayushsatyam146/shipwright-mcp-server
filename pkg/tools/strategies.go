package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	buildv1beta1 "github.com/shipwright-io/build/pkg/apis/build/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/shipwright-io/build/server/pkg/models"
)

func ListBuildStrategies(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[models.ListBuildStrategiesParams]) (*mcp.CallToolResultFor[any], error) {
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

	var strategies []buildv1beta1.BuildStrategy
	for _, strategy := range buildStrategyList.Items {
		if params.Arguments.Prefix == "" || strings.HasPrefix(strategy.Name, params.Arguments.Prefix) {
			strategies = append(strategies, strategy)
		}
	}

	if len(strategies) == 0 {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "No buildstrategies found"}},
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d buildstrategy(ies):\n\n", len(strategies)))

	for _, strategy := range strategies {
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

func ListClusterBuildStrategies(ctx context.Context, cc *mcp.ServerSession, params *mcp.CallToolParamsFor[models.ListClusterBuildStrategiesParams]) (*mcp.CallToolResultFor[any], error) {
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

	var strategies []buildv1beta1.ClusterBuildStrategy
	for _, strategy := range clusterBuildStrategyList.Items {
		if params.Arguments.Prefix == "" || strings.HasPrefix(strategy.Name, params.Arguments.Prefix) {
			strategies = append(strategies, strategy)
		}
	}

	if len(strategies) == 0 {
		return &mcp.CallToolResultFor[any]{
			Content: []mcp.Content{&mcp.TextContent{Text: "No clusterbuildstrategies found"}},
		}, nil
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Found %d clusterbuildstrategy(ies):\n\n", len(strategies)))

	for _, strategy := range strategies {
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
