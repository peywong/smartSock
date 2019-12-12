package smartSock

type SockCliTypes interface {
	Pack(data []byte)[]byte
	ConnHandle(mcf *Mcf, sess *Session)
}
type Mcf struct {
	SessionMaster *SessionCM
	SocketType SockCliTypes
	resource chan []byte
}

func NewMcf(socketType SockCliTypes, ch chan []byte) *Mcf {
	mcf := &Mcf{
		SocketType:  socketType,
		resource: ch,
	}
	mcf.SessionMaster = NewSessionCM(mcf)
	return mcf
}
func (this *Mcf) Dial(address string) {
	go this.SessionMaster.ReConnect(address)
	go this.SessionMaster.TimedSend(20)                      //send heartbeat package each 20s
	go this.SocketType.ConnHandle(this, this.SessionMaster.sessPtr)
	select{}
}
