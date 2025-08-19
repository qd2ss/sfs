package sfs

import (
	"errors"
	"fmt"
	"reflect"
)

func Marshal(v interface{}) (SFSObject, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, errors.New("only structs can be marshaled to SFSObject")
	}

	result := make(SFSObject)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Skip unexported fields
		if !fieldVal.CanInterface() {
			continue
		}

		info, err := parseTag(field)
		if err != nil {
			return nil, err
		}

		// Skip zero value optional fields
		if info.optional && isZero(fieldVal) {
			continue
		}

		sfsValue, err := convertToSFSValue(fieldVal, info.dataType)
		if err != nil {
			return nil, fmt.Errorf("field %s: %v", field.Name, err)
		}

		result[info.name] = sfsValue
	}

	return result, nil
}

func convertToSFSValue(val reflect.Value, dtype DataType) (interface{}, error) {
	// 处理指针类型
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, nil
		}
		val = val.Elem()
	}

	// 自动类型推断
	if dtype == NULL {
		switch val.Kind() {
		case reflect.Bool:
			dtype = BOOL
		case reflect.Int8:
			dtype = BYTE
		case reflect.Uint8:
			dtype = BYTE
		case reflect.Int16:
			dtype = SHORT
		case reflect.Uint16:
			dtype = SHORT
		case reflect.Int32:
			dtype = INT
		case reflect.Uint32:
			dtype = INT
		case reflect.Int, reflect.Int64:
			dtype = LONG
		case reflect.Uint, reflect.Uint64:
			dtype = LONG
		case reflect.Float32:
			dtype = FLOAT
		case reflect.Float64:
			dtype = DOUBLE
		case reflect.String:
			dtype = UTF_STRING
		case reflect.Slice, reflect.Array:
			if val.Type().Elem().Kind() == reflect.Interface {
				return convertInterfaceSliceToSFS(val)
			}
			return convertSliceToSFS(val, dtype)
		case reflect.Struct:
			return Marshal(val.Interface())
		case reflect.Interface:
			// 处理 interface{} 类型
			if val.IsNil() {
				return nil, nil
			}
			return convertToSFSValue(val.Elem(), dtype)
		case reflect.Map:
			// 处理 map 类型
			return convertMapToSFSObject(val)
		default:
			return nil, fmt.Errorf("unsupported type: %s", val.Kind())
		}
	}

	// 处理数组类型
	if dtype >= BOOL_ARRAY && dtype <= UTF_STRING_ARRAY {
		if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
			return nil, fmt.Errorf("expected slice/array for type %d, got %s", dtype, val.Kind())
		}
		// 特殊处理 []interface{} 类型
		if val.Type().Elem().Kind() == reflect.Interface {
			return convertInterfaceSliceToSFS(val)
		}
		return convertSliceToSFS(val, dtype)
	}

	// 基本类型处理
	switch dtype {
	case BOOL:
		return val.Bool(), nil
	case BYTE:
		return byte(val.Uint()), nil
	case SHORT:
		return int16(val.Int()), nil
	case INT:
		return int32(val.Int()), nil
	case LONG:
		return val.Int(), nil
	case FLOAT:
		return float32(val.Float()), nil
	case DOUBLE:
		return val.Float(), nil
	case UTF_STRING, TEXT:
		return val.String(), nil
	case SFS_OBJECT:
		if val.Kind() == reflect.Struct {
			return Marshal(val.Interface())
		}
		return nil, fmt.Errorf("cannot convert %s to SFS_OBJECT", val.Kind())
	case SFS_ARRAY:
		if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			// 特殊处理 []interface{} 类型
			if val.Type().Elem().Kind() == reflect.Interface {
				return convertInterfaceSliceToSFS(val)
			}
			return convertSliceToSFS(val, dtype)
		}
		return nil, fmt.Errorf("cannot convert %s to SFS_ARRAY", val.Kind())
	default:
		return nil, fmt.Errorf("unsupported SFS type: %d", dtype)
	}
}

// convertInterfaceSliceToSFS 将 []interface{} 转换为 SFS 兼容类型
func convertInterfaceSliceToSFS(val reflect.Value) (interface{}, error) {
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return nil, errors.New("value is not a slice or array")
	}

	length := val.Len()
	// if length == 0 {
	// 	return nil, nil
	// }

	// 递归处理每个元素
	arr := make(SFSArray, length)
	for i := 0; i < length; i++ {
		elem := val.Index(i)
		if !elem.IsValid() {
			arr[i] = nil
			continue
		}

		// 处理嵌套的 interface{} 值
		if elem.Kind() == reflect.Interface {
			if elem.IsNil() {
				arr[i] = nil
				continue
			}
			elem = elem.Elem()
		}

		// 根据元素类型递归处理
		var err error
		switch elem.Kind() {
		case reflect.Slice, reflect.Array:
			if elem.Type().Elem().Kind() == reflect.Interface {
				arr[i], err = convertInterfaceSliceToSFS(elem)
			} else {
				arr[i], err = convertSliceToSFS(elem, NULL)
			}
		case reflect.Map:
			arr[i], err = convertMapToSFSObject(elem)
		case reflect.Struct:
			arr[i], err = Marshal(elem.Interface())
		case reflect.Interface:
			// 多层 interface{} 包装
			arr[i], err = convertToSFSValue(elem, NULL)
		default:
			arr[i], err = convertToSFSValue(elem, NULL)
		}

		if err != nil {
			return nil, fmt.Errorf("element %d: %v", i, err)
		}
	}

	return arr, nil
}

// convertMapToSFSObject 将 map 转换为 SFSObject
func convertMapToSFSObject(val reflect.Value) (SFSObject, error) {
	if val.Kind() != reflect.Map {
		return nil, errors.New("value is not a map")
	}

	obj := make(SFSObject)
	for _, key := range val.MapKeys() {
		if key.Kind() != reflect.String {
			return nil, errors.New("map key must be string")
		}

		mapVal := val.MapIndex(key)
		if !mapVal.IsValid() {
			obj[key.String()] = nil
			continue
		}

		var err error
		switch mapVal.Kind() {
		case reflect.Interface:
			if mapVal.IsNil() {
				obj[key.String()] = nil
				continue
			}
			obj[key.String()], err = convertToSFSValue(mapVal.Elem(), NULL)
		case reflect.Slice, reflect.Array:
			if mapVal.Type().Elem().Kind() == reflect.Interface {
				obj[key.String()], err = convertInterfaceSliceToSFS(mapVal)
			} else {
				obj[key.String()], err = convertSliceToSFS(mapVal, NULL)
			}
		case reflect.Map:
			obj[key.String()], err = convertMapToSFSObject(mapVal)
		case reflect.Struct:
			obj[key.String()], err = Marshal(mapVal.Interface())
		default:
			obj[key.String()], err = convertToSFSValue(mapVal, NULL)
		}

		if err != nil {
			return nil, fmt.Errorf("key %s: %v", key.String(), err)
		}
	}

	return obj, nil
}

func convertSliceToSFS(val reflect.Value, dtype DataType) (interface{}, error) {
	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return nil, errors.New("value is not a slice or array")
	}

	length := val.Len()
	if length == 0 {
		return nil, nil
	}

	// 根据指定的SFS数据类型处理
	switch dtype {
	case BOOL_ARRAY:
		arr := make([]bool, length)
		for i := 0; i < length; i++ {
			elem := val.Index(i)
			if elem.Kind() == reflect.Ptr && elem.IsNil() {
				continue
			}
			arr[i] = elem.Bool()
		}
		return arr, nil

	case BYTE_ARRAY:
		if val.Type().Elem().Kind() == reflect.Uint8 {
			return val.Bytes(), nil
		}
		arr := make([]byte, length)
		for i := 0; i < length; i++ {
			elem := val.Index(i)
			if elem.Kind() == reflect.Ptr && elem.IsNil() {
				continue
			}
			arr[i] = byte(elem.Uint())
		}
		return arr, nil

	case SHORT_ARRAY:
		arr := make([]int16, length)
		for i := 0; i < length; i++ {
			elem := val.Index(i)
			if elem.Kind() == reflect.Ptr && elem.IsNil() {
				continue
			}
			arr[i] = int16(elem.Int())
		}
		return arr, nil

	case INT_ARRAY:
		arr := make([]int32, length)
		for i := 0; i < length; i++ {
			elem := val.Index(i)
			if elem.Kind() == reflect.Ptr && elem.IsNil() {
				continue
			}
			arr[i] = int32(elem.Int())
		}
		return arr, nil

	case LONG_ARRAY:
		arr := make([]int64, length)
		for i := 0; i < length; i++ {
			elem := val.Index(i)
			if elem.Kind() == reflect.Ptr && elem.IsNil() {
				continue
			}
			arr[i] = elem.Int()
		}
		return arr, nil

	case FLOAT_ARRAY:
		arr := make([]float32, length)
		for i := 0; i < length; i++ {
			elem := val.Index(i)
			if elem.Kind() == reflect.Ptr && elem.IsNil() {
				continue
			}
			arr[i] = float32(elem.Float())
		}
		return arr, nil

	case DOUBLE_ARRAY:
		arr := make([]float64, length)
		for i := 0; i < length; i++ {
			elem := val.Index(i)
			if elem.Kind() == reflect.Ptr && elem.IsNil() {
				continue
			}
			arr[i] = elem.Float()
		}
		return arr, nil

	case UTF_STRING_ARRAY:
		arr := make([]string, length)
		for i := 0; i < length; i++ {
			elem := val.Index(i)
			if elem.Kind() == reflect.Ptr && elem.IsNil() {
				continue
			}
			arr[i] = elem.String()
		}
		return arr, nil

	case SFS_ARRAY:
		arr := make(SFSArray, length)
		for i := 0; i < length; i++ {
			elem := val.Index(i)
			if elem.Kind() == reflect.Ptr && elem.IsNil() {
				arr[i] = nil
				continue
			}
			var err error
			arr[i], err = convertToSFSValue(elem, NULL)
			if err != nil {
				return nil, err
			}
		}
		return arr, nil

	default:
		// 没有指定具体数组类型，尝试自动推断
		elemType := val.Type().Elem()
		switch elemType.Kind() {
		case reflect.Bool:
			return convertSliceToSFS(val, BOOL_ARRAY)
		case reflect.Uint8:
			return convertSliceToSFS(val, BYTE_ARRAY)
		case reflect.Int16:
			return convertSliceToSFS(val, SHORT_ARRAY)
		case reflect.Int32:
			return convertSliceToSFS(val, INT_ARRAY)
		case reflect.Int, reflect.Int64:
			return convertSliceToSFS(val, LONG_ARRAY)
		case reflect.Float32:
			return convertSliceToSFS(val, FLOAT_ARRAY)
		case reflect.Float64:
			return convertSliceToSFS(val, DOUBLE_ARRAY)
		case reflect.String:
			return convertSliceToSFS(val, UTF_STRING_ARRAY)
		case reflect.Struct, reflect.Interface, reflect.Ptr:
			return convertSliceToSFS(val, SFS_ARRAY)
		default:
			return nil, fmt.Errorf("unsupported slice element type: %s", elemType.Kind())
		}
	}
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		return v.String() == ""
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	case reflect.Struct:
		return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
	default:
		return false
	}
}
