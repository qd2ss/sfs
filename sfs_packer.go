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
	buf                  *bytes.Buffer
	compressionThreshold int
}

func NewPacker() *Packer {
	return &Packer{
		buf:                  new(bytes.Buffer),
		compressionThreshold: DefaultCompressionThreshold,
	}
}

// DefaultCompressionThreshold 为 SDK 默认值，与 JS 实现一致。
const DefaultCompressionThreshold = 1024

// CompressionDisabled 表示服务器通过 Handshake 下发的 ct=2147483647，禁用 SFS 层压缩。
const CompressionDisabled = 2147483647

// SetCompressionThreshold 设置压缩阈值 ct（Handshake 响应字段）。
// 仅当编码后的 payload 长度严格大于 ct 时才压缩；ct 为 CompressionDisabled 时永不压缩。
func (p *Packer) SetCompressionThreshold(ct int) {
	p.compressionThreshold = ct
}

// CompressionThreshold 返回当前压缩阈值。
func (p *Packer) CompressionThreshold() int {
	return p.compressionThreshold
}

// 4 字节长度阈值：与 JS 逻辑保持一致，超过 65335 则使用 4 字节长度并置位 0x08
const lengthThreshold4 = 65335

func (p *Packer) Pack(data SFSObject) ([]byte, error) {
	if err := p.encodeSFSObject(data); err != nil {
		return nil, err
	}

	dataBytes := p.buf.Bytes()
	p.buf = new(bytes.Buffer)

	// 首字节基础标志位设为 0x80，表示采用该协议帧格式（与 JS 保持一致）
	var firstByte byte = 128

	// 超过压缩阈值则进行压缩，并置位压缩标志位 0x20
	if len(dataBytes) > p.compressionThreshold {
		firstByte += 32
		var compressed bytes.Buffer
		w := zlib.NewWriter(&compressed)
		if _, err := w.Write(dataBytes); err != nil {
			w.Close()
			return nil, err
		}
		w.Close()
		dataBytes = compressed.Bytes()
	}

	// 长度使用“最终负载长度”（压缩后长度），超过阈值则置位 0x08 使用 4 字节长度（与 JS 一致为 += 8）
	dataLength := len(dataBytes)
	if dataLength > lengthThreshold4 {
		firstByte += 8
	}

	if err := p.buf.WriteByte(firstByte); err != nil {
		return nil, err
	}

	// 长度字段采用大端序写入：置位 0x08 时写入 uint32，否则写入 uint16
	if (firstByte & 8) > 0 {
		if err := binary.Write(p.buf, binary.BigEndian, uint32(dataLength)); err != nil {
			return nil, err
		}
	} else {
		if err := binary.Write(p.buf, binary.BigEndian, uint16(dataLength)); err != nil {
			return nil, err
		}
	}

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
