package parser

import (
	"log"
	"strconv"

	"github.com/joowonyun/veg/lexer"
)

type Parser struct {
	lex          *lexer.Lexer
	currentToken lexer.Token
	nextToken    lexer.Token
}

func NewParser(lexer *lexer.Lexer) *Parser {
	parser := Parser{lex: lexer}
	parser.readToken()
	parser.readToken()
	return &parser
}

func (p *Parser) readToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.lex.NextToken()
}

func (p *Parser) ParseSvg() *Svg {
	s := Svg{}
	s.Drawables = []drawable{}

	p.expect(lexer.TokenOpen)
	if tagName := p.parseIdentifier(); tagName != "svg" {
		log.Fatalf("expected tag 'svg', got '%s'\n", tagName)
	}

	attributes := p.parseAttributes()
	s.Width = lookupInt(attributes, "width")
	s.Height = lookupInt(attributes, "height")

	p.expect(lexer.TokenClose)

	// Stop parsing shapes when finding "</"
	for p.currentToken.Type != lexer.TokenOpen || p.nextToken.Type != lexer.TokenSlash {
		d := p.parseShape()
		s.Drawables = append(s.Drawables, d)
	}

	// Consume and ignore remaining tokens
	for p.currentToken.Type != lexer.EOF {
		p.readToken()
	}

	return &s
}

func (p *Parser) parseAttributes() map[string]string {
	attributes := make(map[string]string)
	for p.currentToken.Type != lexer.TokenSlash && p.currentToken.Type != lexer.TokenClose {
		k, v := p.parseAttribute()
		attributes[k] = v
	}
	return attributes
}

func (p *Parser) parseAttribute() (key, value string) {
	key = p.parseIdentifier()
	p.expect(lexer.TokenEqual)
	p.expect(lexer.TokenQuote)
	value = p.parseIdentifier()
	p.expect(lexer.TokenQuote)
	return
}

func (p *Parser) expect(tokenType lexer.TokenType) {
	if p.currentToken.Type != tokenType {
		log.Fatalf("expect %q, got %q (%q)", tokenType, p.currentToken.Type, p.currentToken.Literal)
	}
	p.readToken()
}

func (p *Parser) parseIdentifier() string {
	if p.currentToken.Type != lexer.TokenIdentifier {
		log.Printf("expect identifier, got %q (%q)", p.currentToken.Type, p.currentToken.Literal)
	}
	id := p.currentToken.Literal
	p.readToken()
	return id
}

func (p *Parser) parseShape() drawable {
	p.expect(lexer.TokenOpen)
	shapeName := p.parseIdentifier()
	attributes := p.parseAttributes()

	var d drawable
	switch shapeName {
	case "circle":
		d = parseCircle(attributes)
	default:
		log.Printf("unknown shape %s\n", shapeName)
	}

	p.expect(lexer.TokenSlash)
	p.expect(lexer.TokenClose)

	return d
}

func parseCircle(attributes map[string]string) *circle {
	cx := lookupFloat(attributes, "cx")
	cy := lookupFloat(attributes, "cy")
	r := lookupFloat(attributes, "r")
	strokeWidth := lookupFloat(attributes, "stroke-width")
	stroke := attributes["stroke"]
	fill := attributes["fill"]

	return &circle{
		shape: shape{
			x:           cx,
			y:           cy,
			strokeWidth: strokeWidth,
			fillColor:   fill,
			strokeColor: stroke,
		},
		radius: r,
	}
}

func lookupFloat(m map[string]string, key string) float64 {
	if m[key] == "" {
		return 0.0
	}
	fl, err := strconv.ParseFloat(m[key], 64)
	if err != nil {
		panic("cannot parse float: " + key + ":" + m[key])
	}
	return fl
}

func lookupInt(m map[string]string, key string) int {
	if m[key] == "" {
		return 0
	}
	integer, err := strconv.Atoi(m[key])
	if err != nil {
		panic("cannot parse integer: " + key + ":" + m[key])
	}
	return integer
}
