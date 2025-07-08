package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ming-0x0/hexago/internal/shared/dbmocker"
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

// DummyAdapter is a concrete implementation of DummyAdapterInterface.
type DummyAdapter struct{}

func (a DummyAdapter) ToDomain(entity *DummyEntity) (*DummyDomain, error) {
	return &DummyDomain{ID: entity.ID, Name: entity.Name}, nil
}

func (a DummyAdapter) ToEntity(domain *DummyDomain) (*DummyEntity, error) {
	return &DummyEntity{ID: domain.ID, Name: domain.Name}, nil
}

func (a DummyAdapter) ToDomains(entities []*DummyEntity) ([]*DummyDomain, error) {
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

func (a DummyAdapter) ToEntities(domains []*DummyDomain) ([]*DummyEntity, error) {
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

func setupTest(t *testing.T) (*Repository[DummyAdapter, DummyDomain, DummyEntity], *dbmocker.MockedRepository, *gorm.DB, sqlmock.Sqlmock, DummyAdapter, *logrus.Logger) {
	logger := logrus.New()
	mockedDB, err := dbmocker.NewMockedDB()
	if err != nil {
		t.Fatalf("error when creating mock DB: %v", err)
	}
	_, gormDB, sqlMock := mockedDB.DB, mockedDB.GormDB, mockedDB.SqlMock
	adapter := DummyAdapter{}
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
		name      string
		args      args
		setupMock func(sqlmock.Sqlmock)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{domain: &DummyDomain{ID: 1, Name: "Test"}},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `dummy_entities` \\(`name`,`id`\\) VALUES \\(\\?,\\?\\)").WithArgs("Test", 1).WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "Failure_DBError",
			args: args{domain: &DummyDomain{ID: 1, Name: "Test"}},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `dummy_entities` \\(`name`,`id`\\) VALUES \\(\\?,\\?\\)").WithArgs("Test", 1).WillReturnError(gorm.ErrInvalidData)
				mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo, mockedDB, _, sqlMock, _, _ := setupTest(t)
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
		name      string
		args      args
		setupMock func(sqlmock.Sqlmock)
		assertion assert.ErrorAssertionFunc
		expected  []*DummyDomain
	}{
		{
			name: "Success",
			args: args{conditions: map[string]any{"name": "Test"}},
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
			name: "Failure_DBError",
			args: args{conditions: map[string]any{"name": "Test"}},
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
			repo, mockedDB, _, sqlMock, _, _ := setupTest(t)
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
		name      string
		args      args
		setupMock func(sqlmock.Sqlmock)
		assertion assert.ErrorAssertionFunc
		expected  *DummyDomain
	}{
		{
			name: "Success",
			args: args{conditions: map[string]any{"id": 1}},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Test")
				mock.ExpectQuery("SELECT (.+) FROM `dummy_entities`").WillReturnRows(rows)
			},
			assertion: assert.NoError,
			expected:  &DummyDomain{ID: 1, Name: "Test"},
		},
		{
			name: "Failure_NotFound",
			args: args{conditions: map[string]any{"id": 1}},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT (.+) FROM `dummy_entities`").WillReturnError(gorm.ErrRecordNotFound)
			},
			assertion: assert.Error,
			expected:  nil,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo, mockedDB, _, sqlMock, _, _ := setupTest(t)
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
		name      string
		args      args
		setupMock func(sqlmock.Sqlmock)
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "Success",
			args: args{domain: &DummyDomain{ID: 1, Name: "Updated"}},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dummy_entities`").WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			assertion: assert.NoError,
		},
		{
			name: "Failure_DBError",
			args: args{domain: &DummyDomain{ID: 1, Name: "Updated"}},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `dummy_entities`").WillReturnError(gorm.ErrInvalidData)
				mock.ExpectRollback()
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			repo, mockedDB, _, sqlMock, _, _ := setupTest(t)
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
			repo, mockedDB, _, sqlMock, _, _ := setupTest(t)
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
		name      string
		args      args
		setupMock func(sqlmock.Sqlmock)
		assertion assert.ErrorAssertionFunc
		expected  []*DummyDomain
		count     int64
	}{
		{
			name: "Success",
			args: args{
				pagination: map[string]int{"page": 1, "limit": 2},
				conditions: map[string]any{"name": "Test"},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
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
			count: 2,
		},
		{
			name: "Failure_CountError",
			args: args{
				pagination: map[string]int{"page": 1, "limit": 2},
				conditions: map[string]any{"name": "Test"},
			},
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
			setupMock: func(mock sqlmock.Sqlmock) {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(2)
				mock.ExpectQuery("SELECT count\\(\\*\\) FROM `dummy_entities`").WillReturnRows(countRows)
				mock.ExpectQuery("SELECT (.+) FROM `dummy_entities`").WillReturnError(gorm.ErrInvalidField)
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
			repo, mockedDB, _, sqlMock, _, _ := setupTest(t)
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
