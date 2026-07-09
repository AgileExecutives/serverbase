package dbping

import (
	"database/sql"

	"github.com/AgileExecutives/serverbase/pkg/interfaces"
	"gorm.io/gorm"
)

type gormDBPinger struct{ db *gorm.DB }

func NewGormDBPinger(db *gorm.DB) interfaces.DBPinger {
	if db == nil {
		return nil
	}
	return &gormDBPinger{db: db}
}

func (p *gormDBPinger) Ping() error {
	if p == nil || p.db == nil {
		return sql.ErrConnDone
	}
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
