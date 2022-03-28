package mysql

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/sevskii111/microservices-goods/pkg/models"
)

type GoodModel struct {
	DB *sql.DB
}

func (m *GoodModel) Insert(name, description string, price int, inStock bool) (int, error) {
	stmt := `INSERT INTO goods (name, description, price, in_stock, created)
	VALUES(?, ?, ?, ?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, name, description, price, inStock)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *GoodModel) Update(id int, update *models.GoodUpdate) error {
	SQLUpdate := SQLUpdates{}
	if update.Name != nil {
		SQLUpdate.add("name", *update.Name)
	}
	if update.Description != nil {
		SQLUpdate.add("description", *update.Description)
	}
	if update.Price != nil {
		SQLUpdate.add("price", *update.Price)
	}
	if update.InStock != nil {
		SQLUpdate.add("in_stock", *update.InStock)
	}

	stmt := fmt.Sprintf(`UPDATE goods SET %s WHERE id = %d`, SQLUpdate.Assignments(), id)
	_, err := m.DB.Exec(stmt, SQLUpdate.Values()...)
	if err != nil {
		return err
	}

	return nil
}

func (m *GoodModel) Get(id int) (*models.Good, error) {
	stmt := `SELECT name, description, price, in_stock FROM goods
	WHERE id = ?`
	row := m.DB.QueryRow(stmt, id)
	good := &models.Good{}
	err := row.Scan(&good.Name, &good.Description, &good.Price, &good.InStock)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return good, nil
}

func (m *GoodModel) List() ([]*models.Good, error) {
	stmt := `SELECT id, name, description, price, in_stock, created FROM goods`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	goods := []*models.Good{}

	for rows.Next() {
		g := &models.Good{}
		err = rows.Scan(&g.ID, &g.Name, &g.Description, &g.Price, &g.InStock, &g.Created)
		if err != nil {
			return nil, err
		}
		goods = append(goods, g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return goods, nil
}
