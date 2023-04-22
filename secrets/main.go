package secrets

import (
	"fmt"
)

func RunInMemory() {
	v := InMemory("my-fake-key")
	err := v.Set("demo-key", "some value")
	if err != nil {
		panic(err)
	}
	plain, err := v.Get("demo-key")
	if err != nil {
		panic(err)
	}
	fmt.Println("Plain: ", plain)
}
