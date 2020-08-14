package main

import (
	"bytes"
	"fmt"
	phoneDb "gophercises/Phone/db"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "passwood"
	dbname   = "shanghai"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	silence(phoneDb.Reset("postgres", psqlInfo, dbname))

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	silence(phoneDb.Migrate("postgres", psqlInfo))

	db, err := phoneDb.Open("postgres", psqlInfo)
	silence(err)
	defer db.Close()

	silence(db.Seed())

	phones, err := db.AllPhones()
	silence(err)

	for _, phone := range phones {
		fmt.Printf("evedropping on %+v\n", phone)
		number := normalize(phone.Number)
		if number != phone.Number {
			fmt.Println("updating or removing...", number)
			existing, err := db.FindPhone(number)
			silence(err)
			if existing != nil {
				//delete duplicate
				silence(db.DeletePhone(phone.ID))
			} else {
				//update with normalized number
				phone.Number = number
				silence(db.UpdatePhone(&phone))
			}
		} else {
			fmt.Println("no changes required")
		}
	}
}

func normalize(phone string) string {
	var b bytes.Buffer
	for _, c := range phone {
		if c >= '0' && c <= '9' {
			b.WriteRune(c)
		}
	}

	return b.String()
}

func silence(err error) {
	if err != nil {
		panic(err)
	}
}
