package compute

import (
	s "../strategy"
	"fmt"
)

type Computer struct {
	Num1,Num2 int
	strate s.Strategier
}

func (p *Computer) SetStrategy(strate s.Strategier) {
	p.strate = strate
}

func (p Computer) Do() int{
	defer func() {
		if f:=recover();f!=nil {
			fmt.Println(f)
		}
	}()

	if p.strate == nil {
		panic("Strategy is null")
	}

	return p.strate.Compute(p.Num1,p.Num2)
}
