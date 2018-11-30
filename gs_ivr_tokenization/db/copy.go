package db

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"
)

// copyTo  copies the values from an object to another object
func copyTo(toValue interface{}, fromValue interface{}, alias interface{}) (err error) {
	var (
		isSlice bool
		amount  = 1
		from    = indirect(reflect.ValueOf(fromValue))
		to      = indirect(reflect.ValueOf(toValue))
	)

	if !to.CanAddr() {
		return errors.New("copy to value is unaddressable")
	}

	// Return is from value is invalid
	if !from.IsValid() {
		return
	}

	// Just set it if possible to assign
	if from.Type().AssignableTo(to.Type()) {
		to.Set(from)
		return
	}

	fromType := indirectType(from.Type())
	toType := indirectType(to.Type())

	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		return
	}

	if to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			amount = from.Len()
		}
	}

	for i := 0; i < amount; i++ {
		var dest, source reflect.Value

		if isSlice {
			// source
			if from.Kind() == reflect.Slice {
				source = indirect(from.Index(i))
			} else {
				source = indirect(from)
			}

			// dest
			dest = indirect(reflect.New(toType).Elem())
		} else {
			source = indirect(from)
			dest = indirect(to)
		}

		var fieldsTo = deepFields(toType)
		var fieldsFrom = deepFields(fromType)

		var aliasFrom = map[string]string{}
		var _aliasFrom = map[string]string{}
		if ta, ok := alias.(string); ok && ta != "" {
			for _, field := range fieldsFrom {
				name := field.Name
				tag := field.Tag
				alias := tag.Get(ta)
				if alias != "" && alias != "-" {
					aliasFrom[alias] = name
					_aliasFrom[name] = alias
				}
			}
		}

		var aliasTo = map[string]string{}
		var _aliasTo = map[string]string{}
		if ta, ok := alias.(string); ok && ta != "" {
			for _, field := range fieldsTo {
				name := field.Name
				tag := field.Tag
				alias := tag.Get(ta)
				if alias != "" && alias != "-" {
					aliasTo[alias] = name
					_aliasTo[name] = alias
				}
			}
		}

		// Copy from field to field or method
		for _, field := range fieldsFrom {
			name := field.Name

			if n, ok := _aliasFrom[name]; ok {
				name = n
			}

			if fromField := source.FieldByName(name); fromField.IsValid() {
				// has field
				nameTo := name
				if n, ok := aliasTo[name]; ok {
					nameTo = n
				}

				if toField := dest.FieldByName(nameTo); toField.IsValid() {
					if toField.CanSet() {
						if !set(toField, fromField) {
							if err := copyTo(toField.Addr().Interface(), fromField.Interface(), alias); err != nil {
								return err
							}
						}
					}
				} else {
					// try to set to method
					var toMethod reflect.Value
					if dest.CanAddr() {
						toMethod = dest.Addr().MethodByName(nameTo)
					} else {
						toMethod = dest.MethodByName(nameTo)
					}

					if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
						toMethod.Call([]reflect.Value{fromField})
					}
				}
			}
		}

		// Copy from method to field
		for _, field := range fieldsTo {
			name := field.Name

			if n, ok := _aliasTo[name]; ok {
				name = n
			}

			if n, ok := aliasFrom[name]; ok {
				name = n
			}

			var fromMethod reflect.Value
			if source.CanAddr() {
				fromMethod = source.Addr().MethodByName(name)
			} else {
				fromMethod = source.MethodByName(name)
			}

			if fromMethod.IsValid() && fromMethod.Type().NumIn() == 0 && fromMethod.Type().NumOut() == 1 {
				if toField := dest.FieldByName(name); toField.IsValid() && toField.CanSet() {
					values := fromMethod.Call([]reflect.Value{})
					if len(values) >= 1 {
						set(toField, values[0])
					}
				}
			}
		}

		if isSlice {
			if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest.Addr()))
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest))
			}
		}
	}
	return
}

// deepFields returns a list of every fields of the object
func deepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				fields = append(fields, deepFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}
	}

	return fields
}

// indirect returns the value element in case of Ptr value
func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

// indirectType returns the type element in case of Ptr value
func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

// set sets the values from an object to another object
func set(to, from reflect.Value) bool {
	if from.IsValid() {
		if to.Kind() == reflect.Ptr {
			//set `to` to nil if from is nil
			if from.Kind() == reflect.Ptr && from.IsNil() {
				to.Set(reflect.Zero(to.Type()))
				return true
			} else if valuer, ok := from.Addr().Interface().(driver.Valuer); ok {
				if v, e := valuer.Value(); e != nil || v == nil {
					to.Set(reflect.Zero(to.Type()))
					return true
				}
			}

			if to.IsNil() {
				to.Set(reflect.New(to.Type().Elem()))
			}
			to = to.Elem()
		}

		if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))
		} else if from.Kind() == reflect.Ptr {
			return set(to, from.Elem())
		} else if scanner, ok := to.Addr().Interface().(sql.Scanner); ok {
			var i = from.Interface()
			if from.Kind() == reflect.Ptr && from.IsNil() {
				i = nil
			}
			scanner.Scan(i)
		} else if valuer, ok := from.Addr().Interface().(driver.Valuer); ok {
			if v, e := valuer.Value(); e == nil && v != nil {
				var v1 = reflect.ValueOf(v)
				if v1.Type().ConvertibleTo(to.Type()) {
					to.Set(v1.Convert(to.Type()))
				}
			}
		} else {
			return false
		}
	}
	return true
}
