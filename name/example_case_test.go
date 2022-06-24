package name_test

import (
	"fmt"

	"github.com/QisFj/godry/name"
)

func ExampleToCamelCase() {
	for _, s := range []string{
		"a",
		"aa",
		"aa_aa",
		"http_request",
		"battery_life_value",
		"id0_value",
	} {
		fmt.Println(name.ToCamelCase(s))
	}
	// Output:
	// A
	// Aa
	// AaAa
	// HttpRequest
	// BatteryLifeValue
	// Id0Value
}

func ExampleToSnakeCase() {
	for _, s := range []string{
		"A",
		"AA",
		"AaAa",
		"HTTPRequest",
		"BatteryLifeValue",
		"Id0Value",
		"ID0Value",
		"UserID",
		"User.ID",
		"User.Name",
	} {
		fmt.Println(name.ToSnakeCase(s))
	}
	// Output:
	// a
	// aa
	// aa_aa
	// http_request
	// battery_life_value
	// id0_value
	// id0_value
	// user_id
	// user.id
	// user.name
}

func ExampleToSnakeCase_badcase() {
	fmt.Println(name.ToSnakeCase("UserIDs"))
	// Output: user_i_ds
}
