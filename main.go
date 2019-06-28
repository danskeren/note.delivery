package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/packago/config"
	"github.com/danskeren/note.delivery/note"
	"github.com/danskeren/note.delivery/templates"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	limiter "github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors.Handler)

	store := memory.NewStore()
	rate, err := limiter.NewRateFromFormatted("500-H")
	if err != nil {
		panic(err)
	}
	limiterMiddleware := stdlib.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))
	limiterMiddleware.OnLimitReached = rateLimitGET

	r.Group(func(r chi.Router) {
		r.Use(limiterMiddleware.Handler)
		r.Get("/", indexGET)
		r.Get("/protect-your-privacy", protectYourPrivacyGET)
		r.Get("/privacy-policy", privacyPolicyGET)
		r.Post("/", note.CreateNote)
		r.Get("/{noteid}", note.NoteGET)
		r.Post("/{noteid}", note.UnlockNotePOST)
		r.Post("/{noteid}/delete", note.DeleteNotePOST)
		r.NotFound(notFoundGET)
	})

	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	filesDir := filepath.Join(currentDir, "static")
	FileServer(r, "/static", http.Dir(filesDir))

	switch config.File().GetString("environment") {
	case "development":
		panic(http.ListenAndServe(config.File().GetString("port.development"), r))
	case "production":
		panic(http.ListenAndServe(config.File().GetString("port.production"), r))
	default:
		panic("Environment not set")
	}
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}
	fs := http.StripPrefix(path, http.FileServer(root))
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"
	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func notFoundGET(w http.ResponseWriter, r *http.Request) {
	commonData := templates.ReadCommonData(w, r)
	commonData.MetaTitle = "404"
	templates.Render(w, "not-found.html", map[string]interface{}{
		"Common": commonData,
	})
}

func indexGET(w http.ResponseWriter, r *http.Request) {
	commonData := templates.ReadCommonData(w, r)
	templates.Render(w, "index.html", map[string]interface{}{
		"Common": commonData,
	})
}

func protectYourPrivacyGET(w http.ResponseWriter, r *http.Request) {
	commonData := templates.ReadCommonData(w, r)
	templates.Render(w, "protect-your-privacy.html", map[string]interface{}{
		"Common": commonData,
	})
}

func privacyPolicyGET(w http.ResponseWriter, r *http.Request) {
	commonData := templates.ReadCommonData(w, r)
	templates.Render(w, "privacy-policy.html", map[string]interface{}{
		"Common": commonData,
	})
}

func rateLimitGET(w http.ResponseWriter, r *http.Request) {
	commonData := templates.ReadCommonData(w, r)
	commonData.MetaTitle = "Rate Limited"
	templates.Render(w, "rate-limit.html", map[string]interface{}{
		"Common": commonData,
	})
}
