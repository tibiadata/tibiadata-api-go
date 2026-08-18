module github.com/tibiadata/tibiadata-api-go/src/validation

go 1.26.0

replace github.com/tibiadata/tibiadata-api-go/src/tibiamapping => ../tibiamapping

require (
	github.com/projectbarks/cimap v0.1.1
	github.com/stretchr/testify v1.12.0
	github.com/tibiadata/tibiadata-api-go/src/tibiamapping v0.0.0-20260707225952-a0db793b45fa
)

require (
	github.com/go-resty/resty/v2 v2.17.2 // indirect
	golang.org/x/net v0.56.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
