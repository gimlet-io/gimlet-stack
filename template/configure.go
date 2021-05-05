package template

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/enescakir/emoji"
	"github.com/gimlet-io/gimlet-cli/commands/chart/ws"
	"github.com/gimlet-io/gimlet-stack/version"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Component struct {
	Name        string `json:"name,omitempty" yaml:"name"`
	Description string `json:"description,omitempty" yaml:"description"`
	OnePager    string `json:"onePager,omitempty" yaml:"onePager"`
	Schema      string `json:"schema,omitempty" yaml:"schema"`
	UISchema    string `json:"uiSchema,omitempty" yaml:"uiSchema"`
}

type StackDefinition struct {
	Name        string       `json:"name,omitempty" yaml:"name"`
	Description string       `json:"description,omitempty" yaml:"description"`
	Intro       string       `json:"intro,omitempty" yaml:"intro"`
	Components  []*Component `json:"components,omitempty" yaml:"components"`
}

func StackDefinitionFromRepo(repoUrl string) (string, error) {
	stackTemplates, err := cloneStackFromRepo(repoUrl)
	if err != nil {
		return "", err
	}

	return stackTemplates["stack-definition.yaml"], nil
}

func Configure(stackDefinition StackDefinition, existingStackConfig StackConfig) (string, error) {

	port := randomPort()

	workDir, err := ioutil.TempDir(os.TempDir(), "gimlet")
	if err != nil {
		panic(err)
	}
	writeTempFiles(workDir, schema, helmUISchema, string(existingValuesJson))
	defer removeTempFiles(workDir)
	browserClosed := make(chan int, 1)
	r := setupRouter(workDir, browserClosed)
	srv := http.Server{Addr: fmt.Sprintf(":%d", port), Handler: r}

	ctrlC := make(chan os.Signal, 1)
	signal.Notify(ctrlC, os.Interrupt)

	go srv.ListenAndServe()
	fmt.Fprintf(os.Stderr, "%v Configure on http://127.0.0.1:%d\n", emoji.WomanTechnologist, port)
	fmt.Fprintf(os.Stderr, "%v Close the browser when you are done\n", emoji.WomanTechnologist)
	openBrowser(fmt.Sprintf("http://127.0.0.1:%d", port))

	select {
	case <-ctrlC:
	case <-browserClosed:
	}

	fmt.Fprintf(os.Stderr, "%v Generating values..\n\n", emoji.FileFolder)
	srv.Shutdown(context.TODO())

	return "", nil
}

func randomPort() int {
	if version.String() == "idea" {
		return 28000
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(10000) + 20000
}

func writeTempFiles(workDir string, schema string, helmUISchema string, existingValues string) {
	ioutil.WriteFile(filepath.Join(workDir, "values.schema.json"), []byte(schema), 0666)
	ioutil.WriteFile(filepath.Join(workDir, "helm-ui.json"), []byte(helmUISchema), 0666)
	ioutil.WriteFile(filepath.Join(workDir, "values.json"), []byte(existingValues), 0666)
	ioutil.WriteFile(filepath.Join(workDir, "bundle.js"), bundleJs, 0666)
	ioutil.WriteFile(filepath.Join(workDir, "bundle.js.LICENSE.txt"), licenseTxt, 0666)
	ioutil.WriteFile(filepath.Join(workDir, "index.html"), indexHtml, 0666)
}

func removeTempFiles(workDir string) {
	os.Remove(workDir)
}

func setupRouter(workDir string, browserClosed chan int) *chi.Mux {
	r := chi.NewRouter()
	if version.String() == "idea" {
		//r.Use(middleware.Logger)
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:28000", "http://127.0.0.1:28000"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Get("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(browserClosed, w, r)
	})

	r.Post("/saveValues", func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&values)
		w.WriteHeader(200)
		w.Write([]byte("{}"))
	})

	filesDir := http.Dir(workDir)
	fileServer(r, "/", filesDir)

	return r
}

func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}
}

// fileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
