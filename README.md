# Go client and cli for [outline][def]

[def]: https://www.getoutline.com/


# Usage

## Client
```golang
// Create top level client.
cl := outline.New("https://baseurl/api",&http.Client{},"api key")


// Fetch information about a collection:
col,err := cl.Collections().Get("collection id").Do(context.Background())
if err != nil {
	panic(err)
}
fmt.Println(col)


// Fetch information about all collections.
err := cl.Collections().List().Do(context.Background(), func(c *outline.Collection, err error) (bool, error) {
	fmt.Println(c)
	return true, nil
})
if err != nil {
	panic(err)
}
```


## CLI
TBA