package di_demo2_test

import (
	"log"
	"os"
	"testing"

	"learning/test-with-go/di_demo2"
)

func TestThing_SomeFunc(t *testing.T) {
	var thing di_demo2.Thing
	thing.SomeFunc()

	thing = di_demo2.Thing{
		Logger: log.New(os.Stdout, "prefix:", log.Llongfile),
	}
	thing.SomeFunc()
}
