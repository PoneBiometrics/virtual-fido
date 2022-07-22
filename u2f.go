package main

import (
	"bytes"
	"fmt"
)

type U2FCommand uint8

const (
	U2F_COMMAND_REGISTER     U2FCommand = 0x01
	U2F_COMMAND_AUTHENTICATE U2FCommand = 0x02
	U2F_COMMAND_VERSION      U2FCommand = 0x03
)

var u2fCommandDescriptions = map[U2FCommand]string{
	U2F_COMMAND_REGISTER:     "U2F_COMMAND_REGISTER",
	U2F_COMMAND_AUTHENTICATE: "U2F_COMMAND_AUTHENTICATE",
	U2F_COMMAND_VERSION:      "U2F_COMMAND_VERSION",
}

type U2FStatusWord uint16

const (
	U2F_SW_NO_ERROR                 U2FStatusWord = 0x9000
	U2F_SW_CONDITIONS_NOT_SATISFIED U2FStatusWord = 0x6985
	U2F_SW_WRONG_DATA               U2FStatusWord = 0x6A80
	U2F_SW_WRONG_LENGTH             U2FStatusWord = 0x6700
	U2F_SW_CLA_NOT_SUPPORTED        U2FStatusWord = 0x6E00
	U2F_SW_INS_NOT_SUPPORTED        U2FStatusWord = 0x6D00
)

type U2FMessageHeader struct {
	Cla     uint8
	Command U2FCommand
	Param1  uint8
	Param2  uint8
}

func (header U2FMessageHeader) String() string {
	return fmt.Sprintf("U2FMessageHeader{ Cla: 0x%x, Command: %s, Param1: %d, Param2: %d }",
		header.Cla,
		u2fCommandDescriptions[header.Command],
		header.Param1,
		header.Param2)
}

func decodeU2FMessage(messageBytes []byte) (U2FMessageHeader, []byte, uint16) {
	buffer := bytes.NewBuffer(messageBytes)
	header := readBE[U2FMessageHeader](buffer)
	if buffer.Len() == 0 {
		// No reqest length, no reponse length
		return header, []byte{}, 0
	}
	// We should either have a request length or reponse length, so we have at least
	// one '0' byte at the start
	if read(buffer, 1)[0] != 0 {
		panic(fmt.Sprintf("Invalid U2F Payload length: %s %#v", header, messageBytes))
	}
	length := readBE[uint16](buffer)
	if buffer.Len() == 0 {
		// No payload, so length must be the response length
		return header, []byte{}, length
	}
	// length is the request length
	request := read(buffer, uint(length))
	if buffer.Len() == 0 {
		return header, request, 0
	}
	responseLength := readBE[uint16](buffer)
	return header, request, responseLength
}

func processU2FMessage(message []byte) []byte {
	header, _, _ := decodeU2FMessage(message)
	fmt.Printf("U2F MESSAGE: %s\n\n", header)
	switch header.Command {
	case U2F_COMMAND_VERSION:
		response := append([]byte("U2F_V2"), toBE(U2F_SW_NO_ERROR)...)
		fmt.Printf("U2F RESPONSE: %#v\n\n", response)
		return response
	default:
		panic(fmt.Sprintf("Invalid U2F Command: %#v", header))
	}
}