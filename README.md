[![Checks](https://github.com/ioki-mobility/go-outline/actions/workflows/checks.yml/badge.svg)](https://github.com/ioki-mobility/go-outline/actions/workflows/checks.yml)
[![MIT](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/ioki-mobility/go-outline/blob/main/LICENSE)

# go-outline

Go client and cli for [outline](https://www.getoutline.com/).

# Installation

## Installing the client:
```shell
go get github.com/ioki-mobility/go-outline
```

## Installing the CLI:
- Download pre-built binaries from [releases](https://github.com/ioki-mobility/go-outline/releases) page for your platform
- Install via go toolchain:
```shell
go install github.com/ioki-mobility/go-outline/cli@latest
```
NOTE: Very soon CLI binaries would be available for direct download as part of releases.

# Usage

## Client

### Create a client

```go
cl := outline.New("https://server.url", &http.Client{}, "api key")
```

> **Note**: You can create a new API key in your outline **account settings**.

### Get a collection

```go
col, err := cl.Collections().Get("collection id").Do(context.Background())
if err != nil {
	panic(err)
}
fmt.Println(col)
```


### Get all collections

```go
err := cl.Collections().List().Do(context.Background(), func(col *outline.Collection, err error) (bool, error) {
	fmt.Println(col)
	return true, nil
})
if err != nil {
	panic(err)
}
```

### Create a collection

```go
col, err := cl.Collections().Create("collection name").Do(context.Background()) 
if err != nil {
	panic(err)
}
fmt.Println(col)
```

There are also **optional** functions for the `CollectionsCreateClient` available:
```go
colCreateClient := cl.Collections().Create("collection name")
colCreateClient.
	Description("desc"). 
	PermissionRead(). // or PermissionReadWrite()
	Color("#c0c0c0").
	Private(true).
	Do(context.Background())
```

### Document Create

```go
doc := cl.Documents().Create("Document name", "collection id").Do(context.Background())
```

There re also **optional** functions for the `DocumentsCreateClient` available:
```go
docCreateClient := cl.Documents().Create("Document name")
docCreateClient.
	Publish(true). 
	Text("text").
	ParentDocumentID(DocumentId("parent document id")).
	TemplateID(TemplateId("templateId")).
	Template(false).
	Do(context.Background())
```

## CLI

### Build the CLI

```
go build -o outline ./cli
```

### Required flags

Any command requires the flags `server` and `key`.
You can simply add them with `--server https://server.url` 
and `--key sup3rS3cre7Ap1K3Y`.

### Collections

The basic command to work with collections is:
```
outline collections
```

#### Collections fetch

To fetch a collection and display it in a json string use:
```
outline collections fetch [COLLECTION_ID] [flags]
```

#### Collections create

To create a collection and display the created collection as a json string use:
```
outline collections create [COLLECTION_NAME] [flags]
```
