package keyboard

import (
	"fmt"
	"io"
	"text/template"
)

var (
	svgTemplate = `<svg viewBox="{{.ViewBox}}" preserveAspectRatio="xMidYMid meet" xmlns="http://www.w3.org/2000/svg">
<defs>
	<style type="text/css"><![CDATA[
	rect {
		stroke-width:1;
		stroke:black;
	}
	line {
		stroke-width:1;
		stroke:black;
	}
	.key-white {
		fill:white;
	}
	.key-black {
		fill:black;
	}
	rect:hover {
		fill:grey;
	}
	.key-pressed {
		fill:{{.PressedColor}};
	}
	]]></style>
	<clipPath id="canvas">
		<rect x="{{.Canvas.X}}" y="{{.Canvas.Y}}" width="{{.Canvas.Width}}" height="{{.Canvas.Height}}" />
	</clipPath>
</defs>

{{range .WhiteKeys}}<rect id="{{.ID}}" class="{{.Class}}" x="{{.X}}" y="{{.Y}}" width="{{.Width}}" height="{{.Height}}" rx="{{.RX}}" ry="{{.RY}}" clip-path="url(#canvas)" />
{{end}}
{{range .BlackKeys}}<rect id="{{.ID}}" class="{{.Class}}" x="{{.X}}" y="{{.Y}}" width="{{.Width}}" height="{{.Height}}" rx="{{.RX}}" ry="{{.RY}}" clip-path="url(#canvas)" />
{{end}}
<line x1="{{.TopLine.X1}}" y1="{{.TopLine.Y2}}" x2="{{.TopLine.X2}}" y2="{{.TopLine.Y2}}" />
</svg>`

	svgRenderer    = template.Must(template.New("svg").Parse(svgTemplate))
	svgDefaultData = svgTemplateData{
		PressedColor: "#99ccff",
	}
)

type svgTemplateData struct {
	ViewBox      string
	PressedColor string
	Canvas       Rect
	WhiteKeys    []*Rect
	BlackKeys    []*Rect
	TopLine      Line
}

type Rect struct {
	ID     string
	Class  string
	Width  int
	Height int
	X      int
	Y      int
	RX     int
	RY     int
}

type Line struct {
	X1 int
	Y1 int
	X2 int
	Y2 int
}

const (
	startX       = -10
	startY       = -15
	whiteWidth   = 20
	whiteHeight  = 100
	blackWidth   = 14
	blackHeight  = 70
	blackOffset  = -7
	cornerRadius = 5
)

func (k *Keyboard) RenderSVG(w io.Writer) {
	data := svgDefaultData
	data.WhiteKeys = make([]*Rect, 0, len(k.Keys))
	data.BlackKeys = make([]*Rect, 0, len(k.Keys))
	x, y := startX, startY
	for i, key := range k.Keys {
		rect := Rect{
			ID:     fmt.Sprintf("key%d", i),
			X:      x,
			Y:      y,
			Width:  whiteWidth,
			Height: whiteHeight,
			RX:     cornerRadius,
			RY:     cornerRadius,
		}
		if key.isBlack() {
			rect.Class = "key-black"
			rect.X = x + blackOffset
			rect.Width = blackWidth
			rect.Height = blackHeight
			data.BlackKeys = append(data.BlackKeys, &rect)
		} else {
			rect.Class = "key-white"
			data.WhiteKeys = append(data.WhiteKeys, &rect)
			x += whiteWidth
		}
		if key.IsPressed() {
			rect.Class += " key-pressed"
		}
	}
	width := (len(data.WhiteKeys) - 1) * whiteWidth
	height := whiteHeight + startY + 1
	data.ViewBox = fmt.Sprintf("0 0 %d %d", width, height)
	data.Canvas = Rect{Width: width, Height: height}
	data.TopLine = Line{X1: startX, X2: width}
	svgRenderer.Execute(w, data)
}
