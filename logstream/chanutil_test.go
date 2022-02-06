package logstream

import (
	"context"
	"io/ioutil"
	"testing"
	"time"

	"github.com/earthly/earthly/logstream/api"
	"github.com/earthly/earthly/logstream/chanutil"
	"github.com/earthly/earthly/logstream/consumer"
	"github.com/earthly/earthly/logstream/producer"
	"github.com/stretchr/testify/assert"
)

func TestBatcher(t *testing.T) {
	ctx := context.Background()
	p := producer.NewProducer()
	p.SetBuildMetadata(&api.DeltaManifest{
		Status: api.StatusInProgress,
	})
	targetID := "target-1"
	targetName := "+some-target"
	p.SetTargetMetadata(targetID, &api.DeltaTargetManifest{
		Name: targetName,
	})
	w := p.GetWriter(targetID)
	w.Write([]byte("hello"))
	w.Write([]byte(" world"))
	pCh := p.DeltaCh()

	batchedCh := chanutil.IntervalBatcher(ctx, 500*time.Millisecond, pCh)

	onManifestCh := make(chan api.Manifest, 2)
	onManifestChanges := func(m api.Manifest, d api.Delta) {
		onManifestCh <- m
	}
	c := consumer.NewConsumer(ctx, batchedCh, onManifestChanges)
	<-onManifestCh
	m := c.GetManifest()
	assert.Equal(t, api.StatusInProgress, m.Status)
	assert.NotNil(t, m.Targets)
	assert.NotNil(t, m.Targets[targetID])
	assert.Equal(t, targetName, m.Targets[targetID].Name)

	r := c.GetReader(targetID)
	doneRead := make(chan struct{})
	go func() {
		text, err := ioutil.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, "hello world and the entire universe!", string(text))
		close(doneRead)
	}()
	w.Write([]byte(" and the entire"))
	w.Write([]byte(" universe!"))
	err := p.Close()
	assert.NoError(t, err)
	<-doneRead
}

func TestSplitter(t *testing.T) {
	ctx := context.Background()
	p := producer.NewProducer()
	p.SetBuildMetadata(&api.DeltaManifest{
		Status: api.StatusInProgress,
	})
	targetID := "target-1"
	targetName := "+some-target"
	p.SetTargetMetadata(targetID, &api.DeltaTargetManifest{
		Name: targetName,
	})
	w := p.GetWriter(targetID)
	w.Write([]byte("hello"))
	w.Write([]byte(" world"))
	pCh := p.DeltaCh()

	clonesCh := chanutil.Splitter(ctx, pCh, 2, 100)
	ch0 := clonesCh[0]
	ch1 := clonesCh[1]

	onManifestCh0 := make(chan api.Manifest, 2)
	onManifestChanges0 := func(m api.Manifest, d api.Delta) {
		onManifestCh0 <- m
	}
	c0 := consumer.NewConsumer(ctx, ch0, onManifestChanges0)
	<-onManifestCh0
	<-onManifestCh0
	m := c0.GetManifest()
	assert.Equal(t, api.StatusInProgress, m.Status)
	assert.NotNil(t, m.Targets)
	assert.NotNil(t, m.Targets[targetID])
	assert.Equal(t, targetName, m.Targets[targetID].Name)

	r0 := c0.GetReader(targetID)
	doneRead0 := make(chan struct{})
	go func() {
		text, err := ioutil.ReadAll(r0)
		assert.NoError(t, err)
		assert.Equal(t, "hello world and the entire universe!", string(text))
		close(doneRead0)
	}()
	w.Write([]byte(" and the entire"))
	w.Write([]byte(" universe!"))

	onManifestCh1 := make(chan api.Manifest, 2)
	onManifestChanges1 := func(m api.Manifest, d api.Delta) {
		onManifestCh1 <- m
	}
	c1 := consumer.NewConsumer(ctx, ch1, onManifestChanges1)
	<-onManifestCh1
	<-onManifestCh1
	m = c1.GetManifest()
	assert.Equal(t, api.StatusInProgress, m.Status)
	assert.NotNil(t, m.Targets)
	assert.NotNil(t, m.Targets[targetID])
	assert.Equal(t, targetName, m.Targets[targetID].Name)

	r1 := c1.GetReader(targetID)
	doneRead1 := make(chan struct{})
	go func() {
		text, err := ioutil.ReadAll(r1)
		assert.NoError(t, err)
		assert.Equal(t, "hello world and the entire universe!", string(text))
		close(doneRead1)
	}()

	err := p.Close()
	assert.NoError(t, err)
	<-doneRead0
	<-doneRead1
}

func TestFilterManifest(t *testing.T) {
	ctx := context.Background()
	p := producer.NewProducer()
	p.SetBuildMetadata(&api.DeltaManifest{
		Status: api.StatusInProgress,
	})
	targetID := "target-1"
	targetName := "+some-target"
	p.SetTargetMetadata(targetID, &api.DeltaTargetManifest{
		Name: targetName,
	})
	w := p.GetWriter(targetID)
	w.Write([]byte("hello"))
	w.Write([]byte(" world"))
	pCh := p.DeltaCh()

	filterCh := chanutil.Filter(ctx, pCh, true, false, nil)

	onManifestCh := make(chan api.Manifest, 2)
	onManifestChanges := func(m api.Manifest, d api.Delta) {
		onManifestCh <- m
	}
	c := consumer.NewConsumer(ctx, filterCh, onManifestChanges)
	<-onManifestCh
	<-onManifestCh
	m := c.GetManifest()
	assert.Equal(t, api.StatusInProgress, m.Status)
	assert.NotNil(t, m.Targets)
	assert.NotNil(t, m.Targets[targetID])
	assert.Equal(t, targetName, m.Targets[targetID].Name)

	r := c.GetReader(targetID)
	doneRead := make(chan struct{})
	go func() {
		text, err := ioutil.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, "", string(text))
		close(doneRead)
	}()
	w.Write([]byte(" and the entire"))
	w.Write([]byte(" universe!"))
	err := p.Close()
	assert.NoError(t, err)
	<-doneRead
}

func TestFilterAllTargets(t *testing.T) {
	ctx := context.Background()
	p := producer.NewProducer()
	p.SetBuildMetadata(&api.DeltaManifest{
		Status: api.StatusInProgress,
	})
	targetID := "target-1"
	targetName := "+some-target"
	p.SetTargetMetadata(targetID, &api.DeltaTargetManifest{
		Name: targetName,
	})
	w := p.GetWriter(targetID)
	w.Write([]byte("hello"))
	w.Write([]byte(" world"))
	pCh := p.DeltaCh()

	filterCh := chanutil.Filter(ctx, pCh, false, true, nil)

	onManifestChanges := func(m api.Manifest, d api.Delta) {
		t.Fail()
	}
	c := consumer.NewConsumer(ctx, filterCh, onManifestChanges)
	m := c.GetManifest()
	time.Sleep(time.Millisecond * 500)
	assert.Equal(t, api.Status(""), m.Status)

	r := c.GetReader(targetID)
	doneRead := make(chan struct{})
	go func() {
		text, err := ioutil.ReadAll(r)
		assert.NoError(t, err)
		assert.Equal(t, "hello world and the entire universe!", string(text))
		close(doneRead)
	}()
	w.Write([]byte(" and the entire"))
	w.Write([]byte(" universe!"))
	err := p.Close()
	assert.NoError(t, err)
	<-doneRead
}
