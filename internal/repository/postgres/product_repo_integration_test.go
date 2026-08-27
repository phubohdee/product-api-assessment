//go:build integration

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"product-api-assessment/internal/domain"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	_ = godotenv.Load("../../../.env")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	var err error
	testDB, err = sql.Open("postgres", dsn)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer testDB.Close()

	code := m.Run()
	os.Exit(code)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func TestProductRepository_Create_Success(t *testing.T) {
	repo := NewProductRepository(testDB)
	ctx := context.Background()

	desc := "Integration test product"
	product := &domain.Product{
		Name:        "Test Product",
		Description: &desc,
		Price:       199.99,
	}

	result, err := repo.Create(ctx, product)

	require.NoError(t, err)
	assert.NotZero(t, result.ID)
	assert.Equal(t, "Test Product", result.Name)
	assert.Equal(t, &desc, result.Description)
	assert.Equal(t, 199.99, result.Price)
	assert.Nil(t, result.SalePrice)
	assert.NotZero(t, result.CreatedAt)
	assert.NotZero(t, result.UpdatedAt)

	// Cleanup
	testDB.Exec("DELETE FROM products WHERE id = $1", result.ID)
}

func TestProductRepository_Create_WithSalePrice(t *testing.T) {
	repo := NewProductRepository(testDB)
	ctx := context.Background()

	salePrice := 149.99
	product := &domain.Product{
		Name:      "Sale Product",
		SalePrice: &salePrice,
		Price:     199.99,
	}

	result, err := repo.Create(ctx, product)

	require.NoError(t, err)
	assert.NotZero(t, result.ID)
	assert.Equal(t, "Sale Product", result.Name)
	assert.NotNil(t, result.SalePrice)
	assert.Equal(t, 149.99, *result.SalePrice)

	// Cleanup
	testDB.Exec("DELETE FROM products WHERE id = $1", result.ID)
}

func TestProductRepository_Create_WithNullDescription(t *testing.T) {
	repo := NewProductRepository(testDB)
	ctx := context.Background()

	product := &domain.Product{
		Name:  "No Desc Product",
		Price: 50.00,
	}

	result, err := repo.Create(ctx, product)

	require.NoError(t, err)
	assert.NotZero(t, result.ID)
	assert.Equal(t, "No Desc Product", result.Name)
	assert.Nil(t, result.Description)
	assert.Nil(t, result.SalePrice)
	assert.Equal(t, 50.00, result.Price)

	// Cleanup
	testDB.Exec("DELETE FROM products WHERE id = $1", result.ID)
}

func TestProductRepository_Create_AutoIncrementID(t *testing.T) {
	repo := NewProductRepository(testDB)
	ctx := context.Background()

	product1 := &domain.Product{Name: "Product 1", Price: 10.00}
	product2 := &domain.Product{Name: "Product 2", Price: 20.00}

	result1, err := repo.Create(ctx, product1)
	require.NoError(t, err)

	result2, err := repo.Create(ctx, product2)
	require.NoError(t, err)

	assert.Greater(t, result2.ID, result1.ID)

	// Cleanup
	testDB.Exec("DELETE FROM products WHERE id IN ($1, $2)", result1.ID, result2.ID)
}

func TestProductRepository_GetByID_Success(t *testing.T) {
	repo := NewProductRepository(testDB)
	ctx := context.Background()

	created, err := repo.Create(ctx, &domain.Product{Name: "GetByID Test", Price: 99.00})
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, created.ID)
	require.NoError(t, err)
	assert.Equal(t, created.ID, found.ID)
	assert.Equal(t, "GetByID Test", found.Name)

	testDB.Exec("DELETE FROM products WHERE id = $1", created.ID)
}

func TestProductRepository_GetByID_NotFound(t *testing.T) {
	repo := NewProductRepository(testDB)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, 999999)
	assert.ErrorIs(t, err, sql.ErrNoRows)
}

func TestProductRepository_Update_Success(t *testing.T) {
	repo := NewProductRepository(testDB)
	ctx := context.Background()

	created, err := repo.Create(ctx, &domain.Product{Name: "Old Name", Price: 100.00})
	require.NoError(t, err)

	created.Name = "Updated Name"
	created.Price = 150.00
	updated, err := repo.Update(ctx, created)
	require.NoError(t, err)

	assert.NotNil(t, updated)
	assert.Equal(t, "Updated Name", updated.Name)
	assert.Equal(t, 150.00, updated.Price)

	testDB.Exec("DELETE FROM products WHERE id = $1", created.ID)
}
