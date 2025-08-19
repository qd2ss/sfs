package sfs

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

type Unpacker struct {
	buf *bytes.Buffer
}

func NewUnpacker(data []byte) *Unpacker {
	return &Unpacker{buf: bytes.NewBuffer(data)}
}

func (u *Unpacker) Unpack() (interface{}, error) {
	firstByte, err := u.buf.ReadByte()
	if err != nil {
		return nil, err
	}

	compressed := (firstByte & 32) > 0
	lengthIn4Bytes := (firstByte & 8) > 0

	var dataLength uint32
	if lengthIn4Bytes {
		if err := binary.Read(u.buf, binary.BigEndian, &dataLength); err != nil {
			return nil, err
		}
	} else {
		var length uint16
		if err := binary.Read(u.buf, binary.BigEndian, &length); err != nil {
			return nil, err
		}
		dataLength = uint32(length)
	}

	data := make([]byte, dataLength)
	if _, err := io.ReadFull(u.buf, data); err != nil {
		return nil, err
	}

	if compressed {
		r, err := zlib.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		defer r.Close()

		var decompressed bytes.Buffer
		if _, err := io.Copy(&decompressed, r); err != nil {
			return nil, err
		}
		data = decompressed.Bytes()
	}

	u.buf = bytes.NewBuffer(data)
	return u.decodeValue()
}

func (u *Unpacker) decodeValue() (interface{}, error) {
	typeByte, err := u.buf.ReadByte()
	if err != nil {
		return nil, err
	}

	dataType := DataType(typeByte)

	switch dataType {
	case NULL:
		return nil, nil
	case BOOL:
		val, err := u.buf.ReadByte()
		if err != nil {
			return nil, err
		}
		return val != 0, nil
	case BYTE:
		return u.buf.ReadByte()
	case SHORT:
		var val int16
		if err := binary.Read(u.buf, binary.BigEndian, &val); err != nil {
			return nil, err
		}
		return val, nil
	case INT:
		var val int32
		if err := binary.Read(u.buf, binary.BigEndian, &val); err != nil {
			return nil, err
		}
		return val, nil
	case LONG:
		var val int64
		if err := binary.Read(u.buf, binary.BigEndian, &val); err != nil {
			return nil, err
		}
		return val, nil
	case FLOAT:
		var val float32
		if err := binary.Read(u.buf, binary.BigEndian, &val); err != nil {
			return nil, err
		}
		return val, nil
	case DOUBLE:
		var val float64
		if err := binary.Read(u.buf, binary.BigEndian, &val); err != nil {
			return nil, err
		}
		return val, nil
	case UTF_STRING:
		var length uint16
		if err := binary.Read(u.buf, binary.BigEndian, &length); err != nil {
			return nil, err
		}
		strBytes := make([]byte, length)
		if _, err := io.ReadFull(u.buf, strBytes); err != nil {
			return nil, err
		}
		return string(strBytes), nil
	case BOOL_ARRAY:
		var size uint32
		if err := binary.Read(u.buf, binary.BigEndian, &size); err != nil {
			return nil, err
		}
		arr := make([]bool, size)
		for i := uint32(0); i < size; i++ {
			val, err := u.buf.ReadByte()
			if err != nil {
				return nil, err
			}
			arr[i] = val != 0
		}
		return arr, nil
	case BYTE_ARRAY:
		var size uint32
		if err := binary.Read(u.buf, binary.BigEndian, &size); err != nil {
			return nil, err
		}
		arr := make([]byte, size)
		if _, err := io.ReadFull(u.buf, arr); err != nil {
			return nil, err
		}
		return arr, nil
	case SHORT_ARRAY:
		var size uint32
		if err := binary.Read(u.buf, binary.BigEndian, &size); err != nil {
			return nil, err
		}
		arr := make([]int16, size)
		if err := binary.Read(u.buf, binary.BigEndian, &arr); err != nil {
			return nil, err
		}
		return arr, nil
	case INT_ARRAY:
		var size uint32
		if err := binary.Read(u.buf, binary.BigEndian, &size); err != nil {
			return nil, err
		}
		arr := make([]int32, size)
		if err := binary.Read(u.buf, binary.BigEndian, &arr); err != nil {
			return nil, err
		}
		return arr, nil
	case LONG_ARRAY:
		var size uint32
		if err := binary.Read(u.buf, binary.BigEndian, &size); err != nil {
			return nil, err
		}
		arr := make([]int64, size)
		if err := binary.Read(u.buf, binary.BigEndian, &arr); err != nil {
			return nil, err
		}
		return arr, nil
	case FLOAT_ARRAY:
		var size uint32
		if err := binary.Read(u.buf, binary.BigEndian, &size); err != nil {
			return nil, err
		}
		arr := make([]float32, size)
		if err := binary.Read(u.buf, binary.BigEndian, &arr); err != nil {
			return nil, err
		}
		return arr, nil
	case DOUBLE_ARRAY:
		var size uint32
		if err := binary.Read(u.buf, binary.BigEndian, &size); err != nil {
			return nil, err
		}
		arr := make([]float64, size)
		if err := binary.Read(u.buf, binary.BigEndian, &arr); err != nil {
			return nil, err
		}
		return arr, nil
	case UTF_STRING_ARRAY:
		var size uint32
		if err := binary.Read(u.buf, binary.BigEndian, &size); err != nil {
			return nil, err
		}
		arr := make([]string, size)
		for i := uint32(0); i < size; i++ {
			var length uint16
			if err := binary.Read(u.buf, binary.BigEndian, &length); err != nil {
				return nil, err
			}
			strBytes := make([]byte, length)
			if _, err := io.ReadFull(u.buf, strBytes); err != nil {
				return nil, err
			}
			arr[i] = string(strBytes)
		}
		return arr, nil
	case SFS_OBJECT:
		var count uint16
		if err := binary.Read(u.buf, binary.BigEndian, &count); err != nil {
			return nil, err
		}

		obj := make(SFSObject)
		for i := uint16(0); i < count; i++ {
			var keyLen uint16
			if err := binary.Read(u.buf, binary.BigEndian, &keyLen); err != nil {
				return nil, err
			}
			if keyLen < 0 || keyLen > 255 {
				return nil, errors.New("invalid SFSObject key length")
			}

			keyBytes := make([]byte, keyLen)
			if _, err := io.ReadFull(u.buf, keyBytes); err != nil {
				return nil, err
			}
			key := string(keyBytes)

			value, err := u.decodeValue()
			if err != nil {
				return nil, fmt.Errorf("could not decode value for key: %s, error: %v", key, err)
			}
			if value == nil {
				return nil, fmt.Errorf("could not decode value for key: %s", key)
			}

			obj[key] = value
		}
		return obj, nil
	case SFS_ARRAY:
		var count uint16
		if err := binary.Read(u.buf, binary.BigEndian, &count); err != nil {
			return nil, err
		}

		arr := make(SFSArray, count)
		for i := uint16(0); i < count; i++ {
			value, err := u.decodeValue()
			if err != nil {
				return nil, fmt.Errorf("could not decode value for index: %d, error: %v", i, err)
			}
			if value == nil {
				return nil, fmt.Errorf("could not decode value for index: %d", i)
			}

			arr[i] = value
		}
		return arr, nil
	case TEXT:
		var length uint32
		if err := binary.Read(u.buf, binary.BigEndian, &length); err != nil {
			return nil, err
		}
		strBytes := make([]byte, length)
		if _, err := io.ReadFull(u.buf, strBytes); err != nil {
			return nil, err
		}
		return string(strBytes), nil
	default:
		return nil, fmt.Errorf("unknown data type: %d", dataType)
	}
}
