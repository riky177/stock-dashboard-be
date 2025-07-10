package models

import (
	"fmt"
	"stock-dashboard/db"
	"time"
)

type Product struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Price     float64   `json:"price" binding:"required,gt=0"`
	Stock     int       `json:"stock" binding:"required,gt=0"`
	Category  string    `json:"category" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductUpdate struct {
	ID        int64     `json:"id,omitempty"`
	Name      *string   `json:"name,omitempty"`
	Price     *float64  `json:"price,omitempty" binding:"omitempty,gt=0"`
	Stock     *int      `json:"stock,omitempty"`
	Category  *string   `json:"category,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductFilter struct {
	Name      string  `json:"name,omitempty"`
	Category  string  `json:"category,omitempty"`
	MinPrice  float64 `json:"min_price,omitempty"`
	MaxPrice  float64 `json:"max_price,omitempty"`
	MinStock  int     `json:"min_stock,omitempty"`
	MaxStock  int     `json:"max_stock,omitempty"`
	SortOrder string  `json:"sort_order,omitempty"`
	Limit     int     `json:"limit,omitempty"`
	Offset    int     `json:"offset,omitempty"`
}

type ProductListResult struct {
	Products   []Product `json:"products"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"total_pages"`
}

func (p *Product) Get() error {
	query := `
		SELECT id, name, price, stock, category, created_at, updated_at 
		FROM products WHERE id = $1
	`

	row := db.DB.QueryRow(query, p.ID)
	err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.Category, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (p *Product) Save() error {
	query := `
		INSERT INTO products (name, price, stock, category, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	now := time.Now()
	p.CreatedAt = now
	p.UpdatedAt = now

	err := db.DB.QueryRow(query, p.Name, p.Price, p.Stock, p.Category, p.CreatedAt, p.UpdatedAt).Scan(&p.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProductUpdate) Update() error {
	p.UpdatedAt = time.Now()

	query := `UPDATE products SET `
	args := []any{}
	argCount := 1

	if p.Name != nil {
		query += fmt.Sprintf("name = $%d,", argCount)
		args = append(args, *p.Name)
		argCount++
	}
	if p.Price != nil {
		query += fmt.Sprintf("price = $%d,", argCount)
		args = append(args, *p.Price)
		argCount++
	}
	if p.Stock != nil {
		query += fmt.Sprintf("stock = $%d,", argCount)
		args = append(args, *p.Stock)
		argCount++
	}
	if p.Category != nil {
		query += fmt.Sprintf("category = $%d,", argCount)
		args = append(args, *p.Category)
		argCount++
	}

	query += fmt.Sprintf("updated_at = $%d", argCount)
	args = append(args, p.UpdatedAt)
	argCount++

	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, p.ID)

	_, err := db.DB.Exec(query, args...)
	return err
}

func (p *Product) Delete() error {
	query := `DELETE FROM products WHERE id = $1`

	_, err := db.DB.Exec(query, p.ID)
	if err != nil {
		return err
	}

	return nil
}

func SearchProducts(searchTerm string, sortOrder string) ([]Product, error) {
	var products []Product

	order := "DESC"
	if sortOrder == "asc" {
		order = "ASC"
	}

	query := fmt.Sprintf(`
		SELECT id, name, price, stock, category, created_at, updated_at 
		FROM products 
		WHERE name ILIKE $1 OR category ILIKE $1
		ORDER BY created_at %s
	`, order)

	rows, err := db.DB.Query(query, "%"+searchTerm+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock,
			&product.Category, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func GetProductsWithPagination(filter ProductFilter) (*ProductListResult, error) {
	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	fmt.Println("Filter for pagination:", filter.MinPrice)

	page := filter.Offset
	if page <= 0 {
		page = 1
	}

	var args []any
	argCount := 0

	countQuery := `SELECT COUNT(*) FROM products WHERE 1=1`

	dataQuery := `SELECT id, name, price, stock, category, created_at, updated_at FROM products WHERE 1=1`

	filterClause := ""

	if filter.Name != "" {
		argCount++
		filterClause += fmt.Sprintf(" AND name ILIKE $%d", argCount)
		args = append(args, "%"+filter.Name+"%")
	}

	if filter.Category != "" {
		argCount++
		filterClause += fmt.Sprintf(" AND category ILIKE $%d", argCount)
		args = append(args, "%"+filter.Category+"%")
	}

	if filter.MinPrice > 0 {
		argCount++
		filterClause += fmt.Sprintf(" AND price >= $%d", argCount)
		args = append(args, filter.MinPrice)
	}

	if filter.MaxPrice > 0 {
		argCount++
		filterClause += fmt.Sprintf(" AND price <= $%d", argCount)
		args = append(args, filter.MaxPrice)
	}

	if filter.MinStock > 0 {
		argCount++
		filterClause += fmt.Sprintf(" AND stock >= $%d", argCount)
		args = append(args, filter.MinStock)
	}

	if filter.MaxStock > 0 {
		argCount++
		filterClause += fmt.Sprintf(" AND stock <= $%d", argCount)
		args = append(args, filter.MaxStock)
	}

	countQuery += filterClause
	dataQuery += filterClause

	var total int
	err := db.DB.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, err
	}

	sortOrder := "DESC"
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}
	dataQuery += fmt.Sprintf(" ORDER BY created_at %s", sortOrder)

	argCount++
	dataQuery += fmt.Sprintf(" LIMIT $%d", argCount)
	args = append(args, filter.Limit)

	offset := (page - 1) * filter.Limit
	argCount++
	dataQuery += fmt.Sprintf(" OFFSET $%d", argCount)
	args = append(args, offset)

	rows, err := db.DB.Query(dataQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock,
			&product.Category, &product.CreatedAt, &product.UpdatedAt)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	totalPages := (total + filter.Limit - 1) / filter.Limit
	if totalPages == 0 {
		totalPages = 1
	}

	return &ProductListResult{
		Products:   products,
		Total:      total,
		Page:       page,
		Limit:      filter.Limit,
		TotalPages: totalPages,
	}, nil
}
