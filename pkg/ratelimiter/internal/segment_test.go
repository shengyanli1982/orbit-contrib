package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSegment_GetOrCreate(t *testing.T) {
	// Create a new segment
	segment := NewSegment()
	defer segment.Stop()

	// Define a test key and value
	key := "testKey"
	value := "testValue"

	// Define a test function
	fn := func() any {
		return value
	}

	// Call the GetOrCreate function
	result, ok := segment.GetOrCreate(key, fn)

	// Assert that the result is correct
	assert.False(t, ok)
	assert.Equal(t, value, result)

	// Call the GetOrCreate function
	result, ok = segment.GetOrCreate(key, fn)

	// Assert that the result is correct
	assert.True(t, ok)
	assert.Equal(t, value, result)
}

func TestSegment_Get(t *testing.T) {
	// Create a new segment
	segment := NewSegment()
	defer segment.Stop()

	// Define a test key and value
	key := "testKey"
	value := "testValue"

	// Set the value in the segment
	segment.data[key] = value

	// Call the Get function
	result, ok := segment.Get(key)

	// Assert that the result is correct
	assert.True(t, ok)
	assert.Equal(t, value, result)

	// Delete the value from the segment
	delete(segment.data, key)

	// Call the Get function
	result, ok = segment.Get(key)

	// Assert that the result is correct
	assert.False(t, ok)
	assert.Nil(t, result)
}

func TestSegment_Set(t *testing.T) {
	// Create a new segment
	segment := NewSegment()
	defer segment.Stop()

	// Define a test key and value
	key := "testKey"
	value := "testValue"

	// Call the Set function
	segment.Set(key, value)

	// Call the Get function to retrieve the value
	result, ok := segment.Get(key)

	// Assert that the result is correct
	assert.True(t, ok)
	assert.Equal(t, value, result)
}
