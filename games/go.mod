module github.com/ArnaudCalmettes/gohar/games

go 1.27.0

// dex, harmony et synth ne sont pas publiés : le workspace les résout pour build et test,
// ces replace les résolvent aussi pour go mod tidy, qui ignore le workspace.
replace (
	github.com/ArnaudCalmettes/gohar/dex => ../dex
	github.com/ArnaudCalmettes/gohar/harmony => ../harmony
	github.com/ArnaudCalmettes/gohar/synth => ../synth
)

require (
	github.com/ArnaudCalmettes/gohar/dex v0.0.0
	github.com/ArnaudCalmettes/gohar/harmony v0.0.0
	github.com/ArnaudCalmettes/gohar/synth v0.0.0
	github.com/ebitengine/oto/v3 v3.5.1
	github.com/hajimehoshi/ebiten/v2 v2.8.6
	gitlab.com/gomidi/midi/v2 v2.3.24
	golang.org/x/image v0.20.0
)

require (
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.11.0 // indirect
	github.com/go-text/typesetting v0.2.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/jfreymuth/pulse v0.1.3 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)
