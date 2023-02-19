package pagination

import (
	"math"
	"net/http"

	"github.com/seanflannery10/core/internal/helpers"
	"github.com/seanflannery10/core/internal/validator"
)

func New(r *http.Request, v *validator.Validator) Filters {
	return Filters{
		Page:     helpers.ReadIntParam(r.URL.Query(), "page", 1, v),
		PageSize: helpers.ReadIntParam(r.URL.Query(), "page_size", 20, v),
	}
}

type Filters struct {
	Page     int
	PageSize int
}

func (f Filters) Limit() int32 {
	return int32(f.PageSize)
}

func (f Filters) Offset() int32 {
	return int32((f.Page - 1) * f.PageSize)
}

func ValidateFilters(v *validator.Validator, f Filters) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "size must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "size must be a maximum of 100")
}

type Metadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	FirstPage    int `json:"first_page"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
}

func CalculateMetadata(totalRecords int64, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: int(totalRecords),
	}
}
