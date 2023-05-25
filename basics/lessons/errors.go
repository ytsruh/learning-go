package lessons

import "fmt"

func Half(number int) (int,error ) {
	if number%2 != 0 {
		return -1, fmt.Errorf("Cannot half %v",number)
	}
	return number / 2 ,nil
}

func Panic(words string) string {
	fmt.Printf("This works")
	panic("oh no we do the panic")

}