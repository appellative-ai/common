package core

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
)

type Address struct {
	Line1 string
	Line2 string
	City  string
	State string
	Zip   string
}

func ExampleBinary() {
	bytes := []byte("this is a test buffer")
	// []byte -> []byte
	t, status := Unmarshal[[]byte](&Content{Type: ContentTypeBinary, Value: bytes})
	fmt.Printf("test: Unmarshal[[]byte]() -> [%v] [status:%v]\n", string(t), status)

	// []byte -> []byte
	t, status = Marshal[[]byte](&Content{Type: ContentTypeBinary, Value: bytes})
	fmt.Printf("test: Marshal[[]byte]() -> [%v] [status:%v]\n", string(t), status)

	// []byte -> io.Reader
	t2, status2 := Marshal[io.Reader](&Content{Type: ContentTypeBinary, Value: bytes})
	buf, err := io.ReadAll(t2)
	if err != nil {
		fmt.Printf("test: io.ReadAll() -> [err:%v]\n", err)
	}
	fmt.Printf("test: Marshal[io.Reader]() -> [%v] [status:%v]\n", string(buf), status2)

	//Output:
	//test: Unmarshal[[]byte]() -> [this is a test buffer] [status:<nil>]
	//test: Marshal[[]byte]() -> [this is a test buffer] [status:<nil>]
	//test: Marshal[io.Reader]() -> [this is a test buffer] [status:<nil>]

}

func ExampleString() {
	s := "this is a test string"

	// []byte -> string
	t, status := Unmarshal[string](&Content{Type: ContentTypeText, Value: []byte(s)})
	fmt.Printf("test: Unmarshal[string]() -> [%v] [status:%v]\n", t, status)

	// string -> []byte
	t2, status2 := Marshal[[]byte](&Content{Type: ContentTypeText, Value: s})
	fmt.Printf("test: Marshal[[]byte]() -> [%v] [status:%v]\n", string(t2), status2)

	// string -> io.Reader
	t3, status3 := Marshal[io.Reader](&Content{Type: ContentTypeText, Value: s})
	buf, err := io.ReadAll(t3)
	if err != nil {
		fmt.Printf("test: io.ReadAll() -> [err:%v]\n", err)
	}
	fmt.Printf("test: Marshal[io.Reader]() -> [%v] [status:%v]\n", string(buf), status3)

	//Output:
	//test: Unmarshal[string]() -> [this is a test string] [status:<nil>]
	//test: Marshal[[]byte]() -> [this is a test string] [status:<nil>]
	//test: Marshal[io.Reader]() -> [this is a test string] [status:<nil>]

}

func ExampleType() {
	addr := Address{
		Line1: "123 Main",
		Line2: "",
		City:  "Anytown",
		State: "Ohio",
		Zip:   "54321",
	}
	// []byte -> Address
	buf, err := json.Marshal(&addr)
	if err != nil {
		fmt.Printf("test: json.Marshal() -> [err:%v]\n", err)
	}
	t, status := Unmarshal[Address](&Content{Type: ContentTypeJson, Value: buf})
	fmt.Printf("test: Unmarshal[Address]() -> [%v] [status:%v]\n", t, status)

	// Address -> []byte
	t2, status2 := Marshal[[]byte](&Content{Type: ContentTypeJson, Value: addr})
	fmt.Printf("test: Marshal[Address]() -> [%v] [status:%v]\n", string(t2), status2)

	// Address -> io.Reader
	t3, status3 := Marshal[io.Reader](&Content{Type: ContentTypeJson, Value: addr})
	buf, err = io.ReadAll(t3)
	if err != nil {
		fmt.Printf("test: io.ReadAll() -> [err:%v]\n", err)
	}
	fmt.Printf("test: Marshal[io.Reader]() -> [%v] [status:%v]\n", string(buf), status3)

	//Output:
	//test: Unmarshal[Address]() -> [{123 Main  Anytown Ohio 54321}] [status:<nil>]
	//test: Marshal[Address]() -> [{"Line1":"123 Main","Line2":"","City":"Anytown","State":"Ohio","Zip":"54321"}] [status:<nil>]
	//test: Marshal[io.Reader]() -> [{"Line1":"123 Main","Line2":"","City":"Anytown","State":"Ohio","Zip":"54321"}] [status:<nil>]

}

func ExampleMap() {
	m := map[string]string{
		"Line1": "123 Main",
		"Line2": "",
		"City":  "Anytown",
		"State": "Ohio",
		"Zip":   "54321",
	}
	// []byte -> map[string]string
	buf, err := json.Marshal(&m)
	if err != nil {
		fmt.Printf("test: json.Marshal() -> [err:%v]\n", err)
	}
	t, status := Unmarshal[map[string]string](&Content{Type: ContentTypeJson, Value: buf})
	fmt.Printf("test: Unmarshal[map[string]string]() -> [%v] [status:%v]\n", t, status)

	// map[string]string -> []byte
	t2, status2 := Marshal[[]byte](&Content{Type: ContentTypeJson, Value: m})
	fmt.Printf("test: Marshal[[]byte]() -> [%v] [status:%v]\n", string(t2), status2)

	// map[string]string -> io.Reader
	t3, status3 := Marshal[io.Reader](&Content{Type: ContentTypeJson, Value: m})
	buf, err = io.ReadAll(t3)
	if err != nil {
		fmt.Printf("test: io.ReadAll() -> [err:%v]\n", err)
	}
	fmt.Printf("test: Marshal[io.Reader]() -> [%v] [status:%v]\n", string(buf), status3)

	//Output:
	//test: Unmarshal[map[string]string]() -> [map[City:Anytown Line1:123 Main Line2: State:Ohio Zip:54321]] [status:<nil>]
	//test: Marshal[[]byte]() -> [{"City":"Anytown","Line1":"123 Main","Line2":"","State":"Ohio","Zip":"54321"}] [status:<nil>]
	//test: Marshal[io.Reader]() -> [{"City":"Anytown","Line1":"123 Main","Line2":"","State":"Ohio","Zip":"54321"}] [status:<nil>]

}

func _ExampleReader() {
	//bytes := []byte("this is a test io.Reader")
	m := map[string]string{
		"Line1": "123 Main",
		"Line2": "",
		"City":  "Anytown",
		"State": "Ohio",
		"Zip":   "54321",
	}
	buf, err := json.Marshal(&m)
	if err != nil {
		fmt.Printf("test: json.Marshal() -> [err:%v]\n", err)
	}

	t, status := Unmarshal[io.Reader](&Content{Type: ContentTypeJson, Value: buf})
	buf2, err2 := io.ReadAll(t)
	if err2 != nil {
		fmt.Printf("test: io.ReadAll() -> [err:%v]\n", err2)
	}
	fmt.Printf("test: Unmarshal[io.Reader]() -> [%v] [status:%v]\n", string(buf2), status)

	// JSON -> []byte
	t2, status2 := Marshal[[]byte](&Content{Type: ContentTypeJson, Value: m})
	fmt.Printf("test: Marshal[[]byte]() -> [%v] [status:%v]\n", string(t2), status2)

	// JSON -> io.Reader
	t3, status3 := Marshal[io.Reader](&Content{Type: ContentTypeJson, Value: m})
	buf2, err2 = io.ReadAll(t3)
	if err2 != nil {
		fmt.Printf("test: io.ReadAll() -> [err:%v]\n", err2)
	}
	fmt.Printf("test: Marshal[io.Reader]() -> [%v] [status:%v]\n", string(buf2), status3)

	// []byte -> []byte
	t4, status4 := Marshal[[]byte](&Content{Type: ContentTypeBinary, Value: buf})
	fmt.Printf("test: Marshal[[]byte]() -> [%v] [status:%v]\n", string(t4), status4)

	// []byte -> io.Reader
	t5, status5 := Marshal[io.Reader](&Content{Type: ContentTypeBinary, Value: buf})
	buf2, err2 = io.ReadAll(t5)
	if err2 != nil {
		fmt.Printf("test: io.ReadAll() -> [err:%v]\n", err2)
	}
	fmt.Printf("test: Marshal[io.Reader]() -> [%v] [status:%v]\n", string(buf2), status5)

	//Output:
	//test: Unmarshal[io.Reader]() -> [{"City":"Anytown","Line1":"123 Main","Line2":"","State":"Ohio","Zip":"54321"}] [status:<nil>]
	//test: Marshal[[]byte]() -> [{"City":"Anytown","Line1":"123 Main","Line2":"","State":"Ohio","Zip":"54321"}] [status:<nil>]
	//test: Marshal[io.Reader]() -> [{"City":"Anytown","Line1":"123 Main","Line2":"","State":"Ohio","Zip":"54321"}] [status:<nil>]
	//test: Marshal[[]byte]() -> [{"City":"Anytown","Line1":"123 Main","Line2":"","State":"Ohio","Zip":"54321"}] [status:<nil>]
	//test: Marshal[io.Reader]() -> [{"City":"Anytown","Line1":"123 Main","Line2":"","State":"Ohio","Zip":"54321"}] [status:<nil>]

}

func ExampleNew_Type() {
	addr := Address{
		Line1: "123 Main",
		Line2: "",
		City:  "Anytown",
		State: "Ohio",
		Zip:   "54321",
	}
	// Address type
	cnt := Content{Type: ContentTypeJson, Value: addr}
	t, status := New[Address](&cnt)
	fmt.Printf("test: New[Address]() -> %v [type:%v] [status:%v]\n", t, reflect.TypeOf(t), status)

	t2, status2 := New[map[string]string](&cnt)
	fmt.Printf("test: New[map[string]string]() -> %v [type:%v] [status:%v]\n", t2, reflect.TypeOf(t2), status2)

	// String type
	cnt = Content{Type: ContentTypeText, Value: "this is text content"}
	t3, status3 := New[string](&cnt)
	fmt.Printf("test: New[string]() -> %v [type:%v] [status:%v]\n", t3, reflect.TypeOf(t3), status3)

	//Output:
	//test: New[Address]() -> {123 Main  Anytown Ohio 54321} [type:core.Address] [status:<nil>]
	//test: New[map[string]string]() -> map[] [type:map[string]string] [status:content value type: core.Address is not of generic type: map[string]string]
	//test: New[string]() -> this is text content [type:string] [status:<nil>]

}

func ExampleNew_Binary() {
	addr := Address{
		Line1: "123 Main",
		Line2: "",
		City:  "Anytown",
		State: "Ohio",
		Zip:   "54321",
	}
	buf, err := json.Marshal(&addr)
	if err != nil {
		fmt.Printf("test: json.Marshal() -> [err:%v]\n", err)
	}
	// Address
	cnt := Content{Type: ContentTypeJson, Value: buf}
	t, status := New[Address](&cnt)
	fmt.Printf("test: New[Address]() -> %v [type:%v] [status:%v]\n", t, reflect.TypeOf(t), status)

	// String
	cnt = Content{Type: ContentTypeText, Value: []byte("this is a test string")}
	t2, status2 := New[string](&cnt)
	fmt.Printf("test: New[string]() -> %v [type:%v] [status:%v]\n", t2, reflect.TypeOf(t2), status2)

	// Binary
	cnt = Content{Type: ContentTypeBinary, Value: []byte("this is a test string")}
	t3, status3 := New[[]byte](&cnt)
	fmt.Printf("test: New[[]byte]() -> %v [type:%v] [status:%v]\n", string(t3), reflect.TypeOf(t3), status3)

	//Output:
	//test: New[Address]() -> {123 Main  Anytown Ohio 54321} [type:core.Address] [status:<nil>]
	//test: New[string]() -> this is a test string [type:string] [status:<nil>]
	//test: New[[]byte]() -> this is a test string [type:[]uint8] [status:<nil>]

}
