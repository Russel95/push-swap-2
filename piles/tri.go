package piles

func positionMinimum(pile []int) int {
	position := 0
	for i, nombre := range pile {
		if nombre < pile[position] {
			position = i
		}
	}
	return position
}

func Trier(nombres []int) []string {
	if EstTriee(nombres) {
		return nil
	}
	if len(nombres) <= 6 {
		return cheminLePlusCourt(nombres)
	}

	p := Piles{A: nombres}
	var instructions []string
	jouer := func(instruction string) {
		p.Appliquer(instruction)
		instructions = append(instructions, instruction)
	}

	echanger(p.A)
	echangeSuffisant := EstTriee(p.A)
	echanger(p.A)
	if echangeSuffisant {
		return []string{"sa"}
	}

	position := positionMinimum(p.A)
	if EstTriee(p.A[:position]) && EstTriee(p.A[position:]) && p.A[len(p.A)-1] < p.A[0] {
		if position <= len(p.A)/2 {
			for i := 0; i < position; i++ {
				jouer("ra")
			}
		} else {
			for i := position; i < len(p.A); i++ {
				jouer("rra")
			}
		}
		return instructions
	}

	trierParBits(&p, jouer)
	return instructions
}

func trierParBits(p *Piles, jouer func(string)) {
	minimum := p.A[positionMinimum(p.A)]
	var bitsUtiles uint
	for _, nombre := range p.A {
		bitsUtiles |= uint(nombre) - uint(minimum)
	}

	for masque := uint(1); masque != 0 && masque <= bitsUtiles; masque <<= 1 {
		if bitsUtiles&masque == 0 {
			continue
		}
		taille := len(p.A)
		for i := 0; i < taille; i++ {
			ecart := uint(p.A[0]) - uint(minimum)
			if ecart&masque == 0 {
				jouer("pb")
			} else {
				jouer("ra")
			}
		}
		for len(p.B) > 0 {
			jouer("pa")
		}
		if EstTriee(p.A) {
			return
		}
	}
}
