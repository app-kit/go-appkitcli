package appkitcli

import (
	"fmt"
	"os"
	"path"
	"strings"
	//"strconv"

	"github.com/spf13/cobra"
)

func Run() {
	var cli = &cobra.Command{
		Use:   "",
		Short: "Help",
		Long:  `Help`,

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("Use appkit -h to find out how to use the cli.\n")
		},
	}

	var rootPath, repoUrl, backend string
	cmdBootstrap := &cobra.Command{
		Use:   `bootstrap --repo="github.com/m/app" NAME`,
		Short: "Create a new appkit project.",
		Long:  "Create a new appkit project.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Print("Usage: appkit bootstrap NAME")
				os.Exit(1)
			}

			name := args[0]

			if rootPath == "" {
				rootPath, _ = os.Getwd()
				rootPath = path.Join(rootPath, name)
			}

			Bootstrap(rootPath, name, repoUrl, backend)
		},
	}
	cmdBootstrap.Flags().StringVarP(&rootPath, "path", "p", "", "Path to create project at. Defaults to working dir.")
	cmdBootstrap.Flags().StringVarP(&repoUrl, "repo", "r", "", "Repository url of the project.")
	cmdBootstrap.Flags().StringVarP(&backend, "backend", "b", "", "The backend to initialize.")
	cli.AddCommand(cmdBootstrap)

	var resources string
	cmdApp := &cobra.Command{
		Use:   `app --resources="res1,res2" NAME`,
		Short: "Create a new nested app.",
		Long:  "Create a new nested app.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 1 {
				fmt.Print("Usage: appkit app NAME")
				os.Exit(1)
			}
			name := args[0]

			if rootPath == "" {
				rootPath, _ = os.Getwd()
			}

			var resourceList []string
			if resources != "" {
				resourceList = strings.Split(strings.TrimSpace(resources), ",")
			}

			App(rootPath, name, resourceList)
		},
	}
	cmdApp.Flags().StringVarP(&rootPath, "path", "p", "", "Project root path.")
	cmdApp.Flags().StringVarP(&resources, "resources", "r", "", "Name of resources to create inside the app.")
	cli.AddCommand(cmdApp)

	cmdResource := &cobra.Command{
		Use:   `resource APP NAME`,
		Short: "Create a new resource in an app.",
		Long:  "Create a new resource in an app.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < 2 {
				fmt.Print("Usage: appkit app NAME")
				os.Exit(1)
			}
			app := args[0]
			resource := args[1]

			if rootPath == "" {
				rootPath, _ = os.Getwd()
			}

			Resource(rootPath, app, resource)
		},
	}
	cli.AddCommand(cmdResource)

	cli.Execute()
}
