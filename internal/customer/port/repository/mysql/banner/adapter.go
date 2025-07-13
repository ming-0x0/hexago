package banner

import (
	"github.com/ming-0x0/hexago/internal/customer/adapter/repository/mysql/entity"
	"github.com/ming-0x0/hexago/internal/customer/domain/banner"
)

type BannerRepositoryAdapterInterface interface {
	ToDomain(e *entity.Banner) (*banner.Banner, error)
	ToEntity(d *banner.Banner) (*entity.Banner, error)
	ToDomains(es []*entity.Banner) ([]*banner.Banner, error)
	ToEntities(ds []*banner.Banner) ([]*entity.Banner, error)
}
