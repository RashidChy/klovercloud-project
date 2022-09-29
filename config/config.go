package config

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var DbUser string
var DbPass string
var RunMode string
var ClusterEndpoint string
var ServerPort int
var ServerPortStr string

func InitEnVars() error {
	//checking runMode
	RunMode = os.Getenv("RUN_MODE")
	if RunMode == "" {
		RunMode = DEVELOP
	}
	var err error
	log.Println("RUN MODE:", RunMode)

	//loading envArs from .env file is runMode != PRODUCTION
	if RunMode != PRODUCTION {
		err = godotenv.Load()
		if err != nil {
			log.Println("ERROR: ", err.Error())
			return err
		}
	}
	var boolVal bool
	DbUser, boolVal = os.LookupEnv("DB_USER")
	if boolVal == false {
		return errors.New("DB_USER not fount in EnVars")
	}

	DbPass, boolVal = os.LookupEnv("DB_PASS")
	if boolVal == false {
		return errors.New("DB_PASS not fount in EnVars")
	}

	ClusterEndpoint, boolVal = os.LookupEnv("CLUSTER_ENDPOINT")
	if boolVal == false {
		return errors.New("CLUSTER_ENDPOINT not fount in EnVars")
	}

	ServerPortStr, boolVal = os.LookupEnv("SERVER_PORT")
	ServerPort, err = strconv.Atoi(ServerPortStr)
	if err != nil {
		return err
	}

	return nil
}
