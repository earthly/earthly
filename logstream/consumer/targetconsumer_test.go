package consumer

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/earthly/earthly/logstream/producer"
	"github.com/stretchr/testify/assert"
)

func TestTargetConsumer(t *testing.T) {
	ctx := context.Background()
	p := producer.NewProducer()
	targetID := "target-1"
	otherTargetID := "target-2"
	w := p.GetWriter(targetID)
	w.Write([]byte("hello"))
	w.Write([]byte(" world"))
	w2 := p.GetWriter(otherTargetID)
	w2.Write([]byte("SHOULD NOT BE PASSED ALONG"))
	pCh := p.DeltaCh()

	r := NewTargetConsumerReader(ctx, targetID, pCh)

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

func ExampleNewTargetConsumerReader() {
	ctx := context.Background()
	p := producer.NewProducer()
	targetID := "target-1"
	otherTargetID := "target-2"
	w := p.GetWriter(targetID)
	w.Write([]byte("hello"))
	w.Write([]byte(" world"))
	w2 := p.GetWriter(otherTargetID)
	w2.Write([]byte("SHOULD NOT BE PASSED ALONG"))
	pCh := p.DeltaCh()

	r2 := io.TeeReader(NewTargetConsumerReader(ctx, targetID, pCh), os.Stdout)
	go func() {
		ioutil.ReadAll(r2)
	}()

	w.Write([]byte(" and the entire"))
	w.Write([]byte(" universe!"))
	err := p.Close()
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second)
	// Output: hello world and the entire universe!
}
