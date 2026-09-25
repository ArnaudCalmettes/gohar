module github.com/ArnaudCalmettes/gohar/dex

go 1.27.0

// harmony n'est pas publié : le workspace le résout pour build et test,
// ce replace le résout aussi pour go mod tidy, qui ignore le workspace.
replace github.com/ArnaudCalmettes/gohar/harmony => ../harmony

require (
	github.com/ArnaudCalmettes/gohar/harmony v0.0.0
	github.com/stretchr/testify v1.12.1
)

require go.yaml.in/yaml/v3 v3.0.5 // indirect
