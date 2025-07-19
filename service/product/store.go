package product

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

func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]types.Product, 0)
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	return products, nil
}

func (s *Store) CreateProduct(product types.Product) error {
	_, err := s.db.Exec(
		"INSERT INTO products(name, description, price) VALUES(?, ?, ?)",
		product.Name, product.Description, product.Price,
	)
	return err
}

// ✅ Newly added method to implement types.ProductStore
func (s *Store) GetProductsByIDs(ids []int) ([]types.Product, error) {
	if len(ids) == 0 {
		return []types.Product{}, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (%s)", strings.Join(placeholders, ","))
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]types.Product, 0)
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
	_, err := s.db.Exec(`
		UPDATE products 
		SET name = ?, description = ?, image = ?, price = ?, quantity = ?
		WHERE id = ?`,
		product.Name, product.Description, product.Image, product.Price, product.Quantity, product.ID,
	)
	return err
}

func scanRowsIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)

	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return product, nil
}
