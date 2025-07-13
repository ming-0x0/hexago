package banner

import (
	"github.com/ming-0x0/hexago/internal/customer/adapter/repository/mysql/entity"
	"github.com/ming-0x0/hexago/internal/customer/domain/banner"
	bannerPort "github.com/ming-0x0/hexago/internal/customer/port/repository/mysql/banner"
	"github.com/ming-0x0/hexago/internal/shared/repository"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type BannerRepository struct {
	*repository.Repository[bannerPort.BannerRepositoryAdapterInterface, banner.Banner, entity.Banner]
}

func New(
	db *gorm.DB,
	logger *logrus.Logger,
	adapter bannerPort.BannerRepositoryAdapterInterface,
) *BannerRepository {
	return &BannerRepository{
		Repository: repository.NewRepository(db, logger, adapter),
	}
}
