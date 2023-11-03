package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

const (
	host           = "localhost"
	port           = 5432
	user           = "postgres"
	password       = "stanners2020"
	dbname         = "test"
	dbMaxIdleConns = 4
	dbMaxConns     = 100
	totalWorker    = 100
	csvFile        = "majestic_million.csv"
)

var dataHeaders = make([]string, 0)

func main() {

}

func openDbConnection() (*sql.DB, error) {
	log.Println("=> open db connection")

	dbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(dbMaxConns)
	db.SetMaxIdleConns(dbMaxIdleConns)

	return db, nil
}

func openCsvFile() (*csv.Reader, *os.File, error) {
	log.Println("=> open csv file")

	f, err := os.Open(csvFile)
	if err != nil {
		return nil, nil, err
	}

	reader := csv.NewReader(f)
	return reader, f, nil
}
func dispatchWorkers(db *sql.DB, jobs <-chan []interface{}, wg *sync.WaitGroup) {
	for workerIndex := 0; workerIndex <= totalWorker; workerIndex++ {
		go func(workerIndex int, db *sql.DB, jobs <-chan []interface{}, wg *sync.WaitGroup) {
			counter := 0

			for job := range jobs {
				doTheJob(workerIndex, counter, db, job)
				wg.Done()
				counter++
			}
		}(workerIndex, db, jobs, wg)
	}
}
