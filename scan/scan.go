package main

import (
	"fmt"
	"os"
)

func main() {
	argv := os.Args
	list, errs, err := search(argv[1])
	if nil != err {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	for _, fErr := range errs {
		fmt.Printf("%s: %s\n", fErr.Filename, fErr.Error())
	}

	for _, item := range list {
		fmt.Println(item.ResourceMeta, item.AudioMeta)
	}
}
