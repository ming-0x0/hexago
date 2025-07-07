package dbmocker

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MockedRepository struct {
	Logger  *logrus.Logger
	DB      *sql.DB
	GormDB  *gorm.DB
	SqlMock sqlmock.Sqlmock
}

func NewMockedDB() (*MockedRepository, error) {
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}

	sqlMock.MatchExpectationsInOrder(false)

	gormConfig := &gorm.Config{}

	gormDB, err := gorm.Open(
		mysql.New(mysql.Config{
			Conn:                      db,
			SkipInitializeWithVersion: true,
		}),
		gormConfig,
	)
	if err != nil {
		return nil, err
	}

	return &MockedRepository{
		Logger:  logger,
		DB:      db,
		GormDB:  gormDB,
		SqlMock: sqlMock,
	}, nil
}
