package manager

import "fmt"

type Apps []*App

func (apps Apps) List() string {
	if len(apps) == 0 {
		return "No apps"
	}
	out := "Available apps:\n"
	for _, app := range apps {
		out += fmt.Sprintf(" * %s\n", app.Name)
	}
	return out
}

func (apps Apps) String() (out string) {
	for _, app := range apps {
		out += app.String()
	}
	return
}

// Keep track of the apps
var apps = make(Apps, 0)

// TODO save managers independent of apps?

// TODO or call it a schema?
type App struct {
	Name     string
	Managers []Manager // TODO What about other manager types?
}

func (app *App) String() (out string) {
	if len(app.Managers) == 0 {
		return fmt.Sprintf("-- %s (no tables)\n", app.Name)
	} else {
		out = fmt.Sprintf("-- %s\n", app.Name)
		for _, manager := range app.Managers {
			out += fmt.Sprintf("%s\n", manager.Create())
		}
	}
	return out
}

func (app *App) Add(manager Manager) error {
	// TODO return an error if the table or manager already exist?
	app.Managers = append(app.Managers, manager)
	return nil
}

// NewApp creates a new app - a simple way to aggregate managers
func NewApp(name string) *App {
	app := &App{Name: name}
	// TODO error if duplicate?
	apps = append(apps, app)
	return app
}

// Global functions
// ----

func SQL(all bool, names ...string) string {
	if all {
		return apps.String()
	}
	if len(names) == 0 {
		return apps.List()
	}
	var invalid bool
	for _, name := range names {
		invalid = invalid || !HasApp(name)
	}
	if invalid {
		return apps.List()
	}
	out := ""
	for _, name := range names {
		out += fmt.Sprintf("%s\n", GetApp(name))
	}
	return out
}

func AllApps() string {
	return apps.String()
}

func GetApp(name string) *App {
	for _, app := range apps {
		if app.Name == name {
			return app
		}
	}
	return nil
}

func HasApp(name string) bool {
	for _, app := range apps {
		if app.Name == name {
			return true
		}
	}
	return false
}

func GetApps() Apps {
	return apps
}
