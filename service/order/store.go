package cart

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/keshav78-78/ECOM/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateOrder(order types.Order) (int, error) {
	res, err := s.db.Exec("INSERT INTO	orders(userID, total, status, address) VALUES(?,?,?,?)", order.UserID, order.Total, order.Status, order.Address)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s *Store) GetProductsByIDs(ids []int) ([]types.Product, error) {
	// Create a placeholder string for the IN clause, like "(?, ?, ?)"
	placeholders := strings.Repeat("?,", len(ids)-1) + "?"
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (%s)", placeholders)

	// Convert int slice to a slice of interface{} for the query arguments
	args := make([]interface{}, len(ids))
	for i, v := range ids {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []types.Product
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) UpdateProduct(product types.Product) error {
	_, err := s.db.Exec("UPDATE products SET name = ?, description = ?, image = ?, price = ?, quantity = ? WHERE id = ?",
		product.Name, product.Description, product.Image, product.Price, product.Quantity, product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateOrderItem(orderItem types.OrderItem) error {
	_, err := s.db.Exec("INSERT INTO order_items (orderId, productId, quantity, price) VALUES(?,?,?,?)", orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.Price)
	return err

}

func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	p := new(types.Product)
	err := rows.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Image,
		&p.Price,
		&p.Quantity,
		&p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}
