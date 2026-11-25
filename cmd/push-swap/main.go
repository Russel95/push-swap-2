package main

import (
	"fmt"
	"os"

	"pushswap/piles"
)

func main() {
	nombres, erreur := piles.LireNombres(os.Args[1:])
	if erreur != nil {
		fmt.Fprintln(os.Stderr, "Error")
		os.Exit(1)
	}
	for _, instruction := range piles.Trier(nombres) {
		fmt.Println(instruction)
	}
}
