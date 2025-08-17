// Package connect provides ...
package config

import (
	"Template/pkg/utils"
	"Template/pkg/utils/go-utils/database"
	"Template/pkg/utils/go-utils/encryptDecrypt"
	httpUtils "Template/pkg/utils/go-utils/http"
	"fmt"

	"log"
	"net/http"
)

func CreateConnection() {
	username, dbConfigErr := encryptDecrypt.Decrypt(utils.GetEnv("POSTGRES_USERNAME"), utils.GetEnv("SECRET_KEY"))
	if dbConfigErr != nil {
		fmt.Println("error encrypting your classified text: ", dbConfigErr)
	}
	password, dbConfigErr := encryptDecrypt.Decrypt(utils.GetEnv("POSTGRES_PASSWORD"), utils.GetEnv("SECRET_KEY"))
	if dbConfigErr != nil {
		fmt.Println("error encrypting your classified text: ", dbConfigErr)
	}
	host, dbConfigErr := encryptDecrypt.Decrypt(utils.GetEnv("POSTGRES_HOST"), utils.GetEnv("SECRET_KEY"))
	if dbConfigErr != nil {
		fmt.Println("error encrypting your classified text: ", dbConfigErr)
	}
	dbName, dbConfigErr := encryptDecrypt.Decrypt(utils.GetEnv("DATABASE_NAME"), utils.GetEnv("SECRET_KEY"))
	if dbConfigErr != nil {
		fmt.Println("error encrypting your classified text: ", dbConfigErr)
	}
	fmt.Println("username: ", username)
	fmt.Println("password: ", password)
	fmt.Println("host: ", host)
	fmt.Println("dbName: ", dbName)
	httpUtils.Client.New(&http.Client{})
	database.PostgreSQLConnect(
		username,
		password,
		host,
		dbName,
		utils.GetEnv("POSTGRES_PORT"),
		utils.GetEnv("POSTGRES_SSL_MODE"),
		utils.GetEnv("POSTGRES_TIMEZONE"),
	)
	err := database.DBConn.AutoMigrate()

	if err != nil {
		log.Fatal(err.Error())
	}

}
