package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	palabra := os.Args[1]
	output := Scrape(palabra)

	data, _ := json.MarshalIndent(output, "", "    ")
	fmt.Println(string(data))
}
