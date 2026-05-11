package pagination

import "gorm.io/gorm"

type Metadata struct {
	Total           int64 `json:"total"`
	HasNextPage     bool  `json:"has_next_page"`
	HasPreviousPage bool  `json:"has_previous_page"`
	Page            uint  `json:"page"`
	Limit           uint  `json:"limit"`
}

type Pagination[T any] struct {
	Data     []T      `json:"data"`
	Metadata Metadata `json:"metadata"`
}

func NewPagination[T any](queryable *gorm.DB, page uint, limit uint) *Pagination[T] {
	var data []T
	metadata := NewMetadata(queryable, page, limit)
	queryable = queryable.Limit(int(limit)).Offset(int((page - 1) * limit))
	queryable.Find(&data)

	if data == nil {
		data = []T{}
	}

	return &Pagination[T]{
		Data:     data,
		Metadata: metadata,
	}
}

func NewMetadata(queryable *gorm.DB, page uint, limit uint) Metadata {
	var total int64
	queryable.Count(&total)

	return Metadata{
		Total:           total,
		Page:            page,
		Limit:           limit,
		HasNextPage:     total > int64(page*limit),
		HasPreviousPage: page > 1,
	}
}
