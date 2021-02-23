package functionaloptions

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFoobar(t *testing.T) {
	{
		fb, err := NewFoobar()
		require.NoError(t, err)
		require.NotNil(t, fb)
		require.True(t, fb.mutable)
		require.Equal(t, Celsius(37), fb.temperature)
	}
	{
		fb, err := NewFoobar(OptionTemperature(10), OptionReadOnlyFlag)
		require.NoError(t, err)
		require.NotNil(t, fb)
		require.False(t, fb.mutable)
		require.Equal(t, Celsius(10), fb.temperature)
	}
	{
		opts := []OptionFoobar{
			OptionReadOnlyFlag,
			OptionTemperature(10),
		}
		fb, err := NewFoobar(opts...)
		require.NoError(t, err)
		require.NotNil(t, fb)
		require.False(t, fb.mutable)
		require.Equal(t, Celsius(10), fb.temperature)
	}
}

func TestOptionReadOnlyFlag(t *testing.T) {
	fb, err := NewFoobar(OptionReadOnlyFlag)
	require.NoError(t, err)
	require.NotNil(t, fb)
	require.False(t, fb.mutable)
}

func TestOptionTemperature(t *testing.T) {
	fb, err := NewFoobar(OptionTemperature(10))
	require.NoError(t, err)
	require.NotNil(t, fb)
	require.Equal(t, Celsius(10), fb.temperature)
}
