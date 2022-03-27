package events_postgres_repository

import (
	"testing"
)

func TestIpIngo(t *testing.T) {
	db, mock := MockDb()
	defer db.Close()
	//st := NewIpInfoManager(logging.GetLogger("events", "repository"), db)

	mockRows := mock.NewRows([]string{"o_code", "value", "code"}).AddRow(1, 2, 3)
	mock.ExpectQuery("SELECT").WillReturnRows(mockRows)
}
