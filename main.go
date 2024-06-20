package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"time"
)

var conf Config

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(fmt.Errorf("unable to decode into struct: %w", err))
	}
}

func main() {
	o := NewObserver(&conf)
	o.Init()

	duration, err := time.ParseDuration(conf.OpenTransactionNotifier.TimeInterval)
	if err != nil {
		log.Fatalf("error while parsing time duration %s", err)
	}

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			o.Start()
		}
	}
}
