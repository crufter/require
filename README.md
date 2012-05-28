Require
=======

The Require Go package allows you to do file inclusion in any file/string.

Usage
=======

```
package main

import(
	"github.com/opesun/require"
	"fmt"
)

func main() {
	v, err := require.RSimple("/usr", "index.html")
	if err != "" {
		panic(err)
	}
	// If index.html contains "{{require header.html}} text here {{require footer.html}}"
	// and header.html contains "header", footer.html contains "footer"
	// then
	fmt.Println(v)
	// v will be the string "header text here footer"
}
```

The RSimple method is a thin wrapper around the R method.
The R method needs a function with a signature of func(string) ([]byte,error) as its last parameter, where you can provide your very own function which gets the content of the required files,
so you can implement your own caching mechanism, and things like that.