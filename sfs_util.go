package sfs

import (
	"fmt"
	"reflect"
	"strings"
)

func parseTag(field reflect.StructField) (fieldInfo, error) {
	info := fieldInfo{
		name:     field.Name,
		dataType: NULL,
		optional: false,
	}

	tag := field.Tag.Get(tagName)
	if tag == "" {
		return info, nil
	}

	parts := strings.Split(tag, ",")
	if len(parts) > 0 && parts[0] != "" {
		info.name = parts[0]
	}

	for _, part := range parts[1:] {
		if part == "optional" {
			info.optional = true
			continue
		}

		if strings.HasPrefix(part, "type=") {
			typeStr := strings.TrimPrefix(part, "type=")
			dtype, err := parseDataType(typeStr)
			if err != nil {
				return info, err
			}
			info.dataType = dtype
		}
	}

	return info, nil
}

func parseDataType(s string) (DataType, error) {
	switch strings.ToUpper(s) {
	case "NULL":
		return NULL, nil
	case "BOOL":
		return BOOL, nil
	case "BYTE":
		return BYTE, nil
	case "SHORT":
		return SHORT, nil
	case "INT":
		return INT, nil
	case "LONG":
		return LONG, nil
	case "FLOAT":
		return FLOAT, nil
	case "DOUBLE":
		return DOUBLE, nil
	case "UTF_STRING":
		return UTF_STRING, nil
	case "BOOL_ARRAY":
		return BOOL_ARRAY, nil
	case "BYTE_ARRAY":
		return BYTE_ARRAY, nil
	case "SHORT_ARRAY":
		return SHORT_ARRAY, nil
	case "INT_ARRAY":
		return INT_ARRAY, nil
	case "LONG_ARRAY":
		return LONG_ARRAY, nil
	case "FLOAT_ARRAY":
		return FLOAT_ARRAY, nil
	case "DOUBLE_ARRAY":
		return DOUBLE_ARRAY, nil
	case "UTF_STRING_ARRAY":
		return UTF_STRING_ARRAY, nil
	case "SFS_ARRAY":
		return SFS_ARRAY, nil
	case "SFS_OBJECT":
		return SFS_OBJECT, nil
	case "TEXT":
		return TEXT, nil
	default:
		return NULL, fmt.Errorf("unknown data type: %s", s)
	}
}
