package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ayo-ajayi/rest_api_template/model"
	"github.com/joho/godotenv"
	"github.com/twinj/uuid"
)

type DBStruct struct {
	Session *sql.DB
}

func getUuid() string {
	return uuid.NewV4().String()
}

var DBinit = func() (DBStruct , error){

	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			return DBStruct{nil}, err
		}
	}
	db, err := sql.Open("postgres", os.Getenv("NEON_DATABASE_URL"))
	if err != nil {
		return DBStruct{nil}, fmt.Errorf("error opening DB: %s", err.Error())
	}

	if err = db.Ping(); err != nil {
		return DBStruct{nil}, fmt.Errorf("could not ping: %s", err.Error())
	}
	if _, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS choice(
			id VARCHAR(36),
			go BOOL,
			come BOOL
		)`); err != nil {
		return DBStruct{nil}, fmt.Errorf("could not execute table creation: %s", err.Error())
	}

	return DBStruct{db}, nil
}

func (db DBStruct) CheckID(id string) (*model.Choice, error) {
	choice:=model.Choice{}
	//count records with id to be sure it is valid and also exists, else
	//return error here...use validate(playground package) to check if uuid is valid
	row := db.Session.QueryRow(`SELECT id, go, come FROM choice WHERE id= $1 LIMIT 1`, id)
	err := row.Scan(&choice.ID, &choice.Gone, &choice.Come)
	if err != nil{
		return nil, err
	}
	return &choice, nil
}
func (db DBStruct) GetChoice() ([]model.Choice, error) {
	rows, err := db.Session.Query(`SELECT id, go, come FROM choice LIMIT 20`)
	if err != nil {
		log.Println("error immediately after running query:", err)
		return nil, err
	}
	defer rows.Close()
	var choice []model.Choice

	for rows.Next() {
		var c model.Choice
		err := rows.Scan(&c.ID, &c.Gone, &c.Come)

		if err != nil {
			log.Println("error after rows.scan:", err)
			return nil, err
		}
		choice = append(choice, c)

	}

	if err := rows.Err(); err != nil {
		log.Println("error from rows.Err():", err)
		return nil, err
	}
	return choice, nil
}

func (db DBStruct) DeleteChoice(id string) error {
	if _, err := db.Session.Exec(`delete from choice where id = $1`, id); err != nil {
		return fmt.Errorf("500. internal server error. unable to delete from db: %s", err.Error())
	}
	return nil
}

func (db DBStruct) PostChoice(newChoice *model.Choice) error {
	newChoice.ID = getUuid()
	_, err := db.Session.Exec(`insert into choice(id, go, come) values($1, $2, $3)`, newChoice.ID, newChoice.Gone, newChoice.Come)
	if err != nil {
		return fmt.Errorf("DB code could not run with success: %s", err.Error())
	}

	return nil
}

func (db DBStruct) UpdateChoice(updateChoice model.Choice) error {
	_, err := db.Session.Exec(`UPDATE choice SET go=$1, come=$2 WHERE id=$3`, updateChoice.Gone, updateChoice.Come, updateChoice.ID)
	if err != nil {
		return fmt.Errorf(" DB code could not run with success: %s", err.Error())
	}
	return nil
}
func (db DBStruct)Close()error{
	return db.Session.Close()
}