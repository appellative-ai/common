package core

import (
	"fmt"
)

type agent interface {
	Name() string
}

type newAgentFunc func() agent

type address struct {
	Address string
	City    string
	State   string
}

func ExampleNewMap() {
	m := NewSyncMap[string, newAgentFunc]()
	name1 := ""
	t, ok := m.Load("")
	fmt.Printf("test:  Load(\"%v\") -> %v [ok:%v]\n", name1, t, ok)

	name1 = "common:core:ctor/test"
	m.Store(name1, nil)
	t, ok = m.Load(name1)
	fmt.Printf("test:  Load(\"%v\") -> %v [ok:%v]\n", name1, t, ok)

	m.Store(name1, func() agent { return nil })
	t, ok = m.Load(name1)
	fmt.Printf("test:  Load(\"%v\") -> %v [ok:%v]\n", name1, t != nil, ok)

	name1 = "common:core:ctor/invalid"
	t, ok = m.Load(name1)
	fmt.Printf("test:  Load(\"%v\") -> %v [ok:%v]\n", name1, t != nil, ok)

	//Output:
	//test:  Load("") -> <nil> [ok:false]
	//test:  Load("common:core:ctor/test") -> <nil> [ok:true]
	//test:  Load("common:core:ctor/test") -> true [ok:true]
	//test:  Load("common:core:ctor/invalid") -> false [ok:false]

}

func ExampleNewMap_Address() {
	m := NewSyncMap[string, address]()
	name1 := ""
	t, ok := m.Load(name1)
	fmt.Printf("test:  Load(\"%v\") -> %v [ok:%v]\n", name1, t, ok)

	name1 = "common:core:ctor/test"
	m.Store(name1, address{Address: "123 Main", City: "Anytown", State: "Ohio"})
	t, ok = m.Load(name1)
	fmt.Printf("test:  Load(\"%v\") -> %v [ok:%v]\n", name1, t, ok)

	name1 = "common:core:ctor/invalid"
	t, ok = m.Load(name1)
	fmt.Printf("test:  Load(\"%v\") -> %v [ok:%v]\n", name1, t, ok)

	//Output:
	//test:  Load("") -> {  } [ok:false]
	//test:  Load("common:core:ctor/test") -> {123 Main Anytown Ohio} [ok:true]
	//test:  Load("common:core:ctor/invalid") -> {  } [ok:false]

}
