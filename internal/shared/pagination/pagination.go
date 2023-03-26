package pagination

import (
	"math"

	"github.com/go-faster/errors"
	"github.com/seanflannery10/core/internal/oas"
)

const (
	fistPage   = 1
	noRecords  = 0
	pageOffset = 1
)

var ErrPageValueToHigh = errors.New("page must be equal or lower than the last page value")

type (
	Pagination struct {
		Page     int32
		PageSize int32
	}
)

func New(page, pageSize int32) Pagination {
	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

func (p *Pagination) Limit() int32 {
	return p.PageSize
}

func (p *Pagination) Offset() int32 {
	return (p.Page - pageOffset) * p.PageSize
}

func (p *Pagination) CalculateMetadata(totalRecords int64) (oas.MessagesMetadataResponse, error) {
	if totalRecords == noRecords {
		return oas.MessagesMetadataResponse{}, nil
	}

	metadata := oas.MessagesMetadataResponse{
		CurrentPage:  p.Page,
		PageSize:     p.PageSize,
		FirstPage:    fistPage,
		LastPage:     int32(math.Ceil(float64(totalRecords) / float64(p.PageSize))),
		TotalRecords: totalRecords,
	}

	if p.Page > metadata.LastPage {
		return oas.MessagesMetadataResponse{}, ErrPageValueToHigh
	}

	return metadata, nil
}
