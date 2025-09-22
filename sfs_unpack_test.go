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
	hexStr := "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"
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
		bytes, err := NewPacker().Pack(rsp, false)
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
	v, err := packer.Pack(sendData, false)

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
