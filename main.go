package main

import (
	"context"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	buildv1beta1 "github.com/shipwright-io/build/pkg/apis/build/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/shipwright-io/build/server/pkg/tools"
)

var k8sClient client.Client

func main() {
	log.SetOutput(os.Stderr)
	log.Printf("Starting Shipwright Build MCP Server v1.2.0")

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Printf("Not running in cluster, trying kubeconfig...")
		config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			log.Fatalf("Failed to create Kubernetes config: %v", err)
		}
	}

	scheme := runtime.NewScheme()
	if err := buildv1beta1.AddToScheme(scheme); err != nil {
		log.Fatalf("Failed to add build scheme: %v", err)
	}

	k8sClient, err = client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		log.Fatalf("Failed to create Kubernetes client: %v", err)
	}

	tools.SetClient(k8sClient)

	log.Printf("Kubernetes client initialized")

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "shipwright-build-mcp-server",
		Version: "v1.2.0",
	}, nil)

	log.Printf("Registering MCP tools...")

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_builds",
		Description: "List Builds in a namespace with filtering options",
	}, tools.ListBuilds)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_build",
		Description: "Get a specific Build by name",
	}, tools.GetBuild)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_build",
		Description: "Create a new Build resource",
	}, tools.CreateBuild)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_build",
		Description: "Delete a Build resource",
	}, tools.DeleteBuild)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_buildruns",
		Description: "List BuildRuns in a namespace with filtering options",
	}, tools.ListBuildRuns)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_buildrun",
		Description: "Get a specific BuildRun by name",
	}, tools.GetBuildRun)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_buildrun",
		Description: "Create a new BuildRun resource (either from existing Build or inline)",
	}, tools.CreateBuildRun)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "restart_buildrun",
		Description: "Restart a BuildRun by creating a new one",
	}, tools.RestartBuildRun)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "delete_buildrun",
		Description: "Delete a BuildRun resource",
	}, tools.DeleteBuildRun)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_buildstrategies",
		Description: "List BuildStrategies in a namespace with filtering options",
	}, tools.ListBuildStrategies)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_clusterbuildstrategies",
		Description: "List ClusterBuildStrategies with filtering options",
	}, tools.ListClusterBuildStrategies)

	log.Printf("MCP Server listening on stdin/stdout")
	log.Printf("Available tools: list_builds, get_build, create_build, delete_build, list_buildruns, get_buildrun, create_buildrun, restart_buildrun, delete_buildrun, list_buildstrategies, list_clusterbuildstrategies")

	if err := server.Run(context.Background(), mcp.NewStdioTransport()); err != nil {
		log.Fatal(err)
	}
}
