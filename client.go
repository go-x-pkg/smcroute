package smcroute

import (
	"bytes"
	"io"
	"net"
	"time"
)

type Client struct {
	socketPath string
	conn       net.Conn
}

func (c *Client) SetSocketPath(v string) *Client { c.socketPath = v; return c }

func (c *Client) Conn() (net.Conn, *Message) {
	if c.conn == nil {
		conn, err := net.Dial(defaultNetwork, c.socketPath)

		if err != nil {
			return nil, Errorf(ErrorSocketConnect, c.socketPath, err.Error())
		}

		c.conn = conn
	}

	return c.conn, nil
}

func (c *Client) Exec(cmd *Cmd) (*bytes.Buffer, *Message) {
	start := time.Now()

	conn, err := c.Conn()
	if err != nil {
		return nil, err
	}

	// TODO: reuse connection
	//       now:
	//       1. Exec => OK
	//       2. Exec => hangs:
	//            ...
	//            write(3, "(\0\0\0\0\0\0\0j\0\2\0\0\0\0\0eth0.33\000239.255."..., 40) = 40
	//            read(3, 0xc208060200, 255)              = -1 EAGAIN (Resource temporarily unavailable)
	//            epoll_wait(4, {}, 128, 0)               = 0
	//            epoll_wait(4,
	defer func() {
		conn.Close()
		c.conn = nil
	}()

	data, e := cmd.Encode()
	if e != nil {
		return nil, Errorf(ErrorCmdEncode, cmd.StringBash(), e.Error())
	}

	if _, err := io.Copy(conn, data); err != nil {
		return nil, Errorf(ErrorSocketWrite, cmd.StringBash(), c.socketPath, err.Error())
	}

	respBuf := make([]byte, responseBufferSize)
	readed, e := conn.Read(respBuf)
	if e != nil {
		return nil, Errorf(ErrorSocketRead, cmd.StringBash(), c.socketPath, e.Error())
	}
	respBuf = bytes.Trim(respBuf, nullCharacterString)
	resp := bytes.NewBuffer(respBuf)

	latency := time.Since(start)

	log().Infof(`{"cmd": "%s", "response": "%s", "latency": "%s"}`,
		cmd.StringBash(), resp.String(), latency)

	if readed != 1 || resp.Len() != 0 {
		// <match>
		if reErrorDropMembershipFailed99.Match(resp.Bytes()) {
			return resp, Errorf(ErrorDropMembershipFailed99)
		} else if reErrorFailedLeaveNotAMember.Match(resp.Bytes()) {
			return resp, Errorf(ErrorFailedLeaveNotAMember, resp.String())
		}
		// </match>

		return resp, Errorf(ErrorExec, cmd.StringBash(), latency, resp.String())
	}

	return resp, nil
}

func NewClient() *Client {
	c := new(Client)
	c.socketPath = defaultSocketPath
	return c
}
