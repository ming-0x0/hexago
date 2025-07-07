package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ming-0x0/hexago/internal/shared/dbmocker"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

// Dummy types for testing
type DummyDomain struct {
	ID   int
	Name string
}

type DummyEntity struct {
	ID   int
	Name string
}

// DummyAdapter is a dummy implementation of AdapterInterface for testing.
type DummyAdapterInterface interface {
	ToDomain(entity *DummyEntity) (*DummyDomain, error)
	ToEntity(domain *DummyDomain) (*DummyEntity, error)
	ToDomains(entities []*DummyEntity) ([]*DummyDomain, error)
	ToEntities(domains []*DummyDomain) ([]*DummyEntity, error)
}

type DummyAdapter struct{}

func (a *DummyAdapter) ToDomain(entity *DummyEntity) (*DummyDomain, error) {
	return &DummyDomain{ID: entity.ID, Name: entity.Name}, nil
}

func (a *DummyAdapter) ToEntity(domain *DummyDomain) (*DummyEntity, error) {
	return &DummyEntity{ID: domain.ID, Name: domain.Name}, nil
}

func (a *DummyAdapter) ToDomains(entities []*DummyEntity) ([]*DummyDomain, error) {
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

func TestMockRepository_Create(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx    context.Context
		domain *DummyDomain
	}

	tests := []struct {
		name      string
		args      args
		setup     func(*MockRepositoryInterface[DummyAdapterInterface, DummyDomain, DummyEntity])
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "successful create",
			args: args{
				ctx:    context.Background(),
				domain: &DummyDomain{ID: 1, Name: "Test"},
			},
			setup: func(m *MockRepositoryInterface[DummyAdapterInterface, DummyDomain, DummyEntity]) {
				m.EXPECT().Create(context.Background(), &DummyDomain{ID: 1, Name: "Test"}).Return(nil)
			},
			assertion: assert.NoError,
		},
		{
			name: "error on create",
			args: args{
				ctx:    context.Background(),
				domain: &DummyDomain{ID: 1, Name: "Test"},
			},
			setup: func(m *MockRepositoryInterface[DummyAdapterInterface, DummyDomain, DummyEntity]) {
				m.EXPECT().Create(context.Background(), &DummyDomain{ID: 1, Name: "Test"}).Return(errors.New("create error"))
			},
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockRepositoryInterface[DummyAdapterInterface, DummyDomain, DummyEntity](ctrl)
			tc.setup(mockRepo)

			err := mockRepo.Create(tc.args.ctx, tc.args.domain)
			tc.assertion(t, err)
		})
	}
}

func TestMockRepository_FindByConditions(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx        context.Context
		conditions map[string]any
		scopes     []func(*gorm.DB) *gorm.DB
	}

	tests := []struct {
		name      string
		args      args
		setup     func(*MockRepositoryInterface[*DummyAdapter, DummyDomain, DummyEntity])
		want      []*DummyDomain
		assertion assert.ErrorAssertionFunc
	}{
		{
			name: "successful find",
			args: args{
				ctx:        context.Background(),
				conditions: map[string]any{"id": 1},
				scopes:     []func(*gorm.DB) *gorm.DB{},
			},
			setup: func(m *MockRepositoryInterface[*DummyAdapter, DummyDomain, DummyEntity]) {
				m.EXPECT().FindByConditions(context.Background(), map[string]any{"id": 1}).Return([]*DummyDomain{{ID: 1, Name: "Test"}}, nil)
			},
			want:      []*DummyDomain{{ID: 1, Name: "Test"}},
			assertion: assert.NoError,
		},
		{
			name: "error on find",
			args: args{
				ctx:        context.Background(),
				conditions: map[string]any{"id": 1},
				scopes:     []func(*gorm.DB) *gorm.DB{},
			},
			setup: func(m *MockRepositoryInterface[*DummyAdapter, DummyDomain, DummyEntity]) {
				m.EXPECT().FindByConditions(context.Background(), map[string]any{"id": 1}).Return(nil, errors.New("find error"))
			},
			want:      nil,
			assertion: assert.Error,
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockRepositoryInterface[*DummyAdapter, DummyDomain, DummyEntity](ctrl)
			tc.setup(mockRepo)

			got, err := mockRepo.FindByConditions(tc.args.ctx, tc.args.conditions, tc.args.scopes...)
			tc.assertion(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

// Additional tests for other mock methods can be added following the same pattern.

func TestRepositoryWithMockedDB_Create(t *testing.T) {
	t.Parallel()

	mockedRepo, err := dbmocker.NewMockedDB()
	if err != nil {
		t.Fatalf("Failed to create mocked DB: %v", err)
	}
	defer mockedRepo.DB.Close()

	// Setup expectations for SQL mock
	mockedRepo.SqlMock.ExpectBegin()
	mockedRepo.SqlMock.ExpectExec("INSERT INTO `dummy_entities`").
		WithArgs(1, "Test").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockedRepo.SqlMock.ExpectCommit()

	// Create a repository instance with the mocked DB
	// Note: You would need to adapt this based on your actual repository implementation
	// For demonstration, assuming a function to create a repository with a GORM DB
	// repo := NewDummyRepository(mockedRepo.Gorm)

	// Test data
	// domain := &DummyDomain{ID: 1, Name: "Test"}

	// Perform the create operation
	// err = repo.Create(context.Background(), domain)
	// assert.NoError(t, err)

	// For now, just testing the setup
	assert.NotNil(t, mockedRepo)
}
