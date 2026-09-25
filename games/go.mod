module github.com/ArnaudCalmettes/gohar/games

go 1.27.0

// synth n'est pas publié : le workspace le résout pour build et test,
// ce replace le résout aussi pour go mod tidy, qui ignore le workspace.
replace github.com/ArnaudCalmettes/gohar/synth => ../synth

require (
	github.com/ArnaudCalmettes/gohar/synth v0.0.0
	github.com/ebitengine/oto/v3 v3.5.1
	github.com/hajimehoshi/ebiten/v2 v2.8.6
	gitlab.com/gomidi/midi/v2 v2.3.24
)

require (
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.11.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/jfreymuth/pulse v0.1.3 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
)
