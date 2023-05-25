package lessons

import (
	"errors"
)

type Robot interface {
	PowerOn() error
}

type T850 struct {
	Name string
}

func (a *T850) PowerOn() error  {
	return nil
}

type R2D2 struct {
	Broken bool
}

func (r *R2D2) PowerOn() error  {
	if r.Broken {
		return errors.New("R2D2 is broken!")
	} else {
		return nil
	}
}

func Boot (r Robot) error  {
	return r.PowerOn()
}

// func main()  {
// 	sum := types.Addition(7,8)
// 	fmt.Println(sum)

// 	t := T850 {
// 		Name : "Terminator",
// 	}

// 	r := R2D2 {
// 		Broken : true,
// 	}

// 	err := Boot(&r)

// 	if err != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println("Robot is powered on")
// 	}

// 	err2 := Boot(&t)

// 	if err2 != nil {
// 		fmt.Println(err)
// 	} else {
// 		fmt.Println("Robot is powered on")
// 	}
// }