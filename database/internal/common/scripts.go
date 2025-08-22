package common

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"results/errs"
)

func CheckRecord[T any](db *gorm.DB, modelName string, value string, model *T) error {
	if err := db.Where(fmt.Sprintf("%s = ?", modelName), value).First(model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.RecordNotFound
		}
		return err
	}
	return nil
}
