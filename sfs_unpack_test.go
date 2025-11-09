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
	hexStr := "gADYEgADAAFjAgEAAWEDAA0AAXASAAMAAWMIAAdoNS5zcGluAAFyBP////8AAXASAAEABmVudGl0eRIABgAFZGVub20IAAIyMAAMZXh0cmFCZXRUeXBlCAAKTm9FeHRyYUJldAALZ2FtZVN0YXRlSWQIAAEwAAlwbGF5ZXJCZXQIAAI1MAAOYnV5RmVhdHVyZVR5cGUIAARudWxsAApiZXRSZXF1ZXN0EgADAAdiZXRUeXBlCAAHV2F5R2FtZQAJYmV0Q29sdW1uBAAAAAUABndheUJldAQAAAAB"
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
