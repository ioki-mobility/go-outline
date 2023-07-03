module github.com/ioki-mobility/go-outline

go 1.20

require (
	github.com/dghubble/sling v1.4.1
	github.com/rsjethani/secret/v2 v2.3.0
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// The forked repo contains few more features not yet merged into upstream.
replace github.com/dghubble/sling v1.4.1 => github.com/rsjethani/sling v0.0.0-20230703014414-05b42d1f1a76
