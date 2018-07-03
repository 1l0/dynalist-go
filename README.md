# dynalist-go

[![GoDoc](https://godoc.org/github.com/1l0/dynalist-go?status.svg)](https://godoc.org/github.com/1l0/dynalist-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/1l0/dynalist-go)](https://goreportcard.com/report/github.com/1l0/dynalist-go)

Set an env before using:

```bash
export DYNALIST_TOKEN=your_secret_token
```

## example

Get file list:
```go
api, err := dynalist.New()
if err != nil {
	panic(err)
}

res, err := api.FileList()
if err != nil {
	panic(err)
}

fmt.Printf("%+v\n", res)
```