# Comic Info

[![Go Reference](https://pkg.go.dev/badge/github.com/hekmon/go-comicinfo.svg)](https://pkg.go.dev/github.com/hekmon/go-comicinfo)

Easily generate valid `ComicInfo.xml` in Golang based on the [comicinfo schemas](https://github.com/anansi-project/comicinfo) of the [Anansi project](https://anansi-project.github.io/docs/category/comicinfo).


## Versions

It supports:

* [v1](https://anansi-project.github.io/docs/comicinfo/schemas/v1.0)
* [v2](https://anansi-project.github.io/docs/comicinfo/schemas/v2.0)
* [v2.1 DRAFT](https://anansi-project.github.io/docs/comicinfo/schemas/v2.1)

### Used fields

Not all apps will use all fields, here is a list of fields used by some known apps:

* [Komga](https://komga.org/docs/guides/scan-analysis-refresh/#import-metadata-for-cbrcbz-containing-a-comicinfoxml-file) uses some fields from v2.1 draft.
* [Kavita](https://wiki.kavitareader.com/guides/metadata/comics/) uses some fields from v2.1 draft.

## Installation

```bash
go get -u github.com/hekmon/go-comicinfo
```

## Example Usage

See the CBZ creation [example](example/cbz.go) for a full usage example.
