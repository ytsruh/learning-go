package lessons

import (
	"fmt"
	"math"
	"strconv"
)

type Moviay struct {
	Name string
	Rating float64
}

// This is a method. Its the same as a function but it has another parameter & is linked to a Movie struct
func (m *Moviay) summary() string  {
	r := strconv.FormatFloat(m.Rating, 'f', 1, 64)
	return m.Name + ", " + r
}

type Sphere struct {
	Radius float64
}

func (s *Sphere) SurfaceArea() float64  {
	return float64((4) * math.Pi * (s.Radius * s.Radius))
}

func (s *Sphere) Volume() float64  {
	radiusCubed := s.Radius * s.Radius * s.Radius
	return (float64(4) / float64(3)) * math.Pi * radiusCubed
}

type Triangle struct {
	base float64
	height float64
}

func (t *Triangle) area() float64  {
	return 0.5 * t.base * t.height
}

// This method is passed a value reference instead of a pointer reference (asterix above). This method will operate on a copy of the Triangle so will not update the original.
// If you need to modify / mutate the original instantiation then use a pointer. If you need to operate on a struct but do not want to modify the original instantization then use a value
func (t Triangle) changeBase(f float64)  {
	t.base = f
	return
}

// An interface is a blueprint for a method set but does not implement them

func main()  {
	m := Moviay {
		Name: "Lost",
		Rating: 3.22,
	}
	fmt.Println(m.summary())

	// sphere := Sphere {
	// 	Radius : 5,
	// }
	// fmt.Println(sphere.SurfaceArea())
	// fmt.Println(sphere.Volume())

	// t := Triangle{base: 3, height: 1}
	// fmt.Println(t.area())
	// t.changeBase(4)
	// fmt.Println(t.base)
}