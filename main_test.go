package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMarkdown_Complete(t *testing.T) {
	str := `
# My Header

* Fact front #1.
    * Fact back #1.
* Fact front #2.
    * Fact back #2.
* Fact front #3.
	`
	header, facts, err := parseMarkdown(str)
	expected := []*Fact{
		&Fact{Front: "Fact front #1.", Back: "Fact back #1."},
		&Fact{Front: "Fact front #2.", Back: "Fact back #2."},
		&Fact{Front: "Fact front #3."},
	}

	assert.Equal(t, "My Header", header)
	assert.Equal(t, len(expected), len(facts))
	for i, expected := range expected {
		assert.Equal(t, *expected, *facts[i])
	}
	assert.Nil(t, err)
}

func TestParseMarkdown_Empty(t *testing.T) {
	header, facts, err := parseMarkdown("")
	assert.Equal(t, "", header)
	assert.Nil(t, facts)
	assert.Nil(t, err)
}

func TestParseMarkdown_Fact_Empty(t *testing.T) {
	str := `
* 
	`
	header, _, err := parseMarkdown(str)
	assert.Equal(t, "", header)
	assert.Nil(t, err)
}

func TestParseMarkdown_Fact_FrontOnly(t *testing.T) {
	str := `
* Fact front.
	`
	_, facts, err := parseMarkdown(str)
	assert.Equal(t, 1, len(facts))
	assert.Equal(t, Fact{Front: "Fact front."}, *facts[0])
	assert.Nil(t, err)
}

func TestParseMarkdown_Fact_FrontWithBack(t *testing.T) {
	str := `
* Fact front.
    * Fact back.
	`
	_, facts, err := parseMarkdown(str)
	assert.Equal(t, 1, len(facts))
	assert.Equal(t, Fact{Front: "Fact front.", Back: "Fact back."}, *facts[0])
	assert.Nil(t, err)
}

func TestParseMarkdown_Fact_Multiple(t *testing.T) {
	str := `
* Fact front #1.
* Fact front #2.
	`
	_, facts, err := parseMarkdown(str)
	expected := []*Fact{
		&Fact{Front: "Fact front #1."},
		&Fact{Front: "Fact front #2."},
	}

	assert.Equal(t, len(expected), len(facts))
	for i, expected := range expected {
		assert.Equal(t, *expected, *facts[i])
	}
	assert.Nil(t, err)
}

func TestParseMarkdown_Fact_OrphanedBack(t *testing.T) {
	str := `
    * Fact back.
	`
	_, facts, err := parseMarkdown(str)
	assert.Equal(t, 0, len(facts))
	assert.Nil(t, err)
}

func TestParseMarkdown_Header(t *testing.T) {
	str := `
# My Header
	`
	header, _, err := parseMarkdown(str)
	assert.Equal(t, "My Header", header)
	assert.Nil(t, err)
}

func TestParseMarkdown_Header_Empty(t *testing.T) {
	str := `
# 
	`
	header, _, err := parseMarkdown(str)
	assert.Equal(t, "", header)
	assert.Nil(t, err)
}
