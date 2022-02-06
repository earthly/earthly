package consumer

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/earthly/earthly/logstream/api"
	"github.com/earthly/earthly/logstream/producer"
	"github.com/stretchr/testify/assert"
)

func TestConsumer(t *testing.T) {
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

	onManifestCh := make(chan api.Manifest, 2)
	onManifestChanges := func(m api.Manifest, d api.Delta) {
		onManifestCh <- m
	}
	c := NewConsumer(ctx, pCh, onManifestChanges)
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
		assert.Equal(t, "hello world and the entire universe!", string(text))
		close(doneRead)
	}()
	w.Write([]byte(" and the entire"))
	w.Write([]byte(" universe!"))
	err := p.Close()
	assert.NoError(t, err)
	<-doneRead
}
