package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	validator "gopkg.in/go-playground/validator.v9"

	"github.com/FOXCraft40/ezGoApi/src/controller"
)

// Validator strunct
type Validator struct {
	validator *validator.Validate
}

// Validate struct
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

//
func main() {
	initDefaultConfig()
	initService()
}

func initService() {
	e := echo.New()
	e.Validator = &Validator{validator: validator.New()}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: viper.GetStringSlice("server.cors"),
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
	}))

	// init Database connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Set endpoint for "status"
	status := e.Group("status")
	status.GET("", controller.Status)

	port := viper.GetInt("server.port")
	// Start server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}

func initDefaultConfig() {
	// init configuration
	viper.SetConfigName("conf")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/")

	// set defaults for configuration
	viper.SetDefault("server.port", 7001)
	if os.Getenv("DEV_ENV") == "true" {
		dc := []string{"*"}
		viper.Set("server.cors", dc)
	}

	if os.Getenv("DEBUG") == "true" {
		log.Println("** DEBUG MOD ENABLE**")
	}

	if os.Getenv("DEV_ENV") == "true" {
		log.Println("** DEVELOPMENT MOD ENABLE**")
	}

	// try reading in a config
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("could not read configuration: %s", err.Error())
	}
}
