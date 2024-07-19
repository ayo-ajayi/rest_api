package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ayo-ajayi/rest_api_template/model"
	"github.com/joho/godotenv"
	"github.com/twinj/uuid"
)

type DBClient struct {
	Session *sql.DB
}

type DBRepo interface {
	GetChoice(ctx context.Context) ([]model.Choice, error)
	CheckID(ctx context.Context, id string) (*model.Choice, error)
	DeleteChoice(ctx context.Context, id string) error
	PostChoice(ctx context.Context, newChoice *model.Choice) error
	UpdateChoice(ctx context.Context, updateChoice model.Choice) error
}

func getUuid() string {
	return uuid.NewV4().String()
}

var DBinit = func(ctx context.Context) (*DBClient, error) {
	resCh := make(chan *DBClient, 1)
	errCh := make(chan error, 1)
	go func() {
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
			return
		default:
			if _, err := os.Stat(".env"); err == nil {
				if err := godotenv.Load(".env"); err != nil {
					errCh <- fmt.Errorf("could not load .env file: %s", err.Error())
					return
				}
			}
			db, err := sql.Open("postgres", os.Getenv("DB_URI"))
			if err != nil {
				errCh <- fmt.Errorf("could not open db: %s", err.Error())
				return
			}

			if err = db.Ping(); err != nil {
				errCh <- fmt.Errorf("could not ping db: %s", err.Error())
				return
			}
			if _, err = db.Exec(
				`CREATE TABLE IF NOT EXISTS choice(
						id VARCHAR(36),
						go BOOL,
						come BOOL
					)`); err != nil {
				errCh <- fmt.Errorf("could not create table: %s", err.Error())
				return
			}
			resCh <- &DBClient{db}
		}
	}()

	select {
	case res := <-resCh:
		return res, nil
	case err := <-errCh:
		return nil, err
	}
}

func (db DBClient) CheckID(ctx context.Context, id string) (*model.Choice, error) {
	choice := model.Choice{}
	row := db.Session.QueryRowContext(ctx, `SELECT id, go, come FROM choice WHERE id= $1 LIMIT 1`, id)
	err := row.Scan(&choice.ID, &choice.Gone, &choice.Come)
	if err != nil {
		return nil, err
	}
	return &choice, nil
}
func (db DBClient) GetChoice(ctx context.Context) ([]model.Choice, error) {
	rows, err := db.Session.QueryContext(ctx, `SELECT id, go, come FROM choice LIMIT 20`)
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

func (db DBClient) DeleteChoice(ctx context.Context, id string) error {
	if _, err := db.Session.ExecContext(ctx, `delete from choice where id = $1`, id); err != nil {
		return fmt.Errorf("500. internal server error. unable to delete from db: %s", err.Error())
	}
	return nil
}

func (db DBClient) PostChoice(ctx context.Context, newChoice *model.Choice) error {
	newChoice.ID = getUuid()
	_, err := db.Session.ExecContext(ctx, `insert into choice(id, go, come) values($1, $2, $3)`, newChoice.ID, newChoice.Gone, newChoice.Come)
	if err != nil {
		return fmt.Errorf("DB code could not run with success: %s", err.Error())
	}

	return nil
}

func (db DBClient) UpdateChoice(ctx context.Context, updateChoice model.Choice) error {
	_, err := db.Session.ExecContext(ctx, `UPDATE choice SET go=$1, come=$2 WHERE id=$3`, updateChoice.Gone, updateChoice.Come, updateChoice.ID)
	if err != nil {
		return fmt.Errorf(" DB code could not run with success: %s", err.Error())
	}
	return nil
}
func (db DBClient) Close() error {
	return db.Session.Close()
}
