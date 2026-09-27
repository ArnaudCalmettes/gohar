module github.com/ArnaudCalmettes/gohar/charts

go 1.27.0

// harmony n'est pas publié : le workspace le résout pour build et test,
// ce replace le résout aussi pour go mod tidy, qui ignore le workspace.
replace github.com/ArnaudCalmettes/gohar/harmony => ../harmony

require github.com/ArnaudCalmettes/gohar/harmony v0.0.0
