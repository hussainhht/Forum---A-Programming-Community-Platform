package service

import (
	"fmt"
	"forum/Internal/model"
	"forum/Internal/query"
	"strconv"
)

func (service *Service) GetCategories() ([]model.Category, error) {
	categories, err := query.GetCategories(service.DB)
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (service *Service) GetCategoriesNames() []string {
	categories, err := query.GetCategories(service.DB)
	if err != nil {
		return nil
	}
	var names []string
	for _, c := range categories {
		names = append(names, c.Name)
	}
	return names
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
