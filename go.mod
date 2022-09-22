module github.com/kawaway/interceptor

go 1.17

require (
	github.com/pion/logging v0.2.2
	github.com/pion/rtcp v1.2.9
	github.com/pion/rtp v1.7.13
	github.com/stretchr/testify v1.7.1
)

require (
	github.com/davecgh/go-spew v1.1.0 // indirect
	github.com/pion/randutil v0.1.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect
)

//replace github.com/pion/interceptor => github.com/kawaway/interceptor v0.0.3-0.20220916093414-9309326299fb

//replace github.com/kawaway/interceptor => ./
