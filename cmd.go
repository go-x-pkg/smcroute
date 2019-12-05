package smcroute

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"
	"unsafe"
)

type header struct {
	length    int64
	cmd       CmdKind
	argsCount uint16
}

type Cmd struct {
	header

	args     []byte
	argsOrig []string
}

func (cmd *Cmd) StringBash() string {
	return fmt.Sprintf(
		"sudo smcroute -%c %s",
		cmd.cmd,
		strings.Join(cmd.argsOrig, " "),
	)
}

// -j eth0.33 239.255.11.101
//
//  +----+-----+---+-------------------------------+
//  | 40 | 'j' | 2 | "eth0.33\0239.255.11.101\0\0" |
//  +----+-----+---+-------------------------------+
//  ^              ^
//  |              |
//  |              |
//  +-----cmd------+
//
//  sizeof(struct cmd) = 16
//  strlen(args) = 21
//  sizeof(3 NULL_CHARACTERS) = 3
//    2 => after each string arument
//    1 => at the end of arfs bye array
//
//  length = 16 + 21 + 3 = 40 bytes
//
//  strace: write(3, "(\0\0\0\0\0\0\0j\0\2\0\0\0\0\0eth0.33\000239.255.11.101\0\0", 40) = 40
func (cmd *Cmd) Encode() (*bytes.Buffer, error) {
	b := make([]byte, cmd.length)
	buf := bytes.NewBuffer(b)
	buf.Reset()

	if err := binary.Write(buf, binary.LittleEndian, cmd.header); err != nil {
		return nil, err
	}

	// add 0 padding
	// e.g. 1. sizeof(struct{size_t, uint16, uint16}) in C is 16 bytes
	//      2. sizeof(size_t) + sizeof(uint16) + sizeof(uint16) = 8 + 2 + 2 = 12 bytes
	//      3. 16 - 12 = 4 bytes of padding should be added
	structPaddedSize := int(unsafe.Sizeof(cmd.header))
	for structPaddedSize-buf.Len() > 0 {
		if err := buf.WriteByte(nullCharacter); err != nil {
			return nil, err
		}
	}

	if _, err := buf.Write(cmd.args); err != nil {
		return nil, err
	}

	return buf, nil
}

func NewCmd(cmdKind CmdKind, args ...string) *Cmd {
	cmd := new(Cmd)

	cmd.header.cmd = cmdKind
	cmd.header.argsCount = uint16(len(args))

	for _, arg := range args {
		cmd.args = append(cmd.args, []byte(arg)...)
		cmd.args = append(cmd.args, nullCharacter)
	}
	cmd.args = append(cmd.args, nullCharacter)
	cmd.argsOrig = args

	cmd.header.length = int64(unsafe.Sizeof(cmd.header)) + int64(len(cmd.args))

	return cmd
}
