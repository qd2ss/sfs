package sfs

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

type Packer struct {
	buf *bytes.Buffer
}

func NewPacker() *Packer {
	return &Packer{buf: new(bytes.Buffer)}
}

func (p *Packer) Pack(data SFSObject, compress bool) ([]byte, error) {
	// First encode the SFSObject to binary
	if err := p.encodeSFSObject(data); err != nil {
		return nil, err
	}

	dataBytes := p.buf.Bytes()
	p.buf = new(bytes.Buffer) // Reset buffer

	// Set flags in first byte
	var firstByte byte
	if compress {
		firstByte |= 32 // Set compression flag
	}

	// Determine if we need 4-byte length
	dataLength := len(dataBytes)
	if dataLength > math.MaxUint16 {
		firstByte |= 8 // Set 4-byte length flag
	}

	// Write first byte
	if err := p.buf.WriteByte(firstByte); err != nil {
		return nil, err
	}

	// Write length
	if (firstByte & 8) > 0 {
		if err := binary.Write(p.buf, binary.BigEndian, uint32(dataLength)); err != nil {
			return nil, err
		}
	} else {
		if dataLength > math.MaxUint16 {
			return nil, errors.New("data too large for 2-byte length")
		}
		if err := binary.Write(p.buf, binary.BigEndian, uint16(dataLength)); err != nil {
			return nil, err
		}
	}

	// Compress if needed
	if compress {
		var compressed bytes.Buffer
		w := zlib.NewWriter(&compressed)
		if _, err := w.Write(dataBytes); err != nil {
			w.Close()
			return nil, err
		}
		w.Close()
		dataBytes = compressed.Bytes()
	}

	// Write data
	if _, err := p.buf.Write(dataBytes); err != nil {
		return nil, err
	}

	return p.buf.Bytes(), nil
}

func (p *Packer) encodeSFSObject(obj SFSObject) error {
	if err := p.buf.WriteByte(byte(SFS_OBJECT)); err != nil {
		return err
	}

	// Write number of key-value pairs
	if err := binary.Write(p.buf, binary.BigEndian, uint16(len(obj))); err != nil {
		return err
	}

	// Write each key-value pair
	for key, value := range obj {
		// Write key length (UTF-STRING)
		keyBytes := []byte(key)
		if len(keyBytes) > math.MaxUint16 {
			return errors.New("key too long")
		}
		if err := binary.Write(p.buf, binary.BigEndian, uint16(len(keyBytes))); err != nil {
			return err
		}
		// Write key
		if _, err := p.buf.Write(keyBytes); err != nil {
			return err
		}
		// Write value
		if err := p.encodeValue(value); err != nil {
			return err
		}
	}
	return nil
}

func (p *Packer) encodeValue(value interface{}) error {
	if value == nil {
		return p.encodeNull()
	}

	switch v := value.(type) {
	case bool:
		return p.encodeBool(v)
	case byte:
		return p.encodeByte(v)
	case int16:
		return p.encodeShort(v)
	case int32:
		return p.encodeInt(v)
	case int:
		return p.encodeInt(int32(v))
	case int64:
		return p.encodeLong(v)
	case float32:
		return p.encodeFloat(v)
	case float64:
		return p.encodeDouble(v)
	case string:
		return p.encodeUtfString(v)
	case []bool:
		return p.encodeBoolArray(v)
	case []byte:
		return p.encodeByteArray(v)
	case []int16:
		return p.encodeShortArray(v)
	case []int32:
		return p.encodeIntArray(v)
	case []int64:
		return p.encodeLongArray(v)
	case []float32:
		return p.encodeFloatArray(v)
	case []float64:
		return p.encodeDoubleArray(v)
	case []string:
		return p.encodeUtfStringArray(v)
	case map[string]interface{}:
		// 将普通map转换为SFSObject
		sfsObj := make(SFSObject)
		for k, val := range v {
			sfsObj[k] = val
		}
		return p.encodeSFSObject(sfsObj)
	case SFSObject:
		return p.encodeSFSObject(v)
	case SFSArray:
		return p.encodeSFSArray(v)
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
}

func (p *Packer) encodeNull() error {
	return p.buf.WriteByte(byte(NULL))
}

func (p *Packer) encodeBool(v bool) error {
	if err := p.buf.WriteByte(byte(BOOL)); err != nil {
		return err
	}
	var val byte
	if v {
		val = 1
	}
	return p.buf.WriteByte(val)
}

func (p *Packer) encodeByte(v byte) error {
	if err := p.buf.WriteByte(byte(BYTE)); err != nil {
		return err
	}
	return p.buf.WriteByte(v)
}

func (p *Packer) encodeShort(v int16) error {
	if err := p.buf.WriteByte(byte(SHORT)); err != nil {
		return err
	}
	return binary.Write(p.buf, binary.BigEndian, v)
}

func (p *Packer) encodeInt(v int32) error {
	if err := p.buf.WriteByte(byte(INT)); err != nil {
		return err
	}
	return binary.Write(p.buf, binary.BigEndian, v)
}

func (p *Packer) encodeLong(v int64) error {
	if err := p.buf.WriteByte(byte(LONG)); err != nil {
		return err
	}
	return binary.Write(p.buf, binary.BigEndian, v)
}

func (p *Packer) encodeFloat(v float32) error {
	if err := p.buf.WriteByte(byte(FLOAT)); err != nil {
		return err
	}
	return binary.Write(p.buf, binary.BigEndian, v)
}

func (p *Packer) encodeDouble(v float64) error {
	if err := p.buf.WriteByte(byte(DOUBLE)); err != nil {
		return err
	}
	return binary.Write(p.buf, binary.BigEndian, v)
}

func (p *Packer) encodeUtfString(v string) error {
	if err := p.buf.WriteByte(byte(UTF_STRING)); err != nil {
		return err
	}
	strBytes := []byte(v)
	if len(strBytes) > math.MaxUint16 {
		return errors.New("string too long")
	}
	if err := binary.Write(p.buf, binary.BigEndian, uint16(len(strBytes))); err != nil {
		return err
	}
	_, err := p.buf.Write(strBytes)
	return err
}

func (p *Packer) encodeBoolArray(v []bool) error {
	if err := p.buf.WriteByte(byte(BOOL_ARRAY)); err != nil {
		return err
	}
	if err := binary.Write(p.buf, binary.BigEndian, uint32(len(v))); err != nil {
		return err
	}
	for _, b := range v {
		var val byte
		if b {
			val = 1
		}
		if err := p.buf.WriteByte(val); err != nil {
			return err
		}
	}
	return nil
}

func (p *Packer) encodeByteArray(v []byte) error {
	if err := p.buf.WriteByte(byte(BYTE_ARRAY)); err != nil {
		return err
	}
	if err := binary.Write(p.buf, binary.BigEndian, uint32(len(v))); err != nil {
		return err
	}
	_, err := p.buf.Write(v)
	return err
}

func (p *Packer) encodeShortArray(v []int16) error {
	if err := p.buf.WriteByte(byte(SHORT_ARRAY)); err != nil {
		return err
	}
	if err := binary.Write(p.buf, binary.BigEndian, uint16(len(v))); err != nil {
		return err
	}
	return binary.Write(p.buf, binary.BigEndian, v)
}

func (p *Packer) encodeIntArray(v []int32) error {
	if err := p.buf.WriteByte(byte(INT_ARRAY)); err != nil {
		return err
	}
	if err := binary.Write(p.buf, binary.BigEndian, uint16(len(v))); err != nil {
		return err
	}
	return binary.Write(p.buf, binary.BigEndian, v)
}

func (p *Packer) encodeLongArray(v []int64) error {
	if err := p.buf.WriteByte(byte(LONG_ARRAY)); err != nil {
		return err
	}
	if err := binary.Write(p.buf, binary.BigEndian, uint16(len(v))); err != nil {
		return err
	}
	return binary.Write(p.buf, binary.BigEndian, v)
}

func (p *Packer) encodeFloatArray(v []float32) error {
	if err := p.buf.WriteByte(byte(FLOAT_ARRAY)); err != nil {
		return err
	}
	if err := binary.Write(p.buf, binary.BigEndian, uint16(len(v))); err != nil {
		return err
	}
	return binary.Write(p.buf, binary.BigEndian, v)
}

func (p *Packer) encodeDoubleArray(v []float64) error {
	if err := p.buf.WriteByte(byte(DOUBLE_ARRAY)); err != nil {
		return err
	}
	if err := binary.Write(p.buf, binary.BigEndian, uint16(len(v))); err != nil {
		return err
	}
	return binary.Write(p.buf, binary.BigEndian, v)
}

func (p *Packer) encodeUtfStringArray(v []string) error {
	if err := p.buf.WriteByte(byte(UTF_STRING_ARRAY)); err != nil {
		return err
	}
	if err := binary.Write(p.buf, binary.BigEndian, uint16(len(v))); err != nil {
		return err
	}
	for _, s := range v {
		strBytes := []byte(s)
		if len(strBytes) > math.MaxUint16 {
			return errors.New("string too long")
		}
		if err := binary.Write(p.buf, binary.BigEndian, uint16(len(strBytes))); err != nil {
			return err
		}
		if _, err := p.buf.Write(strBytes); err != nil {
			return err
		}
	}
	return nil
}

func (p *Packer) encodeSFSArray(v SFSArray) error {
	if err := p.buf.WriteByte(byte(SFS_ARRAY)); err != nil {
		return err
	}
	if err := binary.Write(p.buf, binary.BigEndian, uint16(len(v))); err != nil {
		return err
	}
	for _, value := range v {
		if err := p.encodeValue(value); err != nil {
			return err
		}
	}
	return nil
}
