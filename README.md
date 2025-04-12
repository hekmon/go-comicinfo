# Comic Info

[![Go Reference](https://pkg.go.dev/badge/github.com/hekmon/go-comicinfo.svg)](https://pkg.go.dev/github.com/hekmon/go-comicinfo)

Easily generate valid `ComicInfo.xml` in Golang based on the [comicinfo schemas](https://github.com/anansi-project/comicinfo) of the [Anansi project](https://anansi-project.github.io/docs/category/comicinfo).

It supports:

* v1
* v2
* v2.1 DRAFT

## Installation

```bash
go get -u github.com/hekmon/go-comicinfo
```

## Example Usage

```go
import "github.com/hekmon/comicinfo"
```

See the CBZ creation example in the `example/cbz.go` file.
