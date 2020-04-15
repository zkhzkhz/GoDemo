package strategy

type Strategier interface {
	Compute(num1, num2 int) int
}
