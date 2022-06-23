# Usage:

Start line is the first line contain `INPLACEGEN_($name)_FROM`, end line is the first line contain `INPLACEGEN_($name)_TO`.

There are 3 parts between start line and end line. Data, template and generated text.

Each line in data part and template part should have the same prefix as start line. All process for those part will ignore the prefix.

For example, start line is: `\t// ;INPLACEGEN_(_)_FROM\n`, prefix is `\t// ;` . 

When generate, only the generated text part will be overwritten, and it will be always overwritten.

## Data

The first line after start line is the first data line.

The first data line is a positive integer n, that means below n lines are data part two.

Each line below would be decode by following rules:

- if start with `*`, this line describe extra data, if not, this line describe normal data
- rest content would be regarded as json, which describe a `[][]string`

## Template

Template would be parsed by `text/template`. And would be execute with data.

### Replace before execute

But before parse:`{{$[i.j]}}` would be replaced by every `j`th element in `i`th normal data line.

And `{{$[i]}}`, `{{$i}}` are short ways of `{{$[i.1]}}`

For example, normal data lines are:

```
*[["A", "a"], ["B", "b"], ["C", "c"]]
[["0", "f", "false"], ["1", "t", "true"]]
```

Template lines are:

```
{{$data := .}}{{range .ex1}}{{.v1}}-{{.v2}};{{$[1.1]}}-{{$[1.2]}}-{{$data.e1.v3}}
{{- end}}
```

Replaced template:

```
{{$data := .}}{{range .ex1}}{{.v1}}-{{.v2}};0-f-{{$data.e1.v3}}
{{-end}}
{{$data := .}}{{range .ex1}}{{.v1}}-{{.v2}};1-t-{{$data.e1.v3}}
{{-end}}
```

There would be `len(normal_data[0])*len(normal_data[1])*...*len(normal_data[n])` replicas to set value. For this example, it's `2`.

### Execute data

Data are converted to map, and use map as data for template execute.

For example, Original data:

```go
NormalData := [][]string{{"0", "f", "false"}, {"1", "t", "true"}}
ExtraData := [][]string{{"A", "a"}, {"B", "b"}, {"C", "c"}}
```

Template Execute Data:

```go
// for 1st replica
Data := map[string]interface{}{
	"e1": map[string]string{
		"v1": "0",
		"v2": "f",
		"v3": "false",
	},
	"ex1": []interface{}{
		map[string]string{"v1": "A", "v2": "a"},
		map[string]string{"v1": "B", "v2": "b"},
	},
}
// for 2ed replica
Data := map[string]interface{}{
	"e1": map[string]string{
		"v1": "1",
		"v2": "t",
		"v3": "true",
	},
	"ex1": []interface{}{
		map[string]string{"v1": "A", "v2": "a"},
		map[string]string{"v1": "B", "v2": "b"},
	},
}
```

## Generated Text

Generated text be replace by template output in-place.

## Import

can use `import $filename`, to treat content in that file as Data and Template. the line after `import $filename` would keep, but not used.

can use `import $filename` in each data line to import file content as Data.

can't import template.