package repository

import (
	"awesomeProject3/models"
	"database/sql"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetByName(name string) (*models.Category, error) {
	var c models.Category
	err := r.db.QueryRow("SELECT id, name, slug FROM categories WHERE name = $1", name).Scan(&c.ID, &c.Name, &c.Slug)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	var c models.Category
	err := r.db.QueryRow("SELECT id, name, slug FROM categories WHERE id = $1", id).Scan(&c.ID, &c.Name, &c.Slug)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) GetBySlug(slug string) (*models.Category, error) {
	var c models.Category
	err := r.db.QueryRow("SELECT id, name, slug FROM categories WHERE slug = $1", slug).Scan(&c.ID, &c.Name, &c.Slug)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) List() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, slug FROM categories ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *CategoryRepository) Create(name, slug string) (inserted bool, err error) {
	res, err := r.db.Exec(
		"INSERT INTO categories (name, slug) VALUES ($1, $2) ON CONFLICT (slug) DO NOTHING",
		name, slug,
	)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n == 1, nil
}

func (r *CategoryRepository) ExistsBySlug(slug string) (bool, error) {
	var n int
	err := r.db.QueryRow("SELECT 1 FROM categories WHERE slug = $1", slug).Scan(&n)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
