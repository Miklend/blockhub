package block

import (
	clickhouseRepo "clickhouse-service/internal/db/click_house"
	"lib/models"
)

func (d *clickhouseRepo.ClickhouseRepo) FetchBlock(table string, hashBlock string) (models.Block, error) {

}
