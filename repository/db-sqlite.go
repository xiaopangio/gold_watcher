package repository

import (
	"database/sql"
	"errors"
	"time"
)

type SQLiteRepository struct {
	Conn *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{Conn: db}
}

func (repo *SQLiteRepository) Migrate() error {
	query := `
	create table if not exists holdings(
		id integer primary key autoincrement,
		amount real not null,
		purchase_date integer not null,
		purchase_price integer not null)
	`
	_, err := repo.Conn.Exec(query)
	return err
}

func (repo *SQLiteRepository) InsertHolding(h Holdings) (*Holdings, error) {
	stmt := "insert into holdings (amount,purchase_date,purchase_price) values (?,?,?)"
	result, err := repo.Conn.Exec(stmt, h.Amount, h.PurchaseDate.Unix(), h.PurchasePrice)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	h.ID = id
	return &h, nil

}

func (repo *SQLiteRepository) AllHoldings() ([]Holdings, error) {
	query := "select id ,amount,purchase_date,purchase_price from holdings order by purchase_date"
	rows, err := repo.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var all []Holdings
	for rows.Next() {
		var h Holdings
		var unixTime int64
		err := rows.Scan(&h.ID, &h.Amount, &unixTime, &h.PurchasePrice)
		if err != nil {
			return nil, err
		}
		h.PurchaseDate = time.Unix(unixTime, 0)
		all = append(all, h)
	}
	//var all []Holdings
	//h := Holdings{
	//	ID:            1,
	//	Amount:        1,
	//	PurchaseDate:  time.Now(),
	//	PurchasePrice: 1000,
	//}
	//all = append(all, h)
	//h = Holdings{
	//	ID:            2,
	//	Amount:        2,
	//	PurchaseDate:  time.Now().Add(1 * time.Second),
	//	PurchasePrice: 1200,
	//}
	//all = append(all, h)
	return all, nil
}

func (repo *SQLiteRepository) GetHoldingByID(id int) (*Holdings, error) {
	query := "select id ,amount,purchase_date,purchase_price from holdings where id = ?"
	row := repo.Conn.QueryRow(query, id)

	var h Holdings
	var unixTime int64
	err := row.Scan(&h.ID, &h.Amount, &unixTime, &h.PurchasePrice)
	if err != nil {
		return nil, err
	}
	h.PurchaseDate = time.Unix(unixTime, 0)
	return &h, nil
}

func (repo *SQLiteRepository) UpdateHolding(id int64, update Holdings) error {
	if id == 0 {
		return errors.New("invalid update id")
	}
	stmt := "update holdings set amount= ?,purchase_date=?,purchase_price=? where id =?"
	res, err := repo.Conn.Exec(stmt, update.Amount, update.PurchaseDate.Unix(), update.PurchasePrice, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errUpdateFailed
	}
	return nil
}

func (repo *SQLiteRepository) DeleteHolding(id int64) error {
	res, err := repo.Conn.Exec("delete from holdings where id = ?", id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errDeleteFailed
	}
	return nil
}
