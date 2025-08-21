package httpx

import (
	"bytes"
	"fmt"
	"time"
)

func ExampleEmptyReader() {
	r := EmptyReader

	buf, err := readAll(r)
	fmt.Printf("test: emptyReader() -> [buf:%v] [err:%v]\n", string(buf), err)

	//Output:
	//test: emptyReader() -> [buf:] [err:<nil>]

}

func ExampleReadAll() {
	s := "this is new content"
	r1 := bytes.NewReader([]byte(s))
	r2 := bytes.NewReader([]byte(s))

	go func() {
		buf, err := readAll(r1)
		fmt.Printf("test: readAll(r1) -> [buf:%v] [err:%v]\n", string(buf), err)
	}()

	buf, err := readAll(r2)
	fmt.Printf("test: readAll(r2) -> [buf:%v] [err:%v]\n", string(buf), err)

	time.Sleep(time.Second * 2)
	//Output:
	//test: readAll(r2) -> [buf:this is new content] [err:<nil>]
	//test: readAll(r1) -> [buf:this is new content] [err:<nil>]

}
