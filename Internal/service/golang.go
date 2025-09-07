package service

import "forum/Internal/query"

func (s *Service) GetCategoryIDByName(name string) (int, error) {
	return query.SelectCategoryIDByName(s.DB, name)
}
