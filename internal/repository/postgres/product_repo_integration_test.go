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

	// Clean up before tests
	testDB.Exec("DELETE FROM products")

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
