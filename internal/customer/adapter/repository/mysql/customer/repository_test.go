package customer

import (
	"testing"

	"github.com/ming-0x0/hexago/internal/shared/dbmocker"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	mockedDB, err := dbmocker.NewMockedDB()
	if err != nil {
		t.Fatalf("error when creating mock DB: %v", err)
	}
	_, gormDB, _ := mockedDB.DB, mockedDB.GormDB, mockedDB.SqlMock

	adapter := &CustomerRepositoryAdapter{}
	repo := New(gormDB, logrus.New(), adapter)
	assert.NotNil(t, repo)
}
