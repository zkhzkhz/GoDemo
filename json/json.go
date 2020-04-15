package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func main() {
	type person struct {
		ID      int    `json:"id" bson:"_id"`
		Name    string `json:"name"`
		Country string `json:"country"`
	}

	zkh := person{1001, "zhaokh", "cn"}

	t := reflect.TypeOf(zkh)

	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)
		fmt.Println(sf.Tag.Get("json"))
		fmt.Println(sf.Tag.Get("bson"))
	}

	fmt.Println("zkh=", zkh)

	result, err := json.Marshal(zkh)
	if err != nil {
		fmt.Println("encoding failed...")
	}
	fmt.Println("JSON RESULT =", result)
	fmt.Println("JSON RESULT(string)=", string(result))

	jdata := []byte(`{"Id":1001,"Name":"liumiao","Country":"China"}`)
	var zkhh person

	errr := json.Unmarshal(jdata, &zkhh)
	if errr != nil {
		fmt.Println("decoding failed")
	}
	fmt.Println("go struct=", zkhh)
}
