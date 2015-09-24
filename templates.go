package appkitcli

var tplMain string = `
package main

import(
	"{{url}}"
)

func main() {
	app := {{pkg}}.BuildApp()
	app.RunCli()
}
`

var tplConfig string = `
dev:
  debug: true

  # Path to store temporary files in.
  tmpDir: tmp
  dataDir: data
  
  # Host to serve on.  
  host: "localhost"
  # Port to serve on.
  port: 8000

  # URL used by the frontend to reach the app.
  url: "localhost:8000"
  
  # Enable server side rendering with phantomjs.
  serverRenderer:
    enabled: true
    cache: fs
    cacheLifetime: 3600
  
  # Configure the frontend template.
  frontend:
    indexTpl: public/index.html

  # Enable crawling of all public app pages to populate the cache.
  crawler:
    onRun: true
    concurrentRequests: 1

  # Backend configuration.
  backends:
    sql:
      url: "postgres://test:test@localhost/test"
`

var tplMainApp string = `
package {{pkg}}

import(
	kit "github.com/theduke/go-appkit"
	"github.com/theduke/go-appkit/app"
	"github.com/theduke/go-appkit/users"

	app_users "{{userPkg}}"
	// APPKIT:APP_IMPORTS - REMOVING THIS COMMENT WILL BREAK THE appkit cli TOOL!
)

func BuildApp() kit.App {
	app := app.NewApp()

	BuildBackends(app)
	BuildApps(app)

	// Configure UserService with profile model.
	userService := users.NewService(nil, app.DefaultBackend(), &app_users.Profile{})
	app.RegisterUserService(userService)

	BuildMigrations(app)
	app.PrepareBackends()

	return app
}

func BuildBackends(app kit.App) {
	// APPKIT:APP_BACKENDS - REMOVING THIS COMMENT WILL BREAK THE appkit cli TOOL!
}

func BuildApps(app kit.App) {
	// APPKIT:APP_APPS - REMOVING THIS COMMENT WILL BREAK THE appkit cli TOOL!
}

func BuildMigrations(app kit.App) {
	// APPKIT:APP_MIGRATIONS - REMOVING THIS COMMENT WILL BREAK THE appkit cli TOOL!
}
`

var tplMigrationsFile string = `
package {{pkg}}

import(
	db "github.com/theduke/go-dukedb"

	kit "github.com/theduke/go-appkit"
	"github.com/theduke/go-appkit/users"

	// APPKIT:{{backend}}_BACKEND_IMPORTS - REMOVING THIS COMMENT WILL BREAK THE appkit cli TOOL!
)

func Build{{upperCaseBackend}}Migrations(backend db.Backend, app kit.App) {
	handler := backend.(db.MigrationBackend).GetMigrationHandler()
	
	userService := app.UserService()
	userMigrations := users.GetUserMigrations(userService)
	handler.Add(userMigrations[0])
	handler.Add(userMigrations[1])

	// APPKIT:{{backend}}_BACKEND_MIGRATIONS - REMOVING THIS COMMENT WILL BREAK THE appkit cli TOOL!
}
`

var tplBuildMigrations string = `
	// Initialize migrations for {{name}} backend.
	Build{{upperCaseName}}Migrations(app.Backend("{{name}}"), app)
`

var tplMigration string = `
{{varName}} := db.Migration{
	Name: "{{name}}",
	Up: func(b db.MigrationBackend) error {
		{{up}}

		return nil
	},
	Down: func(b db.MigrationBackend) error {
		{{down}}

		return nil
	},
}
`

var tplApp string = `
package {{pkg}}

import(
	kit "github.com/theduke/go-appkit"
	"github.com/theduke/go-appkit/resources"
	// APPKIT:{{app}}_APP_IMPORTS - REMOVING THIS COMMENT WILL BREAK THE appkit cli TOOL!
)

// Avoid compile warning if resources not used.
// You can remove this once you have registered a resource.
var _ resources.AllowFindHook = nil 

func Build(app kit.App) {
	// APPKIT:APP_BUILD - REMOVING THIS COMMENT WILL BREAK THE appkit cli TOOL!
}
`

var tplUserProfile string = `
package users

import(
	"github.com/theduke/go-appkit/users"
)

type Profile struct {
	users.IntIDUserProfile

	FirstName string
	LastName  string
}
`

var tplModels string = `
package {{pkg}}

import(
	db "github.com/theduke/go-dukedb"
	"github.com/theduke/go-appkit/users"
)

// Avoid compiler errors.
var _ = db.StrIDModel{}
// You can remove this if you do plan on creating a model connected to a user in this file.
var _ = users.User{}

// APPKIT:APP_MODELS - REMOVING THIS COMMENT WILL BREAK THE appkit cli TOOL!
`

var tplModel string = `
/**
 * {{name}}.
 */

type {{name}} struct {
	db.IntIDModel
}

func (m {{name}}) Collection() string {
	return "{{collection}}" 
}

`

var tplResources string = `
package {{pkg}}

import(
	db "github.com/theduke/go-dukedb"
	kit "github.com/theduke/go-appkit"
)

// Avoid compiler errors.
var _ = db.StrIDModel{}
var _ kit.Model = nil

// APPKIT:RESOURCES - REMOVING THIS COMMENT WILL BREAK THE appkit cli TOOL!
`

var tplResourceHooks string = `
/**
 * Resource hooks for model {{name}}.
 * Check github.com/go-appkit/resources/interfaces.go for the available 
 * hooks you can implement.
 */

type {{name}}Resource struct {}
`

var tplRegisterResource string = `
	// Register {{model}} resource.
	app.RegisterResource(resources.NewResource(&{{model}}{}, &{{hooks}}{}, true))
`

var tplBackendSql string = `
	/**
	 * {{name}} backend.
	 */
	
	// Get backend settings from configuration.	
	url := app.Config().UString("backends.sql.url")
	if url == "" {
		panic("Invalid sql backend configuration: missing backends.sql.url setting.")
	}

	// Create backend.
	{{name}}Backend, err := sql.New("{{name}}", url)
	if err != nil {
		panic("Could not initialize backend {{name}}: " + err.Error())
	}
	{{name}}Backend.SetName("{{name}}")

	// Register backend with app.
	app.RegisterBackend({{name}}Backend)
`
