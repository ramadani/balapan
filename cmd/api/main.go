package main

import (
	"context"
	"flag"
	"github.com/go-redis/redis/v8"
	"github.com/go-zookeeper/zk"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	cmdadapter "github.com/ramadani/balapan/internal/adapter/app/command"
	"github.com/ramadani/balapan/internal/adapter/app/query"
	reposqlx "github.com/ramadani/balapan/internal/adapter/repository/sqlx"
	"github.com/ramadani/balapan/internal/adapter/rest/echo/handler"
	"github.com/ramadani/balapan/internal/app/command"
	"github.com/ramadani/balapan/internal/domain/rewards"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type RaceHandler struct {
	Enabled   bool            `yaml:"enabled"`
	Driver    string          `yaml:"driver"`
	Zookeeper ZookeeperConfig `yaml:"zookeeper"`
	Redis     RedisConfig     `yaml:"redis"`
}

type ZookeeperConfig struct {
	Address        []string      `yaml:"address"`
	SessionTimeout time.Duration `yaml:"sessionTimeout"`
}

type RedisConfig struct {
	Address string        `yaml:"address"`
	SyncIn  time.Duration `yaml:"syncIn"`
}

type Config struct {
	Address     string        `yaml:"address"`
	DB          string        `yaml:"db"`
	RaceHandler RaceHandler   `yaml:"raceHandler"`
	SleepIn     time.Duration `yaml:"sleepIn"`
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

	var claimCommand command.ClaimRewardsCommander

	claimHistoryCommand := command.NewClaimRewardsHistoryCommand(historyRepo)
	claimUsageCommand := command.NewClaimRewardsUsageCommand(rewardsRepo)

	if raceHandler := conf.RaceHandler; raceHandler.Enabled {
		switch conf.RaceHandler.Driver {
		case "zookeeper":
			zkConn, _, err := zk.Connect(raceHandler.Zookeeper.Address, raceHandler.Zookeeper.SessionTimeout)
			if err != nil {
				panic(err)
			}

			claimCommand = command.NewClaimRewardsMiddlewareCommand(claimUsageCommand, claimHistoryCommand)
			claimCommand = cmdadapter.NewClaimRewardsLockerCommand(claimCommand, zkConn)
		case "redis":
			redisClient := redis.NewClient(&redis.Options{Addr: raceHandler.Redis.Address})

			// key formats
			transactionQuotaLimitKeyFormat := "transaction-quota-limit-%s"
			transactionQuotaUsageKeyFormat := "transaction-quota-usage-%s"
			rewardsQuotaLimitKeyFormat := "rewards-quota-limit-%s"
			rewardsQuotaUsageKeyFormat := "rewards-quota-usage-%s"

			usageTtl := time.Duration(0)
			lockExpIn := time.Second * 3
			sleepRetry := 50 * time.Millisecond
			maxRetry := 10

			setNxTransactionQuotaUsageCommand := cmdadapter.NewSetNXRewardsQuotaUsageRedisCommand(
				rewardsRepo,
				redisClient,
				transactionQuotaUsageKeyFormat,
				usageTtl,
				lockExpIn,
				func(r *rewards.Rewards) int64 { return r.TransactionUsage },
			)
			setNxTransactionQuotaUsageCommand = command.NewSetNXRewardsQuotaRetryableCommand(setNxTransactionQuotaUsageCommand, maxRetry, sleepRetry)

			setNxRewardsQuotaUsageCommand := cmdadapter.NewSetNXRewardsQuotaUsageRedisCommand(
				rewardsRepo,
				redisClient,
				rewardsQuotaUsageKeyFormat,
				usageTtl,
				lockExpIn,
				func(r *rewards.Rewards) int64 { return r.RewardsUsage },
			)
			setNxRewardsQuotaUsageCommand = command.NewSetNXRewardsQuotaRetryableCommand(setNxRewardsQuotaUsageCommand, maxRetry, sleepRetry)

			setTransactionQuotaLimitCommand := cmdadapter.NewSetRewardsLimitRedisCommand(redisClient, transactionQuotaLimitKeyFormat, time.Hour*1)
			getTransactionQuotaLimitQuery := query.NewGetRewardsQuotaRedisQueryer(redisClient, transactionQuotaLimitKeyFormat)
			getTransactionQuotaUsageQuery := query.NewGetRewardsQuotaRedisQueryer(redisClient, transactionQuotaUsageKeyFormat)

			setRewardsQuotaLimitCommand := cmdadapter.NewSetRewardsLimitRedisCommand(redisClient, rewardsQuotaLimitKeyFormat, time.Hour*1)
			getRewardsQuotaLimitQuery := query.NewGetRewardsQuotaRedisQueryer(redisClient, rewardsQuotaLimitKeyFormat)
			getRewardsQuotaUsageQuery := query.NewGetRewardsQuotaRedisQueryer(redisClient, rewardsQuotaUsageKeyFormat)

			claimCommand = command.NewClaimRewardsMiddlewareCommand(
				command.NewClaimRewardsSetNXQuotaCommand(setNxRewardsQuotaUsageCommand),
				cmdadapter.NewClaimRewardsQuotaUsageRedisCommand(
					claimHistoryCommand,
					getRewardsQuotaLimitQuery,
					redisClient,
					rewardsQuotaUsageKeyFormat,
				),
			)
			claimCommand = command.NewClaimRewardsMiddlewareCommand(
				command.NewClaimRewardsSetNXQuotaCommand(setNxTransactionQuotaUsageCommand),
				cmdadapter.NewClaimTransactionQuotaUsageRedisCommand(
					claimCommand,
					getTransactionQuotaLimitQuery,
					redisClient,
					transactionQuotaUsageKeyFormat,
				),
			)

			go func() {
				for {
					ctx := context.Background()
					items, err := rewardsRepo.FindAll(ctx)
					if err != nil {
						log.Println(err)
					}

					for _, item := range items {
						if err = setTransactionQuotaLimitCommand.Do(ctx, item.ID, item.TransactionLimit); err != nil {
							log.Println(err)
						}
						if usage, err := getTransactionQuotaUsageQuery.Do(ctx, item.ID); err != nil {
							log.Println(err)
						} else {
							if usage != 0 {
								item.TransactionUsage = usage
							}
						}

						if err = setRewardsQuotaLimitCommand.Do(ctx, item.ID, item.RewardsLimit); err != nil {
							log.Println(err)
						}
						if usage, err := getRewardsQuotaUsageQuery.Do(ctx, item.ID); err != nil {
							log.Println(err)
						} else {
							if usage != 0 {
								item.RewardsUsage = usage
							}
						}

						if err = rewardsRepo.Update(ctx, item); err != nil {
							log.Println(err)
						}
					}

					time.Sleep(raceHandler.Redis.SyncIn)
				}
			}()
		}
	} else {
		claimCommand = command.NewClaimRewardsMiddlewareCommand(claimUsageCommand, claimHistoryCommand)
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
