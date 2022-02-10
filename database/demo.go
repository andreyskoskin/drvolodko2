package database

import (
	"fmt"

	_ "github.com/lib/pq"
)

func Demo() error {
	var db = DB{
		Name:     "postgres",
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
	}

	var conn, err = db.Connect()
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	var orders = Orders{db: conn}

	var anOrder = &Order{
		Title:  "Заказ материалов",
		UserID: 1234,
		Items: []OrderItem{
			{Title: "Лампочки", Quantity: 15},
			{Title: "Гвозди", Quantity: 100},
		},
	}

	printOrder("New", anOrder)

	if err = orders.Save(anOrder); err != nil {
		return fmt.Errorf("can not save order: %w", err)
	}

	printOrder("Saved", anOrder)

	anOrder.Title = "Заказ материалов в электроцех"
	anOrder.Items = anOrder.Items[:len(anOrder.Items)-1] // remove last item
	anOrder.Items[0].Quantity = 12
	anOrder.Items = append(anOrder.Items,
		OrderItem{
			Title:    "Катушка медной проволоки 0.5 мм",
			Quantity: 3,
		}, OrderItem{
			Title:    "Изолента синяя",
			Quantity: 8,
		},
	)

	printOrder("Modified", anOrder)

	if err = orders.Save(anOrder); err != nil {
		return fmt.Errorf("can not save order: %w", err)
	}

	printOrder("Saved", anOrder)

	return nil
}

func printOrder(title string, order *Order) {
	fmt.Printf("%s:\n", title)
	fmt.Printf("    %s %s | user ID: %d, items count: %d\n", idString(order.ID), order.Title, order.UserID, len(order.Items))
	for n, item := range order.Items {
		fmt.Printf("        #%d. %s %s (%d)\n", n+1, idString(item.ID), item.Title, item.Quantity)
	}
	fmt.Println()
}

func idString(id int64) string {
	if id == 0 {
		return "(NEW)"
	}
	return fmt.Sprintf("ID:%d", id)
}
