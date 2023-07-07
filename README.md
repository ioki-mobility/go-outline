# Go client and cli for [outline][def]

[def]: https://www.getoutline.com/


# Usage

## Client
Getting a single document by its id:
```golang
cl := outline.New()
doc, err := cl.Documents().Get().ByID("document id").Do(context.Background())
if err != nil {
	panic(err)
}
fmt.Println(doc)
```

Get all documents for a collection:
```golang
cl := outline.New()
cl.Documents().GetAll().Collection("collection id").Do(context.Background(), func(d *outline.Document, err error) bool {
		fmt.Println(d)
		return true
	})
```


## CLI
TBA