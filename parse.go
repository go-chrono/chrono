package chrono

import (
	"strconv"
)

var rules = []func(v int, sign int, digits uint, conf ParseConfig, state *state, parts *parts) (string, bool){
	func(v int, sign int, digits uint, conf ParseConfig, state *state, parts *parts) (string, bool) {
		if digits == 4 && sign != 0 {
			state.component = componentDate
			state.datePart = datePartYear
			return "%Y", true
		}
		return "", false
	},
	func(v int, sign int, digits uint, conf ParseConfig, state *state, parts *parts) (string, bool) {
		if digits == 4 {
			state.component = componentDate
			state.datePart = datePartYear
			return "%Y", true
		}
		return "", false
	},
}

func evalParseRules(prevTyp rune, buf []rune, conf ParseConfig, state *state, layout *[]string, parts *parts) {
	switch prevTyp {
	case 's', 'w', 'o':
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
			str, ok := rule(v, sign, uint(len(buf)), conf, state, parts)
			if ok {
				*layout = append(*layout, str)
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

type component int

const (
	componentDate component = iota + 1
	componentTime
)

type datePart int

const (
	datePartYear datePart = iota + 1
	datePartMonth
	datePartDay
)

type timePart int

const (
	timePartHour timePart = iota + 1
	timePartMinute
	timePartSecond
)

type parts struct {
	year int
}
