package goimDecoder

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
)

type goimProtocol struct {
	packageLength   uint32
	headerLength    uint16
	protocolVersion uint16
	operation       uint32
	sequenceId      uint32
	body            []byte
}

type frameDecoder interface {
	Decode(conn *net.Conn)
}

type goimDecoder struct {
}

func (decoder *goimDecoder) Decode(conn net.Conn) {
	defer conn.Close()
	for {
		reader := bufio.NewReader(conn)
		for {
			goimPro, err := decodeGoimProtocol(reader)
			if err == io.EOF {
				return
			}

			if err != nil {
				fmt.Println("decode goim protocol error:", err)
				return
			}

			fmt.Println("receive goim protocol message:", *goimPro)
		}
	}
}

func decodeGoimProtocol(reader *bufio.Reader) (*goimProtocol, error) {
	// 包长度
	pkgLenByte, _ := reader.Peek(4)
	pkgLenBuff := bytes.NewBuffer(pkgLenByte)
	var pkgLen uint32
	err := binary.Read(pkgLenBuff, binary.LittleEndian, &pkgLen)
	if err != nil {
		return nil, err
	}

	if uint32(reader.Buffered()) < pkgLen+4 {
		return nil, errors.New("frame length invalid")
	}

	// 根据包长度读取一个完整的包
	pack := make([]byte, int(4+pkgLen))
	_, err = reader.Read(pack)
	if err != nil {
		return nil, err
	}

	headerLength, _ := strconv.Atoi(string(pack[4:6]))
	protocolVersion, _ := strconv.Atoi(string(pack[6:8]))
	operation, _ := strconv.Atoi(string(pack[8:12]))
	sequenceId, _ := strconv.Atoi(string(pack[12:16]))

	return &goimProtocol{
		packageLength:   pkgLen,
		headerLength:    uint16(headerLength),
		protocolVersion: uint16(protocolVersion),
		operation:       uint32(operation),
		sequenceId:      uint32(sequenceId),
		body:            pack[16:],
	}, nil
}
