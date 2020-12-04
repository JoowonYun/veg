package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/fogleman/gg"
	"github.com/joowonyun/veg/lexer"
	"github.com/joowonyun/veg/parser"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Give one argument")
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	input, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	lexer := lexer.NewLexer(string(input))
	parser := parser.NewParser(lexer)
	svg := parser.ParseSvg()

	dc := gg.NewContext(svg.Width, svg.Height)
	// dc.SetColor(color.White)
	// dc.Clear()

	for _, s := range svg.Drawables {
		s.Draw(dc)
	}

	dc.SavePNG("out1.png")
}
