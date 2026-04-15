package converter

import (
	"bytes"
	"embed"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"math"
	"os"
	"regexp"
	"text/template"

	"github.com/tdewolff/canvas"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseFS(templateFS, templatePath)
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}
}

const (
	maxIconSize  = 100
	templatePath = "templates/*.tmpl"
)

type Node struct {
	Height int
	Width  int
	SVG    string
}

var spaceRe = regexp.MustCompile(`[\t\r\n]+`)

func NewNode(svgPath string) (Node, error) {
	svgData, err := os.ReadFile(svgPath)
	if err != nil {
		log.Printf("skipping %s: %v", svgPath, err)
		return Node{}, err
	}
	return NewNodeFromBytes(svgData)
}

func NewNodeFromBytes(svgData []byte) (Node, error) {
	w, h, err := computeSVGBounds(svgData)
	if err != nil {
		return Node{}, err
	}
	svgStr := spaceRe.ReplaceAllString(string(svgData), "")
	return Node{
		Width:  w,
		Height: h,
		SVG:    html.EscapeString(svgStr),
	}, nil
}

type Nodes []Node

func (nodes Nodes) RenderTemplate(w io.Writer) error {
	return tmpl.ExecuteTemplate(w, "base", nodes)
}

type NodeMap map[string]Nodes

// computeSVGBounds calcuates the width and height of a SVG path.
func computeSVGBounds(svg []byte) (int, int, error) {
	// Walk XML to find path
	decoder := xml.NewDecoder(bytes.NewReader(svg))
	var ds []string
	for {
		tok, err := decoder.Token()
		if err != nil {
			break
		}
		el, ok := tok.(xml.StartElement)
		if !ok || el.Name.Local != "path" {
			continue
		}
		for _, attr := range el.Attr {
			if attr.Name.Local == "d" {
				ds = append(ds, attr.Value)
			}
		}
	}

	// Start with an empty rect, inverted infinity so any real point expands it
	bounds := canvas.Rect{X0: math.Inf(1), Y0: math.Inf(1), X1: math.Inf(-1), Y1: math.Inf(-1)}
	for _, d := range ds {
		// Parse the svg path string ("M10 20 L30 40...") into geometry object
		path, err := canvas.ParseSVGPath(d)
		if err != nil {
			continue
		}
		// expand the overall bounds to include this path's bounding box
		bounds = bounds.Add(path.Bounds())
	}
	if math.IsInf(bounds.X0, 1) {
		return 0, 0, fmt.Errorf("no paths found in svg")
	}

	// width/height of the actual drawn content
	w := bounds.X1 - bounds.X0
	h := bounds.Y1 - bounds.Y0

	// Scale to maxIconSize, so all icons generated are around the same size
	ratio := maxIconSize / math.Max(h, w)
	h = math.Floor(h * ratio)
	w = math.Floor(w * ratio)

	return int(w), int(h), nil
}
