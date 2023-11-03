package str

import (
	"context"
	"sync"
)

func greaterBaseCost(a, b []rune) float64 {
	if len(a) > len(b) {
		return float64(len(a))
	}
	return float64(len(b))
}

type charDiff struct {
	baseCost    float64
	perCharCost float64

	cost     float64
	sections []DiffSection
}

func (d *charDiff) calculate() {
	d.cost = 0
	for _, s := range d.sections {
		if s.Type == TypeMatch {
			continue
		}
		d.cost += d.baseCost + d.perCharCost*greaterBaseCost(s.Actual, s.Expected)
	}
}

func (d *charDiff) Cost() float64 {
	return d.cost
}

func (d *charDiff) Sections() []DiffSection {
	return d.sections
}

// broadcast is a type which can broadcast new diffs to multiple subscribers.
type broadcast struct {
	mu sync.Mutex

	closed bool
	curr   Diff
	subs   []chan Diff
}

// subscribe subscribes to an existing broadcast, returning a channel to listen
// for changes on. The current value will be sent on the channel immediately.
func (b *broadcast) subscribe() chan Diff {
	ch := make(chan Diff, 1)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs = append(b.subs, ch)

	if b.curr != nil {
		ch <- b.curr
	}
	if b.closed {
		close(ch)
	}
	return ch
}

// send sends d to all subscribers and updates the current value for new
// subscribers.
func (b *broadcast) send(ctx context.Context, d Diff) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.curr = d

	for _, s := range b.subs {
		select {
		case s <- d:
		case <-ctx.Done():
			return
		}
	}
}

// done signals that b has exhausted all possibilities and all subscribers
// should be closed.
func (b *broadcast) done() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.closed = true
	for _, s := range b.subs {
		close(s)
	}
}

type diffIdx struct {
	aStart, eStart int
}

// CharDiffOpt is an option function for changing the behavior of the
// NewCharDiff constructor.
type CharDiffOpt func(CharDiff) CharDiff

// CharDiffBaseCost is a CharDiff option to set the base cost per diff section.
// Increasing this will reduce the number of diff sections in the output at the
// cost of larger diff sections.
//
// Default is 0.
func CharDiffBaseCost(cost float64) CharDiffOpt {
	return func(d CharDiff) CharDiff {
		d.baseCost = cost
		return d
	}
}

// CharDiffPerCharCost is a CharDiff option to set the cost-per-character of any
// differences returned. Increasing this cost will reduce the size of diff
// sections at the cost of more diff sections.
//
// Default is 1
func CharDiffPerCharCost(cost float64) CharDiffOpt {
	return func(d CharDiff) CharDiff {
		d.perCharCost = cost
		return d
	}
}

// CharDiff is a per-character diff algorithm, meaning that it makes no distinctions
// about word or line boundaries when generating a diff.
type CharDiff struct {
	baseCost    float64
	perCharCost float64
}

func NewCharDiff(opts ...CharDiffOpt) *CharDiff {
	d := CharDiff{
		baseCost:    0,
		perCharCost: 1,
	}
	for _, o := range opts {
		d = o(d)
	}
	return &d
}

func (c *CharDiff) Diffs(ctx context.Context, actual, expected []rune) <-chan Diff {
	ch := make(chan Diff)
	var m sync.Map
	go c.sendDiffs(ctx, ch, &m, actual, expected, 0, 0)
	return ch
}

func (c *CharDiff) sendBestResults(ctx context.Context, ch chan<- Diff, bcast *broadcast, baseSections []DiffSection) {
	defer close(ch)

	subCh := bcast.subscribe()
	var cheapest *charDiff
	for {
		select {
		case subDiff, ok := <-subCh:
			if !ok {
				return
			}
			diff := &charDiff{
				baseCost:    c.baseCost,
				perCharCost: c.perCharCost,
				sections:    append([]DiffSection(nil), baseSections...),
			}
			diff.sections = append(diff.sections, subDiff.Sections()...)
			diff.calculate()
			if cheapest == nil || cheapest.Cost() > diff.Cost() {
				cheapest = diff
				select {
				case ch <- cheapest:
				case <-ctx.Done():
					return
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *CharDiff) runBroadcast(ctx context.Context, bcast *broadcast, ch <-chan Diff, actual, expected []rune, actualStart, expectedStart int) {
	defer bcast.done()

	base := &charDiff{
		baseCost:    c.baseCost,
		perCharCost: c.perCharCost,
		sections: []DiffSection{
			{Type: TypeReplace, Actual: actual[actualStart:], Expected: expected[expectedStart:]},
		},
	}
	base.calculate()
	shortest := Diff(base)
	bcast.send(ctx, shortest)

	if ctx.Err() != nil {
		return
	}

	for diff := range ch {
		if ctx.Err() != nil {
			return
		}
		if diff.Cost() >= shortest.Cost() {
			continue
		}
		shortest = diff
		bcast.send(ctx, shortest)
	}
}

func (c *CharDiff) sendSubDiffs(ctx context.Context, wg *sync.WaitGroup, subCh <-chan Diff, results chan<- Diff, section DiffSection) {
	defer wg.Done()

	for {
		select {
		case subDiff, ok := <-subCh:
			if !ok {
				return
			}
			diff := &charDiff{
				baseCost:    c.baseCost,
				perCharCost: c.perCharCost,
				sections:    append([]DiffSection{section}, subDiff.Sections()...),
			}
			diff.calculate()
			select {
			case results <- diff:
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *CharDiff) sendDiffs(ctx context.Context, ch chan<- Diff, cache *sync.Map, actual, expected []rune, actualStart, expectedStart int) {
	actualEnd, expectedEnd := actualStart, expectedStart
	for actualEnd < len(actual) && expectedEnd < len(expected) && actual[actualEnd] == expected[expectedEnd] {
		actualEnd++
		expectedEnd++
	}
	if actualEnd == len(actual) && expectedEnd == len(expected) {
		if actualEnd-actualStart > 0 || expectedEnd-expectedStart > 0 {
			diff := &charDiff{
				baseCost:    c.baseCost,
				perCharCost: c.perCharCost,
				sections:    []DiffSection{{Type: TypeMatch, Actual: actual[actualStart:actualEnd], Expected: expected[expectedStart:expectedEnd]}},
			}
			select {
			case ch <- diff:
			case <-ctx.Done():
			}
		}
		close(ch)
		return
	}
	bcast := &broadcast{}
	cached, running := cache.LoadOrStore(diffIdx{aStart: actualEnd, eStart: expectedEnd}, bcast)
	bcast = cached.(*broadcast)

	var baseSections []DiffSection
	if actualEnd-actualStart > 0 || expectedEnd-expectedStart > 0 {
		baseSections = []DiffSection{
			{Type: TypeMatch, Actual: actual[actualStart:actualEnd], Expected: expected[expectedStart:expectedEnd]},
		}
	}
	go c.sendBestResults(ctx, ch, bcast, baseSections)

	if running {
		return
	}

	subCh := make(chan Diff)
	go c.runBroadcast(ctx, bcast, subCh, actual, expected, actualEnd, expectedEnd)

	var wg sync.WaitGroup
	for i := actualEnd; i < len(actual); i++ {
		for j := expectedEnd; j < len(expected); j++ {
			if ctx.Err() != nil {
				return
			}
			if actual[i] != expected[j] {
				continue
			}
			subSubCh := make(chan Diff)
			wg.Add(1)
			go c.sendSubDiffs(ctx, &wg, subSubCh, subCh, DiffSection{
				Type:     TypeReplace,
				Actual:   actual[actualEnd:i],
				Expected: expected[expectedEnd:j],
			})
			c.sendDiffs(ctx, subSubCh, cache, actual, expected, i, j)
		}
	}

	go closeAfter(subCh, &wg)
}

func closeAfter(ch chan<- Diff, wg *sync.WaitGroup) {
	wg.Wait()
	close(ch)
}
