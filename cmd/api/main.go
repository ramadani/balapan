package main

import (
	"flag"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	reposqlx "github.com/ramadani/balapan/internal/adapter/repository/sqlx"
	"github.com/ramadani/balapan/internal/adapter/rest/echo/handler"
	"github.com/ramadani/balapan/internal/app/command"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Address string `yaml:"address"`
	Driver  string `yaml:"driver"`
	DB      string `yaml:"db"`
}

func main() {
	address := flag.String("address", ":3000", "server address")

	flag.Parse()

	file, err := os.Open("config.yml")
	if err != nil {
		panic(err)
	}

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	conf := new(Config)

	err = yaml.Unmarshal(fileContent, &conf)
	if err != nil {
		panic(err)
	}

	db, err := sqlx.Connect("postgres", conf.DB)
	if err != nil {
		log.Fatalln(err)
	}

	rewardsRepo := reposqlx.NewRewardsRepository(db)
	usageCommand := command.NewUsageRewardsCommand(rewardsRepo)

	rewardsHandler := handler.NewRewardsHandler(usageCommand)

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.PUT("/rewards/:id/usage", rewardsHandler.Usage)

	// Start server
	e.Logger.Fatal(e.Start(*address))
}
