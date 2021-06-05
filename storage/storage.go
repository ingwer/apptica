package storage

import (
	"apptica/model"

	"context"
	"database/sql"
)

type Storage interface {
	GetTopPositions(ctx context.Context, appID int, countryCode int, date model.Date) (map[model.Category]model.Position, error)
	SaveTopPositions(ctx context.Context, appID int, countryCode int, categoryID model.Category, positions map[model.Date]model.Position) error
}

type storage struct {
	db *sql.DB
}

type Option func(s *storage)

func New(options ...Option) Storage {
	s := &storage{}
	for _, option := range options {
		option(s)
	}

	return s
}

func WithDB(db *sql.DB) Option {
	return func(s *storage) {
		s.db = db
	}
}

func (s storage) GetTopPositions(ctx context.Context, appID int, countryCode int, date model.Date) (map[model.Category]model.Position, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT `category_id`, `position` FROM `top_positions` where app_id=? and country_code = ? and `date` = ?",
		appID,
		countryCode,
		date,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
		_ = rows.Err()
	}()

	result := make(map[model.Category]model.Position)
	for rows.Next() {
		var catID string
		var position int

		if err := rows.Scan(&catID, &position); err != nil {
			return nil, err
		}

		result[model.Category(catID)] = model.Position(position)
	}

	return result, nil
}

func (s storage) SaveTopPositions(ctx context.Context, appID int, countryCode int, categoryID model.Category, positions map[model.Date]model.Position) error {
	ins := "INSERT INTO top_positions(`app_id`, `country_code`, `category_id`, `date`, `position`) VALUES "

	vals := []interface{}{}
	for date, position := range positions {
		ins += "(?,?,?,?,?),"
		vals = append(vals, appID, countryCode, categoryID, date, position)
	}

	ins = ins[0 : len(ins)-1]

	stmt, err := s.db.Prepare(ins)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, vals...)

	return err
}
