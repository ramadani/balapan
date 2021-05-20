package main

import (
	"flag"
	"github.com/go-zookeeper/zk"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	cmdadapter "github.com/ramadani/balapan/internal/adapter/app/command"
	reposqlx "github.com/ramadani/balapan/internal/adapter/repository/sqlx"
	"github.com/ramadani/balapan/internal/adapter/rest/echo/handler"
	"github.com/ramadani/balapan/internal/app/command"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type RaceHandler struct {
	Enabled bool   `yaml:"enabled"`
	Driver  string `yaml:"driver"`
}

type ZookeeperConfig struct {
	Address        []string      `yaml:"address"`
	SessionTimeout time.Duration `yaml:"sessionTimeout"`
}

type Config struct {
	Address     string          `yaml:"address"`
	DB          string          `yaml:"db"`
	RaceHandler RaceHandler     `yaml:"raceHandler"`
	Zookeeper   ZookeeperConfig `yaml:"zookeeper"`
	SleepIn     time.Duration   `yaml:"sleepIn"`
}

func main() {
	address := flag.String("address", "", "server address")

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

	if addr := *address; strings.TrimSpace(addr) != "" {
		conf.Address = addr
	}

	db, err := sqlx.Connect("postgres", conf.DB)
	if err != nil {
		panic(err)
	}

	rewardsRepo := reposqlx.NewRewardsRepository(db)
	historyRepo := reposqlx.NewHistoryRepository(db)

	claimHistoryCommand := command.NewClaimRewardsHistoryCommand(historyRepo)
	claimUsageCommand := command.NewClaimRewardsUsageCommand(rewardsRepo)
	claimCommand := command.NewClaimRewardsMiddlewareCommand(claimUsageCommand, claimHistoryCommand)

	if conf.RaceHandler.Enabled {
		switch conf.RaceHandler.Driver {
		case "zookeeper":
			zkConn, _, err := zk.Connect(conf.Zookeeper.Address, conf.Zookeeper.SessionTimeout)
			if err != nil {
				panic(err)
			}
			claimCommand = cmdadapter.NewClaimRewardsLockerCommand(claimCommand, zkConn)
		}
	}

	claimCommand = command.NewClaimRewardsSleeperCommand(claimCommand, conf.SleepIn)
	claimCommand = command.NewClaimRewardsLoggerCommand(claimCommand)

	rewardsHandler := handler.NewRewardsHandler(claimCommand)

	e := echo.New()
	e.HideBanner = true

	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.PUT("/rewards/:id/claim", rewardsHandler.Claim)

	// Start server
	e.Logger.Fatal(e.Start(conf.Address))
}
