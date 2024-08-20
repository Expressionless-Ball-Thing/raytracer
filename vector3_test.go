package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCross(t *testing.T) {

	assert.Equal(t, Cross(&vec3{3, -3, 1}, &vec3{4, 9, 2}), &vec3{-15, -2, 39})

}

func TestDot(t *testing.T) {

	assert.Equal(t, Dot(&vec3{3, -3, 1}, &vec3{4, 9, 2}), float64(-13))
}
