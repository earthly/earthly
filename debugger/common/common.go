package common

import (
	"encoding/binary"
	"io"
	"io/ioutil"
)

// PtyStream is a magic number to indicate the yamux stream is for pesudo-tty data
const PtyStream = 0x01

// WinChangeStream is a magic number to indicate the yamux stream is for window size change data
const WinChangeStream = 0x02

// ReadUint16PrefixedData first reads a uint16 indicating how many bytes to read
// followed by reading that many bytes
func ReadUint16PrefixedData(r io.Reader) ([]byte, error) {
	var l uint16
	err := binary.Read(r, binary.LittleEndian, &l)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(io.LimitReader(r, int64(l)))
}

// WriteUint16PrefixedData writes a uint16 indicating the length of the data payload
// followed by the data payload
func WriteUint16PrefixedData(w io.Writer, data []byte) error {
	length := uint16(len(data))
	err := binary.Write(w, binary.LittleEndian, length)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}
