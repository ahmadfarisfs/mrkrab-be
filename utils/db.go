package utils

import (
	"log"

	"gorm.io/gorm"
)

func AppendCommonRequest(db *gorm.DB, req CommonRequest) *gorm.DB {

	for k, v := range req.Filter {
		if k == "created_at" {
			continue
		}
		if k == "start_date" {
			db = db.Where("created_at > ?", v)
			continue
		}
		if k == "end_date" {
			db = db.Where("created_at < ?", v)
			continue
		}
		switch v.(type) {
		case float64:
			db = db.Where(k+" = ?", v)
			break
		case []float64:
			db = db.Where(k+" IN (?)", v)
			break
		case int:
			db = db.Where(k+" = ?", v)
			break
		case []int:
			db = db.Where(k+" IN (?)", v)
			break
		case string:
			db = db.Where(k+" = ?", v)
			break
		case []string:
			db = db.Where(k+" IN (?)", v)
			break
		}
	}

	for k, v := range req.Search {
		switch v.(type) {
		case string:
			db = db.Where(k+" LIKE ?", "%"+v.(string)+"%")
		}
	}

	if req.SortBy != "" && req.SortType != "" {
		sortQuery := req.SortBy + " " + req.SortType
		log.Println("sortQuery")

		log.Println(sortQuery)
		db = db.Order(sortQuery)
	}
	var endIndex int
	if req.EndIndex != 0 {
		endIndex = req.EndIndex
	}
	//	if req.EndIndex != 0 {
	//[0,9]
	//[10,19] -> id ke 10 - 19
	db = db.Offset(req.StartIndex).Limit(endIndex + 1 - req.StartIndex)
	//	}
	return db
}
