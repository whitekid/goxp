package goxp

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type testData struct {
	Name string `json:"name" xml:"name" yaml:"name"`
	Age  int    `json:"age" xml:"age" yaml:"age"`
}

func TestReadJSON(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		jsonData := `{"name": "Alice", "age": 30}`
		reader := strings.NewReader(jsonData)

		result, err := ReadJSON[testData](reader)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, "Alice", result.Name)
		require.Equal(t, 30, result.Age)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		invalidJSON := `{"name": "Alice", "age":}`
		reader := strings.NewReader(invalidJSON)

		result, err := ReadJSON[testData](reader)
		require.Error(t, err)
		require.NotNil(t, result) // Function always returns pointer even on error
		require.Empty(t, result.Name) // But values should be zero values
	})

	t.Run("empty reader", func(t *testing.T) {
		reader := strings.NewReader("")

		result, err := ReadJSON[testData](reader)
		require.Error(t, err)
		require.NotNil(t, result) // Function always returns pointer even on error
		require.Empty(t, result.Name) // But values should be zero values
	})

	t.Run("wrong type", func(t *testing.T) {
		jsonData := `{"name": 123, "age": "invalid"}`
		reader := strings.NewReader(jsonData)

		result, err := ReadJSON[testData](reader)
		require.Error(t, err)
		require.NotNil(t, result) // Function always returns pointer even on error
	})
}

func TestReadXML(t *testing.T) {
	t.Run("valid XML", func(t *testing.T) {
		xmlData := `<testData><name>Bob</name><age>25</age></testData>`
		reader := strings.NewReader(xmlData)

		result, err := ReadXML[testData](reader)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, "Bob", result.Name)
		require.Equal(t, 25, result.Age)
	})

	t.Run("invalid XML", func(t *testing.T) {
		invalidXML := `<testData><name>Bob</name><age>25</testData>` // Missing closing tag
		reader := strings.NewReader(invalidXML)

		result, err := ReadXML[testData](reader)
		require.Error(t, err)
		require.NotNil(t, result) // Function always returns pointer even on error
	})

	t.Run("empty reader", func(t *testing.T) {
		reader := strings.NewReader("")

		result, err := ReadXML[testData](reader)
		require.Error(t, err)
		require.NotNil(t, result) // Function always returns pointer even on error
	})
}

func TestReadYAML(t *testing.T) {
	t.Run("valid YAML", func(t *testing.T) {
		yamlData := `name: Charlie
age: 35`
		reader := strings.NewReader(yamlData)

		result, err := ReadYAML[testData](reader)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, "Charlie", result.Name)
		require.Equal(t, 35, result.Age)
	})

	t.Run("invalid YAML", func(t *testing.T) {
		invalidYAML := `name: Charlie
  age: 35
invalid: [unclosed`
		reader := strings.NewReader(invalidYAML)

		result, err := ReadYAML[testData](reader)
		require.Error(t, err)
		require.NotNil(t, result) // Function always returns pointer even on error
	})

	t.Run("empty reader", func(t *testing.T) {
		reader := strings.NewReader("")

		result, err := ReadYAML[testData](reader)
		require.Error(t, err)
		require.NotNil(t, result) // Function always returns pointer even on error
	})
}

func TestWriteJSON(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		data := testData{Name: "David", Age: 40}
		var buf bytes.Buffer

		err := WriteJSON(&buf, data)
		require.NoError(t, err)

		result := buf.String()
		require.Contains(t, result, `"name":"David"`)
		require.Contains(t, result, `"age":40`)
	})

	t.Run("nil data", func(t *testing.T) {
		var buf bytes.Buffer
		var nilData *testData

		err := WriteJSON(&buf, nilData)
		require.NoError(t, err)

		result := buf.String()
		require.Equal(t, "null\n", result)
	})

	t.Run("invalid data", func(t *testing.T) {
		var buf bytes.Buffer
		invalidData := make(chan int) // Channels can't be JSON encoded

		err := WriteJSON(&buf, invalidData)
		require.Error(t, err)
	})
}

func TestWriteXML(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		data := testData{Name: "Eve", Age: 28}
		var buf bytes.Buffer

		err := WriteXML(&buf, data)
		require.NoError(t, err)

		result := buf.String()
		require.Contains(t, result, "<name>Eve</name>")
		require.Contains(t, result, "<age>28</age>")
	})

	t.Run("invalid data", func(t *testing.T) {
		var buf bytes.Buffer
		invalidData := make(chan int) // Channels can't be XML encoded

		err := WriteXML(&buf, invalidData)
		require.Error(t, err)
	})
}

func TestWriteYAML(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		data := testData{Name: "Frank", Age: 33}
		var buf bytes.Buffer

		err := WriteYAML(&buf, data)
		require.NoError(t, err)

		result := buf.String()
		require.Contains(t, result, "name: Frank")
		require.Contains(t, result, "age: 33")
	})

	t.Run("invalid data", func(t *testing.T) {
		var buf bytes.Buffer
		invalidData := make(chan int) // Channels can't be YAML encoded

		// YAML encoder panics instead of returning error
		require.Panics(t, func() {
			WriteYAML(&buf, invalidData)
		})
	})
}

func TestMustMarshalJson(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		data := testData{Name: "Grace", Age: 27}

		result := MustMarshalJson(data)
		require.NotEmpty(t, result)

		resultStr := string(result)
		require.Contains(t, resultStr, `"name":"Grace"`)
		require.Contains(t, resultStr, `"age":27`)
	})

	t.Run("nil data", func(t *testing.T) {
		result := MustMarshalJson(nil)
		require.Equal(t, []byte("null"), result)
	})

	t.Run("panic on invalid data", func(t *testing.T) {
		invalidData := make(chan int) // Channels can't be JSON encoded

		require.Panics(t, func() {
			MustMarshalJson(invalidData)
		})
	})
}

// Test generic type inference
func TestEncodingGenerics(t *testing.T) {
	t.Run("different types", func(t *testing.T) {
		// Test with map
		mapData := map[string]int{"count": 42}
		jsonStr := `{"count": 42}`
		reader := strings.NewReader(jsonStr)

		result, err := ReadJSON[map[string]int](reader)
		require.NoError(t, err)
		require.Equal(t, mapData, *result)

		// Test with slice
		sliceData := []string{"a", "b", "c"}
		jsonStr = `["a", "b", "c"]`
		reader = strings.NewReader(jsonStr)

		resultSlice, err := ReadJSON[[]string](reader)
		require.NoError(t, err)
		require.Equal(t, sliceData, *resultSlice)
	})
}

// Test round-trip encoding/decoding
func TestEncodingRoundTrip(t *testing.T) {
	original := testData{Name: "Henry", Age: 45}

	t.Run("JSON round-trip", func(t *testing.T) {
		var buf bytes.Buffer
		
		// Encode
		err := WriteJSON(&buf, original)
		require.NoError(t, err)

		// Decode
		result, err := ReadJSON[testData](&buf)
		require.NoError(t, err)
		require.Equal(t, original, *result)
	})

	t.Run("XML round-trip", func(t *testing.T) {
		var buf bytes.Buffer
		
		// Encode
		err := WriteXML(&buf, original)
		require.NoError(t, err)

		// Decode
		result, err := ReadXML[testData](&buf)
		require.NoError(t, err)
		require.Equal(t, original, *result)
	})

	t.Run("YAML round-trip", func(t *testing.T) {
		var buf bytes.Buffer
		
		// Encode
		err := WriteYAML(&buf, original)
		require.NoError(t, err)

		// Decode
		result, err := ReadYAML[testData](&buf)
		require.NoError(t, err)
		require.Equal(t, original, *result)
	})
}