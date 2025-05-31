//go:build embed_dashboard

package dashboard

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

//go:embed all:out
var dashboardFiles embed.FS

// Handler returns an http.Handler that serves the embedded dashboard files
func Handler() http.Handler {
	// Get the embedded filesystem starting from the 'out' directory
	dashboardFS, err := fs.Sub(dashboardFiles, "out")
	if err != nil {
		panic("failed to create dashboard filesystem: " + err.Error())
	}

	return &dashboardHandler{fs: dashboardFS}
}

type dashboardHandler struct {
	fs fs.FS
}

func (h *dashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Clean the path and remove leading slash
	filePath := strings.TrimPrefix(path.Clean(r.URL.Path), "/")

	// If the path is empty, serve index.html
	if filePath == "" || filePath == "." {
		filePath = "index.html"
	}

	// Check if file exists
	file, err := h.fs.Open(filePath)
	if err != nil {
		// For SPA routing, serve index.html for non-existent files
		// unless the request is for an asset (has file extension)
		if !strings.Contains(filePath, ".") {
			filePath = "index.html"
			file, err = h.fs.Open(filePath)
			if err != nil {
				http.NotFound(w, r)
				return
			}
		} else {
			http.NotFound(w, r)
			return
		}
	}
	defer file.Close()

	// Set content type based on file extension
	setContentType(w, filePath)

	// Add no-cache headers for dashboard requests
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Read file content and serve it
	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Serve the content
	w.Header().Set("Content-Length", string(rune(len(content))))
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

func setContentType(w http.ResponseWriter, filePath string) {
	ext := path.Ext(filePath)
	switch ext {
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case ".json":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	case ".woff":
		w.Header().Set("Content-Type", "font/woff")
	case ".woff2":
		w.Header().Set("Content-Type", "font/woff2")
	case ".ttf":
		w.Header().Set("Content-Type", "font/ttf")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}
}
