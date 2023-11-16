package statsstreamparser

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/alexcb/binarystream"
	"github.com/containerd/go-runc"
)

type Parser struct {
	buf                 *bytes.Buffer
	bsr                 *binarystream.BinaryStream
	readProtocolVersion bool
}

func New() *Parser {
	buf := bytes.NewBuffer(nil)
	return &Parser{
		buf: buf,
		bsr: binarystream.NewReader(buf, binary.LittleEndian),
	}
}

func (ssp *Parser) Parse(b []byte) ([]*runc.Stats, error) {
	_, err := ssp.buf.Write(b)
	if err != nil {
		return nil, err
	}
	var stats []*runc.Stats
	for {
		if !ssp.readProtocolVersion {
			protocolVersion, err := ssp.bsr.ReadUint8()
			if err != nil {
				if err == binarystream.ErrBufferUnderflow {
					break
				}
				return nil, err
			}
			if protocolVersion != 1 {
				return nil, fmt.Errorf("unexpected stats stream protocol version %d", protocolVersion)
			}
			ssp.readProtocolVersion = true
		}
		statsStreamJSON, err := ssp.bsr.ReadUint32PrefixedString()
		if err != nil {
			if err == binarystream.ErrBufferUnderflow {
				break
			}
			return nil, err
		}
		var runcStat runc.Stats
		err = json.Unmarshal([]byte(statsStreamJSON), &runcStat)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &runcStat)
		ssp.readProtocolVersion = false
	}
	return stats, nil
}
