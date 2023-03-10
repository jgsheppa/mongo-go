package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth"
	"github.com/jgsheppa/mongo-go/auth"
	"github.com/jgsheppa/mongo-go/controllers"
	middlewares "github.com/jgsheppa/mongo-go/middlewares"
	"github.com/jgsheppa/mongo-go/models"
	"github.com/spf13/viper"
)

var TokenAuth *jwtauth.JWTAuth

func init() {
	viper.SetConfigName("config")               // name of config file (without extension)
	viper.SetConfigType("yaml")                 // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/workspaces/mongo-go") // call multiple times to add many search paths
	viper.AddConfigPath("/etc/secrets")
	viper.AddConfigPath(".") // optionally look for config in the working directory
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			fmt.Println("warning: config.yaml not found")
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("fatal error loading config file: %w", err))
		}
	}

	Secret := viper.GetString("JWT_SECRET")
	TokenAuth = jwtauth.New("HS256", []byte(Secret), nil)
}

func main() {
	MONGO_URI := viper.GetString("mongodb")
	PASSWORD_PEPPER := viper.GetString("PASSWORD_PEPPER")
	// Inject secrets
	auth.TokenAuth = TokenAuth
	models.PasswordPepper = PASSWORD_PEPPER

	s := CreateNewServer()
	s.MountHandlers(MONGO_URI)
	http.ListenAndServe(":3000", s.Router)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type Server struct {
	Router   *chi.Mux
	Services *models.Services
}

func CreateNewServer() *Server {
	s := &Server{}
	s.Router = chi.NewRouter()
	return s
}

func (s *Server) MountHandlers(mongoURI string) {
	services, err := models.NewServices(mongoURI)
	if err != nil {
		panic(err)
	}
	must(err)

	magazineController := controllers.NewMagazine(services.Magazine)
	userController := controllers.NewUser(services.User)

	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Timeout(time.Minute * 3))
	s.Router.Use(middleware.StripSlashes)
	// TODO: improve CORS once API has frontend
	s.Router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Enable httprate request limiter of 100 requests per minute.
	s.Router.Use(httprate.Limit(100, 1*time.Minute, httprate.WithKeyFuncs(httprate.KeyByIP), httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
		// We can send custom responses for the rate limited requests, e.g. a JSON message
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error": "Too many requests"}`))
	})))

	s.Router.Get("/", HelloWorld)

	s.Router.Route("/magazines", func(r chi.Router) {
		r.Get("/", magazineController.GetAllMagazines)
		r.Get("/slug/{magazineSlug:[a-zA-Z ]+}", magazineController.MagazineBySlug)

		// Protected update routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(auth.TokenAuth))
			r.Use(jwtauth.Authenticator)

			r.Post("/{title}/{price}", magazineController.CreateMagazine)
			r.Put("/{id}/{title}/{price}", magazineController.UpdateMagazine)
		})

		r.Route("/search", func(r chi.Router) {
			r.Get("/{field:[a-zA-Z ]+}/{term:[a-zA-Z ]+}", magazineController.SearchMagazines)
		})

		r.Route("/{magazineId}", func(r chi.Router) {
			r.Get("/", magazineController.MagazineById)

			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(auth.TokenAuth))
				r.Use(jwtauth.Authenticator)

				r.Delete("/", magazineController.DeleteMagazine)
			})
		})

		r.Route("/aggregations", func(r chi.Router) {
			r.Get("/price/{price}", magazineController.AggregateMagazinePrice)
		})
	})

	s.Router.Route("/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(auth.TokenAuth))
			r.Use(middlewares.Authenticator)

			r.Get("/me", userController.GetUser)
		})

		r.Group(func(r chi.Router) {
			r.Post("/login", userController.Login)
			r.Post("/logout", userController.Logout)
		})
	})
}

// HelloWorld api Handler
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
