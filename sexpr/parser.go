package sexpr

import "errors"

// ErrParser is the error value returned by the Parser if the string is not a
// valid term.
// See also https://golang.org/pkg/errors/#New
// and // https://golang.org/pkg/builtin/#error
var ErrParser = errors.New("parser error")

//
// <sexpr>       ::= <atom> | <pars> | QUOTE <sexpr>
// <atom>        ::= NUMBER | SYMBOL
// <pars>        ::= LPAR <dotted_list> RPAR | LPAR <proper_list> RPAR
// <dotted_list> ::= <proper_list> <sexpr> DOT <sexpr>
// <proper_list> ::= <sexpr> <proper_list> | \epsilon
//

/*
	DISCUSSION GRAMMAR:
	<sexpr> => NUM | SYM | LPAR <list> RPAR | QUOTE <sexpr>
	<list> => <sexpr> <tail> | eps
	<tail> => <sexpr> <tail> | DOT <sexpr> | eps

	NT    | FIRST                           | FOLLOW
	sexpr | NUM, SYM, LPAR, QUOTE           | NUM, SYM, LPAR, QUOTE, DOT, LPAR, RPAR, $ 
	list  | NUM, SYM, LPAR, QUOTE, EPS      | RPAR
	tail  | NUM, SYM, LPAR, QUOTE, DOT, EPS | RPAR

	PARSE TABLE
	NT    | NUM                | SYM                | LPAR                     | RPAR     | QUOTE                | DOT           | $
	sexpr | sexpr -> NUM       | sexpr -> SYM       | sexpr -> LPAR list RPAR |          | sexpr -> QUOTE sexpr |               |
	list  | list -> sexpr tail | list -> sexpr tail | list -> sexpr tail       | L -> eps | list -> sexpr tail   |               |
	tail  | tail -> sexpr tail | tail -> sexpr tail | tail -> sexpr tail       | T -> eps | tail -> sexpr tail   | tail -> dot S /
*/

type Parser interface {
	Parse(string) (*SExpr, error)
}

// Implement the Parser interface.
type ParserImpl struct {
	lex     *lexer
	peekTok *Token
}

// NewParser creates a struct of a type that satisfies the Parser interface.
func NewParser() Parser {
	return &ParserImpl{}
}

// Helper function which returns the next token.
func (p *ParserImpl) nextToken() (*Token, error) {
	if tok := p.peekTok; tok != nil {
		p.peekTok = nil
		return tok, nil
	}

	tok, err := p.lex.next()
	if err != nil {
		return nil, ErrParser
	}

	return tok, nil
}

// Helper function which puts a token back as the next token.
func (p *ParserImpl) backToken(tok *Token) {
	p.peekTok = tok
}

// Helper function to peek the next token.
func (p *ParserImpl) peekToken() (*Token, error) {
	tok, err := p.nextToken()
	if err != nil {
		return nil, ErrParser
	}

	p.backToken(tok)

	return tok, nil
}


func (p *ParserImpl) Parse(input string) (*SExpr, error) {
	return p.startNT(input)
}

func (p *ParserImpl) startNT(input string) (*SExpr, error) {
	p.lex = newLexer(input)

	// apply the sexprNT rule 
	sexpr, err := sexprNT()
	if err != nil {
		return nil, ErrParser
	}

	// check that next token is the endmarker $, there should be nothing left after parsing S
	if nextTok, err := p.nextToken(); err != nil || nextTok.typ != tokenEOF {
		return nil, ErrParser
	}

	return sexpr, nil
}

func (p *ParserImpl) sexprNT() (*SExpr, error) {
	// we don't know which rule to use so, we peek
	tok, err := p.nextToken()
	if err != nil {
		return nil, ErrParser
	}

	var sexpr *SExpr
	var er error

	// figure out which rule to use
	switch tok.tokenType {
	// form an atom from these two tokens
	case tokenNumber, tokenSymbol:
		sexpr = &SExpr{
			atom: mkTokenSymbol(tok.literal),
			car: nil,
			cdr: nil,
		}
	
	// apply rule S -> LPAR <list> RPAR
	case tokenLpar:
		// parse list
		list, err := p.listNT()
		if err != nil {
			return nil, ErrParser
		}

		// check for closing parenthesis
		closeParen, err := p.nextToken()
		if err != nil || closeParen.typ != tokenRpar {
	 		return nil, ErrParser
		}
		sexpr = list
	
	default:
		return nil, ErrParser
	}

	return sexpr, nil
}

func (p *ParserImpl) listNT() (*SExpr, error) {
	// no tokens are being consumed, so we only peek
	tok, err := p.peekToken()
	if err != nil {
		return nil, ErrParser
	}
	switch tok.tokenType {
	// apply list -> sexpr tail
	case tokenNumber, tokenSymbol, tokenLpar, tokenQuote:
		sexpr, err := p.sexprNT()
		if err != nil {
			return nil, ErrParser
		}

		tail, err := p.tailNT()
		if err != nil {
			return nil, ErrParser
		}

		return &SExpr{
			atom: nil,
			car: sexpr,
			cdr: tail,
		}, nil
	// apply list -> eps
	case tokenRpar:
		return nil, nil
	default:
		return nil, ErrParser
	}
}

func (p *ParserImpl) tailNT(first *SExpr) (*SExpr, error) {
// no tokens are being consumed, so we only peek
	tok, err := p.peekToken()
	if err != nil {
		return nil, ErrParser
	}
	switch tok.tokenType {
	// apply tail -> sexpr tail
	case tokenNumber, tokenSymbol, tokenLpar, tokenQuote:
		sexpr, err := p.sexprNT()
		if err != nil {
			return nil, ErrParser
		}

		tail, err := p.tailNT()
		if err != nil {
			return nil, ErrParser
		}

		return &SExpr{
			atom: nil,
			car: sexpr,
			cdr: tail,
		}, nil
	// apply tail -> eps
	case tokenRpar:
		return nil, nil

	/*
		FIGURE OUT HOW TO WRITE THIS DOT CELL!!
	*/
	case tokenDot:
		tok, _ := p.nextToken()

		sexpr, err := p.sexprNT()
		if err != nil {
			return nil, ParserErr
		}

		return &SExpr{
			atom: nil,
			car: first,
			cdr: sexpr,
		}, nil

	default:
		return nil, ErrParser
	}
}