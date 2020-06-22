package utilities

import (
	"context"
	"log"
	"math"

	"gorm.io/gorm"
)

// Param 分页参数
type Param struct {
	DB      *gorm.DB
	Page    int
	Limit   int
	OrderBy []string
	ShowSQL bool
}

// Paginator 分页返回
type Paginator struct {
	TotalRecord int         `json:"total_record"`
	TotalPage   int         `json:"total_page"`
	Records     interface{} `json:"records"`
	//	Offset      int         `json:"offset"`
	//	Limit       int         `json:"limit"`
	//	Page        int         `json:"page"`
	//	PrevPage    int         `json:"prev_page"`
	//	NextPage    int         `json:"next_page"`
}

func CalculatePage(limit int, totalRecord int) int {
	return int(math.Ceil(float64(totalRecord) / float64(limit)))
}

// Paging 分页
func Paging(ctx context.Context, p *Param, result interface{}) *Paginator {
	db := p.DB.WithContext(ctx)

	//	if p.ShowSQL {
	db = db.Debug()
	//	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}

	var offset int
	//	done := make(chan bool, 1)
	var paginator Paginator
	var count int64
	//	var offset int
	//	db.SkipDefaultTransaction = true
	//  db = db.Order("id DESC")
	//  go countRecords(db, result, done, &count)
	db.Model(result).Count(&count)

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}
	//offset = 0
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
		}
	}
	db.Limit(p.Limit).Offset(offset).Find(result)

	paginator.TotalRecord = int(count)
	paginator.Records = result
	//	paginator.Page = p.Page

	//	paginator.Offset = offset
	//	paginator.Limit = p.Limit
	//	paginator.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))
	paginator.TotalPage = CalculatePage(p.Limit, int(count))
	log.Println("len res")
	log.Println(paginator.TotalPage)

	/*	if p.Page > 1 {
			paginator.PrevPage = p.Page - 1
		} else {
			paginator.PrevPage = p.Page
		}

		if p.Page == paginator.TotalPage {
			paginator.NextPage = p.Page
		} else {
			paginator.NextPage = p.Page + 1
		}
	*/return &paginator
}

func countRecords(db *gorm.DB, anyType interface{}, done chan bool, count *int64) {
	db.Model(anyType).Count(count)
	done <- true
}
