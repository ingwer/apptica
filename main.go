package main

import (
	"apptica/commander"
	"apptica/fetcher"
	"apptica/grpcserver"
	"apptica/httpserver"
	"apptica/storage"

	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const httpShutdownTimeout = 15 * time.Second

func main() {
	db, err := initMySQL()
	if err != nil {
		log.Fatal("failed to connect to mysql")
	}

	defer func() {
		if err := db.Close(); err != nil {
			fmt.Println("failed to close DB connection")
		}
	}()

	st := storage.New(
		storage.WithDB(db),
	)

	endpoint := getConfigValue("API_ENDPOINT")
	token := getConfigValue("API_TOKEN")
	f := fetcher.New(
		fetcher.WithEndpoint(endpoint),
		fetcher.WithToken(token),
		fetcher.WithHTTPClient(&http.Client{}),
	)

	cmd := commander.New(
		commander.WithFetcher(f),
		commander.WithStorage(st),
	)

	httpServer := httpserver.NewServer(cmd)
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			log.Println("failed to shutdown http server")
		}
	}()

	grpcServer := grpcserver.New(":9000", cmd)
	defer grpcServer.GracefulStop()

	log.Printf("server started")
	loop()
	log.Printf("server stopped")
}

func initMySQL() (*sql.DB, error) {
	host := getConfigValue("MYSQL_HOST")
	port := getConfigValue("MYSQL_PORT")
	user := getConfigValue("MYSQL_USER")
	password := getConfigValue("MYSQL_PASSWORD")
	dbname := getConfigValue("MYSQL_DATABASE_NAME")

	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	return sql.Open("mysql", dataSource)
}

func loop() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	<-signals
}

func getConfigValue(key string) string {
	return os.Getenv(key)
}
