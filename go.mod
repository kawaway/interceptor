module github.com/kawaway/interceptor

go 1.15

require (
	github.com/pion/logging v0.2.2
	github.com/pion/rtcp v1.2.9
	github.com/pion/rtp v1.7.13
	github.com/stretchr/testify v1.7.1
)

replace github.com/pion/interceptor => github.com/kawaway/interceptor v0.0.3-0.20220916093414-9309326299fb

replace github.com/kawaway/interceptor => ./
