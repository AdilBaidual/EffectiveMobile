package main

import (
	"EffectiveMobile/internal/kafka"
	"EffectiveMobile/internal/repository"
	"EffectiveMobile/internal/service"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
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

	repository := repository.NewRepository(db)
	service := service.NewService(repository)

	messageProcessor, err := kafka.NewKafkaMessageProcessor(viper.GetString("brokerAddr"), "test", service)
	if err != nil {
		logrus.Fatalf("Kafka init error: ", err.Error())
		panic(err)
	}
	messageProcessor.Start()

	time.Sleep(time.Minute * 5)
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
