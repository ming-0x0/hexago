package transaction

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ming-0x0/hexago/shared/dbmocker"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (*Transaction, *dbmocker.MockedRepository, *gorm.DB, sqlmock.Sqlmock) {
	mockedDB, err := dbmocker.NewMockedDB()
	if err != nil {
		t.Fatalf("error when creating mock DB: %v", err)
	}
	_, gormDB, sqlMock := mockedDB.DB, mockedDB.GormDB, mockedDB.SqlMock

	tx := NewTransaction(gormDB)
	return tx, mockedDB, gormDB, sqlMock
}

func teardownTest(mockedDB *dbmocker.MockedRepository) {
	if mockedDB != nil {
		mockedDB.DB.Close()
	}
}

func TestTransaction_Do(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupMock func(sqlmock.Sqlmock)
		fn        func(ctx context.Context) error
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectCommit()
			},
			fn: func(ctx context.Context) error {
				return nil
			},
			assertion: assert.NoError,
		},
		{
			name: "Failure_TransactionError",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectRollback()
			},
			fn: func(ctx context.Context) error {
				return gorm.ErrInvalidTransaction
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tx, mockedDB, _, sqlMock := setupTest(t)
			defer teardownTest(mockedDB)

			tc.setupMock(sqlMock)
			err := tx.Do(context.Background(), tc.fn)
			tc.assertion(t, err)
		})
	}
}

func TestTransactionFromContext(t *testing.T) {
	t.Parallel()

	t.Run("TransactionInContext", func(t *testing.T) {
		t.Parallel()
		_, mockedDB, gormDB, _ := setupTest(t)
		defer teardownTest(mockedDB)

		// Create a mock transaction
		tx := gormDB.Begin()
		ctx := context.WithValue(context.Background(), Tx, tx)

		// Get transaction from context
		retrievedTx, ok := TransactionFromContext(ctx)

		// Verify that the transaction was retrieved correctly
		assert.True(t, ok)
		assert.Equal(t, tx, retrievedTx)
	})

	t.Run("NoTransactionInContext", func(t *testing.T) {
		t.Parallel()
		_, mockedDB, _, _ := setupTest(t)
		defer teardownTest(mockedDB)

		// Use a context without a transaction
		ctx := context.Background()

		// Get transaction from context
		retrievedTx, ok := TransactionFromContext(ctx)

		// Verify that no transaction was retrieved
		assert.False(t, ok)
		assert.Nil(t, retrievedTx)
	})
}
