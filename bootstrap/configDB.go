package bootstrap

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	logger "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"time"
)

type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Database        string        `yaml:"database"`
	UserName        string        `yaml:"user_name"`
	Password        string        `yaml:"password"`
	MaxOpenConn     int           `yaml:"max_open_conn"`
	MaxIdleConn     int           `yaml:"max_idle_conn"`
	ConnMaxLifeTime time.Duration `yaml:"conn_max_life_time"`
}

var Conn *sql.DB

func InitOptionalDB() {

	dbConfig := DatabaseConfig{}

	yamlFile, err := ioutil.ReadFile("./config/database_optional.yaml")

	if err != nil {
		log.Fatal("unable to read ./config/database_optional.yaml file ", err)
	}

	err = yaml.Unmarshal(yamlFile, &dbConfig)

	if err != nil {
		log.Fatal("unable to unmarshal ./config/database_optional.yaml file ", err)
	}

	logger.Info("config/database_optional.yaml loaded...")

	connString := dbConfig.UserName + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + fmt.Sprintf(`:%d`, dbConfig.Port) + ")/" + dbConfig.Database

	Conn, err = sql.Open("mysql", connString)

	if err != nil {
		log.Fatal("unable to connect to the database")
		panic(err.Error())
	}

	Conn.SetMaxIdleConns(dbConfig.MaxIdleConn)
	Conn.SetMaxOpenConns(dbConfig.MaxOpenConn)
	Conn.SetConnMaxLifetime(dbConfig.ConnMaxLifeTime)
}

func CloseOptionalDB() {
	err := Conn.Close()
	if err != nil {
		logger.Error("unable to close database")
	}
}
