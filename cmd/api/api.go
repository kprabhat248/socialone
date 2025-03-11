package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"socialone/internal/store"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	// http-swagger middleware
)



type application struct{
	config config
	store store.Storage



}
type config struct {
	addr string
	db dbConfig
	env string

}
type dbConfig struct{
	addr string
	maxOpenConns int
	maxidealConns int
	maxidealTime string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthcheckHandler)
		docsURL:= fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)
				r.Get("/", app.getUserHandler)

				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)


			})
			r.Group(func(r chi.Router) {
				r.Get("/feed",app.getUserFeedHandler)
			})
		})
	})

	return r
}
func (app *application) run(Mux http.Handler) error {
	//Docs
	//docs.SwaggerInfo.Version = version
	srv:=&http.Server{
		Addr: app.config.addr,
		Handler: Mux,
		WriteTimeout: 30 * time.Second,
		ReadTimeout: 10 * time.Second,
		IdleTimeout: time.Minute,
		ErrorLog: log.New(os.Stderr, "", log.LstdFlags),
	}
	log.Printf("Starting server on %s", app.config.addr)
	return srv.ListenAndServe()
}