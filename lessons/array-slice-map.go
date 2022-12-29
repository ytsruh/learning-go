package lessons

import "fmt"

func CreateCheeseArray() [2]string {
	//Go arrays have to be instantiated with the length & type
	var cheeses [2]string
	cheeses[0] = "Halloumi"
	cheeses[1] = "Gouda"
	return cheeses
}

func CreateCheeseSlice() []string {
	//Use 'make' keyword to create a slice. Slice is similar to array but can be added to.
	var cheeses = make([]string, 2)
	cheeses[0] = "Halloumi"
	cheeses[1] = "Gouda"
	return cheeses
}

func RemoveSlow(slice []string, s int) []string {
    return append(slice[:s], slice[s+1:]...)
}

func RemoveFast(s []string, i int) []string {
	// This is faster but does not retain the order of the slice
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
}

func CopySlice(original []string)[]string{
	var subSlice = make([]string,2)
	copy(subSlice,original)
	return subSlice
}

func CreateMay() {
	var players = make(map[string] int)
	players["Hurst"] = 36
	players["Bairstow"] = 40
	fmt.Println(players)
	delete(players, "Bairstow")
	fmt.Println(players)
}

/*
func main()  {
	// ary := createCheeseArray()
	// fmt.Println(ary)
	// slc := createCheeseSlice()
	// slc = append(slc, "Philly")
	// slc = append(slc, "Mozzarella","Brie")
	// fmt.Println(slc)
	// fmt.Println("Length of cheeses is : ",len(slc))
	// removedSlow := removeSlow(slc, 1)
	// fmt.Println(removedSlow)
	// removedFast := removeFast(slc,3)
	// fmt.Println(removedFast)
	// fmt.Println(copySlice(slc))
	//createMay()
}
*/