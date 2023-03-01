package pagination

import (
	"fmt"
	"math"
	"net/http"

	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/validator"
)

type Pagination struct {
	Page      int
	PageSize  int
	Validator *validator.Validator
}

type Metadata struct {
	CurrentPage  int `json:"current_page"`
	PageSize     int `json:"page_size"`
	FirstPage    int `json:"first_page"`
	LastPage     int `json:"last_page"`
	TotalRecords int `json:"total_records"`
}

func New(r *http.Request) Pagination {
	v := validator.New()

	return Pagination{
		Page:      helpers.ReadIntParam(r.URL.Query(), "page", 1, v),
		PageSize:  helpers.ReadIntParam(r.URL.Query(), "page_size", 20, v),
		Validator: v,
	}
}

func (p *Pagination) Limit() int32 {
	return int32(p.PageSize)
}

func (p *Pagination) Offset() int32 {
	return int32((p.Page - 1) * p.PageSize)
}

func (p *Pagination) CalculateMetadata(totalRecords int64) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	metadata := Metadata{
		CurrentPage:  p.Page,
		PageSize:     p.PageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(p.PageSize))),
		TotalRecords: int(totalRecords),
	}

	if p.Page > metadata.LastPage {
		msg := fmt.Sprintf("must be equal or lower than the last page value of %d", metadata.LastPage)
		p.Validator.AddError("page", msg)

		return Metadata{}
	}

	return metadata
}

func (p *Pagination) Validate() {
	p.Validator.Check(p.Page > 0, "page", "must be greater than zero")
	p.Validator.Check(p.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	p.Validator.Check(p.PageSize > 0, "page_size", "size must be greater than zero")
	p.Validator.Check(p.PageSize <= 100, "page_size", "size must be a maximum of 100")
}
