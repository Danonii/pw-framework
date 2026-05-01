package framework

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
)

const (
	PW_TEMPLATE_FILE_TYPE = ".html"
	PW_TEMPLATE_DIR       = "templates/"
)

type Gate func(http.ResponseWriter, *http.Request) bool
type Middleware func(writer http.ResponseWriter, request *http.Request) bool
type HandleFunc func(http.ResponseWriter, *http.Request, *PWFramework)
type TemplateData map[string]any

/*
*
* This should not be needed, most of the time the better option should be to dynamically create a struct
* and feed it to the renderTemplate function.
*
 */
func Convert(a any) TemplateData {
	templateData, err := convertStructToMap(a)
	if err != nil {
		panic(err)
	}
	for field, val := range templateData {
		fmt.Println("KV Pair2: ", field, val)
	}
	return templateData
}

func GetLink(regexp_str string, request_path string) map[string]any {
	rgx, err := regexp.Compile(regexp_str)
	if err != nil {
		panic(err)
	}
	batata := rgx.FindStringSubmatch(request_path)
	fmt.Print(batata)
	return nil
}

type PWFramework struct {
	templates_path []string
	templates      *template.Template
	gates          map[string]Gate
}

type PWFrameworkInitData struct {
	Templates_init_func   func(*PWFramework)
	File_server_init_func func(*PWFramework)
	Routes_init_func      func(*PWFramework)
	Gates_init_func       func(*PWFramework)
}

func (framework *PWFramework) AddTemplate(path string) {
	framework.templates_path = append(framework.templates_path, PW_TEMPLATE_DIR+path)
}

func (framework *PWFramework) cacheBufferedTemplates() {
	if len(framework.templates_path) > 0 {
		framework.templates = template.Must(template.ParseFiles(framework.templates_path...))
	}
}

/*
* INIT FUNCTIONS
*
*
*
 */
func (framework *PWFramework) initTemplates(fn func(*PWFramework)) {
	fn(framework)
	framework.cacheBufferedTemplates()
}

func (framework *PWFramework) initFileServer(fn func(*PWFramework)) {
	fn(framework)
}

func (framework *PWFramework) initRoutes(fn func(*PWFramework)) {
	fn(framework)
}

func (framework *PWFramework) initGates(fn func(*PWFramework)) {
	framework.gates =
		make(map[string]Gate)
	fn(framework)
}

func (framework *PWFramework) RenderTemplate(w http.ResponseWriter, tmpl string, data any) {
	err := framework.templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (framework *PWFramework) AddRoute(path string, middlewares []Middleware, handleFunc HandleFunc) {
	finalHandler := func(w http.ResponseWriter, r *http.Request) {
		for _, mw := range middlewares {
			if !mw(w, r) {
				http.Error(w, "No permission", http.StatusUnauthorized)
				return
			}
		}
		handleFunc(w, r, framework)
	}
	http.HandleFunc(path, finalHandler)
}

/*
*
*
* GATES
*
*
 */

func (framework *PWFramework) GetGate(name string) Gate {
	g := (framework.gates)[name]
	if g == nil {
		fmt.Printf("The gate %s does not exist.\n", name)
	}
	return g
}

func (framework *PWFramework) Gate(name string, w http.ResponseWriter, r *http.Request) bool {
	gate := framework.GetGate(name)
	if gate == nil {
		fmt.Printf("Failed to execute gate %s.", name)
		return false
	}
	fmt.Printf("Executing gate %s.\n", name)
	return gate(w, r)

}

func (framework *PWFramework) AddGate(name string, gate Gate) {
	if (framework.gates)[name] != nil {
		fmt.Printf("Gate %s already exists.\n", name)
		return
	}
	(framework.gates)[name] = gate
}

/*
*
*
* MAIN FUNCTIONS
*
*
 */
func (framework *PWFramework) Init(data *PWFrameworkInitData) {
	framework.initGates(data.Gates_init_func)
	framework.initTemplates(data.Templates_init_func)
	framework.initFileServer(data.File_server_init_func)
	framework.initRoutes(data.Routes_init_func)
}

func (framework *PWFramework) Serve(address string) {
	log.Fatal(http.ListenAndServe(address, nil))
}
