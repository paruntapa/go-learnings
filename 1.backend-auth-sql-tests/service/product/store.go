package product

import (
	"backend-apis/types"
	"database/sql"
	"fmt"
	"strings"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CheckDuplicateProducts(product *types.Product) (*types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products WHERE name = ?", product.Name)
	if err != nil {
		fmt.Errorf("Error in Finding Product", err)
	}

	p := new(types.Product)

	for rows.Next() {
		p, err = scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
	}

	if p.ID != 0 {
		return nil, fmt.Errorf("Product Already Exists")
	}

	return nil, nil
}

func (s *Store) CreateProducts(product types.Product) error {
	_, err := s.db.Exec("INSERT INTO products (name, description, image, price, quantity) VALUES (?, ?, ?, ?, ?)", product.Name, product.Description, product.Image, product.Price, product.Quantity)

	if err != nil {
		fmt.Errorf("Product details are not correct")
		return err
	}

	return nil
}

func (s *Store) GetProducts() ([]*types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}

	products := make([]*types.Product, 0)
	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func (s *Store) UpdateProduct(product types.Product) error {
	_, err := s.db.Exec("UPDATE products SET name = ?, description = ?, image = ?, price = ?, quantity = ? WHERE id = ?", product.Name, product.Description, product.Image, product.Price, product.Quantity, product.ID)

	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetProductsByIDs(productIDs []int) ([]types.Product, error) {
	placeholders := strings.Repeat(",?", len(productIDs)-1)
	query := fmt.Sprintf("SELECT * FROM products WHERE id IN (?%s)", placeholders)

	args := make([]interface{}, len(productIDs))
	for i, v := range productIDs {
		args[i] = v
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	products := []types.Product{}

	for rows.Next() {
		p, err := scanRowsIntoProduct(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)

	}

	return products, nil
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
