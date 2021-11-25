package meander

import "fmt"

type Facade interface {
	Public() interface{}
}

func Public(o interface{}) interface{} {
	if p, ok := o.(Facade); ok {
		return p.Public()
	}
	fmt.Println("hoghoee")
	return o
}
