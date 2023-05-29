package main

import (
	"bufio"
	"bytes"
	"os"

	"github.com/nlundbo/sqltostruct/sqltostruct"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var blob bytes.Buffer
	for scanner.Scan() {
		t := scanner.Text()
		blob.WriteString(t)

	}
	sqltostruct.Convert(blob.String())
}
