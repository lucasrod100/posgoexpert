package database

import (
	"database/sql"

	"github.com/lucasrod100/posgoexpert/CleanArch/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{Db: db}
}

func (r *OrderRepository) Save(order *entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.ID, order.Price, order.Tax, order.FinalPrice)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepository) List() ([]entity.Order, error) {
	orders, err := r.Db.Query("SELECT * FROM orders")
	if err != nil {
		return nil, err
	}
	var result []entity.Order
	for orders.Next() {
		var order entity.Order
		err = orders.Scan(&order.ID, &order.Price, &order.Tax, &order.FinalPrice)
		if err != nil {
			return nil, err
		}
		result = append(result, order)
	}
	return result, nil
}
