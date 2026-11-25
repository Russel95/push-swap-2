package piles

type Piles struct {
	A []int
	B []int
}

func EstTriee(pile []int) bool {
	for i := 1; i < len(pile); i++ {
		if pile[i-1] > pile[i] {
			return false
		}
	}
	return true
}

func echanger(pile []int) {
	if len(pile) >= 2 {
		pile[0], pile[1] = pile[1], pile[0]
	}
}

func pousser(source, destination []int) ([]int, []int) {
	if len(source) == 0 {
		return source, destination
	}
	sommet := source[0]
	return source[1:], append([]int{sommet}, destination...)
}

func tourner(pile []int) []int {
	if len(pile) < 2 {
		return pile
	}
	return append(pile[1:], pile[0])
}

func tournerEnSensInverse(pile []int) []int {
	if len(pile) < 2 {
		return pile
	}
	dernier := len(pile) - 1
	return append([]int{pile[dernier]}, pile[:dernier]...)
}

func (p *Piles) Appliquer(instruction string) bool {
	switch instruction {
	case "sa":
		echanger(p.A)
	case "sb":
		echanger(p.B)
	case "ss":
		echanger(p.A)
		echanger(p.B)
	case "pa":
		p.B, p.A = pousser(p.B, p.A)
	case "pb":
		p.A, p.B = pousser(p.A, p.B)
	case "ra":
		p.A = tourner(p.A)
	case "rb":
		p.B = tourner(p.B)
	case "rr":
		p.A = tourner(p.A)
		p.B = tourner(p.B)
	case "rra":
		p.A = tournerEnSensInverse(p.A)
	case "rrb":
		p.B = tournerEnSensInverse(p.B)
	case "rrr":
		p.A = tournerEnSensInverse(p.A)
		p.B = tournerEnSensInverse(p.B)
	default:
		return false
	}
	return true
}
