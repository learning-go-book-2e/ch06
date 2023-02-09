# Exercise 2

## Question
Write two functions. The `UpdateSlice` function takes in a `[]string` and a `string`. It sets the last position in the passed-in slice to the passed-in `string`. Print the slice after making the change. The `GrowSlice` function also takes in a `[]string` and a `string`. It appends the `string` onto the slice. Print the slice after making the change. Call these functions from `main`. Print out the slice before each function is called and after each function is called.

## Solution

The logic to update the last element of a slice to a new value isn't difficult, but unlike Python, there's no shortcut built into the language.

```go
func UpdateSlice(s []string, val string) {
	s[len(s)-1] = val
	fmt.Println("in UpdateSlice:", s)
}

func GrowSlice(s []string, val string) {
	s = append(s, val)
	fmt.Println("in GrowSlice:", s)
}

func main() {
	s := []string{"a", "b", "c"}
	UpdateSlice(s, "d")
	fmt.Println("in main after UpdateSlice:", s)
	GrowSlice(s, "e")
	fmt.Println("in main, after GrowSlice:", s)
}
```

Running this code produces:

```shell
$ go build
$ ./ex2 
in UpdateSlice: [a b d]
in main after UpdateSlice: [a b d]
in GrowSlice: [a b d e]
in main, after GrowSlice: [a b d]
```

You can see how changes to a slice's contents are visible after it is passed to a function, but changes that extend past its length are not.