package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func main() {
	for {
		data, err := ioutil.ReadFile("/vault/secrets/mysecret.txt")
		if err != nil {
			log.Println("Error reading secret:", err)
		} else {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				fmt.Println(line)
			}
		}
		fmt.Println("-----")
		// sleep 5s
		time.Sleep(5 * time.Second)
	}
}
