package database

import (
	"context"
	"database/sql"
)

type Orders struct {
	db *sql.DB
}

type (
	Order struct {
		ID     int64  `json:"id"`
		Title  string `json:"title"`
		UserID int64  `json:"user_id"`

		Items []OrderItem `json:"items"`
	}

	OrderItem struct {
		ID       int64  `json:"id"`
		Title    string `json:"title"`
		Quantity int64  `json:"quantity"`
	}
)

func (r *Orders) Save(order *Order) error {
	return WithTransaction(context.Background(), r.db, func(tx Transaction) (err error) {
		if order.ID == 0 {
			err = r.insertOrderTx(tx, order)
		} else {
			err = r.updateOrderTx(tx, order)
		}

		if err != nil {
			return err
		}

		if err = r.deleteOrderItemsTx(tx, order.ID); err != nil {
			return err
		}

		if err := r.insertOrderItemsTx(tx, order.ID, order.Items); err != nil {
			return err
		}

		return nil
	}, nil)
}

func (r *Orders) insertOrderTx(tx Transaction, order *Order) (err error) {
	return tx.QueryRow(`
			INSERT INTO orders (title, user_id)
			VALUES ($1, $2) RETURNING id
		`,
		order.Title,
		order.UserID,
	).
		Scan(&order.ID)
}

func (r *Orders) updateOrderTx(tx Transaction, order *Order) (err error) {
	_, err = tx.Exec(`
		UPDATE orders SET
			title = $2,
			user_id = $3
		WHERE
			id = $1
		`,
		order.ID,
		order.Title,
		order.UserID,
	)
	return err
}

func (r *Orders) deleteOrderItemsTx(tx Transaction, orderID int64) (err error) {
	_, err = tx.Exec(`DELETE FROM order_items WHERE order_id = $1`, orderID)
	return err
}

func (r *Orders) insertOrderItemsTx(tx Transaction, orderID int64, items []OrderItem) (err error) {
	for i, item := range items {
		err = tx.QueryRow(`
			INSERT INTO order_items (order_id, title, quantity)
			VALUES ($1, $2, $3) RETURNING id
		`,
			orderID,
			item.Title,
			item.Quantity,
		).
			Scan(&items[i].ID)

		if err != nil {
			return err
		}
	}
	return nil
}
