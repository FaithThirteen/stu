package test

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {

	l := DefaultLru()

	l.Set("key1","value1")
	l.Set("key2","value2")
	l.Set("key3","value3")
	l.Set("key4","value4")
	l.Set("key5","value5")

	l.Remove("key4")

	fmt.Println(l.Get("key1"))
	fmt.Println(l.Get("key2"))
	fmt.Println(l.Get("key3"))
	fmt.Println(l.Get("key4"))
	fmt.Println(l.Get("key5"))

	/**
		=== RUN   TestName
		<nil> <nil>
		value2 <nil>
		value3 <nil>
		<nil> <nil>
		value5 <nil>
		--- PASS: TestName (0.00s)
		PASS
	 */




}