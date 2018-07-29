package HttpUtil

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var sh spaHandler

type spaHandler struct {
	fs   http.Handler
	path string
}

func init() {
	sh.path = "."
}

func InitFileServer(path string) (*spaHandler, error) {
	if _, err := os.Stat(path); err != nil {
		return nil, err
	}
	sh.path = path
	sh.fs = http.FileServer(http.Dir(path))
	return &sh, nil
}

// HandleSPA : Serves files when html/css/js is requested
func (s *spaHandler) HandleSPA(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, `/`)
	lastPart := pathParts[len(pathParts)-1]
	if strings.Contains(lastPart, ".") {
		sh.fs.ServeHTTP(w, r)
	} else {
		path := fmt.Sprintf("%s/index.html", sh.path)
		http.ServeFile(w, r, path)
	}
}

// CORSWrap : Handles cross origin authentication
func CORSWrap(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if origin := req.Header.Get("Origin"); origin != "" {
			rw.Header().Set("Access-Control-Allow-Origin", origin)
			rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			rw.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Auth-Token")
		}
		// Stop here if its Preflighted OPTIONS request
		if req.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(rw, req)
	})
}
