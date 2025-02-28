package stat

import (
	"time"
	"url-shortener/pkg/db"

	"gorm.io/datatypes"
)

type StatRepository struct {
	*db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{
		Db: db,
	}
}

func (repo *StatRepository) AddClick(linkID uint) {
	// Если за сегодня нет статистики - создаем
	// если есть - увеличиваем на +1
	currentDate := datatypes.Date(time.Now())

	var stat Stat
	repo.Db.DB.Find(&stat, "link_id = ? and date = ?", linkID, currentDate)

	if stat.ID == 0 {
		repo.Db.DB.Create(&Stat{
			LinkID: linkID,
			Clicks: 1,
			Date:   currentDate,
		})
	}

	stat.Clicks += 1
	repo.Db.DB.Save(&stat)
}

func (repo *StatRepository) GetStats(by string, from time.Time, to time.Time) []GetStatResponse {
	var stats []GetStatResponse
	var selectQuery string

	switch by {
	case GroupByDay:
		selectQuery = "to_char(date, 'YYYY-MM-DD') as period, sum(clicks)"
	case GroupByMonth:
		selectQuery = "to_char(date, 'YYYY-MM') as period, sum(clicks)"
	}

	repo.Db.DB.
		Table("stats").
		Select(selectQuery).
		Where("date BETWEEN ? AND ?", from, to).
		Group("period").
		Group("period").
		Scan(&stats)

	return stats
}
