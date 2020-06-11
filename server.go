package main

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/mitchellh/go-ps"
)

// Hello just outputs "hello, world" using an HTML template
func Hello(c echo.Context) error {
	return c.Render(http.StatusOK, "hello", "World")
}

type nomachineStatus struct {
	HostName         string
	NoMachineRunning bool
	ClientAttached   bool
}

func getFirstProcessByName(name string) (int, error) {
	processes, err := ps.Processes()
	if err != nil {
		return -1, err
	}

	for _, process := range processes {
		if process.Executable() == name {
			return process.Pid(), nil
		}
	}

	return -1, errors.New("Could not find process " + name)
}

// Index writes the status of NoMachine
func Index(c echo.Context) error {
	hostName, err := os.Hostname()
	if err != nil {
		c.Logger().Fatal(err)
	}
	status := nomachineStatus{HostName: hostName}

	noMachinePid, _ := getFirstProcessByName("nxserver.bin")
	noMachineClientPid, _ := getFirstProcessByName("nxexec")

	if noMachinePid >= 0 {
		status.NoMachineRunning = true
	}
	if noMachineClientPid >= 0 {
		status.ClientAttached = true
	}

	return c.Render(http.StatusOK, "index", status)
}

func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}

	e := echo.New()
	e.Renderer = t
	e.GET("/hello", Hello)
	e.GET("/", Index)

	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })
	e.Logger.Fatal(e.Start(":1323"))
}

// Template struct
type Template struct {
	templates *template.Template
}

// Render is used by Echo handler funcs
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
