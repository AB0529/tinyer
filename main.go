package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Config the configuration JSON structure
type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

var config Config

func main() {
	file, _ := ioutil.ReadFile("./config.json")
	json.Unmarshal(file, &config)

	fmt.Println(config.Host)
}
