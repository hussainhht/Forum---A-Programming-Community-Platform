package query

import (
	"database/sql"
	"forum/Internal/model"
	"strings"
)

const selectCategoryIDByNameQ = `SELECT id FROM categories WHERE LOWER(name)=LOWER(?) LIMIT 1`
const selectCategories = `SELECT * FROM categories`

func SelectCategoryIDByName(db *sql.DB, name string) (int, error) {
	var id int
	err := db.QueryRow(selectCategoryIDByNameQ, name).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, err
}

func GetCategories(db *sql.DB) ([]model.Category, error) {
	rows, err := db.Query(selectCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var name string
		var id int
		if err := rows.Scan(&id, &name); err != nil {
			return nil, err
		}
		categories = append(categories, model.Category{Name: strings.ToLower(name), ID: id})
	}
	return categories, nil
}
