package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(writer, "Hello World, %s!", request.URL.Path[1:])
}

type httpDir string

func (hd httpDir) Open(name string) (http.File, error) {
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return nil, errors.New("http: invalid character in file path")
	}

	dir := string(hd)
	if dir == "" {
		dir = "."
	}
	fullName := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name)))
	if !strings.HasSuffix(fullName, ".html") {
		if _, err := os.Stat(fullName); os.IsNotExist(err) {
			fullName = fullName + ".html"
		}
	}
	f, err := os.Open(fullName)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// Server An HTTP Web Server
type Server struct {
	mux    *http.ServeMux
	server *http.Server
}

// NewServer Creates a new Server
func NewServer() *Server {
	return &Server{mux: http.NewServeMux()}
}

// Serve Begins serving pages
func (s *Server) Serve() {
	s.mux.HandleFunc("/echo/", handler)

	files := http.FileServer(httpDir("docs"))
	s.mux.Handle("/developer-docs/", http.StripPrefix("/developer-docs/", files))
	s.mux.Handle("/", files)

	s.server = &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: s.mux,
	}
	s.server.ListenAndServe()
}

// Serve Starts the web server.
func Serve() {
	NewServer().Serve()
}
