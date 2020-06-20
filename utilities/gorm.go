package utilities

import (
	"context"
	"math"

	"github.com/jinzhu/gorm"
)

//CreateWithContext is the same with gorm Create but respecting context, return nil if timed out
func CreateWithContext(ctx context.Context, value interface{}, g *gorm.DB) (*gorm.DB, error) {
	resChan := make(chan *gorm.DB)
	go func(resChan chan<- *gorm.DB, value interface{}) {
		resChan <- g.Create(value)
	}(resChan, value)
	return waitContextOrGORMError(ctx, resChan)
}

//SaveWithContext is the same with gorm Save but respecting context, return nil if timed out
func SaveWithContext(ctx context.Context, value interface{}, g *gorm.DB) (*gorm.DB, error) {
	resChan := make(chan *gorm.DB)
	go func(resChan chan<- *gorm.DB, value interface{}) {
		resChan <- g.Save(value)
	}(resChan, value)
	return waitContextOrGORMError(ctx, resChan)
}

//CountWithContext is the same with gorm count but respecting context, return nil if timed out
func CountWithContext(ctx context.Context, g *gorm.DB, value interface{}) (*gorm.DB, error) {
	resChan := make(chan *gorm.DB)
	go func(resChan chan<- *gorm.DB, value interface{}) {
		resChan <- g.Count(value)

		//todo: handle send to closed channel
		//cant do : https://medium.com/@rocketlaunchr.cloud/canceling-mysql-in-go-827ed8f83b30
		//on mysql side it will be still executing, kill the process is undesired
	}(resChan, value)
	return waitContextOrGORMError(ctx, resChan)
}

//FirstWithContext is the same with gorm First but respecting context, return nil if timed out
func FirstWithContext(ctx context.Context, g *gorm.DB, out interface{}, where ...interface{}) (*gorm.DB, error) {
	resChan := make(chan *gorm.DB)

	go func(resChan chan<- *gorm.DB, out interface{}, where ...interface{}) {
		resChan <- g.First(out, where) //how to close ongoing stuff ? so it wont send to closed chan

	}(resChan, out, where)
	return waitContextOrGORMError(ctx, resChan)
}

//FindWithContext will find gorm and respecting context, return nil if context timed out
func FindWithContext(ctx context.Context, g *gorm.DB, out interface{}, where ...interface{}) (*gorm.DB, error) {
	//encapsulate with context timeout, this is not the cleanest way to do this, as resource maybe dangling inside gorm process
	resChan := make(chan *gorm.DB) //finger crossed to GC

	go func(resChan chan<- *gorm.DB, out interface{}, where ...interface{}) {
		resChan <- g.Find(out, where)
	}(resChan, out, where)
	return waitContextOrGORMError(ctx, resChan)
}

//ScanWithContext will Scan gorm and respecting context, return nil if context timed out
func ScanWithContext(ctx context.Context, g *gorm.DB, dest interface{}) (*gorm.DB, error) {
	//encapsulate with context timeout, this is not the cleanest way to do this, as resource maybe dangling inside gorm process
	resChan := make(chan *gorm.DB)
	go func(resChan chan<- *gorm.DB, dest interface{}) {
		resChan <- g.Scan(dest)
	}(resChan, dest)
	return waitContextOrGORMError(ctx, resChan)
}

func waitContextOrGORMError(ctx context.Context, resChan chan *gorm.DB) (*gorm.DB, error) {
	select {
	case <-ctx.Done():
		//timeout happen
		return nil, ctx.Err()
	case res := <-resChan:
		//returned ok but error in sql itself
		if res.Error != nil {
			return res, res.Error
		}
		//good things happen
		return res, nil
	}
}

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
	Offset      int         `json:"offset"`
	Limit       int         `json:"limit"`
	Page        int         `json:"page"`
	PrevPage    int         `json:"prev_page"`
	NextPage    int         `json:"next_page"`
}

// Paging
func Paging(ctx context.Context, p *Param, result interface{}) *Paginator {
	db := p.DB

	if p.ShowSQL {
		db = db.Debug()
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(o)
		}
	}

	done := make(chan bool, 1)
	var paginator Paginator
	var count int
	var offset int

	go countRecords(ctx, db, result, done, &count)

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}

	//	db.Limit(p.Limit).Offset(offset).Find(result)
	FindWithContext(ctx, db.Limit(p.Limit).Offset(offset), result)
	<-done

	paginator.TotalRecord = count
	paginator.Records = result
	paginator.Page = p.Page

	paginator.Offset = offset
	paginator.Limit = p.Limit
	paginator.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	} else {
		paginator.PrevPage = p.Page
	}

	if p.Page == paginator.TotalPage {
		paginator.NextPage = p.Page
	} else {
		paginator.NextPage = p.Page + 1
	}
	return &paginator
}

func countRecords(ctx context.Context, db *gorm.DB, anyType interface{}, done chan bool, count *int) {
	CountWithContext(ctx, db.Model(anyType), count)
	done <- true
}
