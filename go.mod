module github.com/jonhadfield/soba

go 1.19

require (
	github.com/carlescere/scheduler v0.0.0-20170109141437-ee74d2f83d82
	github.com/jonhadfield/githosts-utils v0.0.0-20221002110647-bbd9e109c7b9
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.8.0
	golang.org/x/exp v0.0.0-20221002003631-540bb7301a08
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

//replace github.com/jonhadfield/githosts-utils => ../githosts-utils
