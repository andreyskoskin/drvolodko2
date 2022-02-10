package database

import (
	"database/sql"
	"fmt"
)

const (
	initOrders = `
		CREATE TABLE IF NOT EXISTS orders
		(
			id      SERIAL,
			title   VARCHAR,
			user_id INT,
			
			PRIMARY KEY (id)
		)
	`
	initOrderItems = `
		CREATE TABLE IF NOT EXISTS order_items
		(
			id       SERIAL,
			title    VARCHAR,
			quantity INT,
			order_id INT,

			PRIMARY KEY (id)
		)
	`
)

type DB struct {
	Name     string
	Host     string
	Port     int
	User     string
	Password string
}

func (c DB) Connect() (db *sql.DB, err error) {
	var connStr = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name)

	if db, err = sql.Open("postgres", connStr); err != nil {
		return nil, err
	}

	if _, err = db.Exec(initOrders); err != nil {
		return nil, disconnect(db, err)
	}

	if _, err = db.Exec(initOrderItems); err != nil {
		return nil, disconnect(db, err)
	}

	return db, nil
}

func disconnect(db *sql.DB, err error) error {
	_ = db.Close()
	return err
}
