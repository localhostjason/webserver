package uv

import (
	"gorm.io/gorm"
	"strings"
	"time"
)

const LAST_PAGE = 0

type PagingIn struct {
	Page    uint `form:"page"`
	PerPage uint `form:"per_page"`
	NoCount bool `form:"no_count"`
}

func (pi *PagingIn) HasPage() bool {
	return pi.PerPage != 0
}

type Order struct {
	// 列名,只支持一个 TODO validation
	SortBy string `form:"sort_by"`
	// 顺序or倒叙 : "" or "desc"
	Order string `form:"order"`
}

func (s *Order) HasOrder() bool {
	return s.SortBy != ""
}

func (s *Order) ToExp() string {
	return strings.Join([]string{s.SortBy, s.Order}, " ")
}

type PagingOut struct {
	Page    uint `json:"page"`
	PerPage uint `json:"per_page"`
	Pages   uint `json:"pages"`
	Total   uint `json:"total"`
}

type PagingItems struct {
	Items interface{} `json:"items"`
	PagingOut
}

type Limit struct {
	Offset uint
	Limit  uint
}

func (l *Limit) ToPagingOut(total uint) *PagingOut {
	return &PagingOut{
		Page:    1 + (l.Offset / l.Limit),
		PerPage: l.Limit,
		Pages:   (total + l.Limit - 1) / l.Limit,
		Total:   total,
	}
}

func (pi *PagingIn) ToLimit(total uint) *Limit {
	var limitIn *Limit
	if pi.PerPage != 0 {
		if pi.Page != LAST_PAGE {
			limitIn = &Limit{
				Limit:  pi.PerPage,
				Offset: pi.PerPage * (pi.Page - 1),
			}
		} else {
			limitIn = &Limit{
				Limit:  pi.PerPage,
				Offset: total / pi.PerPage * pi.PerPage,
			}
		}
	}
	return limitIn
}

// PagingFind 这里不能用pdb
// TODO return error, 不要panic
func PagingFind(db *gorm.DB, o interface{}, pi *PagingIn, order *Order) (po *PagingOut, err error) {
	if order.HasOrder() {
		db = db.Order(order.ToExp())
	}

	// 不用分页
	if pi == nil || !pi.HasPage() {
		err = db.Find(o).Error
		return
	}

	var q = db
	var total int64
	if pi.NoCount && pi.Page != LAST_PAGE {
		total = 0
	} else {
		err = db.Count(&total).Error
		if err != nil {
			return
		}
	}

	var li = pi.ToLimit(uint(total))
	db = q.Limit(int(li.Limit)).Offset(int(li.Offset))
	err = db.Find(o).Error
	if err != nil {
		return
	}
	po = li.ToPagingOut(uint(total))
	return
}

func PagingScan(db *gorm.DB, o interface{}, pi *PagingIn, order *Order) *PagingOut {
	if order.HasOrder() {
		db = db.Order(order.ToExp())
	}

	if pi == nil || !pi.HasPage() {
		db.Find(o)
		return nil
	}
	var q = db
	var total int64

	if pi.NoCount && pi.Page != LAST_PAGE {
		total = 0
	} else {
		PDB(db.Count(&total))
	}

	var li = pi.ToLimit(uint(total))
	db = q.Limit(int(li.Limit)).Offset(int(li.Offset))
	db.Scan(o)
	return li.ToPagingOut(uint(total))
}

type PagedItemList struct {
	Items interface{} `json:"items"`
	PagingOut
}

func PagedOut(o interface{}, po *PagingOut) interface{} {
	if po == nil {
		return o
	} else {
		return &PagedItemList{
			Items:     o,
			PagingOut: *po,
		}
	}
}

func QueryTime(db *gorm.DB, expr string, t time.Time) *gorm.DB {
	if t.IsZero() {
		return db
	}
	return db.Where(expr, t)
}

// PDB panic if db error
func PDB(db *gorm.DB) {
	var err = db.Error
	if err == nil {
		return
	}
	PEIf(E_DB_ERROR, err)
}
