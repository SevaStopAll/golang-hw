package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	err := Copy("D:\\GolandProjects\\golang-hw\\hw07_file_copying\\testdata\\out_offset0_limit10.txt", "D:\\GolandProjects\\golang-hw\\hw07_file_copying\\testdata\\test.txt", 0, 1024)
	if err != nil {
		fmt.Println("error")
	}
}
