package chrono

import (
	"fmt"
	"strconv"

	"github.com/davecgh/go-spew/spew"
)

var rules = []func(v int, sign int, digits uint, conf ParseConfig, state *state, parts *parts, rev bool) (string, bool){
	// [2006]-01-02 => %Y
	func(v int, sign int, digits uint, conf ParseConfig, state *state, parts *parts, rev bool) (string, bool) {
		if digits == 4 {
			fmt.Println("rule 0", v)
			state.component = componentDate
			state.datePart = datePartYear
			// TODO sign
			return "%Y", true
		}
		return "", false
	},
	// 01-[02] => %d
	func(v int, sign int, digits uint, conf ParseConfig, state *state, parts *parts, rev bool) (string, bool) {
		if digits == 2 && state.component == componentDate && state.d(rev) == datePartNone {
			fmt.Println("rule 1", v)
			state.component = componentDate
			state.datePart = datePartDay
			return "%d", true
		}
		return "", false
	},
	// 2006-[01]-02 => %m
	func(v int, sign int, digits uint, conf ParseConfig, state *state, parts *parts, rev bool) (string, bool) {
		if digits == 2 && state.component == componentDate && state.d(rev) != datePartMonth {
			fmt.Println("rule 2", v)
			state.component = componentDate
			state.datePart = datePartMonth
			return "%m", true
		}
		return "", false
	},
	// 2006-01-[02] %d
	func(v int, sign int, digits uint, conf ParseConfig, state *state, parts *parts, rev bool) (string, bool) {
		if digits == 2 && state.component == componentDate && state.d(rev) != datePartDay {
			fmt.Println("rule 3", v)
			state.component = componentDate
			state.datePart = datePartDay
			return "%d", true
		}
		return "", false
	},
}

func evalParseRules(prevTyp rune, buf []rune, conf ParseConfig, state *state, layout *[]string, parts *parts, ambiguous *[]chunk) {
	switch prevTyp {
	case 'w', 'o':
		*layout = append(*layout, string(buf))
	case 's':
		switch buf[0] {
		case '-':
			state.component = componentDate
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

		for _, rule := range rules {
			spew.Dump(state)
			if str, ok := rule(v, sign, uint(len(buf)), conf, state, parts, false); ok {
				*layout = append(*layout, str)

				if len(*ambiguous) != 0 {
					(*ambiguous)[len(*ambiguous)-1].state = *state
				}

				return
			}
		}

		*layout = append(*layout, "")
		*ambiguous = append(*ambiguous, chunk{
			pos:    uint(len(*layout) - 1),
			v:      v,
			sign:   sign,
			digits: uint(len(buf)),
		})
	}
}

func evalAmbiguous(chunks []chunk, conf ParseConfig, layout []string, parts *parts) {
	fmt.Println("========= evalAmbiguous")
	for i := len(chunks) - 1; i >= 0; i-- {
		spew.Dump(chunks[i])

		for _, rule := range rules {
			str, ok := rule(chunks[i].v, chunks[i].sign, chunks[i].digits, conf, &chunks[i].state, parts, true)
			layout[chunks[i].pos] = str
			if ok {
				break
			}
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
	} else {
		return s.datePart
	}
}

func (s state) t(rev bool) timePart {
	return s.timePart // TODO
}

type component rune

const (
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

type parts struct {
	year int
}

type chunk struct {
	pos    uint
	v      int
	sign   int
	digits uint
	state  state
}
