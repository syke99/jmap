package jmap

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Map struct {
	m map[string]val
}

func NewMap() *Map {
	return &Map{m: make(map[string]val)}
}

type Option int8

const (
	def Option = iota
	Omit
	OmitEmpty
	Null
)

type val struct {
	v   any
	opt Option
}

func (j *Map) Len() int {
	return len(j.m)
}

func (j *Map) Get(key string) any {
	return j.m[key].v
}

func (j *Map) Set(key string, value any, opt ...Option) {
	var vl any

	switch value.(type) {
	case Map:
		x := value.(Map)
		vl = x.m
	default:
		vl = value
	}

	j.m[key] = val{
		v:   vl,
		opt: def,
	}

	if len(opt) != 0 {
		v := j.m[key]
		v.opt = opt[0]
		j.m[key] = v
	}
}

func (j *Map) Delete(key string) {
	delete(j.m, key)
}

func (j *Map) MarshalJSON() ([]byte, error) {
	var str strings.Builder

	str.WriteString("{")

	for k, v := range j.m {
		v := v

		section := make(map[string]any)

		switch v.opt {
		case Omit:
			continue
		case Null:
			section = map[string]any{
				k: nil,
			}
		case OmitEmpty:
			if v.v == nil {
				continue
			}
			fallthrough
		default:
			section = map[string]any{
				k: v.v,
			}
		}

		jstr, err := json.Marshal(section)
		if err != nil {
			return nil, err
		}

		str.WriteString(fmt.Sprintf("%s,", string(jstr)[1:len(jstr)-1]))
	}

	res := str.String()

	res = fmt.Sprintf("%s}", res[:len(res)-1])

	return []byte(res), nil
}

func (j *Map) UnmarshalJSON(data []byte) error {
	bufMap := make(map[string]any)

	err := json.Unmarshal(data, &bufMap)
	if err != nil {
		return err
	}

	for k, v := range bufMap {
		x := j.m[k]

		switch x.opt {
		case Omit:
			if x.v != nil {
				v = x.v
			} else {
				continue
			}
		case OmitEmpty:
			// TODO: determine v type and skip if it's empty
		case Null:
			if x.v != nil {
				v = x.v
			} else {
				v = nil
			}
		}

		if isMap(v) {
			var b []byte

			b, err = json.Marshal(v)
			if err != nil {
				return err
			}

			x := NewMap()

			err = x.UnmarshalJSON(b)

			y := any(x)

			j.m[k] = val{v: y}

			continue
		}

		y := v

		j.m[k] = val{v: y}
	}

	return nil
}

func isMap(v any) bool {
	var res bool

	switch v.(type) {
	case map[string]any:
		res = true
	default:
		res = false
	}

	return res
}
