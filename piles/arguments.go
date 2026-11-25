package piles

import (
	"errors"
	"strconv"
	"strings"
)

func LireNombres(arguments []string) ([]int, error) {
	var nombres []int
	dejaVus := make(map[int]bool)

	for _, argument := range arguments {
		mots := strings.Fields(argument)
		if len(mots) == 0 {
			return nil, errors.New("argument vide")
		}

		for _, mot := range mots {
			nombre, erreur := strconv.Atoi(mot)
			if erreur != nil {
				return nil, erreur
			}
			if dejaVus[nombre] {
				return nil, errors.New("nombre en double")
			}
			dejaVus[nombre] = true
			nombres = append(nombres, nombre)
		}
	}
	return nombres, nil
}
