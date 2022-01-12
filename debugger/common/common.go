package common

import (
	"bytes"
	"encoding/binary"
	"io"
)

//******************************************************************************************
// Magic numbers used in the shell repeater network protocol

// Protocol handshake:
// <byte:connection_type> The first byte identifies the type of connection (i.e. is it a terminal or the shell)
// next comes any number of data packets of the form:
// <byte:data_packet_type><uint16:n><n bytes of data>

////////////////////////////////////////////////////////////////////
// connection_type identifiers

// TermID is a magic byte to identify the connection is from the terminal
const TermID = 0x01

// ShellID is a magic byte to identify the connection is from the shell
const ShellID = 0x02

////////////////////////////////////////////////////////////////////
// data packet identifiers

// StartShellSession identifies the start of a shell session data packet
const StartShellSession = 0x01

// EndShellSession identifies the end of a shell session data packet
const EndShellSession = 0x02

// PtyData identifies the psuedo terminal (pty) data payload packet
const PtyData = 0x03

// WinSizeData identifies the terminal window data payload packet
const WinSizeData = 0x04

// End of network protocol magic numbers
//******************************************************************************************

func readUint16PrefixedData(r io.Reader) ([]byte, error) {
	var l uint16
	err := binary.Read(r, binary.LittleEndian, &l)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(io.LimitReader(r, int64(l)))
}

// ReadDataPacket decodes a data packet from the reader
func ReadDataPacket(r io.Reader) (int, []byte, error) {
	var connDataType uint16
	err := binary.Read(r, binary.LittleEndian, &connDataType)
	if err != nil {
		return 0, nil, err
	}
	data, err := readUint16PrefixedData(r)
	if err != nil {
		return 0, nil, err
	}
	return int(connDataType), data, nil
}

// WriteDataPacket writes a data packet to the writer
func WriteDataPacket(w io.Writer, n int, data []byte) error {
	err := binary.Write(w, binary.LittleEndian, uint16(n))
	if err != nil {
		return err
	}
	return writeUint16PrefixedData(w, data)
}

func writeUint16PrefixedData(w io.Writer, data []byte) error {
	length := uint16(len(data))
	err := binary.Write(w, binary.LittleEndian, length)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// SerializeDataPacket returns a serialized a data packet
func SerializeDataPacket(payloadID int, data []byte) ([]byte, error) {
	var b bytes.Buffer
	err := WriteDataPacket(&b, payloadID, data)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
