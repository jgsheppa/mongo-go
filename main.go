package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jgsheppa/mongo-go/controllers"
	"github.com/jgsheppa/mongo-go/models"
	"github.com/spf13/viper"
)

type Magazine struct {
	Title string `json:"title"`
	Price string `json:"price"`
}

func init() {
	viper.SetConfigName("config")         // name of config file (without extension)
	viper.SetConfigType("yaml")           // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME/mongo-go") // call multiple times to add many search paths
	viper.AddConfigPath(".")              // optionally look for config in the working directory
	err := viper.ReadInConfig()           // Find and read the config file
	if err != nil {                       // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func main() {
	MONGO_URI := viper.GetString("mongodb")

	services, err := models.NewServices(MONGO_URI)
	if err != nil {
		panic(err)
	}
	must(err)

	magazineController := controllers.NewMagazine(services.Magazine)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/magazines", func(r chi.Router) {
		r.Get("/", magazineController.GetAllMagazines)
		r.Post("/{title}/{price}", magazineController.CreateMagazine)

		r.Route("/{magazineId}", func(r chi.Router){
			r.Get("/", magazineController.MagazineById)
			r.Delete("/", magazineController.DeleteMagazine)
		})
	})
	http.ListenAndServe(":3000", r)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
