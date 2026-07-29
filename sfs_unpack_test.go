package sfs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"testing"
)

type Data struct {
	Index int32 `sfs:"index"`
}

func TestArrayRoundtrip(t *testing.T) {
	obj := SFSObject{
		"bools": []bool{true, false, true, false},
		"longs": []int64{1, 2, 3, 4, 5},
	}
	p := NewPacker()
	data, err := p.Pack(obj)
	if err != nil {
		t.Fatal(err)
	}
	u := NewUnpacker(data)
	v, err := u.Unpack()
	if err != nil {
		t.Fatal(err)
	}
	sfo, ok := v.(SFSObject)
	if !ok {
		t.Fatalf("unexpected type: %T", v)
	}
	bools, ok := sfo["bools"].([]bool)
	if !ok || len(bools) != 4 {
		t.Fatalf("bools wrong type/len: %T len=%d", sfo["bools"], len(bools))
	}
	longs, ok := sfo["longs"].([]int64)
	if !ok || len(longs) != 5 {
		t.Fatalf("longs wrong type/len: %T len=%d", sfo["longs"], len(longs))
	}
}

type Respond struct {
	Code int32  `sfs:"code"`
	Data Data   `sfs:"data"`
	Msg  string `sfs:"msg"`
}

type P struct {
	Code   string `sfs:"code"`
	Entity []byte `sfs:"entity"`
}

func TestUpack(t *testing.T) {
	hexStr := "gADhEgADAAFjAgEAAWEDAA0AAXASAAMAAWMIAAhiYXIuc3BpbgABcgT/////AAFwEgADAARjb2RlCAAIYmFyLnNwaW4AAnNuBQAAuEUvq3U0AAZlbnRpdHkSAAMABWRlbm9tBAAAA+gACXBsYXllckJldAUAAAAAAAAAAAADYmV0DQAMAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	data, err := base64.StdEncoding.DecodeString(hexStr)
	if err != nil {
		fmt.Println("base64 decode error:", err)
		return
	}
	// fmt.Println("base64 decode:", data)

	//data, err := hex.DecodeString(hexStr)
	unpacker := NewUnpacker(data)
	v, err := unpacker.Unpack()
	if err != nil {
		t.Fatal(err)
	}
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}

	log.Println(string(jsonBytes))
}

func TestPackByStruct(t *testing.T) {
	var moudleType int16 = 1
	var commandType uint8 = 13

	hexStr := "eyJiYWxhbmNlIjo4MDAwMDAwMDMxMC40MDAsInNwaW5SZXN1bHQiOnsicGxheWVyVG90YWxXaW4iOjExMCwiYmFzZUdhbWVSZXN1bHQiOnsidXNlZFRhYmxlSW5kZXgiOjAsInNjcmVlblN5bWJvbCI6W1s2LDAsMywxMiw2XSxbOCw3LDMsNiwxMV0sWzgsNywwLDEwLDExXV0sImJhc2VHYW1lVG90YWxXaW4iOjExMCwid2F5c0dhbWVSZXN1bHQiOnsicGxheWVyV2luIjoxMTAsIndheXNSZXN1bHQiOlt7InN5bWJvbElEIjo2LCJoaXREaXJlY3Rpb24iOiJMZWZ0VG9SaWdodCIsImhpdE51bWJlciI6NSwiY291bnQiOjEsImhpdE9kZHMiOjEwMCwic3ltYm9sV2luIjoxMDAsInNjcmVlbkhpdERhdGEiOltbdHJ1ZSx0cnVlLGZhbHNlLGZhbHNlLHRydWVdLFtmYWxzZSxmYWxzZSxmYWxzZSx0cnVlLGZhbHNlXSxbZmFsc2UsZmFsc2UsdHJ1ZSxmYWxzZSxmYWxzZV1dfSx7InN5bWJvbElEIjo4LCJoaXREaXJlY3Rpb24iOiJMZWZ0VG9SaWdodCIsImhpdE51bWJlciI6MywiY291bnQiOjIsImhpdE9kZHMiOjUsInN5bWJvbFdpbiI6MTAsInNjcmVlbkhpdERhdGEiOltbZmFsc2UsdHJ1ZSxmYWxzZSxmYWxzZSxmYWxzZV0sW3RydWUsZmFsc2UsZmFsc2UsZmFsc2UsZmFsc2VdLFt0cnVlLGZhbHNlLHRydWUsZmFsc2UsZmFsc2VdXX1dfSwic3BlY2lhbEZlYXR1cmVSZXN1bHQiOlt7InNwZWNpYWxIaXRJbmZvIjoibm9TcGVjaWFsSGl0Iiwic3BlY2lhbE9wZXJhdGlvbnMiOltdLCJzcGVjaWFsU2NyZWVuSGl0RGF0YSI6W1tmYWxzZSxmYWxzZSxmYWxzZSxmYWxzZSxmYWxzZV0sW2ZhbHNlLGZhbHNlLGZhbHNlLGZhbHNlLGZhbHNlXSxbZmFsc2UsZmFsc2UsZmFsc2UsZmFsc2UsZmFsc2VdXSwic3BlY2lhbFNjcmVlbldpbiI6MH0seyJzcGVjaWFsSGl0SW5mbyI6Im5vU3BlY2lhbEhpdCIsInNwZWNpYWxPcGVyYXRpb25zIjpbXSwic3BlY2lhbFNjcmVlbkhpdERhdGEiOltbZmFsc2UsZmFsc2UsZmFsc2UsZmFsc2UsZmFsc2VdLFtmYWxzZSxmYWxzZSxmYWxzZSxmYWxzZSxmYWxzZV0sW2ZhbHNlLGZhbHNlLGZhbHNlLGZhbHNlLGZhbHNlXV0sInNwZWNpYWxTY3JlZW5XaW4iOjB9LHsic3BlY2lhbEhpdEluZm8iOiJub1NwZWNpYWxIaXQiLCJzcGVjaWFsT3BlcmF0aW9ucyI6W10sInNwZWNpYWxTY3JlZW5IaXREYXRhIjpbW2ZhbHNlLGZhbHNlLGZhbHNlLGZhbHNlLGZhbHNlXSxbZmFsc2UsZmFsc2UsZmFsc2UsZmFsc2UsZmFsc2VdLFtmYWxzZSxmYWxzZSxmYWxzZSxmYWxzZSxmYWxzZV1dLCJzcGVjaWFsU2NyZWVuV2luIjowfV0sImRpc3BsYXlJbmZvIjp7ImRpc3BsYXlNZXRob2QiOltbZmFsc2UsZmFsc2UsZmFsc2UsZmFsc2UsZmFsc2VdXSwiYmlnV2luVHlwZSI6Im5vcm1hbCIsImRhbXBJbmZvIjpbWzksMTAsMywxMiw3XSxbNSwwLDQsMTAsNl1dfSwiZXh0ZW5kSW5mb0ZvcmJhc2VHYW1lUmVzdWx0Ijp7InJlYWR5SGFuZEZsYWciOltmYWxzZSxmYWxzZSxmYWxzZSxmYWxzZSxmYWxzZV0sInNldExpZ2h0RmxhZyI6W1tmYWxzZSxmYWxzZSxmYWxzZSxmYWxzZSxmYWxzZV0sW2ZhbHNlLGZhbHNlLGZhbHNlLGZhbHNlLGZhbHNlXSxbZmFsc2UsZmFsc2UsZmFsc2UsZmFsc2UsZmFsc2VdXSwiZXh0ZW5kUGxheWVyV2luIjowfX0sImdhbWVGbG93UmVzdWx0Ijp7IklzQm9hcmRFbmRGbGFnIjp0cnVlLCJwZXJtaXNzaW9uT3BlcmF0aW9ucyI6WyJiYXNlR2FtZSJdfX0sImdhbWVTZXEiOjc0OTk3MzY0NDQ3NjksInRzIjoxNzU1MDgyODkyODc4fQ=="
	entity, err := base64.StdEncoding.DecodeString(hexStr)
	if err != nil {
		panic(err)
	}

	rsp, err := Marshal(&P{
		Code:   "spinResponse",
		Entity: entity,
	})

	if err == nil {
		t.Logf("pack by struct: %v", rsp)
		bytes, err := NewPacker().Pack(rsp)
		if err == nil {
			t.Logf("pack by struct: %v", bytes)
		}
	}

	data := SFSObject{
		"c": "h5.spinResponse",
		"p": rsp,
	}

	sendData := SFSObject{
		"a": commandType,
		"c": moudleType,
		"p": data,
	}

	packer := NewPacker()
	v, err := packer.Pack(sendData)

	fmt.Println(base64.StdEncoding.EncodeToString(v))

	// rsp, err := Marshal(Respond{
	// 	Code: 200,
	// 	Data: Data{Index: 1},
	// 	Msg:  "success",
	// })
	// if err == nil {
	// 	t.Logf("pack by struct: %v", rsp)
	// 	bytes, err := NewPacker().Pack(rsp, false)
	// 	if err == nil {
	// 		t.Logf("pack by struct: %v", bytes)
	// 	}
	// }

	// data := &Respond{}
	// err = Unmarshal(rsp, data)
	// if err == nil {
	// 	t.Logf("pack by struct: %v", data)
	// } else {
	// 	t.Logf("pack by struct: %v", err)
	// }

}

func TestPackCompressed(t *testing.T) {
	obj := SFSObject{
		"a": uint8(0),
		"c": int16(0),
		"p": SFSObject{
			"api": "1.7.15",
			"cl":  "JavaScript",
		},
	}
	p := NewPacker()
	data, err := p.Pack(obj)
	if err != nil {
		t.Fatal(err)
	}
	u := NewUnpacker(data)
	v, err := u.Unpack()
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := v.(SFSObject); !ok {
		t.Fatalf("unexpected type: %T", v)
	}
}

func TestCompressionThreshold(t *testing.T) {
	large := make([]byte, 1100)
	for i := range large {
		large[i] = byte('a' + (i % 26))
	}
	obj := SFSObject{"data": large}

	p := NewPacker()
	packed, err := p.Pack(obj)
	if err != nil {
		t.Fatal(err)
	}
	if packed[0]&0x20 == 0 {
		t.Fatalf("expected compressed packet with default ct, got header 0x%02x", packed[0])
	}

	p.SetCompressionThreshold(CompressionDisabled)
	uncompressed, err := p.Pack(obj)
	if err != nil {
		t.Fatal(err)
	}
	if uncompressed[0]&0x20 != 0 {
		t.Fatalf("expected uncompressed packet with ct=%d, got header 0x%02x", CompressionDisabled, uncompressed[0])
	}

	// 与抓包场景一致：1530 字节 payload 在 ct=2147483647 时不压缩
	capturedB64 := "gAX6EgADAAFjAgEAAWEDAA0AAXASAAMAAWMIAAlnYW1lTG9naW4AAXIE/////wABcBIAFwADdWlkCAAXNjI5ZGVtb2lucjAyMzExc2xvdEBBV0MACGdhbWVUeXBlBAAAAAgAC21hY2hpbmVUeXBlBAAAH1EABmJhbmtJZAgAAAAMc3RhcnRCYWxhbmNlBwAAAAAAAAAAAAVkZWJ1ZwEAAAdnYW1lVWlkCAAXNjI5ZGVtb2lucjAyMzExc2xvdEBBV0MACGdhbWVQYXNzCAAHYzAwZDExMQAIdXNlck5hbWUIABc2MjlkZW1vaW5yMDIzMTFzbG90QEFXQwAKc2Vzc2lvbklEMAgAAAAKc2Vzc2lvbklEMQgAAAAKc2Vzc2lvbklEMggAAAAKc2Vzc2lvbklEMwgEEENGNDQxNEMwREIxQzc4MTg2QTUzMjk2NzE2RDhFQzAxMjkxNDNDMjNBRkMxODRBMDcyQUVEMTJGQjc3NTNFODc4NEEwNTRGRjVCN0RBRTk4NUVGRDhDMUUxM0U5RjgyNjBFMTc3ODUzMTcyNENDMEEzMjEwMTQ3MzQ2MTczODEzMkUzMzlFRENBMUFBMDAxMkZBQzQ4OUM1OUI1MjAxMjdEOTgyQTg5NzZCQTQ1MzhCREZBQjg3MjY4RUE1RTYzNzREMjU5Q0MxOTA1NTkwOTlFMTIzNURGRTVCODRDQjQ0MDUwOEZBRUI0NEEzOTZCNDA4MkEzQjcyMEUyQTgyMUU0RjJFMUQ2QkI4NkRCMjU0QTEyRDcwODIxODY2NjE1RTFGMjRDMEU3NjdDNUREN0NGNzczMTI0NTUxRTYyNzFGOEQ3Q0NEQTJEMTYyQUREQ0EwNTU1MkJGQUM3N0Y1M0ZBQjc0MzJFQ0VFMEQ1RjQ4MjIwNEQ4RjE1Q0RBNDM1RTJEMTczQzkzMzlEMEExQzM1NDVGQUU2RDU1MEVBNEUwNjBCNUIwOTVCOUI5NTk2ODU4REEzOEE1RTYyNTMzMTI0QTkyOTFFOTI5RDc5RDNDNTEyRUNGMDI4NEI0M0EzRTZEMkNBOTdEQkQxOUVBNzcwMTY1NTU4RDkyMjA2MDgzRkZBNTJFQUM3RDlGNzkzOTdDRDcxMTAyMTE3NjNEQjJGODlBRkIxMTU2NTc3QUUyQjIxQzQ1RDgzMDlCMkQzNEUwNUM3NDkzOEU0RjY1RURDRjVCNjkwREEwQ0M4NkE5M0FGRDQ0NDZEMUNDMkU0NzczMjg4QzlEODVFQ0Y2OUUxNTBDRERCQjY1QzNEMTQ1NTczNjMzNzU3RDk0ODA2MzUxNTkzRjZBNkZEOTE5QTY2MzlCQzBDMzEzQkZFMzY0QTY4Qjc3NDhGNjNCQTdFREE5Q0VENjE3Q0M4MjcwN0NBRDIyNkMwOEYwNzkyNDM1QzNGMkJDRkE0MEUzRERFMUJFQzkwMzAyMUU2OUU3RjJEOENDRDhENUQ2N0ZCNEU4NTM5QThCODI3MTUwOTRDMEQzMjRCQjkwMjY4NjRGMEYzOThCQ0FGRTg2RkVGOUM2RURDMzdBMjVBMTMxOEE0QjRFRTNCOTJEM0IyN0YwNEMyMjA0QjA2N0Y2RURDNTM1RjJFRThBRDlDMTI4MzkxODcyNTQxQzA2OUNFQzI5QUNFRDRFMjE3Q0ZBQTIyNTMwMjhEOURFMkVFRkYwN0U4NTc0OTI1RDUyMEIxNjM2MUNDN0U2RjRFMzI1NEFBRkRDOTA1NEU2QzE3NDgyMTE1NjJCNEZDQjhBMTg5NUM4ODNGNzA0NzZFM0ZEOTY3MUZDAApzZXNzaW9uSUQ0CAAAAAZ1c2VTU0wBAQAIcGFzc3dvcmQIAAFhAApjbGllbnRUeXBlCAADV2ViAAF0CAAKamRiMTY4Lm5ldAANZ2FtZUxvZ2luTmFtZQgACWdhbWVMb2dpbgAEem9uZQgADUpEQl9aT05FX0dBTUUABHBvcnQEAAABuwAEaG9zdAgAD3N0MDMuamRiNzExLmNvbQAIem9uZU5hbWUIAA1KREJfWk9ORV9HQU1F"
	captured, err := base64.StdEncoding.DecodeString(capturedB64)
	if err != nil {
		t.Fatal(err)
	}
	if captured[0] != 0x80 {
		t.Fatalf("captured packet should be uncompressed, got 0x%02x", captured[0])
	}
	u := NewUnpacker(captured)
	v, err := u.Unpack()
	if err != nil {
		t.Fatal(err)
	}
	p.SetCompressionThreshold(CompressionDisabled)
	repacked, err := p.Pack(v.(SFSObject))
	if err != nil {
		t.Fatal(err)
	}
	if repacked[0] != 0x80 {
		t.Fatalf("repacked with server ct should stay uncompressed, got 0x%02x", repacked[0])
	}
}
