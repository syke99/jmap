package jmap

import (
	"encoding/json"
	"strings"
)

type Map[V any] interface {
	Set(key string, val V, opt ...Option)
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}

func NewMap[V any]() Map[V] {
	return jmap[V]{}
}

type jmap[V any] struct {
	m map[string]val[V]
}

type Option int8

const (
	Omit Option = iota
	OmitEmpty
	Null
)

type val[V any] struct {
	v   *V
	opt Option
}

func (j jmap[V]) Set(key string, value V, opt ...Option) {
	j.m[key] = val[V]{
		v: &value,
	}

	if len(opt) != 0 {
		v := j.m[key]
		v.opt = opt[0]
		j.m[key] = v
	}
}

func (j jmap[V]) MarshalJSON() ([]byte, error) {
	var str strings.Builder

	str.WriteString("{")

	for k, v := range j.m {
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
		str.WriteString(string(jstr)[1 : len(jstr)-1])
	}

	str.WriteString("}")

	return []byte(str.String()), nil
}

func (j jmap[V]) UnmarshalJSON([]byte) error {
	// TODO: implement JSON unmarshalling .-.
	return nil
}
