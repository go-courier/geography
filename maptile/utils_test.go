package maptile

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type LayerXXAttrs struct {
	Class string  `name:"class"`
	Ele   float64 `name:"ele"`
	Is    bool    `name:"bool"`
}

func TestStructToFields(t *testing.T) {
	require.Equal(t, map[string]FieldType{
		"class": FieldTypeString,
		"ele":   FieldTypeNumber,
		"bool":  FieldTypeBoolean,
	}, StructToFields(LayerXXAttrs{}))
}

func TestStructToProperties(t *testing.T) {
	require.Equal(t, map[string]interface{}{
		"class": "test",
		"ele":   float64(11),
		"bool":  true,
	}, StructToProperties(LayerXXAttrs{
		Class: "test",
		Ele:   11,
		Is:    true,
	}))
}

type LayerXXAttrsWithOmitempty struct {
	Class string  `name:"class,omitempty"`
	Ele   float64 `name:"ele,omitempty"`
	Is    bool    `name:"bool,omitempty"`
}

func TestStructToPropertiesWithOmitEmpty(t *testing.T) {
	require.Equal(t, map[string]interface{}{}, StructToProperties(LayerXXAttrsWithOmitempty{
	}))
}
