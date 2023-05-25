package lessons

type Movie struct {
	Name string
	Rating float32
}

type SuperHero struct {
	Name string
	Age	int
	Address Address
}

type Address struct {
	Number int
	Street string
	City string
}

type Alarm struct {
	Time string
	Sound string
}

func NewAlarm(time string) Alarm {
	a := Alarm {
		Time: time,
		Sound: "klaxon",
	} 
	return a
}

type Drink struct {
	Name string
	Ice bool 
}

// func main()  {
// 	sum := lessons.Addition(7,8)
// 	fmt.Println(sum)
// 	m := Movie{
// 		Name: "citizen kane",
// 		Rating: 10,
// 	}
// 	fmt.Println(m.Name, m.Rating)
// 	m.Rating = 5
// 	fmt.Println(m.Name, m.Rating)
// 	a := NewAlarm("7:00")
// 	b := NewAlarm("8:00")
// 	if b == a {
// 		fmt.Println("B & A are the same")
// 	} else {
// 		fmt.Println("B & A are not the same")
// 	}
// 	fmt.Println(reflect.TypeOf(a))

// 	c := Drink {
// 		Name: "Lemonade",
// 		Ice: true,
// 	}

// 	d := c // Value reference
// 	d.Ice = false
// 	fmt.Printf("%+v\n",d)
// 	fmt.Printf("%+v\n",c)
// 	fmt.Printf("%p\n",&c)
// 	fmt.Printf("%p\n",&d)
// 	fmt.Println("-------")
// 	e := &c // Pointer reference - the memory assigned to c is also updated when e is updated. c & e will have same memory address
// 	e.Ice = false
// 	fmt.Printf("%+v\n",e)
// 	fmt.Printf("%+v\n",c)
// 	fmt.Printf("%p\n",&c)
// 	fmt.Printf("%p\n",&e)
// }