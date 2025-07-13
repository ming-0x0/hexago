package banner

import (
	"github.com/ming-0x0/hexago/internal/customer/adapter/repository/mysql/entity"
	"github.com/ming-0x0/hexago/internal/customer/domain/banner"
	"github.com/ming-0x0/hexago/internal/shared/repository"
)

type BannerRepositoryInterface interface {
	repository.RepositoryInterface[BannerRepositoryAdapterInterface, banner.Banner, entity.Banner]
}
