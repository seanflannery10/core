package pagination

import (
	"math"
	"net/http"

	"github.com/seanflannery10/core/internal/helpers"
	"github.com/seanflannery10/core/internal/validator"
)

type Pagination struct {
	Page     int
	PageSize int
}

type Metadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	FirstPage    int `json:"first_page"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
}

func New(r *http.Request, v *validator.Validator) Pagination {
	return Pagination{
		Page:     helpers.ReadIntParam(r.URL.Query(), "page", 1, v),
		PageSize: helpers.ReadIntParam(r.URL.Query(), "page_size", 20, v),
	}
}

func (f *Pagination) Limit() int32 {
	return int32(f.PageSize)
}

func (f *Pagination) Offset() int32 {
	return int32((f.Page - 1) * f.PageSize)
}

func (f *Pagination) CalculateMetadata(totalRecords int64) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  f.Page,
		PageSize:     f.PageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(f.PageSize))),
		TotalRecords: int(totalRecords),
	}
}

func ValidatePagination(v *validator.Validator, p Pagination) {
	v.Check(p.Page > 0, "page", "must be greater than zero")
	v.Check(p.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(p.PageSize > 0, "page_size", "size must be greater than zero")
	v.Check(p.PageSize <= 100, "page_size", "size must be a maximum of 100")
}
