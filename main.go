package main

import (
	"book-store/bootstrap"
	"book-store/http/routers"
	_ "github.com/go-sql-driver/mysql"
)

var err error

func main() {

	routers.Router()
	bootstrap.InitOptionalDB()

	defer bootstrap.CloseOptionalDB()

}
