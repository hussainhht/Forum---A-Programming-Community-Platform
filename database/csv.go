package database

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// ImportUsersCSV reads users.csv and UPSERTs into users table by unique email.
func ImportUsersCSV(db *sql.DB, csvPath string) error {
	f, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true

	// header
	if _, err := r.Read(); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt := `
	INSERT INTO users (id, username, email, password, created_at)
	VALUES (?, ?, ?, ?, ?)
	ON CONFLICT(email) DO UPDATE SET
		username=excluded.username,
		password=excluded.password,
		created_at=excluded.created_at
	`
	ps, err := tx.Prepare(stmt)
	if err != nil {
		return err
	}
	defer ps.Close()

	for {
		rec, e := r.Read()
		if e == io.EOF {
			break
		}
		if e != nil {
			return e
		}
		if len(rec) < 5 {
			return errors.New("invalid CSV row: need 5 columns (id,username,email,password,created_at)")
		}

		// Trim cells
		for i := range rec {
			rec[i] = strings.TrimSpace(rec[i])
		}

		if _, err := ps.Exec(rec[0], rec[1], rec[2], rec[3], strings.Trim(rec[4], `"`)); err != nil {
			return err
		}
	}

	return tx.Commit()
}



func ExportUsersCSV(db *sql.DB, csvPath string) error {
    rows, err := db.Query(`SELECT id, username, email, password, created_at FROM users ORDER BY id ASC`)
    if err != nil { return err }
    defer rows.Close()

    tmp := csvPath + ".tmp"
    f, err := os.Create(tmp)
    if err != nil { return err }
    defer f.Close()

    w := csv.NewWriter(f)
    defer w.Flush()
    _ = w.Write([]string{"id","username","email","password","created_at"})

    for rows.Next() {
        var id, user, email, pass, created string
        if err := rows.Scan(&id, &user, &email, &pass, &created); err != nil { return err }
        if err := w.Write([]string{id, user, email, pass, created}); err != nil { return err }
    }
    if err := rows.Err(); err != nil { return err }
    return os.Rename(tmp, csvPath) 
}


func AutoSyncCSVToDB(csvPath string, interval time.Duration) {
	go func() {
		var lastMod time.Time
		var lastSize int64

		for {
			fi, err := os.Stat(csvPath)
			if err == nil {
				mod := fi.ModTime()
				sz := fi.Size()

				if mod.After(lastMod) || sz != lastSize {
					time.Sleep(300 * time.Millisecond)
					if err := ImportUsersCSV(DB, csvPath); err != nil {
						log.Println("[csv-sync] import error:", err)
					} else {
						lastMod = mod
						lastSize = sz
						log.Println("[csv-sync] re-imported users.csv")
					}
				}
			}
			time.Sleep(interval)
		}
	}()
}
