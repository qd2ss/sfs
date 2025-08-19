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

func TestUpack(t *testing.T) {
	hexStr := "gACIEgADAAFwEgACAAFwEgAGAAlsb2dpblJvb20IAAlTTE9UX1JPT00ABGRhdGEBAQAHYmFsYW5jZQdAwVwAAAAAAAAIdGVzdE1vZGUBAAAIc2VydmVySWQIAAIwMQACdHMHQnmLvH0ZwAAAAWMIAA9nYW1lTG9naW5SZXR1cm4AAWEHQCoAAAAAAAAAAWMHP/AAAAAAAAA="
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
	rsp, err := Marshal(Respond{
		Code: 200,
		Data: Data{Index: 1},
		Msg:  "success",
	})
	if err == nil {
		t.Logf("pack by struct: %v", rsp)
		bytes, err := NewPacker().Pack(rsp, false)
		if err == nil {
			t.Logf("pack by struct: %v", bytes)
		}
	}

	data := &Respond{}
	err = Unmarshal(rsp, data)
	if err == nil {
		t.Logf("pack by struct: %v", data)
	} else {
		t.Logf("pack by struct: %v", err)
	}

}
