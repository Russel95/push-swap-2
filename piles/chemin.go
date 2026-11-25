package piles

import "fmt"

type essai struct {
	piles  Piles
	chemin []string
}

func cheminLePlusCourt(nombres []int) []string {
	depart := Piles{A: nombres}
	file := []essai{{piles: depart}}
	dejaVus := map[string]bool{fmt.Sprint(depart): true}
	mouvements := []string{"pb", "ra", "sa", "pa", "rra", "sb", "ss", "rb", "rr", "rrb", "rrr"}

	for indice := 0; indice < len(file); indice++ {
		actuel := file[indice]
		for _, mouvement := range mouvements {
			suivant := Piles{
				A: append([]int(nil), actuel.piles.A...),
				B: append([]int(nil), actuel.piles.B...),
			}
			suivant.Appliquer(mouvement)
			cle := fmt.Sprint(suivant)
			if dejaVus[cle] {
				continue
			}
			dejaVus[cle] = true
			chemin := append([]string(nil), actuel.chemin...)
			chemin = append(chemin, mouvement)
			if len(suivant.B) == 0 && EstTriee(suivant.A) {
				return chemin
			}
			file = append(file, essai{piles: suivant, chemin: chemin})
		}
	}
	return nil
}
