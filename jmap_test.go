package jmap

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var loc = "us"

func TestIsMap(t *testing.T) {
	// Arrange
	m := map[string]any{}

	// Act
	res := isMap(m)

	// Asset
	assert.True(t, res)
}

func TestIsMap_NotMap(t *testing.T) {
	// Arrange
	s := "hello"

	// Act
	res := isMap(s)

	// Assert
	assert.False(t, res)
}

func TestNewMap(t *testing.T) {
	// Arrange
	expected := &Map{m: make(map[string]val)}

	// Act
	actual := NewMap()

	// Assert
	assert.NotNil(t, actual)
	assert.Equal(t, *expected, *actual)
}

func TestSet(t *testing.T) {
	// Arrange
	m := NewMap()

	// Assert
	assert.Equal(t, 0, len(m.m))

	// Act
	m.Set("greeting", "hello")

	// Assert
	assert.Equal(t, 1, len(m.m))
}

func TestLen(t *testing.T) {
	// Arrange
	m := NewMap()

	// Act
	l := m.Len()

	// Assert
	assert.Equal(t, 0, l)

	// Arrange
	m.Set("greeting", "hello")

	// Act
	l = m.Len()

	// Assert
	assert.Equal(t, 1, l)
}

func TestGet(t *testing.T) {
	// Arrange
	greeting := "hello"

	m := NewMap()

	m.Set("greeting", greeting)

	// Act
	hello := m.Get("greeting")

	// Assert
	assert.Equal(t, greeting, hello)
}

func TestDelete(t *testing.T) {
	// Arrange
	m := NewMap()

	l := m.Len()

	// Assert
	assert.Equal(t, 0, l)

	// Arrange
	m.Set("greeting", "hello")

	m.Set("location", loc)

	l = m.Len()

	// Assert
	assert.Equal(t, 2, l)

	// Act
	m.Delete("greeting")

	// Arrange
	l = m.Len()

	// Assert
	assert.Equal(t, 1, l)
}

var greeting = "hello"
var jsonstring = fmt.Sprintf("{\"greeting\":\"hello\",\"nested\":{\"location\":\"%s\",\"season\":\"fall\"}}", loc)

func TestJSONMarshal(t *testing.T) {
	// Arrange
	m := NewMap()

	m.Set("greeting", greeting)

	n := NewMap()

	n.Set("location", loc)

	n.Set("season", "fall")

	m.Set("nested", n)

	// Act
	str, err := json.Marshal(m)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, jsonstring, string(str))
}

func TestJSONMarshalOmit(t *testing.T) {
	// Arrange
	m := NewMap()

	m.Set("greeting", greeting)

	n := NewMap()

	n.Set("location", loc)

	n.Set("season", "fall")

	m.Set("nested", n, Omit)

	// Act
	str, err := json.Marshal(m)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "{\"greeting\":\"hello\"}", string(str))
}

func TestJSONMarshalNull(t *testing.T) {
	// Arrange
	m := NewMap()

	m.Set("greeting", greeting, Null)

	// Act
	str, err := json.Marshal(m)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "{\"greeting\":null}", string(str))
}

func TestJSONUnmarshal(t *testing.T) {
	// Arrange
	m := NewMap()

	m.Set("greeting", "howdy")

	// Assert
	assert.Equal(t, "howdy", m.Get("greeting"))

	// Act
	err := json.Unmarshal([]byte(jsonstring), m)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, greeting, m.Get("greeting"))
	assert.IsType(t, m, m.Get("nested"))
	assert.Equal(t, loc, m.Get("nested").(*Map).Get("location"))
	assert.Equal(t, "fall", m.Get("nested").(*Map).Get("season"))
}

func TestJSONUnmarshalOmit(t *testing.T) {
	// Arrange
	m := NewMap()

	m.Set("greeting", "howdy")

	m.Set("nested", NewMap(), Omit)

	// Assert
	assert.Equal(t, "howdy", m.Get("greeting"))

	// Act
	err := json.Unmarshal([]byte(jsonstring), m)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, greeting, m.Get("greeting"))
	assert.Equal(t, 0, m.Get("nested").(*Map).Len())
}

func TestJSONUnmarshalNull(t *testing.T) {
	// Arrange
	m := NewMap()

	m.Set("greeting", "howdy")

	n := NewMap()

	m.Set("nested", n, Null)

	// Assert
	assert.Equal(t, "howdy", m.Get("greeting"))

	// Act
	err := json.Unmarshal([]byte(jsonstring), m)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, greeting, m.Get("greeting"))
	assert.Equal(t, n, m.Get("nested"))
}
