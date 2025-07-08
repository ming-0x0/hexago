package dbmocker

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
)

type MockedRepository struct {
	DB      *sql.DB
	GormDB  *gorm.DB
	SqlMock sqlmock.Sqlmock
	Logger  *logrus.Logger
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
	gormConfig.Logger = gormLog.New(logger, gormLog.Config{
		LogLevel: gormLog.Warn,
		Colorful: true,
	})

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
		DB:      db,
		GormDB:  gormDB,
		SqlMock: sqlMock,
		Logger:  logger,
	}, nil
}
