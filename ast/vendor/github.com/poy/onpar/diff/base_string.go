package diff

import "github.com/poy/onpar/diff/str"

type baseDiff struct {
	actual, expected []rune

	sections []str.DiffSection
	cost     float64
}

func (d *baseDiff) equal() bool {
	if len(d.actual) != len(d.expected) {
		return false
	}
	for i, ar := range d.actual {
		if ar != d.expected[i] {
			return false
		}
	}
	return true
}

func (d *baseDiff) calculate() {
	if d.sections != nil {
		return
	}
	d.sections = []str.DiffSection{{Actual: d.actual, Expected: d.expected}}
	if d.equal() {
		d.sections[0].Type = str.TypeMatch
		return
	}
	d.sections[0].Type = str.TypeReplace
	if len(d.actual) > len(d.expected) {
		d.cost = float64(len(d.actual))
		return
	}
	d.cost = float64(len(d.expected))
}

func (d *baseDiff) Cost() float64 {
	d.calculate()
	return d.cost
}

func (d *baseDiff) Sections() []str.DiffSection {
	d.calculate()
	return d.sections
}

func baseStringDiff(actual, expected []rune) str.Diff {
	return &baseDiff{actual: actual, expected: expected}
}
