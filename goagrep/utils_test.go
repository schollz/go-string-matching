package goagrep

import "fmt"

func ExampleUtils1() {
	fmt.Println(removeDuplicates([]int{1, 2, 3, 4, 5, 5, 5, 6, 7, 7, 8, 10}))
	// Output: [1 2 3 4 5 6 7 8 10]
}

func ExampleUtils2() {
	fmt.Println(abs(-1), abs(30))
	// Output: 1 30
}

func ExampleUtils3() {
	fmt.Println(stringInSlice("hello", []string{"hello", "hell", "hello there"}), stringInSlice("hello", []string{"hell", "hello there"}))
	// Output: true false
}

func ExampleUtils4() {
	fmt.Println(lineCount("../example/testlist"))
	// Output: 1007
}
