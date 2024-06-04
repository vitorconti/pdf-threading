package seed

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type AirlineDelay struct {
	id             uint16
	airport_code   string
	airport_name   string
	time_label     string
	delay_late     string
	delay_security string
}

func PopulateDatabase() {
	os.Remove("sqlite-database.db")
	var wg sync.WaitGroup
	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
	createTable(sqliteDatabase)
	defer sqliteDatabase.Close()

	file, err = os.Open("../seed/airlines.csv")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	encodedFile := csv.NewReader(file)
	records, err := encodedFile.ReadAll()
	if err != nil {
		panic(err)
	}
	for i, r := range records {
		if i == 0 {
			continue
		}
		airlineDelay := AirlineDelay{
			airport_code:   r[0],
			airport_name:   r[1],
			time_label:     r[2],
			delay_late:     r[7],
			delay_security: r[9],
		}
		wg.Add(1)
		go insertAirline(&airlineDelay, sqliteDatabase, &wg)

	}

}
func createTable(db *sql.DB) {
	createAirlineDelaySQL := `CREATE TABLE airline_delay (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"airport_code" TEXT,
		"airport_name" TEXT,
		"time_label" TEXT,
		"delay_late" TEXT,
		"delay_security" TEXT		
	  );`

	statement, err := db.Prepare(createAirlineDelaySQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("created table airline")
}
func insertAirline(airline *AirlineDelay, db *sql.DB, wg *sync.WaitGroup) {
	insertAirlineSQL := fmt.Sprintf(`
	INSERT INTO airline_delay(airport_code,airport_name,time_label,delay_late,delay_security)
	VALUES("%s","%s","%s","%s","%s")`, airline.airport_code, airline.airport_name, airline.time_label, airline.delay_late, airline.delay_security)
	statement, err := db.Prepare(insertAirlineSQL)
	if err != nil {
		log.Fatal(err.Error(), "deu b.o aqui")
	}
	statement.Exec() // Execute SQL Statements
	log.Println("inserted succecced")
	wg.Done()
}
