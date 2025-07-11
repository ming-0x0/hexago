package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ming-0x0/hexago/internal/shared/dbmocker"
	"github.com/ming-0x0/hexago/internal/shared/errors"
	"github.com/ming-0x0/hexago/internal/shared/transaction"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type DummyEntity struct {
	ID   int
	Name string
}

type DummyDomain struct {
	ID   int
	Name string
}

// DummyAdapterInterface defines the methods for converting between entity and domain.
type DummyAdapterInterface interface {
	ToDomain(*DummyEntity) (*DummyDomain, error)
	ToEntity(*DummyDomain) (*DummyEntity, error)
	ToDomains([]*DummyEntity) ([]*DummyDomain, error)
	ToEntities([]*DummyDomain) ([]*DummyEntity, error)
}

// DummyAdapter is a concrete implementation of DummyAdapterInterface with configurable error injection for testing.
type DummyAdapter struct {
	// Flags to simulate errors in adapter methods
	ShouldFailToDomain   bool
	ShouldFailToEntity   bool
	ShouldFailToDomains  bool
	ShouldFailToEntities bool
}

func (a *DummyAdapter) ToDomain(entity *DummyEntity) (*DummyDomain, error) {
	if a.ShouldFailToDomain {
		return nil, errors.NewDomainError(errors.Validation, "failed to convert to domain")
	}
	return &DummyDomain{ID: entity.ID, Name: entity.Name}, nil
}

func (a *DummyAdapter) ToEntity(domain *DummyDomain) (*DummyEntity, error) {
	if a.ShouldFailToEntity {
		return nil, errors.NewDomainError(errors.Validation, "failed to convert to entity")
	}
	return &DummyEntity{ID: domain.ID, Name: domain.Name}, nil
}

func (a *DummyAdapter) ToDomains(entities []*DummyEntity) ([]*DummyDomain, error) {
	if a.ShouldFailToDomains {
		return nil, errors.NewDomainError(errors.Validation, "failed to convert to domains")
	}

	domains := make([]*DummyDomain, len(entities))
	for i, e := range entities {
		d, err := a.ToDomain(e)
		if err != nil {
			return nil, err
		}
		domains[i] = d
	}
	return domains, nil
}

func (a *DummyAdapter) ToEntities(domains []*DummyDomain) ([]*DummyEntity, error) {
	if a.ShouldFailToEntities {
		return nil, errors.NewDomainError(errors.Validation, "failed to convert to entities")
	}

	entities := make([]*DummyEntity, len(domains))
	for i, d := range domains {
		e, err := a.ToEntity(d)
		if err != nil {
			return nil, err
		}
		entities[i] = e
	}
	return entities, nil
}

func setupTest(t *testing.T, adapterConfig DummyAdapter) (*Repository[*DummyAdapter, DummyDomain, DummyEntity], *dbmocker.MockedRepository, *gorm.DB, sqlmock.Sqlmock, *DummyAdapter, *logrus.Logger) {
	logger := logrus.New()
	mockedDB, err := dbmocker.NewMockedDB()
	if err != nil {
		t.Fatalf("error when creating mock DB: %v", err)
	}
	_, gormDB, sqlMock := mockedDB.DB, mockedDB.GormDB, mockedDB.SqlMock

	adapter := &adapterConfig
	repo := NewRepository(gormDB, logger, adapter)
	return repo, mockedDB, gormDB, sqlMock, adapter, logger
}

func teardownTest(mockedDB *dbmocker.MockedRepository) {
	if mockedDB != nil {
		mockedDB.DB.Close()
	}
}

func TestRepository_Create(t *testing.T) {
	t.Parallel()

	type args struct {
		domain *DummyDomain
	}

	tests := []struct {
		name          string
		args          args
		adapterConfig DummyAdapter
		setupMock     func(sqlmock.Sqlmock)
		assertion     assert.ErrorAssertionFunc
	}{
		{
			name:          "Success",
			args:          args{domain: &DummyDomain{ID: 1, Name: "Test"}},
			adapterConfig: DummyAdapter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `dummy_entities`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name:          "Failure_DBError",
			args:          args{domain: &DummyDomain{ID: 1, Name: "Test"}},
			adapterConfig: DummyAdapter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `dummy_entities`").WillReturnError(gorm.ErrInvalidData)
				mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name:          "Failure_AdapterToEntityError",
			args:          args{domain: &DummyDomain{ID: 1, Name: "Test"}},
			adapterConfig: DummyAdapter{ShouldFailToEntity: true},
			setupMock:     func(mock sqlmock.Sqlmock) {},
			assertion:     assert.Error,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo, mockedDB, _, sqlMock, _, _ := setupTest(t, tc.adapterConfig)
			defer teardownTest(mockedDB)

			tc.setupMock(sqlMock)
			err := repo.Create(context.Background(), tc.args.domain)
			tc.assertion(t, err)
		})
	}
}

func TestRepository_FindByConditions(t *testing.T) {
	t.Parallel()

	type args struct {
		conditions map[string]any
	}

	tests := []struct {
		name          string
		args          args
		adapterConfig DummyAdapter
		setupMock     func(sqlmock.Sqlmock)
		assertion     assert.ErrorAssertionFunc
		expected      []*DummyDomain
	}{
		{
			name:          "Success",
			args:          args{conditions: map[string]any{"name": "Test"}},
			adapterConfig: DummyAdapter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Test1").
					AddRow(2, "Test2")
				mock.ExpectQuery("SELECT (.+) FROM `dummy_entities`").WillReturnRows(rows)
			},
			assertion: assert.NoError,
			expected:  []*DummyDomain{{ID: 1, Name: "Test1"}, {ID: 2, Name: "Test2"}},
		},
		{
			name:          "Failure_DBError",
			args:          args{conditions: map[string]any{"name": "Test"}},
			adapterConfig: DummyAdapter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `dummy_entities`").WillReturnError(gorm.ErrInvalidField)
			},
			assertion: assert.Error,
			expected:  nil,
		},
		{
			name:          "Failure_AdapterToDomainsError",
			args:          args{conditions: map[string]any{"name": "Test"}},
			adapterConfig: DummyAdapter{ShouldFailToDomains: true},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Test1").
					AddRow(2, "Test2")
				mock.ExpectQuery("SELECT (.+) FROM `dummy_entities`").WillReturnRows(rows)
			},
			assertion: assert.Error,
			expected:  nil,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo, mockedDB, _, sqlMock, _, _ := setupTest(t, tc.adapterConfig)
			defer teardownTest(mockedDB)

			tc.setupMock(sqlMock)
			data, err := repo.FindByConditions(context.Background(), tc.args.conditions)
			tc.assertion(t, err)
			if err == nil {
				assert.Equal(t, tc.expected, data)
			}
		})
	}
}

func TestRepository_TakeByConditions(t *testing.T) {
	t.Parallel()

	type args struct {
		conditions map[string]any
	}

	tests := []struct {
		name          string
		args          args
		adapterConfig DummyAdapter
		setupMock     func(sqlmock.Sqlmock)
		assertion     assert.ErrorAssertionFunc
		expected      *DummyDomain
	}{
		{
			name:          "Success",
			args:          args{conditions: map[string]any{"id": 1}},
			adapterConfig: DummyAdapter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Test")
				mock.ExpectQuery("SELECT (.+) FROM `dummy_entities`").WillReturnRows(rows)
			},
			assertion: assert.NoError,
			expected:  &DummyDomain{ID: 1, Name: "Test"},
		},
		{
			name:          "Failure_NotFound",
			args:          args{conditions: map[string]any{"id": 1}},
			adapterConfig: DummyAdapter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `dummy_entities`").WillReturnError(gorm.ErrRecordNotFound)
			},
			assertion: assert.Error,
			expected:  nil,
		},
		{
			name:          "Failure_AdapterToDomainError",
			args:          args{conditions: map[string]any{"id": 1}},
			adapterConfig: DummyAdapter{ShouldFailToDomain: true},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Test")
				mock.ExpectQuery("SELECT (.+) FROM `dummy_entities`").WillReturnRows(rows)
			},
			assertion: assert.Error,
			expected:  nil,
		},
		{
			name:          "Failure_DBError",
			args:          args{conditions: map[string]any{"id": 1}},
			adapterConfig: DummyAdapter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `dummy_entities`").WillReturnError(gorm.ErrInvalidField)
			},
			assertion: assert.Error,
			expected:  nil,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo, mockedDB, _, sqlMock, _, _ := setupTest(t, tc.adapterConfig)
			defer teardownTest(mockedDB)

			tc.setupMock(sqlMock)
			data, err := repo.TakeByConditions(context.Background(), tc.args.conditions)
			tc.assertion(t, err)
			if err == nil {
				assert.Equal(t, tc.expected, data)
			}
		})
	}
}

func TestRepository_Save(t *testing.T) {
	t.Parallel()

	type args struct {
		domain *DummyDomain
	}

	tests := []struct {
		name          string
		args          args
		adapterConfig DummyAdapter
		setupMock     func(sqlmock.Sqlmock)
		assertion     assert.ErrorAssertionFunc
	}{
		{
			name:          "Success",
			args:          args{domain: &DummyDomain{ID: 1, Name: "Updated"}},
			adapterConfig: DummyAdapter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dummy_entities`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name:          "Failure_DBError",
			args:          args{domain: &DummyDomain{ID: 1, Name: "Updated"}},
			adapterConfig: DummyAdapter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dummy_entities`").WillReturnError(gorm.ErrInvalidData)
				mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
		{
			name:          "Failure_AdapterToEntityError",
			args:          args{domain: &DummyDomain{ID: 1, Name: "Updated"}},
			adapterConfig: DummyAdapter{ShouldFailToEntity: true},
			setupMock:     func(mock sqlmock.Sqlmock) {},
			assertion:     assert.Error,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo, mockedDB, _, sqlMock, _, _ := setupTest(t, tc.adapterConfig)
			defer teardownTest(mockedDB)

			tc.setupMock(sqlMock)
			err := repo.Save(context.Background(), tc.args.domain)
			tc.assertion(t, err)
		})
	}
}

func TestRepository_DeleteByConditions(t *testing.T) {
	t.Parallel()

	type args struct {
		conditions map[string]any
	}

	tests := []struct {
		name      string
		args      args
		setupMock func(sqlmock.Sqlmock)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{conditions: map[string]any{"id": 1}},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `dummy_entities`").WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "Failure_DBError",
			args: args{conditions: map[string]any{"id": 1}},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `dummy_entities`").WillReturnError(gorm.ErrInvalidField)
				mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo, mockedDB, _, sqlMock, _, _ := setupTest(t, DummyAdapter{})
			defer teardownTest(mockedDB)

			tc.setupMock(sqlMock)
			err := repo.DeleteByConditions(context.Background(), tc.args.conditions)
			tc.assertion(t, err)
		})
	}
}

func TestRepository_FindByConditionsWithPagination(t *testing.T) {
	t.Parallel()

	type args struct {
		pagination map[string]int
		conditions map[string]any
	}

	tests := []struct {
		name          string
		args          args
		adapterConfig DummyAdapter
		setupMock     func(sqlmock.Sqlmock)
		assertion     assert.ErrorAssertionFunc
		expected      []*DummyDomain
		count         int64
	}{
		{
			name: "Success_DefaultPagination",
			args: args{
				pagination: map[string]int{}, // Empty map to test default pagination values
				conditions: map[string]any{"name": "Test"},
			},
			adapterConfig: DummyAdapter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(5)
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Test1").
					AddRow(2, "Test2")
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `dummy_entities`").WillReturnRows(countRows)
				mock.ExpectQuery("SELECT (.+) FROM `dummy_entities`").WillReturnRows(rows)
			},
			assertion: assert.NoError,
			expected: []*DummyDomain{
				{ID: 1, Name: "Test1"},
				{ID: 2, Name: "Test2"},
			},
			count: 5,
		},
		{
			name: "Success_CustomPagination",
			args: args{
				pagination: map[string]int{"page": 2, "limit": 3},
				conditions: map[string]any{"name": "Test"},
			},
			adapterConfig: DummyAdapter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(10)
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(4, "Test4").
					AddRow(5, "Test5").
					AddRow(6, "Test6")
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `dummy_entities`").WillReturnRows(countRows)
				mock.ExpectQuery("SELECT (.+) FROM `dummy_entities`").WillReturnRows(rows)
			},
			assertion: assert.NoError,
			expected: []*DummyDomain{
				{ID: 4, Name: "Test4"},
				{ID: 5, Name: "Test5"},
				{ID: 6, Name: "Test6"},
			},
			count: 10,
		},
		{
			name: "Success_ZeroPageAndLimit",
			args: args{
				pagination: map[string]int{"page": 0, "limit": 0}, // Test handling of zero or invalid values
				conditions: map[string]any{"name": "Test"},
			},
			adapterConfig: DummyAdapter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(5)
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Test1").
					AddRow(2, "Test2")
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `dummy_entities`").WillReturnRows(countRows)
				mock.ExpectQuery("SELECT (.+) FROM `dummy_entities`").WillReturnRows(rows)
			},
			assertion: assert.NoError,
			expected: []*DummyDomain{
				{ID: 1, Name: "Test1"},
				{ID: 2, Name: "Test2"},
			},
			count: 5,
		},
		{
			name: "Failure_CountError",
			args: args{
				pagination: map[string]int{"page": 1, "limit": 2},
				conditions: map[string]any{"name": "Test"},
			},
			adapterConfig: DummyAdapter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `dummy_entities`").WillReturnError(gorm.ErrInvalidField)
			},
			assertion: assert.Error,
			expected:  nil,
			count:     0,
		},
		{
			name: "Failure_QueryError",
			args: args{
				pagination: map[string]int{"page": 1, "limit": 2},
				conditions: map[string]any{"name": "Test"},
			},
			adapterConfig: DummyAdapter{},
			setupMock: func(mock sqlmock.Sqlmock) {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `dummy_entities`").WillReturnRows(countRows)
				mock.ExpectQuery("SELECT (.+) FROM `dummy_entities`").WillReturnError(gorm.ErrInvalidField)
			},
			assertion: assert.Error,
			expected:  nil,
			count:     0,
		},
		{
			name: "Failure_AdapterToDomainsError",
			args: args{
				pagination: map[string]int{"page": 1, "limit": 2},
				conditions: map[string]any{"name": "Test"},
			},
			adapterConfig: DummyAdapter{ShouldFailToDomains: true},
			setupMock: func(mock sqlmock.Sqlmock) {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "Test1").
					AddRow(2, "Test2")
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `dummy_entities`").WillReturnRows(countRows)
				mock.ExpectQuery("SELECT (.+) FROM `dummy_entities`").WillReturnRows(rows)
			},
			assertion: assert.Error,
			expected:  nil,
			count:     0,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo, mockedDB, _, sqlMock, _, _ := setupTest(t, tc.adapterConfig)
			defer teardownTest(mockedDB)

			tc.setupMock(sqlMock)
			data, count, err := repo.FindByConditionsWithPagination(context.Background(), tc.args.pagination, tc.args.conditions)
			tc.assertion(t, err)
			if err == nil {
				assert.Equal(t, tc.expected, data)
				assert.Equal(t, tc.count, count)
			}
		})
	}
}

func TestRepository_DB_WithTransaction(t *testing.T) {
	t.Parallel()

	t.Run("TransactionInContext", func(t *testing.T) {
		t.Parallel()
		repo, mockedDB, gormDB, _, _, _ := setupTest(t, DummyAdapter{})
		defer teardownTest(mockedDB)

		// Create a mock transaction
		tx := gormDB.Begin()
		ctx := context.WithValue(context.Background(), transaction.Tx, tx)

		// Get DB from context with transaction
		db := repo.DB(ctx)

		// Verify that the returned DB is the transaction
		assert.Equal(t, tx, db)
	})

	t.Run("NoTransactionInContext", func(t *testing.T) {
		t.Parallel()
		repo, mockedDB, gormDB, _, _, _ := setupTest(t, DummyAdapter{})
		defer teardownTest(mockedDB)

		// Use a context without a transaction
		ctx := context.Background()

		// Get DB from context without transaction
		db := repo.DB(ctx)

		// Verify that the returned DB is the default DB with context
		assert.NotEqual(t, gormDB, db) // It's a new instance with context
		assert.Equal(t, ctx, db.Statement.Context)
	})
}
