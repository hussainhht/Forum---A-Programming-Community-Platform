package service

import (
	"fmt"
	"forum/Internal/model"
	"forum/Internal/query"
	"strconv"
)

func (service *Service) GetCategories() ([]model.Category, error) {
	rows, err := service.DB.Query("SELECT id, name FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var cat model.Category
		if err := rows.Scan(&cat.ID, &cat.Name); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

func (service *Service) validateCategories(categories []string) error {
	for _, v := range categories {
		cat, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("invalid category ID: %s", v)
		}
		q, err := query.SelectCategoryByID(service.DB, cat)
		if err != nil {
			return fmt.Errorf("failed to select category: %s", v)
		}
		if q == nil {
			return fmt.Errorf("invalid category ID: %s", v)
		}
	}

	return nil
}
