package binding

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func newRequest(values url.Values) *http.Request {
	return &http.Request{URL: &url.URL{RawQuery: values.Encode()}}
}

func TestQueryBinding(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		type Req struct {
			Int     int
			Uint    uint
			Float   float64
			String  string
			Boolean bool
		}
		var req Req
		err := QueryBinding{}.Bind(newRequest(url.Values{
			"int":     {"1"},
			"uint":    {"1"},
			"float":   {"1"},
			"string":  {"1"},
			"boolean": {"true"},
		}), &req)
		require.NoError(t, err)
		require.Equal(t, Req{
			Int:     1,
			Uint:    1,
			Float:   1,
			String:  "1",
			Boolean: true,
		}, req)
	})
	t.Run("struct", func(t *testing.T) {
		type structure struct {
			Int     int
			Uint    uint
			Float   float64
			String  string
			Boolean bool
			inner   int // ignore
		}
		type Req struct {
			structure // anonymous, can be accessed
			Nested    structure
		}
		var req Req
		err := QueryBinding{}.Bind(newRequest(url.Values{
			"int":            {"1"},
			"uint":           {"1"},
			"float":          {"1"},
			"string":         {"1"},
			"boolean":        {"true"},
			"inner":          {"1"},
			"nested.int":     {"1"},
			"nested.uint":    {"1"},
			"nested.float":   {"1"},
			"nested.string":  {"1"},
			"nested.boolean": {"true"},
			"nested.inner":   {"1"},
		}), &req)
		require.NoError(t, err)
		expStruct := structure{
			Int:     1,
			Uint:    1,
			Float:   1,
			String:  "1",
			Boolean: true,
			inner:   0, // still 0
		}
		require.Equal(t, expStruct, req.structure)
		require.Equal(t, expStruct, req.Nested)
	})
	t.Run("multiple values", func(t *testing.T) {
		// use the first
		// for no value, never happen
		type Req struct{ Value int }
		var req Req
		err := QueryBinding{}.Bind(newRequest(url.Values{
			"value": {"1", "2"},
		}), &req)
		require.NoError(t, err)
		require.Equal(t, Req{Value: 1}, req)
	})
	t.Run("simple slice", func(t *testing.T) {
		type Req struct{ Values []int }
		var req Req
		err := QueryBinding{}.Bind(newRequest(url.Values{
			"values": {"1", "2"},
		}), &req)
		require.NoError(t, err)
		require.Equal(t, Req{Values: []int{1, 2}}, req)
	})
	t.Run("name", func(t *testing.T) {
		// name would be convert to snake case
		// if converted name duplicated:
		// for fields, set value for the first field
		// for url.Values, merge them with a random order (because url.Values doesn't have order)
		type Req1 struct {
			IntValues0  []int
			Int_Values0 []int // duplicated, and not first, won't set
		}
		type Req2 struct {
			Int_Values0 []int
			IntValues0  []int // duplicated, and not first, won't set
		}
		var req1 Req1
		var req2 Req2
		req := newRequest(url.Values{
			"int_values0": {"1", "2"},
			"IntValues0":  {"3", "4"}, // duplicate, merge
		})
		checheValue := func(values []int) {
			// values should be 1,2,3,4 or 3,4,1,2
			require.Len(t, values, 4)
			require.Equal(t, 4, values[0]+values[2])
			require.Equal(t, 6, values[1]+values[3])
		}
		err := QueryBinding{}.Bind(req, &req1)
		require.NoError(t, err)
		checheValue(req1.IntValues0)
		require.Nil(t, req1.Int_Values0)

		err = QueryBinding{}.Bind(req, &req2)
		require.NoError(t, err)
		checheValue(req2.Int_Values0)
		require.Nil(t, req2.IntValues0)
	})
	t.Run("ptr", func(t *testing.T) {
		type Req struct {
			PValue  *int
			PPValue **int
			PsValue []*int
			PValues *[]int
		}
		var req Req
		err := QueryBinding{}.Bind(newRequest(url.Values{
			"p_value":  {"1"},
			"pp_value": {"1"},
			"ps_value": {"1", "2"},
			"p_values": {"1", "2"},
		}), &req)
		require.NoError(t, err)
		pi := func(i int) *int { return &i }
		require.Equal(t, Req{
			PValue:  pi(1),
			PPValue: func(i *int) **int { return &i }(pi(1)),
			PsValue: []*int{pi(1), pi(2)},
			PValues: func(is []int) *[]int { return &is }([]int{1, 2}),
		}, req)
	})
	t.Run("type alias", func(t *testing.T) {
		type Value int
		type Req struct {
			Value    Value
			Values   []Value
			ValuePtr *Value
		}
		var req Req
		err := QueryBinding{}.Bind(newRequest(url.Values{
			"value":     {"1"},
			"values":    {"1", "2"},
			"value_ptr": {"1"},
		}), &req)
		require.NoError(t, err)
		require.Equal(t, Req{
			Value:    Value(1),
			Values:   []Value{Value(1), Value(2)},
			ValuePtr: func(v Value) *Value { return &v }(Value(1)),
		}, req)
	})
	t.Run("not support complex slice", func(t *testing.T) {
		type Value struct{ Value int }
		type Req struct {
			Values []Value
		}
		var req Req
		err := QueryBinding{}.Bind(newRequest(url.Values{
			"values.value": []string{"1", "2", "3"},
		}), &req)
		// won't return error
		// because Values[*].Value actual has no name, nothing match them
		require.NoError(t, err)
		require.Equal(t, 0, len(req.Values))
		err = QueryBinding{}.Bind(newRequest(url.Values{
			"values": []string{"1", "2", "3"},
		}), &req)
		// return error
		// because Values match the name of Values
		require.Error(t, err)
	})
	t.Run("not support map", func(t *testing.T) {
		type Value struct{ Value int }
		type Req struct {
			Values map[string]Value
		}
		var req Req
		err := QueryBinding{}.Bind(newRequest(url.Values{
			"values.value":          []string{"1", "2", "3"},
			"values.anything.value": []string{"1", "2", "3"},
		}), &req)
		// won't return error
		// because Values[*].Value actual has no name, nothing match them
		require.NoError(t, err)
		require.Equal(t, 0, len(req.Values))
		err = QueryBinding{}.Bind(newRequest(url.Values{
			"values": []string{"1", "2", "3"},
		}), &req)
		// return error
		// because values match the name of Values
		require.Error(t, err)
	})
}
