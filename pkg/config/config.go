package config

import (
	"log"
	"os"
	"sync"

	"github.com/Plaenkler/BatteryHistory/pkg/handler"
	yaml "gopkg.in/yaml.v3"
)

var (
	instance *Config
	once     sync.Once
)

type Config struct {
	ServerAddress  string `yaml:"serverAddress"`
	ServerPort     string `yaml:"serverPort"`
	ServerUser     string `yaml:"serverUser"`
	ServerPassword string `yaml:"serverPassword"`
	WebPort        string `yaml:"webPort"`
}

func GetConfig() *Config {
	defer handler.HandlePanic("config")

	once.Do(func() {
		err := initConfig()
		if err != nil {
			log.Fatalf("[config] initialization failed - error: %s", err.Error())
		}
	})

	return instance
}

func initConfig() error {
	if _, err := os.Stat("./config/config.yaml"); err != nil {
		err = createConfig()
		if err != nil {
			return err
		}
	}

	file, err := os.Open("./config/config.yaml")
	if err != nil {
		return err
	}
	defer file.Close()

	err = yaml.NewDecoder(file).Decode(&instance)
	if err != nil {
		return err
	}

	return nil
}

func createConfig() error {
	config := Config{
		ServerAddress:  "127.0.0.1",
		ServerPort:     "8550",
		ServerUser:     "user",
		ServerPassword: "password",
		WebPort:        "9000",
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	_, err = os.Stat("./config")
	if os.IsNotExist(err) {
		err = os.Mkdir("./config", os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile("./config/config.yaml", data, 0600)
	if err != nil {
		return err
	}

	log.Println("[config] created default configuration file exiting...")
	os.Exit(0)
	return nil
}
