package dal

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"

	"github.com/puzzles/services/orders/gen"
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
	status TEXT NOT NULL,
	shipping_info_id TEXT NOT NULL,
	payment_info_id TEXT NOT NULL,
	FOREIGN KEY (shipping_info_id) REFERENCES shipping_info(id) ON DELETE CASCADE,
	FOREIGN KEY (payment_info_id) REFERENCES payment_info(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS puzzles (
	id TEXT PRIMARY KEY,
	order_info_id TEXT NOT NULL,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	price DECIMAL NOT NULL,
	puzzle_type puzzle_enum NOT NULL,
	FOREIGN KEY (order_info_id) REFERENCES order_info(id) ON DELETE CASCADE
);
`

func NewSQLDal(host string, port int, user string, password string, dbname string, timeout time.Duration) (*SQLDal, error) {
	// Initialize database connection
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Create schema
	result := db.MustExec(schema)
	if _, err := result.RowsAffected(); err != nil {
		return nil, err
	}

	// Return DAL
	return &SQLDal{
		db: db,
	}, nil
}

func (d *SQLDal) Close() error {
	return d.db.Close()
}

func (d *SQLDal) DeleteOrder(id string) *gen.Error {
	tx, err := d.db.Beginx()
	if err != nil {
		return &gen.Error{Code: 500, Message: err.Error()}
	}

	_, err = tx.Exec("DELETE FROM order_info WHERE id = $1", id)
	if err == sql.ErrNoRows {
		tx.Rollback()
		return &gen.Error{Code: 404, Message: "Order not found"}
	} else if err != nil {
		tx.Rollback()
		return &gen.Error{Code: 500, Message: err.Error()}
	}

	err = tx.Commit()
	if err != nil {
		return &gen.Error{Code: 500, Message: err.Error()}
	}

	return nil
}

func (d *SQLDal) GetOrder(id string) (*OrderInfo, *PaymentInfo, *ShippingInfo, *[]Puzzle, *gen.Error) {
	var order OrderInfo
	err := d.db.Get(&order, "SELECT * FROM order_info WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, nil, nil, nil, &gen.Error{Code: 404, Message: "Order not found"}
	} else if err != nil {
		return nil, nil, nil, nil, &gen.Error{Code: 500, Message: err.Error()}
	}

	var payment PaymentInfo
	err = d.db.Get(&payment, "SELECT * FROM payment_info WHERE id = $1", order.PaymentInfoId)
	if err != nil {
		return nil, nil, nil, nil, &gen.Error{Code: 500, Message: err.Error()}
	}

	var shipping ShippingInfo
	err = d.db.Get(&shipping, "SELECT * FROM shipping_info WHERE id = $1", order.ShippingInfoId)
	if err != nil {
		return nil, nil, nil, nil, &gen.Error{Code: 500, Message: err.Error()}
	}

	var puzzles []Puzzle
	err = d.db.Select(&puzzles, "SELECT * FROM puzzles WHERE order_info_id = $1", order.Id)
	if err != nil {
		return nil, nil, nil, nil, &gen.Error{Code: 500, Message: err.Error()}
	}

	return &order, &payment, &shipping, &puzzles, nil
}

func (d *SQLDal) CreateOrder(order OrderInfo, payment PaymentInfo, shipping ShippingInfo, puzzles []Puzzle) *gen.Error {
	tx, err := d.db.Beginx()
	if err != nil {
		return &gen.Error{Code: 500, Message: err.Error()}
	}

	_, err = tx.NamedExec("INSERT INTO payment_info (id, card_number, expiration, cvv, area_code) VALUES (:id, :card_number, :expiration, :cvv, :area_code)", payment)
	if err != nil {
		tx.Rollback()
		return &gen.Error{Code: 500, Message: err.Error()}
	}

	_, err = tx.NamedExec("INSERT INTO shipping_info (id, name, address, area_code, city, state, country) VALUES (:id, :name, :address, :area_code, :city, :state, :country)", shipping)
	if err != nil {
		tx.Rollback()
		return &gen.Error{Code: 500, Message: err.Error()}
	}

	_, err = tx.NamedExec("INSERT INTO order_info (id, name, status, shipping_info_id, payment_info_id) VALUES (:id, :name, :status, :shipping_info_id, :payment_info_id)", order)
	if err != nil {
		tx.Rollback()
		return &gen.Error{Code: 500, Message: err.Error()}
	}

	_, err = tx.NamedExec("Insert into puzzles (id, order_info_id, name, description, price, puzzle_type) VALUES (:id, :order_info_id, :name, :description, :price, :puzzle_type)", puzzles)
	if err != nil {
		tx.Rollback()
		return &gen.Error{Code: 500, Message: err.Error()}
	}

	err = tx.Commit()
	if err != nil {
		return &gen.Error{Code: 500, Message: err.Error()}
	}

	return nil
}
