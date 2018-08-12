# basenconv
---

## Introduction

`basenconv` extends the functionality of the `FormatUint` and `ParseUint` functions in the Go `strconv` package. It extends the functionality by allowing you to format a number in `base-n`, so long as you can provide a character set of length `n`.

The library includes convenience methods to format and parse values to and from base62, hex, octal and binary.

## Usage

```go
package main

import(
    "github.com/w-k-s/basenconv"
    "fmt"
)

func main() {
	fmt.Println(basenconv.FormatBase62(20000))  //Prints 5Ca
	fmt.Println(basenconv.FormatHex(20000))     //Prints 4E20
	fmt.Println(basenconv.FormatOctal(20000))   //Prints 47040
	fmt.Println(basenconv.FormatBinary(20000))  //Prints 100111000100000
	fmt.Println(basenconv.FormatUint(2, "ABC")) // Prints C

	fmt.Println(basenconv.ParseBase62("5Ca"))             //Prints 20000
	fmt.Println(basenconv.ParseHex("4E20"))               //Prints 20000
	fmt.Println(basenconv.ParseOctal("47040"))            //Prints 20000
	fmt.Println(basenconv.ParseBinary("100111000100000")) //Prints 20000
	fmt.Println(basenconv.ParseUint("B", "ABC"))          // Prints 1
}
```

## Author

* [W.K.S](https://stackoverflow.com/users/821110/w-k-s)
