package filter

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const endSymbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	rulee
	rulee1
	rulee2
	rulevalue
	ruleand
	ruleor
	rulenot
	ruleopen
	ruleclose
	rulegt
	rulelt
	ruleeq
	rulene
	rulequote
	ruleexpression
	rulecomparator
	ruleidentifier
	rulep_string
	rulep_number
	ruleparameter
	rulesp
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	rulePegText
	ruleAction8
	ruleAction9
	ruleAction10
)

var rul3s = [...]string{
	"Unknown",
	"e",
	"e1",
	"e2",
	"value",
	"and",
	"or",
	"not",
	"open",
	"close",
	"gt",
	"lt",
	"eq",
	"ne",
	"quote",
	"expression",
	"comparator",
	"identifier",
	"p_string",
	"p_number",
	"parameter",
	"sp",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"PegText",
	"Action8",
	"Action9",
	"Action10",
}

type token32 struct {
	pegRule
	begin, end uint32
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v", rul3s[t.pegRule], t.begin, t.end)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(pretty bool, buffer string) {
	var print func(node *node32, depth int)
	print = func(node *node32, depth int) {
		for node != nil {
			for c := 0; c < depth; c++ {
				fmt.Printf(" ")
			}
			rule := rul3s[node.pegRule]
			quote := strconv.Quote(string(([]rune(buffer)[node.begin:node.end])))
			if !pretty {
				fmt.Printf("%v %v\n", rule, quote)
			} else {
				fmt.Printf("\x1B[34m%v\x1B[m %v\n", rule, quote)
			}
			if node.up != nil {
				print(node.up, depth+1)
			}
			node = node.next
		}
	}
	print(node, 0)
}

func (node *node32) Print(buffer string) {
	node.print(false, buffer)
}

func (node *node32) PrettyPrint(buffer string) {
	node.print(true, buffer)
}

type tokens32 struct {
	tree []token32
}

func (t *tokens32) Trim(length uint32) {
	t.tree = t.tree[:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) AST() *node32 {
	type element struct {
		node *node32
		down *element
	}
	tokens := t.Tokens()
	var stack *element
	for _, token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	if stack != nil {
		return stack.node
	}
	return nil
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	t.AST().Print(buffer)
}

func (t *tokens32) PrettyPrintSyntaxTree(buffer string) {
	t.AST().PrettyPrint(buffer)
}

func (t *tokens32) Add(rule pegRule, begin, end, index uint32) {
	if tree := t.tree; int(index) >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	t.tree[index] = token32{
		pegRule: rule,
		begin:   begin,
		end:     end,
	}
}

func (t *tokens32) Tokens() []token32 {
	return t.tree
}

type FilterPeg struct {
	AST

	Buffer string
	buffer []rune
	rules  [34]func() bool
	parse  func(rule ...int) error
	reset  func()
	Pretty bool
	tokens32
}

func (p *FilterPeg) Parse(rule ...int) error {
	return p.parse(rule...)
}

func (p *FilterPeg) Reset() {
	p.reset()
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer []rune, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p   *FilterPeg
	max token32
}

func (e *parseError) Error() string {
	tokens, error := []token32{e.max}, "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.buffer, positions)
	format := "parse error near %v (line %v symbol %v - line %v symbol %v):\n%v\n"
	if e.p.Pretty {
		format = "parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n"
	}
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf(format,
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			strconv.Quote(string(e.p.buffer[begin:end])))
	}

	return error
}

func (p *FilterPeg) PrintSyntaxTree() {
	if p.Pretty {
		p.tokens32.PrettyPrintSyntaxTree(p.Buffer)
	} else {
		p.tokens32.PrintSyntaxTree(p.Buffer)
	}
}

func (p *FilterPeg) Execute() {
	buffer, _buffer, text, begin, end := p.Buffer, p.buffer, "", 0, 0
	for _, token := range p.Tokens() {
		switch token.pegRule {

		case rulePegText:
			begin, end = int(token.begin), int(token.end)
			text = string(_buffer[begin:end])

		case ruleAction0:
			p.AddOperator(OpAnd)
		case ruleAction1:
			p.AddOperator(OpOr)
		case ruleAction2:
			p.AddOperator(OpNot)
		case ruleAction3:
			p.AddComparator(CmpGt)
		case ruleAction4:
			p.AddComparator(CmpLt)
		case ruleAction5:
			p.AddComparator(CmpEq)
		case ruleAction6:
			p.AddComparator(CmpNe)
		case ruleAction7:
			p.AddExpression()
		case ruleAction8:
			p.AddIdentifier(buffer[begin:end])
		case ruleAction9:
			p.AddArgument(TypeString, buffer[begin:end])
		case ruleAction10:
			p.AddArgument(TypeInt, buffer[begin:end])

		}
	}
	_, _, _, _, _ = buffer, _buffer, text, begin, end
}

func (p *FilterPeg) Init() {
	var (
		max                  token32
		position, tokenIndex uint32
		buffer               []rune
	)
	p.reset = func() {
		max = token32{}
		position, tokenIndex = 0, 0

		p.buffer = []rune(p.Buffer)
		if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != endSymbol {
			p.buffer = append(p.buffer, endSymbol)
		}
		buffer = p.buffer
	}
	p.reset()

	_rules := p.rules
	tree := tokens32{tree: make([]token32, math.MaxInt16)}
	p.parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokens32 = tree
		if matches {
			p.Trim(tokenIndex)
			return nil
		}
		return &parseError{p, max}
	}

	add := func(rule pegRule, begin uint32) {
		tree.Add(rule, begin, position, tokenIndex)
		tokenIndex++
		if begin != position && position > max.end {
			max = token32{rule, begin, position}
		}
	}

	matchDot := func() bool {
		if buffer[position] != endSymbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 e <- <(sp e1 !.)> */
		func() bool {
			position0, tokenIndex0 := position, tokenIndex
			{
				position1 := position
				if !_rules[rulesp]() {
					goto l0
				}
				if !_rules[rulee1]() {
					goto l0
				}
				{
					position2, tokenIndex2 := position, tokenIndex
					if !matchDot() {
						goto l2
					}
					goto l0
				l2:
					position, tokenIndex = position2, tokenIndex2
				}
				add(rulee, position1)
			}
			return true
		l0:
			position, tokenIndex = position0, tokenIndex0
			return false
		},
		/* 1 e1 <- <(e2 ((and e2 Action0) / (or e2 Action1))*)> */
		func() bool {
			position3, tokenIndex3 := position, tokenIndex
			{
				position4 := position
				if !_rules[rulee2]() {
					goto l3
				}
			l5:
				{
					position6, tokenIndex6 := position, tokenIndex
					{
						position7, tokenIndex7 := position, tokenIndex
						{
							position9 := position
							{
								position10, tokenIndex10 := position, tokenIndex
								if buffer[position] != rune('A') {
									goto l11
								}
								position++
								if buffer[position] != rune('N') {
									goto l11
								}
								position++
								if buffer[position] != rune('D') {
									goto l11
								}
								position++
								if !_rules[rulesp]() {
									goto l11
								}
								goto l10
							l11:
								position, tokenIndex = position10, tokenIndex10
								if buffer[position] != rune('a') {
									goto l8
								}
								position++
								if buffer[position] != rune('n') {
									goto l8
								}
								position++
								if buffer[position] != rune('d') {
									goto l8
								}
								position++
								if !_rules[rulesp]() {
									goto l8
								}
							}
						l10:
							add(ruleand, position9)
						}
						if !_rules[rulee2]() {
							goto l8
						}
						{
							add(ruleAction0, position)
						}
						goto l7
					l8:
						position, tokenIndex = position7, tokenIndex7
						{
							position13 := position
							{
								position14, tokenIndex14 := position, tokenIndex
								if buffer[position] != rune('O') {
									goto l15
								}
								position++
								if buffer[position] != rune('R') {
									goto l15
								}
								position++
								if !_rules[rulesp]() {
									goto l15
								}
								goto l14
							l15:
								position, tokenIndex = position14, tokenIndex14
								if buffer[position] != rune('o') {
									goto l6
								}
								position++
								if buffer[position] != rune('r') {
									goto l6
								}
								position++
								if !_rules[rulesp]() {
									goto l6
								}
							}
						l14:
							add(ruleor, position13)
						}
						if !_rules[rulee2]() {
							goto l6
						}
						{
							add(ruleAction1, position)
						}
					}
				l7:
					goto l5
				l6:
					position, tokenIndex = position6, tokenIndex6
				}
				add(rulee1, position4)
			}
			return true
		l3:
			position, tokenIndex = position3, tokenIndex3
			return false
		},
		/* 2 e2 <- <((not value Action2) / value)> */
		func() bool {
			position17, tokenIndex17 := position, tokenIndex
			{
				position18 := position
				{
					position19, tokenIndex19 := position, tokenIndex
					{
						position21 := position
						{
							position22, tokenIndex22 := position, tokenIndex
							if buffer[position] != rune('N') {
								goto l23
							}
							position++
							if buffer[position] != rune('O') {
								goto l23
							}
							position++
							if buffer[position] != rune('T') {
								goto l23
							}
							position++
							if !_rules[rulesp]() {
								goto l23
							}
							goto l22
						l23:
							position, tokenIndex = position22, tokenIndex22
							if buffer[position] != rune('n') {
								goto l20
							}
							position++
							if buffer[position] != rune('o') {
								goto l20
							}
							position++
							if buffer[position] != rune('t') {
								goto l20
							}
							position++
							if !_rules[rulesp]() {
								goto l20
							}
						}
					l22:
						add(rulenot, position21)
					}
					if !_rules[rulevalue]() {
						goto l20
					}
					{
						add(ruleAction2, position)
					}
					goto l19
				l20:
					position, tokenIndex = position19, tokenIndex19
					if !_rules[rulevalue]() {
						goto l17
					}
				}
			l19:
				add(rulee2, position18)
			}
			return true
		l17:
			position, tokenIndex = position17, tokenIndex17
			return false
		},
		/* 3 value <- <((expression sp) / (open e1 close))> */
		func() bool {
			position25, tokenIndex25 := position, tokenIndex
			{
				position26 := position
				{
					position27, tokenIndex27 := position, tokenIndex
					{
						position29 := position
						{
							position30 := position
							{
								position31 := position
								{
									position32, tokenIndex32 := position, tokenIndex
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l33
									}
									position++
									goto l32
								l33:
									position, tokenIndex = position32, tokenIndex32
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l28
									}
									position++
								}
							l32:
							l34:
								{
									position35, tokenIndex35 := position, tokenIndex
									{
										switch buffer[position] {
										case '_':
											if buffer[position] != rune('_') {
												goto l35
											}
											position++
											break
										case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l35
											}
											position++
											break
										case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l35
											}
											position++
											break
										default:
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l35
											}
											position++
											break
										}
									}

									goto l34
								l35:
									position, tokenIndex = position35, tokenIndex35
								}
								add(rulePegText, position31)
							}
							{
								add(ruleAction8, position)
							}
							add(ruleidentifier, position30)
						}
						if !_rules[rulesp]() {
							goto l28
						}
						{
							position38 := position
							{
								switch buffer[position] {
								case 'n':
									{
										position40 := position
										if buffer[position] != rune('n') {
											goto l28
										}
										position++
										if buffer[position] != rune('e') {
											goto l28
										}
										position++
										if !_rules[rulesp]() {
											goto l28
										}
										{
											add(ruleAction6, position)
										}
										add(rulene, position40)
									}
									break
								case 'e':
									{
										position42 := position
										if buffer[position] != rune('e') {
											goto l28
										}
										position++
										if buffer[position] != rune('q') {
											goto l28
										}
										position++
										if !_rules[rulesp]() {
											goto l28
										}
										{
											add(ruleAction5, position)
										}
										add(ruleeq, position42)
									}
									break
								case 'l':
									{
										position44 := position
										if buffer[position] != rune('l') {
											goto l28
										}
										position++
										if buffer[position] != rune('t') {
											goto l28
										}
										position++
										if !_rules[rulesp]() {
											goto l28
										}
										{
											add(ruleAction4, position)
										}
										add(rulelt, position44)
									}
									break
								default:
									{
										position46 := position
										if buffer[position] != rune('g') {
											goto l28
										}
										position++
										if buffer[position] != rune('t') {
											goto l28
										}
										position++
										if !_rules[rulesp]() {
											goto l28
										}
										{
											add(ruleAction3, position)
										}
										add(rulegt, position46)
									}
									break
								}
							}

							add(rulecomparator, position38)
						}
						if !_rules[rulesp]() {
							goto l28
						}
						{
							position48 := position
							{
								position49, tokenIndex49 := position, tokenIndex
								{
									position51 := position
									if !_rules[rulequote]() {
										goto l50
									}
									{
										position52 := position
										{
											position55, tokenIndex55 := position, tokenIndex
											if buffer[position] != rune('\'') {
												goto l55
											}
											position++
											goto l50
										l55:
											position, tokenIndex = position55, tokenIndex55
										}
										if !matchDot() {
											goto l50
										}
									l53:
										{
											position54, tokenIndex54 := position, tokenIndex
											{
												position56, tokenIndex56 := position, tokenIndex
												if buffer[position] != rune('\'') {
													goto l56
												}
												position++
												goto l54
											l56:
												position, tokenIndex = position56, tokenIndex56
											}
											if !matchDot() {
												goto l54
											}
											goto l53
										l54:
											position, tokenIndex = position54, tokenIndex54
										}
										add(rulePegText, position52)
									}
									if !_rules[rulequote]() {
										goto l50
									}
									{
										add(ruleAction9, position)
									}
									add(rulep_string, position51)
								}
								goto l49
							l50:
								position, tokenIndex = position49, tokenIndex49
								{
									position58 := position
									{
										position59 := position
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l28
										}
										position++
									l60:
										{
											position61, tokenIndex61 := position, tokenIndex
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l61
											}
											position++
											goto l60
										l61:
											position, tokenIndex = position61, tokenIndex61
										}
										add(rulePegText, position59)
									}
									{
										add(ruleAction10, position)
									}
									add(rulep_number, position58)
								}
							}
						l49:
							add(ruleparameter, position48)
						}
						{
							add(ruleAction7, position)
						}
						add(ruleexpression, position29)
					}
					if !_rules[rulesp]() {
						goto l28
					}
					goto l27
				l28:
					position, tokenIndex = position27, tokenIndex27
					{
						position64 := position
						if buffer[position] != rune('(') {
							goto l25
						}
						position++
						if !_rules[rulesp]() {
							goto l25
						}
						add(ruleopen, position64)
					}
					if !_rules[rulee1]() {
						goto l25
					}
					{
						position65 := position
						if buffer[position] != rune(')') {
							goto l25
						}
						position++
						if !_rules[rulesp]() {
							goto l25
						}
						add(ruleclose, position65)
					}
				}
			l27:
				add(rulevalue, position26)
			}
			return true
		l25:
			position, tokenIndex = position25, tokenIndex25
			return false
		},
		/* 4 and <- <(('A' 'N' 'D' sp) / ('a' 'n' 'd' sp))> */
		nil,
		/* 5 or <- <(('O' 'R' sp) / ('o' 'r' sp))> */
		nil,
		/* 6 not <- <(('N' 'O' 'T' sp) / ('n' 'o' 't' sp))> */
		nil,
		/* 7 open <- <('(' sp)> */
		nil,
		/* 8 close <- <(')' sp)> */
		nil,
		/* 9 gt <- <('g' 't' sp Action3)> */
		nil,
		/* 10 lt <- <('l' 't' sp Action4)> */
		nil,
		/* 11 eq <- <('e' 'q' sp Action5)> */
		nil,
		/* 12 ne <- <('n' 'e' sp Action6)> */
		nil,
		/* 13 quote <- <'\''> */
		func() bool {
			position75, tokenIndex75 := position, tokenIndex
			{
				position76 := position
				if buffer[position] != rune('\'') {
					goto l75
				}
				position++
				add(rulequote, position76)
			}
			return true
		l75:
			position, tokenIndex = position75, tokenIndex75
			return false
		},
		/* 14 expression <- <(identifier sp comparator sp parameter Action7)> */
		nil,
		/* 15 comparator <- <((&('n') ne) | (&('e') eq) | (&('l') lt) | (&('g') gt))> */
		nil,
		/* 16 identifier <- <(<(([a-z] / [A-Z]) ((&('_') '_') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*)> Action8)> */
		nil,
		/* 17 p_string <- <(quote <(!'\'' .)+> quote Action9)> */
		nil,
		/* 18 p_number <- <(<[0-9]+> Action10)> */
		nil,
		/* 19 parameter <- <(p_string / p_number)> */
		nil,
		/* 20 sp <- <(' ' / '\t')*> */
		func() bool {
			{
				position84 := position
			l85:
				{
					position86, tokenIndex86 := position, tokenIndex
					{
						position87, tokenIndex87 := position, tokenIndex
						if buffer[position] != rune(' ') {
							goto l88
						}
						position++
						goto l87
					l88:
						position, tokenIndex = position87, tokenIndex87
						if buffer[position] != rune('\t') {
							goto l86
						}
						position++
					}
				l87:
					goto l85
				l86:
					position, tokenIndex = position86, tokenIndex86
				}
				add(rulesp, position84)
			}
			return true
		},
		/* 22 Action0 <- <{ p.AddOperator(OpAnd) }> */
		nil,
		/* 23 Action1 <- <{ p.AddOperator(OpOr) }> */
		nil,
		/* 24 Action2 <- <{ p.AddOperator(OpNot) }> */
		nil,
		/* 25 Action3 <- <{ p.AddComparator(CmpGt) }> */
		nil,
		/* 26 Action4 <- <{ p.AddComparator(CmpLt) }> */
		nil,
		/* 27 Action5 <- <{ p.AddComparator(CmpEq) }> */
		nil,
		/* 28 Action6 <- <{ p.AddComparator(CmpNe) }> */
		nil,
		/* 29 Action7 <- <{ p.AddExpression() }> */
		nil,
		nil,
		/* 31 Action8 <- <{ p.AddIdentifier(buffer[begin:end]) }> */
		nil,
		/* 32 Action9 <- <{ p.AddArgument(TypeString, buffer[begin:end]) }> */
		nil,
		/* 33 Action10 <- <{ p.AddArgument(TypeInt, buffer[begin:end]) }> */
		nil,
	}
	p.rules = _rules
}
