package pagination

import (
	"fmt"
	"math"
	"net/http"

	"github.com/seanflannery10/core/internal/pkg/helpers"
	"github.com/seanflannery10/core/internal/pkg/validator"
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
	pageSizeMin     = 0
	noLength        = 0
)

type (
	Pagination struct {
		Validator *validator.Validator
		Page      int
		PageSize  int
	}

	Metadata struct {
		CurrentPage  int `json:"current_page"`
		PageSize     int `json:"page_size"`
		FirstPage    int `json:"first_page"`
		LastPage     int `json:"last_page"`
		TotalRecords int `json:"total_records"`
	}
)

func New(r *http.Request) Pagination {
	v := validator.New()

	return Pagination{
		Page:      helpers.ReadIntParam(r.URL.Query(), "page", defaultPage, v),
		PageSize:  helpers.ReadIntParam(r.URL.Query(), "page_size", defaultPageSize, v),
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
	if totalRecords == noLength {
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
		msg := fmt.Sprintf("must be equal or lower than the last page value of %d", metadata.LastPage)
		p.Validator.AddError("page", msg)

		return Metadata{}
	}

	return metadata
}

func (p *Pagination) Validate() {
	p.Validator.Check(p.Page > pageMin, keyPage, "must be greater than zero")
	p.Validator.Check(p.Page <= pageMax, keyPage, "must be a maximum of 10 million")
	p.Validator.Check(p.PageSize > pageSizeMin, keyPageSize, "size must be greater than zero")
	p.Validator.Check(p.PageSize <= pageSizeMax, keyPageSize, "size must be a maximum of 100")
}
