package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/TwoBlueCats/diceRolls"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data := scanner.Bytes()
		result, err := diceRolls.Parser(string(data))
		if err != nil {
			fmt.Println("get error ", err)
		} else {
			fmt.Printf("get result %v; description %v\n", result.Value(), result.Description())
		}
	}
}
