package appkitcli

import (
	"fmt"
	"os"
	"path"
	"strings"

	db "github.com/theduke/go-dukedb"

	"github.com/theduke/go-appkit/utils"
)

func checkErr(err error) {
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func StringReplace(content string, data map[string]string) string {
	for key, val := range data {
		content = strings.Replace(content, "{{"+key+"}}", val, -1)
	}

	return content
}

func FileReplace(path string, data map[string]string) {
	rawContent, err := utils.ReadFile(path)
	checkErr(err)

	rawContent = []byte(StringReplace(string(rawContent), data))
	checkErr(utils.WriteFile(path, rawContent, false))
}

func ReplaceToken(token, content, newContent string) (string, bool) {
	tokenIndex := strings.Index(content, token)
	if tokenIndex == -1 {
		return "", false
	}

	tokenStartIndex := -1
	for i := tokenIndex; i >= 0; i-- {
		if content[i] == '\n' {
			tokenStartIndex = i + 1
			break
		}
	}

	tokenEndIndex := -1
	for i := tokenIndex; i < len(content); i++ {
		if content[i] == '\n' {
			tokenEndIndex = i - 1
			break
		}
	}

	if tokenStartIndex == -1 || tokenEndIndex == -1 {
		return "", false
	}

	// Ensure double newline at end of new content.
	if newContent[len(newContent)-2] != '\n' {
		if newContent[len(newContent)-1] != '\n' {
			newContent += "\n\n"
		} else {
			newContent += "\n"
		}
	}

	content = content[:tokenStartIndex] + newContent + content[tokenStartIndex:]

	return content, true
}

func FileReplaceToken(path, token, newContent string) bool {
	contents, err := utils.ReadFile(path)
	if err != nil {
		return false
	}

	newContent, ok := ReplaceToken(token, string(contents), newContent)
	if !ok {
		return false
	}

	contents = []byte(newContent)

	if err := utils.WriteFile(path, []byte(contents), false); err != nil {
		return false
	}

	return true
}

func determinePkgPath(rootPath string) string {
	gopath := os.Getenv("GOPATH") + "/"
	if gopath == "/" {
		fmt.Print("Could not determine $GOPATH")
		os.Exit(1)
	}

	if !strings.Contains(rootPath, gopath) {
		fmt.Printf("Error: Project root is not inside $GOPATH\n")
		os.Exit(1)
	}

	rootPath = strings.Replace(rootPath, path.Join(gopath, "src")+"/", "", 1)

	if rootPath[0] == '/' {
		rootPath = rootPath[1:]
	}

	return rootPath
}

func Backend(rootPath, pkg, backend string) {
	if ok, err := utils.FileExists(path.Join(rootPath, "app.go")); !ok || err != nil {
		fmt.Printf("Could not find app.go. Are you at the root of your project?")
		os.Exit(1)
	}

	var buildImports []string
	var build string

	switch backend {
	case "postgres", "postgresql":
		buildImports = []string{`_ "github.com/lib/pq"`, `"github.com/theduke/go-dukedb/backends/sql"`}
		build = StringReplace(tplBackendSql, map[string]string{
			"name":             "postgres",
			"upperCaseBackend": strings.ToUpper(string(backend[0])) + backend[1:],
		})
	case "mysql":
	case "memory":
	default:
		fmt.Printf("Unknown backend type: %v\n", backend)
		os.Exit(1)
	}

	migrations := StringReplace(tplMigrationsFile, map[string]string{
		"pkg":              pkg,
		"backend":          backend,
		"upperCaseBackend": strings.ToUpper(string(backend[0])) + backend[1:],
	})

	buildMigrations := StringReplace(tplBuildMigrations, map[string]string{
		"name":          backend,
		"upperCaseName": strings.ToUpper(string(backend[0])) + backend[1:],
	})

	appPath := path.Join(rootPath, "app.go")

	imports := ""
	for _, val := range buildImports {
		imports += "\t" + val + "\n"
	}

	// Add imports to app.go.
	if !FileReplaceToken(appPath, "APPKIT:APP_IMPORTS", imports) {
		fmt.Printf("APPKIT:APP_IMPORTS token not found in %v\n", appPath)
		os.Exit(1)
	}

	// Add build code to BuildBackends() function in app.go
	if !FileReplaceToken(appPath, "APPKIT:APP_BACKENDS", build) {
		fmt.Printf("APPKIT:APP_BACKENDS token not found in %v\n", appPath)
		os.Exit(1)
	}

	// Add code to BuildMigratins() function in app.go.
	if !FileReplaceToken(appPath, "APPKIT:APP_MIGRATIONS", buildMigrations) {
		fmt.Printf("APPKIT:APP_MIGRATIONS token not found in %v\n", appPath)
		os.Exit(1)
	}

	// Create migrations file.
	utils.WriteFile(path.Join(rootPath, backend+"_"+"migrations.go"), []byte(migrations), false)
}

func App(rootPath, appName string, resources []string) {
	if ok, err := utils.FileExists(path.Join(rootPath, "app.go")); !ok || err != nil {
		fmt.Printf("Could not find app.go. Are you at the root of your project?")
		os.Exit(1)
	}

	pkgPath := determinePkgPath(rootPath)

	// Build appName.go file.
	tpl := StringReplace(tplApp, map[string]string{
		"pkg": appName,
	})
	utils.WriteFile(path.Join(rootPath, "apps", appName, appName+".go"), []byte(tpl), true)

	// Build models.go file.
	tpl = StringReplace(tplModels, map[string]string{
		"pkg": appName,
	})
	utils.WriteFile(path.Join(rootPath, "apps", appName, "models.go"), []byte(tpl), true)

	// Build resources.go file.
	tpl = StringReplace(tplResources, map[string]string{
		"pkg": appName,
	})
	utils.WriteFile(path.Join(rootPath, "apps", appName, "resources.go"), []byte(tpl), true)

	// Add pkg import and build call to app.go.
	appPath := path.Join(rootPath, "app.go")
	FileReplaceToken(appPath, "APPKIT:APP_IMPORTS", fmt.Sprintf(`	app_%v "%v/apps/%v"`, appName, pkgPath, appName))
	FileReplaceToken(appPath, "APPKIT:APP_APPS", fmt.Sprintf("	app_%v.Build(app)", appName))

	// Now, build resources.
	for _, resource := range resources {
		Resource(rootPath, appName, resource)
	}
}

func Resource(rootPath, app, resourceName string) {
	if ok, err := utils.FileExists(path.Join(rootPath, "app.go")); !ok || err != nil {
		fmt.Printf("Could not find app.go. Are you at the root of your project?")
		os.Exit(1)
	}

	collection := db.Pluralize(db.CamelCaseToUnderscore(resourceName))
	if collection[len(collection)-1] != 's' {
		collection += "s"
	}

	// Write the model to the apps/APP/models.go file.
	tpl := StringReplace(tplModel, map[string]string{
		"name":       resourceName,
		"collection": collection,
	})
	FileReplaceToken(path.Join(rootPath, "apps", app, "models.go"), "APPKIT:APP_MODELS", tpl)

	// Write the resource hook struct to apps/APP/resources.go .
	tpl = StringReplace(tplResourceHooks, map[string]string{
		"name": resourceName,
	})
	FileReplaceToken(path.Join(rootPath, "apps", app, "resources.go"), "APPKIT:RESOURCES", tpl)

	// Now, finally, write the .RegisterResource to apps/APP/appname.go Build().
	tpl = StringReplace(tplRegisterResource, map[string]string{
		"model": resourceName,
		"hooks": resourceName + "Resource",
	})
	FileReplaceToken(path.Join(rootPath, "apps", app, app+".go"), "APPKIT:APP_BUILD", tpl)
}

func Bootstrap(rootPath, name, pkgUrl, backend string) bool {
	pkgPath := determinePkgPath(rootPath)

	appPath := path.Join(rootPath, "apps")

	if err := os.MkdirAll(appPath, 0777); err != nil {
		fmt.Printf("Could not create root directory %v: %v", rootPath, err)
		return false
	}

	// Create app.go file.
	tpl := StringReplace(tplMainApp, map[string]string{
		"pkg":     name,
		"url":     pkgPath,
		"userPkg": pkgPath + "/apps/users",
	})
	utils.WriteFile(path.Join(rootPath, "app.go"), []byte(tpl), false)

	// Create the main.go file.
	tpl = StringReplace(tplMain, map[string]string{
		"pkg": name,
		"url": pkgPath,
	})
	utils.WriteFile(path.Join(rootPath, name, "main.go"), []byte(tpl), true)

	// Create users app.
	tpl = StringReplace(tplUserProfile, map[string]string{})
	utils.WriteFile(path.Join(rootPath, "apps", "users", "models.go"), []byte(tpl), true)

	// Create backend, if any.
	if backend != "" {
		Backend(rootPath, name, backend)
	}

	// Write example config.
	utils.WriteFile(path.Join(rootPath, name, "config.yaml"), []byte(tplConfig), false)

	return true
}
