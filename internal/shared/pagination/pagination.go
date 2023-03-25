package pagination

import (
	"math"
)

const (
	defaultPage     = 1
	defaultPageSize = 20
	fistPage        = 1
	keyPage         = "page"
	keyPageSize     = "page_size"
	pageMax         = 10000000
	pageMin         = 0
	pageSizeMax     = 100
	pageSizeMin     = 5
	noRecords       = 0
)

type (
	Pagination struct {
		Page     int
		PageSize int
	}

	Metadata struct {
		CurrentPage  int `json:"current_page"`
		PageSize     int `json:"page_size"`
		FirstPage    int `json:"first_page"`
		LastPage     int `json:"last_page"`
		TotalRecords int `json:"total_records"`
	}
)

func New(page, pageSize int) Pagination {
	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p *Pagination) Limit() int32 {
	return int32(p.PageSize)
}

func (p *Pagination) Offset() int32 {
	return int32((p.Page - 1) * p.PageSize)
}

func (p *Pagination) CalculateMetadata(totalRecords int64) Metadata {
	if totalRecords == noRecords {
		return Metadata{}
	}

	metadata := Metadata{
		CurrentPage:  p.Page,
		PageSize:     p.PageSize,
		FirstPage:    fistPage,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(p.PageSize))),
		TotalRecords: int(totalRecords),
	}

	if p.Page > metadata.LastPage {
		// msg := fmt.Sprintf("must be equal or lower than the last page value of %d", metadata.LastPage)
		// p.Validator.AddError("page", msg)

		return Metadata{}
	}

	return metadata
}
