package name

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type caseTest struct {
	camel, snake string
}

func TestToCamelCase(t *testing.T) {
	for _, c := range []caseTest{
		{snake: "", camel: ""},
		{snake: "a", camel: "A"},
		{snake: "aa", camel: "Aa"},
		{snake: "aa_aa", camel: "AaAa"},
		{snake: "http_request", camel: "HttpRequest"},
		{snake: "battery_life_value", camel: "BatteryLifeValue"},
		{snake: "id0_value", camel: "Id0Value"},
	} {
		require.Equal(t, ToCamelCase(c.snake), c.camel)
		require.Equal(t, ToCamelCase(c.camel), c.camel)
	}
}

func TestToSnakeCase(t *testing.T) {
	for _, c := range []caseTest{
		{camel: "", snake: ""},
		{camel: "A", snake: "a"},
		{camel: "AA", snake: "aa"},
		{camel: "AaAa", snake: "aa_aa"},
		{camel: "HTTPRequest", snake: "http_request"},
		{camel: "BatteryLifeValue", snake: "battery_life_value"},
		{camel: "Id0Value", snake: "id0_value"},
		{camel: "ID0Value", snake: "id0_value"},
		{camel: "UserID", snake: "user_id"},
		// ! camel: UserIDs -> snake: user_i_ds
		{camel: "User.ID", snake: "user.id"},
		{camel: "User.Name", snake: "user.name"},
	} {
		require.Equal(t, c.snake, ToSnakeCase(c.snake))
		require.Equal(t, c.snake, ToSnakeCase(c.camel))
	}
}
