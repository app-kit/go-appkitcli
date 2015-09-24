# go-appkit-cli

Scaffolding cli for [go-appkit](https://github.com/theduke/go-appkit)

This command line tool allows you to very conveniently create new go-appkit setups and
to create new apps and resources.

## Usage

First, install the command.

`go install github.com/theduke/go-appkitcli/appkit`


### Create a new project.

`appkit bootstrap --backend="postgres" mynewapp`

### Create new apps inside your project.

`appkit app --resources="Todo,TodoProject" todo`

### Create a new resource within an app.

`appkit resource myapp MyNewResource`
