package domain

import (
	"github.com/mpppk/imagine/domain/model"
	"github.com/mpppk/imagine/util"
)

// NewNumbersFromStringSlice create new Numbers with numbers from string slice
func NewNumbersFromStringSlice(strNumbers []string) (model.Numbers, error) {
	rawNumbers, err := util.ConvertStringSliceToIntSlice(strNumbers)
	if err != nil {
		return nil, err
	}
	return model.NewNumbers(rawNumbers), nil
}


