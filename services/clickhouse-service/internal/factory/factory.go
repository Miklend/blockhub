package factory

import (
	dbpkg "clickhouse-service/internal/db"
	clickhouseRepo "clickhouse-service/internal/db/click_house"
	clientsDB "lib/clients/db"
	"lib/utils/logging"
)

const (
	DBClickHouse = "clickhouse"
)

// New creates concrete db implementations by type
func New(dbType string, chClient clientsDB.ClickhouseClient, logger *logging.Logger) dbpkg.DB {
	switch dbType {
	case DBClickHouse:
		return clickhouseRepo.NewClickhouseService(chClient, logger)
	default:
		return clickhouseRepo.NewClickhouseService(chClient, logger)
	}
}

