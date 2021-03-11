package binding

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueryBinding(t *testing.T) {
	type Req struct {
		Int        int
		Ints       []int
		IntPtr     *int
		Uint       uint
		Uints      []uint
		UintPtr    *uint
		Float      float64
		Floats     []float64
		FloatPtr   *float64
		String     string
		Strings    []string
		StringPtr  *string
		Boolean    bool
		Booleans   []bool
		BooleanPtr *bool
	}
	type ReqWithNestedReq struct {
		Req
		Nested Req
	}
	var req ReqWithNestedReq
	newRequest := func(values url.Values) *http.Request {
		return &http.Request{URL: &url.URL{RawQuery: values.Encode()}}
	}
	err := QueryBinding{}.Bind(newRequest(url.Values{
		"int":               {"1", "2", "3"},
		"ints":              {"1", "2", "3"},
		"int_ptr":           {"1", "2", "3"},
		"uint":              {"1", "2", "3"},
		"uints":             {"1", "2", "3"},
		"uint_ptr":          {"1", "2", "3"},
		"float":             {"1", "2", "3"},
		"floats":            {"1", "2", "3"},
		"float_ptr":         {"1", "2", "3"},
		"string":            {"1", "2", "3"},
		"strings":           {"1", "2", "3"},
		"string_ptr":        {"1", "2", "3"},
		"boolean":           {"1", "t", "true"},
		"booleans":          {"1", "t", "true"},
		"boolean_ptr":       {"1", "t", "true"},
		"nested.Int":        {"1", "2", "3"},
		"nested.Ints":       {"1", "2", "3"},
		"nested.IntPtr":     {"1", "2", "3"},
		"nested.Uint":       {"1", "2", "3"},
		"nested.Uints":      {"1", "2", "3"},
		"nested.UintPtr":    {"1", "2", "3"},
		"nested.Float":      {"1", "2", "3"},
		"nested.Floats":     {"1", "2", "3"},
		"nested.FloatPtr":   {"1", "2", "3"},
		"nested.String":     {"1", "2", "3"},
		"nested.Strings":    {"1", "2", "3"},
		"nested.StringPtr":  {"1", "2", "3"},
		"nested.Boolean":    {"1", "t", "true"},
		"nested.Booleans":   {"1", "t", "true"},
		"nested.BooleanPtr": {"1", "t", "true"},
	}), &req)
	require.NoError(t, err)
	pi := func(v int) *int { return &v }
	pu := func(v uint) *uint { return &v }
	pf := func(v float64) *float64 { return &v }
	ps := func(v string) *string { return &v }
	pb := func(v bool) *bool { return &v }
	require.Equal(t, ReqWithNestedReq{
		Req: Req{
			Int:        1,
			Ints:       []int{1, 2, 3},
			IntPtr:     pi(1),
			Uint:       1,
			Uints:      []uint{1, 2, 3},
			UintPtr:    pu(1),
			Float:      1,
			Floats:     []float64{1, 2, 3},
			FloatPtr:   pf(1),
			String:     "1",
			Strings:    []string{"1", "2", "3"},
			StringPtr:  ps("1"),
			Boolean:    true,
			Booleans:   []bool{true, true, true},
			BooleanPtr: pb(true),
		},
		Nested: Req{
			Int:        1,
			Ints:       []int{1, 2, 3},
			IntPtr:     pi(1),
			Uint:       1,
			Uints:      []uint{1, 2, 3},
			UintPtr:    pu(1),
			Float:      1,
			Floats:     []float64{1, 2, 3},
			FloatPtr:   pf(1),
			String:     "1",
			Strings:    []string{"1", "2", "3"},
			StringPtr:  ps("1"),
			Boolean:    true,
			Booleans:   []bool{true, true, true},
			BooleanPtr: pb(true),
		},
	}, req)
}

func TestQueryBinding_TypeAlias(t *testing.T) {
	// test dose QueryBinding work well for type alias
	type Int int
	type Ints []int
	type IntPtr *int
	type Uint uint
	type Uints []uint
	type UintPtr *uint
	type Float float64
	type Floats []float64
	type FloatPtr *float64
	type String string
	type Strings []string
	type StringPtr *string
	type Boolean bool
	type Booleans []bool
	type BooleanPtr *bool
	type Req struct {
		Int        Int
		Ints       Ints
		IntPtr     IntPtr
		Uint       Uint
		Uints      Uints
		UintPtr    UintPtr
		Float      Float
		Floats     Floats
		FloatPtr   FloatPtr
		String     String
		Strings    Strings
		StringPtr  StringPtr
		Boolean    Boolean
		Booleans   Booleans
		BooleanPtr BooleanPtr
	}
	type AliasedReq Req
	type ReqWithNestedReq struct {
		AliasedReq
		Nested AliasedReq
	}
	var req ReqWithNestedReq
	newRequest := func(values url.Values) *http.Request {
		return &http.Request{URL: &url.URL{RawQuery: values.Encode()}}
	}
	err := QueryBinding{}.Bind(newRequest(url.Values{
		"int":               {"1", "2", "3"},
		"ints":              {"1", "2", "3"},
		"int_ptr":           {"1", "2", "3"},
		"uint":              {"1", "2", "3"},
		"uints":             {"1", "2", "3"},
		"uint_ptr":          {"1", "2", "3"},
		"float":             {"1", "2", "3"},
		"floats":            {"1", "2", "3"},
		"float_ptr":         {"1", "2", "3"},
		"string":            {"1", "2", "3"},
		"strings":           {"1", "2", "3"},
		"string_ptr":        {"1", "2", "3"},
		"boolean":           {"1", "t", "true"},
		"booleans":          {"1", "t", "true"},
		"boolean_ptr":       {"1", "t", "true"},
		"nested.Int":        {"1", "2", "3"},
		"nested.Ints":       {"1", "2", "3"},
		"nested.IntPtr":     {"1", "2", "3"},
		"nested.Uint":       {"1", "2", "3"},
		"nested.Uints":      {"1", "2", "3"},
		"nested.UintPtr":    {"1", "2", "3"},
		"nested.Float":      {"1", "2", "3"},
		"nested.Floats":     {"1", "2", "3"},
		"nested.FloatPtr":   {"1", "2", "3"},
		"nested.String":     {"1", "2", "3"},
		"nested.Strings":    {"1", "2", "3"},
		"nested.StringPtr":  {"1", "2", "3"},
		"nested.Boolean":    {"1", "t", "true"},
		"nested.Booleans":   {"1", "t", "true"},
		"nested.BooleanPtr": {"1", "t", "true"},
	}), &req)
	require.NoError(t, err)
	pi := func(v int) *int { return &v }
	pu := func(v uint) *uint { return &v }
	pf := func(v float64) *float64 { return &v }
	ps := func(v string) *string { return &v }
	pb := func(v bool) *bool { return &v }
	require.Equal(t, ReqWithNestedReq{
		AliasedReq: AliasedReq{
			Int:        Int(1),
			Ints:       Ints([]int{1, 2, 3}),
			IntPtr:     IntPtr(pi(1)),
			Uint:       Uint(1),
			Uints:      Uints([]uint{1, 2, 3}),
			UintPtr:    UintPtr(pu(1)),
			Float:      Float(1),
			Floats:     Floats([]float64{1, 2, 3}),
			FloatPtr:   FloatPtr(pf(1)),
			String:     String("1"),
			Strings:    Strings([]string{"1", "2", "3"}),
			StringPtr:  StringPtr(ps("1")),
			Boolean:    Boolean(true),
			Booleans:   Booleans([]bool{true, true, true}),
			BooleanPtr: BooleanPtr(pb(true)),
		},
		Nested: AliasedReq{
			Int:        Int(1),
			Ints:       Ints([]int{1, 2, 3}),
			IntPtr:     IntPtr(pi(1)),
			Uint:       Uint(1),
			Uints:      Uints([]uint{1, 2, 3}),
			UintPtr:    UintPtr(pu(1)),
			Float:      Float(1),
			Floats:     Floats([]float64{1, 2, 3}),
			FloatPtr:   FloatPtr(pf(1)),
			String:     String("1"),
			Strings:    Strings([]string{"1", "2", "3"}),
			StringPtr:  StringPtr(ps("1")),
			Boolean:    Boolean(true),
			Booleans:   Booleans([]bool{true, true, true}),
			BooleanPtr: BooleanPtr(pb(true)),
		},
	}, req)
}