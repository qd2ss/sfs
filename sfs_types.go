package sfs

type DataType byte

const (
	NULL             DataType = 0
	BOOL             DataType = 1
	BYTE             DataType = 2
	SHORT            DataType = 3
	INT              DataType = 4
	LONG             DataType = 5
	FLOAT            DataType = 6
	DOUBLE           DataType = 7
	UTF_STRING       DataType = 8
	BOOL_ARRAY       DataType = 9
	BYTE_ARRAY       DataType = 10
	SHORT_ARRAY      DataType = 11
	INT_ARRAY        DataType = 12
	LONG_ARRAY       DataType = 13
	FLOAT_ARRAY      DataType = 14
	DOUBLE_ARRAY     DataType = 15
	UTF_STRING_ARRAY DataType = 16
	SFS_ARRAY        DataType = 17
	SFS_OBJECT       DataType = 18
	TEXT             DataType = 20
)

type SFSObject map[string]interface{}
type SFSArray []interface{}

type fieldInfo struct {
	name     string
	dataType DataType
	optional bool
}

const tagName = "sfs"
