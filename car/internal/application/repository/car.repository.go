package repository

import (
	"car-dealership/car/internal/application/model"
	"context"
	"errors"
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
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (id) DO UPDATE SET 
		brand = EXCLUDED.brand, 
		model = EXCLUDED.model, 
		year = EXCLUDED.year`)
	q, _, err := query.ToPgsql()
	if err != nil {
		return fmt.Errorf("failed to convert query to pgsql: %w", err)
	}
	_, err = r.connectionPool.Exec(ctx, q, car.ID, car.Brand, car.Model, car.Year)
	if err != nil {
		return fmt.Errorf("failed to upsert car: %w", err)
	}
	return nil
}

func (r *PostgresCarRepository) ListCars(ctx context.Context) ([]model.Car, error) {
	fmt.Println("LIST CARS REPOSITORY FUNC IS REACHED")
	return nil, errors.New("Not Implemented")
}
