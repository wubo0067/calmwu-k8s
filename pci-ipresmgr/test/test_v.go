package main

import "fmt"

var time string
var version string

func main() {
	fmt.Println(time)
	fmt.Println(version)
}

// CGO_ENABLED=1 go build  -ldflags="-X 'main.time=`date`' -X main.version=1.0.2 -linkmode external -extldflags -static" test_v.go
