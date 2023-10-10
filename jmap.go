package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Map[V any] struct {
	m map[string]val[V]
	comparer
}

func NewMap[V any]() Map[V] {
	return Map[V]{m: make(map[string]val[V]), comparer: comparer{}}
}

type Option int8

const (
	def Option = iota
	Omit
	OmitEmpty
	Null
)

type val[V any] struct {
	v   *V
	opt Option
}

func (j Map[V]) Len() int {
	return len(j.m)
}

func (j Map[V]) Get(key string) V {
	return *j.m[key].v
}

func (j Map[V]) Set(key string, value V, opt ...Option) {
	j.m[key] = val[V]{
		v:   &value,
		opt: def,
	}

	if len(opt) != 0 {
		v := j.m[key]
		v.opt = opt[0]
		j.m[key] = v
	}
}

func (j Map[V]) Delete(key string) {
	delete(j.m, key)
}

func (j Map[V]) Nil(key string) {
	v := j.m[key]
	v.v = nil
	j.m[key] = v
}

func (j Map[V]) MarshalJSON() ([]byte, error) {
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
				k: v.v,
			}
		case OmitEmpty:
			if v.v == nil {
				continue
			}
			fallthrough
		default:
			section = map[string]any{
				k: &v.v,
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

func (j Map[V]) UnmarshalJSON(data []byte) error {
	bufMap := make(map[string]interface{})

	err := json.Unmarshal(data, &bufMap)
	if err != nil {
		return err
	}

	for k, v := range bufMap {
		if j.comparer.compare(v) {
			var b []byte
			b, err = json.Marshal(v)
			if err != nil {
				return err
			}

			x := NewMap[V]()

			err = x.UnmarshalJSON(b)

			y := any(x).(V)

			j.m[k] = val[V]{v: &y}

			continue
		}

		x := v.(V)

		j.m[k] = val[V]{v: &x}
	}

	return nil
}

type comparer struct{}

func (c comparer) compare(v any) bool {
	var res bool

	switch v.(type) {
	case map[string]interface{}:
		res = true
	default:
		res = false
	}

	return res
}

var jsonstring = `{"greeting": "hello", "nested": {"destination": "Colorado", "season": "fall"}}`

func main() {
	m := NewMap[any]()

	m.Set("greeting", "hello")

	m.Set("nested", map[string]any{
		"location": "Colorado",
		"season":   "fall",
	})

	str, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	println(string(str))

	//err := json.Unmarshal([]byte(jsonstring), &m)
	//if err != nil {
	//	panic(err)
	//}
	//
	//v := m.Get("nested")
	//
	//x := v.(Map[any]).Get("destination")
	//
	//println(fmt.Sprintln(x))
}
