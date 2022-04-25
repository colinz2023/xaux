package server

import (
	"encoding/json"
	"fmt"
	"net"
	"xaux/pkg/x"
)

var _ ISessionMaker = (*FakeSessionMaker)(nil)
var _ ISession = (*FakeSession)(nil)

type FakeSessionMaker struct {
	cnt uint32
}

type FakeSession struct {
	id uint32
}

func NewFakeSessionMaker() *FakeSessionMaker {
	return &FakeSessionMaker{}
}

func (f *FakeSessionMaker) MakeSession() (ISession, error) {
	f.cnt++
	return &FakeSession{id: f.cnt}, nil
}

func (f *FakeSession) ID() uint32 {
	return f.id
}

func (f *FakeSession) CommandCb(conn net.Conn, allResponse *x.AllRequest) error {
	fmt.Println("client addr =", conn.RemoteAddr().String(), ":")
	buf, _ := json.MarshalIndent(allResponse, "", " ")
	fmt.Println(string(buf))

	cmd := allResponse.Cmd
	var rspBuf []byte
	if cmd == x.CmdStart {
		startRsp := x.StartResponse{
			Cmd:       cmd,
			SessionID: f.id,
			UDPPort:   x.UDPPort,
		}
		rspBuf, _ = json.Marshal(&startRsp)
	} else if cmd == x.CmdEnd {
		endRsp := x.EndResponse{
			Cmd: cmd,
			Msg: "session end!",
		}
		rspBuf, _ = json.Marshal(&endRsp)
	}
	_, err := conn.Write(rspBuf)
	if err != nil {
		panic(err)
	}
	return nil
}
func (f *FakeSession) DataCb(data []byte, seq uint32) error {
	fmt.Println("get seq=", seq, ", data len=", len(data))
	return nil
}
