package common

import (
	"math"

	modelCommon "orchid-starter/internal/common/model"

	"github.com/yudhiana/logos"
)

func CalculatePagination(total int64, limit int, page int) (pagination modelCommon.PaginationResult) {

	totalPage := math.Ceil(float64(total) / float64(limit))
	hasNext := false

	if float64(page) < totalPage {
		hasNext = true
	}

	totalPageInt, errConvert := ConvertFloat64ToInt(totalPage)
	if errConvert != nil {
		logos.NewLogger().Error("Failed calculate pagination at convert totalPage to int", "error", errConvert)
		return
	}

	pagination = modelCommon.PaginationResult{
		Page:        page,
		TotalPage:   totalPageInt,
		TotalItems:  total,
		PerPage:     limit,
		HasNext:     hasNext,
		HasPrevious: page > 1,
	}
	return
}
