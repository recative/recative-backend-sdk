package mock

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
)

func Mock(db *gorm.DB, force bool, values ...interface{}) error {
	var onConflict = func(isForce bool) clause.OnConflict {
		if isForce {
			return clause.OnConflict{UpdateAll: true}
		} else {
			return clause.OnConflict{DoNothing: true}
		}
	}(force)

	for _, value := range values {
		reflectValue := reflect.Indirect(reflect.ValueOf(value))
		switch reflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < reflectValue.Len(); i++ {
				err := db.Select("*").Clauses(onConflict).Create(toStructPtr(reflectValue.Index(i).Interface())).Error
				if err != nil {
					return err
				}
			}
		default:
			if reflectValue.IsZero() {
				continue
			}
			err := db.Select("*").Clauses(onConflict).Create(toStructPtr(value)).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// https://groups.google.com/g/golang-nuts/c/KB3_Yj3Ny4c
func toStructPtr(obj interface{}) interface{} {
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		return obj
	}
	// Create a new instance of the underlying type
	vp := reflect.New(reflect.TypeOf(obj))

	vp.Elem().Set(reflect.ValueOf(obj))

	// NOTE: `vp.Elem().Set(reflect.ValueOf(&obj).Elem())` does not work
	return vp.Interface()
}
