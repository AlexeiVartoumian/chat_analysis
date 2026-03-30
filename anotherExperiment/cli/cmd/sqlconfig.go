package cmd

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDb() (*sql.DB, error) {

	fmt.Println("Attemption connection to db")
	err := godotenv.Load(".env")

	if err != nil {
		return nil, ErrorHandler(err, "env variables did not load ")
	}

	user := os.Getenv("db_user")
	password := os.Getenv("db_password")
	host := os.Getenv("db_host")
	database := os.Getenv("db_database")
	dbport := os.Getenv("db_port")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, dbport, user, password, database)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("SET client_encoding TO 'UTF8'")
	if err != nil {
		panic(err)
	}
	var encoding string
	db.QueryRow("SHOW client_encoding").Scan(&encoding)
	fmt.Println("Client encoding:", encoding)

	fmt.Println("Connected to db")
	return db, nil
}
