package main

import "fmt"

func main() {
	firstExample()
	secondExample()
	thirdExample()
	fourthExample()
	fifthExample()
}

func firstExample() {
	var x int32 = 10
	var y bool = true
	pointerX := &x
	pointerY := &y
	var pointerZ *string

	fmt.Println(x, y, pointerX, pointerY, pointerZ)
}

func secondExample() {
	x := 10
	pointerToX := &x
	fmt.Println(pointerToX)  // prints a memory address
	fmt.Println(*pointerToX) // prints 10
	z := 5 + *pointerToX
	fmt.Println(z) // prints 15
}

func thirdExample() {
	// chapter 9 explains panic and recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	var x *int
	fmt.Println(x == nil) // prints true
	fmt.Println(*x)       // panics
}

func fourthExample() {
	var x = new(int)
	fmt.Println(x == nil) // prints false
	fmt.Println(*x)       // prints 0
}

type person struct {
	FirstName  string
	MiddleName *string
	LastName   string
}

func stringp(s string) *string {
	return &s
}

func fifthExample() {
	p := person{
		FirstName:  "Pat",
		MiddleName: stringp("Perry"), // This works
		LastName:   "Peterson",
	}
	fmt.Println(p)
}
