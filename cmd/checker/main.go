package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"pushswap/piles"
)

func verifier(nombres []int, entree io.Reader) string {
	p := piles.Piles{A: nombres}
	lecteur := bufio.NewReader(entree)
	for {
		ligne, erreur := lecteur.ReadString('\n')
		if erreur == io.EOF && ligne == "" {
			break
		}
		if erreur != nil || !p.Appliquer(strings.TrimSuffix(ligne, "\n")) {
			return "Error"
		}
	}
	if piles.EstTriee(p.A) && len(p.B) == 0 {
		return "OK"
	}
	return "KO"
}

func main() {
	if len(os.Args) == 1 {
		return
	}
	nombres, erreur := piles.LireNombres(os.Args[1:])
	if erreur != nil {
		fmt.Fprintln(os.Stderr, "Error")
		os.Exit(1)
	}
	resultat := verifier(nombres, os.Stdin)
	if resultat == "Error" {
		fmt.Fprintln(os.Stderr, resultat)
		os.Exit(1)
	}
	fmt.Println(resultat)
}
