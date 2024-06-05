package seed

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type AirlineDelay struct {
	Id             uint16
	Airport_code   string
	Airport_name   string
	Time_label     string
	Delay_late     string
	Delay_security string
}

func PopulateDatabase() {
	os.Remove("sqlite-database.db")
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
	var insertString string
	for i, r := range records {
		if i == 0 {
			continue
		}
		airlineDelay := AirlineDelay{
			Airport_code:   r[0],
			Airport_name:   r[1],
			Time_label:     r[2],
			Delay_late:     r[7],
			Delay_security: r[9],
		}
		insertAirlineSQL := fmt.Sprintf(
			`INSERT INTO airline_delay(airport_code,airport_name,time_label,delay_late,delay_security) VALUES('%s',"%s",'%s','%s','%s');`,
			airlineDelay.Airport_code,
			airlineDelay.Airport_name,
			airlineDelay.Time_label,
			airlineDelay.Delay_late,
			airlineDelay.Delay_security,
		)
		insertString = fmt.Sprintf("%s %s", insertString, insertAirlineSQL)
	}

	insertAirline(insertString, sqliteDatabase)
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
func insertAirline(sqlInsertString string, db *sql.DB) {

	_, err := db.Exec(sqlInsertString)
	if err != nil {
		log.Fatal(err.Error(), "deu b.o aqui")
	}

	log.Println("inserted succecced")
}
func RetriveAirlineDelays() (*sql.Rows, error) {
	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
	defer sqliteDatabase.Close()
	selectStatement := "SELECT * FROM airline_delay"
	rs, err := sqliteDatabase.Query(selectStatement)
	if err != nil {
		return nil, err
	}
	return rs, err
}
