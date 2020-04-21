package smartSock

import "time"

type SockCliTypes interface {
	Pack(data []byte) []byte
	ConnHandle(mcf *Mcf, sess *Session)
}

type Mcf struct {
	SessionMaster *SessionCM
	SocketType    SockCliTypes
	targetCh      chan []byte
	HeartSpan     time.Duration
	DialAddr      string
}

func NewMcf(socketType SockCliTypes, ch chan []byte, addr string, sp time.Duration) *Mcf {
	mcf := &Mcf{
		SocketType: socketType,
		targetCh:   ch,
		HeartSpan:  sp,
		DialAddr:   addr,
	}
	mcf.SessionMaster = NewSessionCM(mcf)
	return mcf
}
func (this *Mcf) Dial(address string) {
	go this.SessionMaster.ReConnect(this.DialAddr)
	go this.SessionMaster.TimedSend(this.HeartSpan) //send heartbeat package each 20s
	go this.SocketType.ConnHandle(this, this.SessionMaster.sessPtr)
	select {}
}
