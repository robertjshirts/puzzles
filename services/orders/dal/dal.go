package dal

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
)

type SQLDal struct {
	db *sqlx.DB
}

var schema = `
DO $$ BEGIN
	CREATE TYPE puzzle_enum AS ENUM ('clock', 'megaminx', 'pyraminx', '2x2', '3x3', '4x4', '5x5', '6x6', '7x7', '8x8+', 'skewb', 'sqaure-1', 'other');
EXCEPTION
	WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
	CREATE TYPE order_enum AS ENUM ('cancelled', 'placed');
EXCEPTION
	WHEN duplicate_object THEN null;
END $$;

DROP TABLE IF EXISTS puzzles CASCADE;

CREATE TABLE IF NOT EXISTS shipping_info (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	address TEXT NOT NULL,
	area_code TEXT NOT NULL,
	city TEXT NOT NULL,
	state TEXT NOT NULL,
	country TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS payment_info (
	id TEXT PRIMARY KEY,
	card_number TEXT NOT NULL,
	expiration TEXT NOT NULL,
	cvv TEXT NOT NULL,
	area_code TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS order_info (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	shipping_info_id TEXT NOT NULL,
	payment_info_id TEXT NOT NULL,
	FOREIGN KEY (shipping_info_id) REFERENCES shipping_info(id),
	FOREIGN KEY (payment_info_id) REFERENCES payment_info(id)
);

CREATE TABLE IF NOT EXISTS puzzles (
	id TEXT PRIMARY KEY,
	order_info_id TEXT NOT NULL,
	description TEXT NOT NULL,
	price DECIMAL NOT NULL,
	puzzle_type puzzle_enum NOT NULL,
	FOREIGN KEY (order_info_id) REFERENCES order_info(id)
);
`

func NewSQLDal(host string, port int, user string, password string, dbname string, timeout time.Duration) (*SQLDal, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	fmt.Println(dsn)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	result := db.MustExec(schema)
	if _, err := result.RowsAffected(); err != nil {
		return nil, err
	}

	return &SQLDal{
		db: db,
	}, nil
}

func (d *SQLDal) Close() error {
	return d.db.Close()
}

func (d *SQLDal) GetOrderInfo(id string) (*OrderInfo, error) {
	var order OrderInfo
	err := d.db.Get(&order, "SELECT * FROM orders WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &order, err
}

func (d *SQLDal) CreateOrderInfo(order OrderInfo, payment PaymentInfo, shipping ShippingInfo, puzzles []Puzzle) error {
	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	_, err = tx.NamedExec("INSERT INTO payment_info (id, card_number, expiration, cvv, area_code) VALUES (:id, :card_number, :expiration, :cvv, :area_code)", payment)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.NamedExec("INSERT INTO shipping_info (id, name, address, area_code, city, state, country) VALUES (:id, :name, :address, :area_code, :city, :state, :country)", shipping)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.NamedExec("INSERT INTO order_info (id, name, shipping_info_id, payment_info_id) VALUES (:id, :name, :shipping_info_id, :payment_info_id)", order)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.NamedExec("Insert into puzzles (id, order_info_id, description, price, puzzle_type) VALUES (:id, :order_info_id, :description, :price, :puzzle_type)", puzzles)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
