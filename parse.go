//go:build parse

package chrono

import (
	"strconv"
	"strings"
	"unicode"
)

type ParseConfig struct {
	DayFirst bool
}

type Chronological interface {
	String() string
	Parse(layout, value string) error
	get() (dv, tv, ov *int64)
	set(dv, tv, ov int64)
}

// func Parse(value string, conf ParseConfig) (Chronological, error) {
// 	// pick type

// 	return nil, nil
// }

// ParseToLayout attempts to parse the input string as C and returns
// the applicable layout string that would parse that string, or format C to that string.
// Any valid and non-ambiguous ISO 8601 string will be parsed correctly,
// and any other string will be parsed with a best effort attempt, although the resulting
// layout string may not be valid.
func ParseToLayout(value string, conf ParseConfig, c Chronological) (string, error) {
	var (
		typ, prevTyp rune // a = alpha, n = numeric, s = separator, w = whitespace, o = other
		sign         int  // -1 = negative, 0 = no sign, 1 = positive
		buf          []rune

		state     state
		layout    []string
		parts     parts
		ambiguous []chunk
	)

	var date, time, offset *int64
	if c != nil {
		date, time, offset = c.get()
	}

	var err error
	if date != nil {
		if parts.year, parts.month, parts.day, err = fromDate(*date); err != nil {
			return "", err
		}

		if parts.isoYear, parts.isoWeek, err = getISOWeek(*date); err != nil {
			return "", err
		}
	}

	if time != nil {
		parts.hour, parts.min, parts.sec, parts.nsec = fromTime(*time)
		_, parts.isAfternoon = convert24To12HourClock(parts.hour)
	}

	if offset != nil {
		parts.offset = *offset
	}

	_ = sign // TODO

	for i, c := range []rune(value) {
		switch {
		case (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z'):
			typ = 'a'
		case c >= '0' && c <= '9':
			typ = 'n'
		case (c == '+' || c == '-') && (prevTyp == 0 || prevTyp != 'n'):
			typ = 'n'
		case c == '/' || c == '-' || c == '.' || c == ',' || c == ':':
			typ = 's'
		case unicode.IsSpace(c):
			typ = 'w'
		default:
			typ = 'o'
		}

		if typ == prevTyp {
			buf = append(buf, c)

			if i != len(value)-1 {
				continue
			}
		}

		evalParseRules(prevTyp, buf, conf, &state, &layout, &parts, &ambiguous)

		prevTyp = typ
		typ = 0
		sign = 0
		buf = []rune{c}
	}

	//fmt.Println("=>", layout)

	evalAmbiguous(ambiguous, conf, layout, &parts)

	layoutStr := strings.Join(layout, "")

	if c != nil {
		if err := applyParts(parts, date, time, offset); err != nil {
			return layoutStr, err
		}

		c.set(*date, *time, *offset)
	}

	return layoutStr, nil
}

func evalParseRules(prevTyp rune, buf []rune, conf ParseConfig, state *state, layout *[]string, parts *parts, ambiguous *[]chunk) {
	switch prevTyp {
	case 'w', 'o':
		*layout = append(*layout, string(buf))
	case 's':
		switch buf[0] {
		case '-':
			state.component = componentDate
		case '/':
			state.component = componentDate
			// TODO conf
		case ':':
			state.component = componentTime
		}
		*layout = append(*layout, string(buf))
	case 'n':
		var sign int
		switch buf[0] {
		case '+':
			sign = 1
		case '-':
			sign = -1
		}

		v, err := strconv.Atoi(string(buf))
		if err != nil {
			panic(err)
		}

		if str, ok := eval(v, sign, uint(len(buf)), conf, state, parts, false); ok {
			*layout = append(*layout, str)

			if len(*ambiguous) != 0 {
				(*ambiguous)[len(*ambiguous)-1].state = *state
			}

			return
		}

		*layout = append(*layout, string(buf))
		*ambiguous = append(*ambiguous, chunk{
			pos:    uint(len(*layout) - 1),
			v:      v,
			sign:   sign,
			digits: uint(len(buf)),
		})
	}
}

func evalAmbiguous(chunks []chunk, conf ParseConfig, layout []string, parts *parts) {
	for i := len(chunks) - 1; i >= 0; i-- {
		if str, ok := eval(chunks[i].v, chunks[i].sign, chunks[i].digits, conf, &chunks[i].state, parts, true); ok {
			layout[chunks[i].pos] = str
		}
	}
}

type state struct {
	component
	datePart
	timePart
}

func (s state) String() string {
	return string(s.component) + string(s.datePart) + string(s.timePart)
}

func (s state) d(rev bool) datePart {
	if rev && s.datePart == datePartYear {
		return datePartDay
	} else if rev && s.datePart == datePartDay {
		return datePartYear
	}
	return s.datePart
}

func (s state) t(rev bool) timePart {
	return s.timePart // TODO
}

type component rune

const (
	componentNone           = 0
	componentDate component = 'd'
	componentTime component = 't'
)

type datePart rune

const (
	datePartNone           = 0
	datePartYear  datePart = 'y'
	datePartMonth datePart = 'm'
	datePartDay   datePart = 'd'
)

type timePart rune

const (
	timePartHour   timePart = 'h'
	timePartMinute timePart = 'm'
	timePartSecond timePart = 's'
)

type chunk struct {
	pos    uint
	v      int
	sign   int
	digits uint
	state  state
}

func eval(v int, sign int, digits uint, conf ParseConfig, state *state, parts *parts, rev bool) (str string, ok bool) {
	// fmt.Println("========= eval")
	// fmt.Println("v:", v)
	// fmt.Println("sign:", sign)
	// fmt.Println("digits:", digits)
	// fmt.Println("conf:", conf)
	// fmt.Println("state:", *state)
	// fmt.Println("datePart:", string(state.d(rev)))
	// fmt.Println("parts:", *parts)
	// fmt.Println("rev:", rev)

	// defer func() {
	// 	fmt.Println("out:", ok, str)
	// }()

	switch state.component {
	case componentNone:
		switch digits {
		case 4: // [2006](-01-02) => %Y
			if sign == 0 {
				sign = 1
			}

			state.component = componentDate
			state.datePart = datePartYear
			parts.year = v * sign
			return "%Y", true
		}
	case componentDate:
		switch d := state.d(rev); d {
		case datePartNone:
			switch digits {
			case 2:
				switch conf.DayFirst {
				case false:
					switch sign {
					case 0: // 01-[02] => %d
						state.component = componentDate
						state.datePart = datePartDay
						parts.day = v
						return "%d", true
					}
				case true:
					switch sign {
					case 0: // 02-[01](-2006) => %m
						state.component = componentDate
						state.datePart = datePartMonth
						parts.month = v
						return "%m", true
					}
				}
			case 4: // 02-[2006] => %Y
				if sign == 0 {
					sign = 1
				}

				state.component = componentDate
				state.datePart = datePartYear
				parts.year = v * sign
				return "%Y", true
			}
		case datePartYear:
			switch digits {
			case 2:
				switch conf.DayFirst {
				case false:
					switch sign {
					case 0: // 2006-[01](-02) => %m
						state.component = componentDate
						state.datePart = datePartMonth
						parts.month = v
						return "%m", true
					}
				case true:
					switch sign {
					case 0: // 2006-[02](-01) => %d
						state.component = componentDate
						state.datePart = datePartDay
						parts.day = v
						return "%d", true
					}
				}
			}
		case datePartMonth:
			switch digits {
			case 2:
				switch conf.DayFirst {
				case false:
					switch sign {
					case 0: // (2006)-01-[02] => %d
						state.component = componentDate
						state.datePart = datePartDay
						parts.day = v
						return "%d", true
					}
				case true:
					switch sign {
					case 0: // // [02]-01(-2006) => %m
						state.component = componentDate
						state.datePart = datePartDay
						parts.day = v
						return "%d", true
					}
				}
			}
		case datePartDay:
			switch digits {
			case 2:
				switch conf.DayFirst {
				case true:
					switch sign {
					case 0: // 2006-02-[01] => %m
						state.component = componentDate
						state.datePart = datePartMonth
						parts.month = v
						return "%m", true
					}
				case false:
					switch sign {
					case 0: // [01]-2006 => %m
						state.component = componentDate
						state.datePart = datePartMonth
						parts.month = v
						return "%m", true
					}
				}
			}
		}
	}

	return "", false
}
