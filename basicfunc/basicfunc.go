package basicfunc

func GetGrade(score int) string {
	switch {
	case score < 60:
		return "D"
	case score <= 70:
		return "C"
	case score <= 80:
		return "B"
	case score <= 90:
		return "A"
	default:
		return "Undefined"
	}
}

func Add(num1 int, num2 int) int {
	return num1 + num2
}
