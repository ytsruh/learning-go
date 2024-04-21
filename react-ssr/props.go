package main

import (
	"encoding/json"
	"fmt"
)

type PageProps struct {
	Props map[string]interface{}
}

func CreateProps(props map[string]interface{}) string {
	initialProps := PageProps{
		Props: props,
	}
	jsonProps, err := json.Marshal(initialProps.Props)
	if err != nil {
		fmt.Println("Error marshalling props")
		return ""
	}
	return string(jsonProps)
}
