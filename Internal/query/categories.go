package query

import "database/sql"

const selectCategoryIDByNameQ = `SELECT id FROM categories WHERE LOWER(name)=LOWER(?) LIMIT 1`

func SelectCategoryIDByName(db *sql.DB, name string) (int, error) {
	var id int
	err := db.QueryRow(selectCategoryIDByNameQ, name).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	return id, err
}
