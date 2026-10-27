package grpcserver

import (
	carpb "car-dealership/car/gen"
	"car-dealership/car/internal/application"
	"car-dealership/car/internal/application/model"
	"context"
	"net"
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type mockRepository struct {
	upsertCar func(context.Context, model.Car) error
	listCars  func(context.Context, *model.CarFilter, *model.Pagination) ([]model.Car, error)
}

func (m *mockRepository) UpsertCar(ctx context.Context, car model.Car) error {
	return m.upsertCar(ctx, car)
}

func (m *mockRepository) ListCars(ctx context.Context, filter *model.CarFilter, pagination *model.Pagination) ([]model.Car, error) {
	return m.listCars(ctx, filter, pagination)
}

func TestUpsertCar(t *testing.T) {
	ctx := context.Background()
	carID := uuid.NewV4()

	var storedCar []model.Car
	repository := &mockRepository{
		upsertCar: func(ctx context.Context, car model.Car) error {
			storedCar = append(storedCar, car)
			return nil
		},
		listCars: func(ctx context.Context, filter *model.CarFilter, pagination *model.Pagination) ([]model.Car, error) {
			return storedCar, nil
		},
	}

	app := application.NewHandler(repository)
	mockService := NewServer(app)

	expectedCar := model.Car{
		ID:    carID,
		Brand: "Toyota",
		Model: "Corolla",
		Year:  2022,
	}

	response, err := mockService.UpsertCar(ctx, &carpb.UpsertCarRequest{
		Id:    carID.String(),
		Brand: "Toyota",
		Model: "Corolla",
		Year:  2022,
	})

	require.NoError(t, err)
	require.NotNil(t, response)
	require.NotNil(t, response.Car)
	assert.Equal(t, expectedCar.ID.String(), response.Car.Id)
	assert.Equal(t, expectedCar.Brand, response.Car.Brand)
	assert.Equal(t, expectedCar.Model, response.Car.Model)
	assert.Equal(t, int32(expectedCar.Year), response.Car.Year)

}

func TestUpsertCar_InvalidIDRejected(t *testing.T) {
	ctx := context.Background()

	var storedCar []model.Car
	repository := &mockRepository{
		upsertCar: func(ctx context.Context, car model.Car) error {
			assert.Fail(t, "UpsertCar should not be called with invalid ID")
			return nil
		},
		listCars: func(ctx context.Context, filter *model.CarFilter, pagination *model.Pagination) ([]model.Car, error) {
			assert.Fail(t, "ListCars should not be called with invalid ID")
			return storedCar, nil
		},
	}

	app := application.NewHandler(repository)
	mockService := NewServer(app)
	client := newTestClient(t, mockService)

	// Attempt to upsert a car with an invalid ID using the new test client
	response, err := client.UpsertCar(ctx, &carpb.UpsertCarRequest{
		Id:    "invalid-uuid",
		Brand: "Toyota",
		Model: "Corolla",
		Year:  2022,
	})

	require.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestUpsertCar_ValidIDPropagates(t *testing.T) {
	ctx := context.Background()
	carID := uuid.NewV4().String()

	var storedCar []model.Car
	repository := &mockRepository{
		upsertCar: func(ctx context.Context, car model.Car) error {
			storedCar = append(storedCar, car)
			assert.Equal(t, carID, car.ID.String())
			return nil
		},
		listCars: func(ctx context.Context, filter *model.CarFilter, pagination *model.Pagination) ([]model.Car, error) {
			return storedCar, nil
		},
	}

	app := application.NewHandler(repository)
	mockService := NewServer(app)
	client := newTestClient(t, mockService)

	response, err := client.UpsertCar(ctx, &carpb.UpsertCarRequest{
		Id:    carID,
		Brand: "Toyota",
		Model: "Corolla",
		Year:  2022,
	})
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "Toyota", response.Car.Brand)
	assert.Equal(t, "Corolla", response.Car.Model)
	assert.Equal(t, int32(2022), response.Car.Year)

}

func TestListCars_Filtering(t *testing.T) {
	ctx := context.Background()
	searchText := "BMW"
	searchYear := int32(2022)

	var storedCar []model.Car
	repository := &mockRepository{
		upsertCar: func(ctx context.Context, car model.Car) error {
			storedCar = append(storedCar, car)
			return nil
		},
		listCars: func(ctx context.Context, filter *model.CarFilter, pagination *model.Pagination) ([]model.Car, error) {
			assert.NotNil(t, filter)
			if filter.Search != nil {
				assert.NotEmpty(t, *filter.Search)
				assert.Equal(t, searchText, *filter.Search)
			}
			if filter.Year != nil {
				assert.Greater(t, *filter.Year, 0)
				assert.Equal(t, int(searchYear), *filter.Year)
			}
			return storedCar, nil
		},
	}

	app := application.NewHandler(repository)
	mockService := NewServer(app)
	client := newTestClient(t, mockService)

	response, err := client.ListCars(ctx, &carpb.ListCarsRequest{
		Search: &searchText,
		Year:   &searchYear,
	})
	require.NoError(t, err)
	assert.NotNil(t, response)
}

func TestListCars_Pagination(t *testing.T) {
	ctx := context.Background()
	var storedCar []model.Car
	repository := &mockRepository{
		upsertCar: func(ctx context.Context, car model.Car) error {
			storedCar = append(storedCar, car)
			return nil
		},
		listCars: func(ctx context.Context, filter *model.CarFilter, pagination *model.Pagination) ([]model.Car, error) {
			assert.NotNil(t, pagination)
			return storedCar, nil
		},
	}

	app := application.NewHandler(repository)
	mockService := NewServer(app)
	client := newTestClient(t, mockService)

	cars := []model.Car{
		{
			ID:    uuid.NewV4(),
			Brand: "Toyota",
			Model: "Corolla",
			Year:  2022,
		},
		{
			ID:    uuid.NewV4(),
			Brand: "BMW",
			Model: "X5",
			Year:  2022,
		},
		{
			ID:    uuid.NewV4(),
			Brand: "Audi",
			Model: "A4",
			Year:  2022,
		},
	}
	for _, car := range cars {
		client.UpsertCar(ctx, &carpb.UpsertCarRequest{
			Id:    car.ID.String(),
			Brand: car.Brand,
			Model: car.Model,
			Year:  int32(car.Year),
		})
	}

	response, err := client.ListCars(ctx, &carpb.ListCarsRequest{
		PageSize:  2,
		PageToken: "",
	})
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 2, len(response.Cars))
	assert.NotEmpty(t, response.NextPageToken)

	response, err = client.ListCars(ctx, &carpb.ListCarsRequest{
		PageSize:  3,
		PageToken: "",
	})
	require.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 3, len(response.Cars))
	assert.Empty(t, response.NextPageToken)
}

// Test helper to create a new gRPC client connected to the test server using an in-memory listener
// It sets up an in-memory gRPC server and returns a client connected to it.
func newTestClient(t *testing.T, handler carpb.CarServiceServer) carpb.CarServiceClient {
	t.Helper() // Mark this function as a test helper

	listener := bufconn.Listen(1024 * 1024) // Create an in-memory listener for testing
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			LoggingInterceptor,
			InvalidUUIDInterceptor,
		),
	)
	carpb.RegisterCarServiceServer(grpcServer, handler)

	serverErrorChan := make(chan error, 1)
	// Start the gRPC server in a separate goroutine so it doesn't block the test execution
	go func(c chan error) {
		c <- grpcServer.Serve(listener)
	}(serverErrorChan)

	// Ensure the gRPC server is stopped and the listener is closed after the test finishes
	t.Cleanup(func() {
		grpcServer.Stop()
		_ = listener.Close()

		if err := <-serverErrorChan; err != nil {
			t.Fatalf("Failed to start gRPC server: %v", err)
		}
	})

	conn, err := grpc.NewClient(
		"passthrough:///bufconn",
		grpc.WithContextDialer(
			func(ctx context.Context, s string) (net.Conn, error) {
				return listener.Dial()
			},
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = conn.Close()
	})

	return carpb.NewCarServiceClient(conn)

}
