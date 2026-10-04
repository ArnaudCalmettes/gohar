module github.com/ArnaudCalmettes/gohar/games

go 1.27.0

// charts, dex, harmony et synth ne sont pas publiés : le workspace les résout pour build et test,
// ces replace les résolvent aussi pour go mod tidy, qui ignore le workspace.
replace (
	github.com/ArnaudCalmettes/gohar/charts => ../charts
	github.com/ArnaudCalmettes/gohar/dex => ../dex
	github.com/ArnaudCalmettes/gohar/harmony => ../harmony
	github.com/ArnaudCalmettes/gohar/synth => ../synth
)

require (
	github.com/ArnaudCalmettes/gohar/charts v0.0.0
	github.com/ArnaudCalmettes/gohar/dex v0.0.0
	github.com/ArnaudCalmettes/gohar/harmony v0.0.0
	github.com/ArnaudCalmettes/gohar/synth v0.0.0
	github.com/BurntSushi/toml v1.6.0
	github.com/ebitengine/oto/v3 v3.5.1
	github.com/hajimehoshi/ebiten/v2 v2.8.6
	github.com/nicksnyder/go-i18n/v2 v2.6.1
	gitlab.com/gomidi/midi/v2 v2.3.24
	golang.org/x/image v0.20.0
	golang.org/x/text v0.32.0
)

require (
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.11.0 // indirect
	github.com/go-text/typesetting v0.2.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/jfreymuth/pulse v0.1.3 // indirect
	github.com/sinshu/go-meltysynth v0.0.0-20230205031334-05d311382fc4 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
)
