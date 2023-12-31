package main

import (
	"EffectiveMobile/internal/graph"
	"EffectiveMobile/internal/http"
	"EffectiveMobile/internal/kafka"
	"EffectiveMobile/internal/repository"
	"EffectiveMobile/internal/service"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"sync"
	"time"
)

func main() {
	fmt.Println("Server started", "done")

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("env error: %s", err.Error())
	}
	//ожидание инициализации контейнера с postgres
	time.Sleep(10 * time.Second)
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("db initialize error: %s", err.Error())
	}

	cache, err := repository.NewRedisClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: "",
		DB:       0,
	})
	if err != nil {
		logrus.Fatalf("redis initialize error: %s", err.Error())
	}

	repository := repository.NewRepository(db, cache)
	service := service.NewService(repository)

	var wg sync.WaitGroup

	messageProcessor, err := kafka.NewKafkaMessageProcessor(viper.GetString("brokerAddr"), "test", service)
	if err != nil {
		logrus.Fatalf("Kafka init error: ", err.Error())
	}

	wg.Add(1)
	go messageProcessor.Start(&wg)

	httpHandlers := http.NewHandler(service)
	httpServer := new(http.Server)

	wg.Add(1)
	go httpServer.Run(viper.GetString("port"), httpHandlers.InitRoute(), &wg)

	resolver := graph.NewResolver(repository)
	gqlServer := new(graph.Server)

	wg.Add(1)
	go gqlServer.Run(resolver, &wg)

	wg.Wait()
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
