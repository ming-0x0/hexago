package repository

import (
	"context"
	"errors"

	sharedErrors "github.com/ming-0x0/hexago/shared/errors"
	"github.com/ming-0x0/hexago/shared/transaction"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 20
)

type AdapterInterface[D, E any] interface {
	// ToDomain converts the entity to its domain representation.
	ToDomain(*E) (*D, error)
	// ToEntity converts the domain to its entity representation.
	ToEntity(*D) (*E, error)
	// ToDomains converts the entities to their domain representations.
	ToDomains([]*E) ([]*D, error)
	// ToEntities converts the domains to their entity representations.
	ToEntities([]*D) ([]*E, error)
}

//go:generate go tool mockgen -destination mock/repository.go -package mock github.com/ming-0x0/hexago/shared/repository RepositoryInterface
type RepositoryInterface[A AdapterInterface[D, E], D, E any] interface {
	Create(
		ctx context.Context,
		domain *D,
	) error
	FindByConditions(
		ctx context.Context,
		conditions map[string]any,
		scopes ...func(*gorm.DB) *gorm.DB,
	) ([]*D, error)
	TakeByConditions(
		ctx context.Context,
		conditions map[string]any,
		scopes ...func(*gorm.DB) *gorm.DB,
	) (*D, error)
	Save(
		ctx context.Context,
		domain *D,
	) error
	DeleteByConditions(
		ctx context.Context,
		conditions map[string]any,
	) error
	FindByConditionsWithPagination(
		ctx context.Context,
		pageData map[string]int,
		conditions map[string]any,
		scopes ...func(*gorm.DB) *gorm.DB,
	) ([]*D, int64, error)
}

type Repository[A AdapterInterface[D, E], D, E any] struct {
	db      *gorm.DB
	logger  *logrus.Logger
	adapter A
}

func NewRepository[A AdapterInterface[D, E], D, E any](
	db *gorm.DB,
	logger *logrus.Logger,
	adapter A,
) *Repository[A, D, E] {
	return &Repository[A, D, E]{
		db:      db,
		logger:  logger,
		adapter: adapter,
	}
}

func (r *Repository[A, D, E]) DB(ctx context.Context) *gorm.DB {
	if tx, ok := transaction.TransactionFromContext(ctx); ok {
		return tx
	}

	return r.db.WithContext(ctx)
}

func (r *Repository[A, D, E]) pagination(pageData map[string]int) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page := DefaultPage
		if valPage, ok := pageData["page"]; ok && valPage > 0 {
			page = valPage
		}

		pageSize := DefaultPerPage
		if valPageSize, ok := pageData["limit"]; ok && valPageSize > 0 {
			pageSize = valPageSize
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func (r *Repository[A, D, E]) Create(
	ctx context.Context,
	domain *D,
) error {
	entity, err := r.adapter.ToEntity(domain)
	if err != nil {
		return err
	}

	err = r.DB(ctx).Create(entity).Error
	if err != nil {
		return sharedErrors.NewDomainError(sharedErrors.System, err.Error())
	}

	return nil
}

func (r *Repository[A, D, E]) FindByConditions(
	ctx context.Context,
	conditions map[string]any,
	scopes ...func(*gorm.DB) *gorm.DB,
) ([]*D, error) {
	var entities []*E
	if err := r.DB(ctx).Scopes(scopes...).Where(conditions).Find(&entities).Error; err != nil {
		return nil, sharedErrors.NewDomainError(sharedErrors.System, err.Error())
	}

	return r.adapter.ToDomains(entities)
}

func (r *Repository[A, D, E]) TakeByConditions(
	ctx context.Context,
	conditions map[string]any,
	scopes ...func(*gorm.DB) *gorm.DB,
) (*D, error) {
	entity := new(E)
	err := r.DB(ctx).Scopes(scopes...).Where(conditions).Take(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedErrors.NewDomainError(sharedErrors.NotFound, err.Error())
		}

		return nil, sharedErrors.NewDomainError(sharedErrors.System, err.Error())
	}

	return r.adapter.ToDomain(entity)
}

func (r *Repository[A, D, E]) Save(
	ctx context.Context,
	domain *D,
) error {
	entity, err := r.adapter.ToEntity(domain)
	if err != nil {
		return err
	}

	err = r.DB(ctx).Save(entity).Error
	if err != nil {
		return sharedErrors.NewDomainError(sharedErrors.System, err.Error())
	}

	return nil
}

func (r *Repository[A, D, E]) DeleteByConditions(
	ctx context.Context,
	conditions map[string]any,
) error {
	entity := new(E)
	err := r.DB(ctx).Where(conditions).Delete(entity).Error
	if err != nil {
		return sharedErrors.NewDomainError(sharedErrors.System, err.Error())
	}

	return nil
}

func (r *Repository[A, D, E]) FindByConditionsWithPagination(
	ctx context.Context,
	pageData map[string]int,
	conditions map[string]any,
	scopes ...func(*gorm.DB) *gorm.DB,
) ([]*D, int64, error) {
	cdb := r.DB(ctx)

	var entities []*E
	var count int64

	countBuilder := cdb.Model(&entities)
	queryBuilder := cdb.Scopes(r.pagination(pageData))

	err := countBuilder.Scopes(scopes...).Where(conditions).Count(&count).Error
	if err != nil {
		return []*D{}, 0, sharedErrors.NewDomainError(sharedErrors.System, err.Error())
	}

	err = queryBuilder.Scopes(scopes...).Where(conditions).Find(&entities).Error
	if err != nil {
		return []*D{}, 0, sharedErrors.NewDomainError(sharedErrors.System, err.Error())
	}

	domains, err := r.adapter.ToDomains(entities)
	if err != nil {
		return []*D{}, 0, err
	}

	return domains, count, nil
}
