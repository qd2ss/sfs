package sfs

import (
	"errors"
	"fmt"
	"math"
	"reflect"
)

func Unmarshal(data SFSObject, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return errors.New("must pass a pointer to a struct")
	}

	val = val.Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		info, err := parseTag(field)
		if err != nil {
			return err
		}

		sfsValue, exists := data[info.name]
		if !exists {
			if info.optional {
				continue
			}
			return fmt.Errorf("required field %s not found", info.name)
		}

		if err := convertFromSFSValue(fieldVal, sfsValue, info.dataType); err != nil {
			return fmt.Errorf("field %s: %v", field.Name, err)
		}
	}

	return nil
}

func convertFromSFSValue(field reflect.Value, sfsValue interface{}, dtype DataType) error {
	if sfsValue == nil {
		field.Set(reflect.Zero(field.Type()))
		return nil
	}

	if dtype == NULL {
		return autoConvert(field, sfsValue)
	}

	if dtype >= BOOL_ARRAY && dtype <= UTF_STRING_ARRAY {
		if field.Kind() != reflect.Slice {
			return fmt.Errorf("cannot convert array to non-slice type %s", field.Kind())
		}
		return convertArrayToField(field, sfsValue, dtype)
	}

	switch dtype {
	case BOOL:
		if b, ok := sfsValue.(bool); ok {
			field.SetBool(b)
			return nil
		}
	case BYTE:
		if b, ok := sfsValue.(byte); ok {
			field.SetUint(uint64(b))
			return nil
		}
	case SHORT:
		if s, ok := sfsValue.(int16); ok {
			field.SetInt(int64(s))
			return nil
		}
	case INT:
		if i, ok := sfsValue.(int32); ok {
			field.SetInt(int64(i))
			return nil
		}
	case LONG:
		if l, ok := sfsValue.(int64); ok {
			field.SetInt(l)
			return nil
		}
	case FLOAT:
		if f, ok := sfsValue.(float32); ok {
			field.SetFloat(float64(f))
			return nil
		}
	case DOUBLE:
		if d, ok := sfsValue.(float64); ok {
			field.SetFloat(d)
			return nil
		}
	case UTF_STRING, TEXT:
		if s, ok := sfsValue.(string); ok {
			field.SetString(s)
			return nil
		}
	case SFS_OBJECT:
		if obj, ok := sfsValue.(SFSObject); ok {
			if field.Kind() == reflect.Struct {
				if field.CanAddr() && field.Addr().CanInterface() {
					return Unmarshal(obj, field.Addr().Interface())
				}
			} else if field.Kind() == reflect.Ptr {
				if field.IsNil() {
					field.Set(reflect.New(field.Type().Elem()))
				}
				return Unmarshal(obj, field.Interface())
			}
		}
	case SFS_ARRAY:
		return convertSliceToField(field, sfsValue)
	}

	return fmt.Errorf("cannot convert %T to %s with type %d",
		sfsValue, field.Type(), dtype)
}

func convertArrayToField(field reflect.Value, sfsValue interface{}, dtype DataType) error {
	sliceType := field.Type()
	elemType := sliceType.Elem()

	switch dtype {
	case BOOL_ARRAY:
		if arr, ok := sfsValue.([]bool); ok {
			slice := reflect.MakeSlice(sliceType, len(arr), len(arr))
			for i, val := range arr {
				slice.Index(i).SetBool(val)
			}
			field.Set(slice)
			return nil
		}

	case BYTE_ARRAY:
		if arr, ok := sfsValue.([]byte); ok {
			if elemType.Kind() == reflect.Uint8 {
				field.SetBytes(arr)
				return nil
			}
			slice := reflect.MakeSlice(sliceType, len(arr), len(arr))
			for i, val := range arr {
				slice.Index(i).SetUint(uint64(val))
			}
			field.Set(slice)
			return nil
		}

	case SHORT_ARRAY:
		if arr, ok := sfsValue.([]int16); ok {
			slice := reflect.MakeSlice(sliceType, len(arr), len(arr))
			for i, val := range arr {
				slice.Index(i).SetInt(int64(val))
			}
			field.Set(slice)
			return nil
		}

	case INT_ARRAY:
		if arr, ok := sfsValue.([]int32); ok {
			slice := reflect.MakeSlice(sliceType, len(arr), len(arr))
			for i, val := range arr {
				slice.Index(i).SetInt(int64(val))
			}
			field.Set(slice)
			return nil
		}

	case LONG_ARRAY:
		if arr, ok := sfsValue.([]int64); ok {
			slice := reflect.MakeSlice(sliceType, len(arr), len(arr))
			for i, val := range arr {
				slice.Index(i).SetInt(val)
			}
			field.Set(slice)
			return nil
		}

	case FLOAT_ARRAY:
		if arr, ok := sfsValue.([]float32); ok {
			slice := reflect.MakeSlice(sliceType, len(arr), len(arr))
			for i, val := range arr {
				slice.Index(i).SetFloat(float64(val))
			}
			field.Set(slice)
			return nil
		}

	case DOUBLE_ARRAY:
		if arr, ok := sfsValue.([]float64); ok {
			slice := reflect.MakeSlice(sliceType, len(arr), len(arr))
			for i, val := range arr {
				slice.Index(i).SetFloat(val)
			}
			field.Set(slice)
			return nil
		}

	case UTF_STRING_ARRAY:
		if arr, ok := sfsValue.([]string); ok {
			slice := reflect.MakeSlice(sliceType, len(arr), len(arr))
			for i, val := range arr {
				slice.Index(i).SetString(val)
			}
			field.Set(slice)
			return nil
		}

	case SFS_ARRAY:
		if arr, ok := sfsValue.(SFSArray); ok {
			slice := reflect.MakeSlice(sliceType, len(arr), len(arr))
			for i, val := range arr {
				elem := reflect.New(elemType).Elem()
				if err := autoConvert(elem, val); err != nil {
					return err
				}
				slice.Index(i).Set(elem)
			}
			field.Set(slice)
			return nil
		}
	}

	return fmt.Errorf("cannot convert %T to %s array type %d", sfsValue, elemType, dtype)
}

func autoConvert(field reflect.Value, sfsValue interface{}) error {
	switch field.Kind() {
	case reflect.Bool:
		if b, ok := sfsValue.(bool); ok {
			field.SetBool(b)
			return nil
		}

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch v := sfsValue.(type) {
		case int64:
			// 检查是否在目标类型的范围内
			if field.OverflowInt(v) {
				return fmt.Errorf("value %d overflows %s", v, field.Type())
			}
			field.SetInt(v)
			return nil
		case int32:
			if field.Kind() == reflect.Int32 || field.Kind() == reflect.Int {
				field.SetInt(int64(v))
				return nil
			}
		case int16:
			if field.Kind() == reflect.Int16 {
				field.SetInt(int64(v))
				return nil
			}
		case int8:
			if field.Kind() == reflect.Int8 {
				field.SetInt(int64(v))
				return nil
			}
		case int:
			if field.Kind() == reflect.Int {
				field.SetInt(int64(v))
				return nil
			}
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch v := sfsValue.(type) {
		case int64:
			if v >= 0 {
				u := uint64(v)
				if field.OverflowUint(u) {
					return fmt.Errorf("value %d overflows %s", v, field.Type())
				}
				field.SetUint(u)
				return nil
			}
		case uint64:
			if field.OverflowUint(v) {
				return fmt.Errorf("value %d overflows %s", v, field.Type())
			}
			field.SetUint(v)
			return nil
		case uint32:
			if field.Kind() == reflect.Uint32 || field.Kind() == reflect.Uint {
				field.SetUint(uint64(v))
				return nil
			}
		case uint16:
			if field.Kind() == reflect.Uint16 {
				field.SetUint(uint64(v))
				return nil
			}
		case uint8:
			if field.Kind() == reflect.Uint8 {
				field.SetUint(uint64(v))
				return nil
			}
		case uint:
			if field.Kind() == reflect.Uint {
				field.SetUint(uint64(v))
				return nil
			}
		}

	case reflect.Float32, reflect.Float64:
		switch v := sfsValue.(type) {
		case float32:
			if field.Kind() == reflect.Float32 {
				field.SetFloat(float64(v))
				return nil
			}
		case float64:
			if field.Kind() == reflect.Float64 || field.Kind() == reflect.Float32 {
				if field.Kind() == reflect.Float32 && (v > math.MaxFloat32 || v < -math.MaxFloat32) {
					return fmt.Errorf("value %f overflows float32", v)
				}
				field.SetFloat(v)
				return nil
			}
		}

	case reflect.String:
		if s, ok := sfsValue.(string); ok {
			field.SetString(s)
			return nil
		}

	case reflect.Struct:
		if obj, ok := sfsValue.(SFSObject); ok {
			if field.CanAddr() && field.Addr().CanInterface() {
				return Unmarshal(obj, field.Addr().Interface())
			}
		}

	case reflect.Ptr:
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		return autoConvert(field.Elem(), sfsValue)

	case reflect.Slice:
		// 处理 []interface{} 到具体切片类型的转换
		if srcSlice, ok := sfsValue.([]interface{}); ok {
			return convertSliceToField(field, srcSlice)
		}
		return convertSliceToField(field, sfsValue)

	case reflect.Interface:
		// 如果目标字段是 interface{} 类型，直接设置值
		if field.Type().NumMethod() == 0 { // 空接口
			field.Set(reflect.ValueOf(sfsValue))
			return nil
		}
	}

	return fmt.Errorf("cannot auto-convert %T to %s", sfsValue, field.Type())
}

// convertInterfaceSliceToField 将 []interface{} 转换为目标切片类型
// func convertInterfaceSliceToField(field reflect.Value, srcSlice []interface{}) error {
// 	dstType := field.Type()
// 	if dstType.Kind() != reflect.Slice {
// 		return fmt.Errorf("target is not a slice type: %s", dstType)
// 	}

// 	elemType := dstType.Elem()
// 	dstSlice := reflect.MakeSlice(dstType, len(srcSlice), len(srcSlice))

// 	for i, elem := range srcSlice {
// 		elemValue := reflect.New(elemType).Elem()
// 		if err := autoConvert(elemValue, elem); err != nil {
// 			return fmt.Errorf("element %d: %v", i, err)
// 		}
// 		dstSlice.Index(i).Set(elemValue)
// 	}

// 	field.Set(dstSlice)
// 	return nil
// }

func convertSliceToField(field reflect.Value, sfsValue interface{}) error {
	sliceType := field.Type()
	elemType := sliceType.Elem()

	var slice reflect.Value

	switch v := sfsValue.(type) {
	case []bool:
		if elemType.Kind() == reflect.Bool {
			slice = reflect.MakeSlice(sliceType, len(v), len(v))
			for i, val := range v {
				slice.Index(i).SetBool(val)
			}
			field.Set(slice)
			return nil
		}
	case []byte:
		if elemType.Kind() == reflect.Uint8 {
			field.SetBytes(v)
			return nil
		}
	case []int16:
		if elemType.Kind() == reflect.Int16 {
			slice = reflect.MakeSlice(sliceType, len(v), len(v))
			for i, val := range v {
				slice.Index(i).SetInt(int64(val))
			}
			field.Set(slice)
			return nil
		}
	case []int32:
		if elemType.Kind() == reflect.Int32 {
			slice = reflect.MakeSlice(sliceType, len(v), len(v))
			for i, val := range v {
				slice.Index(i).SetInt(int64(val))
			}
			field.Set(slice)
			return nil
		}
	case []int64:
		if elemType.Kind() == reflect.Int64 {
			slice = reflect.MakeSlice(sliceType, len(v), len(v))
			for i, val := range v {
				slice.Index(i).SetInt(val)
			}
			field.Set(slice)
			return nil
		}
	case []float32:
		if elemType.Kind() == reflect.Float32 {
			slice = reflect.MakeSlice(sliceType, len(v), len(v))
			for i, val := range v {
				slice.Index(i).SetFloat(float64(val))
			}
			field.Set(slice)
			return nil
		}
	case []float64:
		if elemType.Kind() == reflect.Float64 {
			slice = reflect.MakeSlice(sliceType, len(v), len(v))
			for i, val := range v {
				slice.Index(i).SetFloat(val)
			}
			field.Set(slice)
			return nil
		}
	case []string:
		if elemType.Kind() == reflect.String {
			slice = reflect.MakeSlice(sliceType, len(v), len(v))
			for i, val := range v {
				slice.Index(i).SetString(val)
			}
			field.Set(slice)
			return nil
		}
	case SFSArray:
		slice = reflect.MakeSlice(sliceType, len(v), len(v))
		for i, val := range v {
			elem := reflect.New(elemType).Elem()
			if err := autoConvert(elem, val); err != nil {
				return err
			}
			slice.Index(i).Set(elem)
		}
		field.Set(slice)
		return nil
	}

	return fmt.Errorf("cannot convert %T to %s", sfsValue, sliceType)
}
