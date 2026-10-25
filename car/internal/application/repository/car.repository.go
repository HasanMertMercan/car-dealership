package repository

import (
	"car-dealership/car/internal/application/model"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nullism/bqb"
)

type PostgresCarRepository struct {
	connectionPool *pgxpool.Pool
}

func NewCarRepository(connectionPool *pgxpool.Pool) *PostgresCarRepository {
	return &PostgresCarRepository{
		connectionPool: connectionPool,
	}
}

func (r *PostgresCarRepository) UpsertCar(ctx context.Context, car model.Car) error {
	query := bqb.New(`
		INSERT INTO car (id, brand, model, year) 
		VALUES (?, ?, ?, ?) 
		ON CONFLICT (id) DO UPDATE SET 
		brand = EXCLUDED.brand, 
		model = EXCLUDED.model, 
		year = EXCLUDED.year`, car.ID, car.Brand, car.Model, car.Year)
	q, params, err := query.ToPgsql()
	if err != nil {
		return fmt.Errorf("failed to convert query to pgsql: %w", err)
	}
	_, err = r.connectionPool.Exec(ctx, q, params...)
	if err != nil {
		return fmt.Errorf("failed to upsert car: %w", err)
	}
	return nil
}

func (r *PostgresCarRepository) ListCars(ctx context.Context, filter *model.CarFilter, pagination *model.Pagination) ([]model.Car, error) {
	where := bqb.Optional("WHERE")
	if filter != nil {
		if filter.Search != nil {
			where = where.And("(brand ILIKE ? OR model ILIKE ?)", "%"+*filter.Search+"%", "%"+*filter.Search+"%")
		}
		if filter.Year != nil {
			where = where.And("year = ?", *filter.Year)
		}
	}

	where = where.And("id > ?", pagination.PageToken)
	orderBy := bqb.Optional("ORDER BY").Comma("id " + string(model.SortDirectionAsc)) // Default sort by ID ascending
	limit := bqb.Optional("LIMIT").And("?", pagination.PageSize)

	query := bqb.New(`
		SELECT id, brand, model, year
		FROM car
		?
		?
		?
	`, where, orderBy, limit)

	queryString, params, err := query.ToPgsql()
	if err != nil {
		return nil, fmt.Errorf("failed to convert query to pgsql: %w", err)
	}

	rows, err := r.connectionPool.Query(ctx, queryString, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	cars := make([]model.Car, 0)
	for rows.Next() {
		var car model.Car
		err := rows.Scan(&car.ID, &car.Brand, &car.Model, &car.Year)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		cars = append(cars, car)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return cars, nil
}
