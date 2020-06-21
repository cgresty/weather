package main

import (
	"fmt"
	"go.gresty.dev/weather/common"
)

func main() {
	p := common.DefaultProvider()
	var f = p.OneDayForecast(0)
	fmt.Printf("%s\n", f.Description)
}
