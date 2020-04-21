package smartSock

import (
	"net"
)

type SocketTypes interface {
	ConnHandle(msf *Msf, sess *Session)
	Pack(data []byte) []byte
}

type Msf struct {
	EventPool     *RoutersMap
	SessionMaster *SessionSM
	SocketType    SocketTypes
	//add destination chan []byte correspond to communication mode, 1 in 1 out specially.
	des chan []byte
}

func NewMsf(socketType SocketTypes, ch chan []byte) *Msf {
	msf := &Msf{
		EventPool:  NewRoutersMap(),
		SocketType: socketType,
		des:        ch,
	}
	msf.SessionMaster = NewSessionSM(msf)
	return msf
}

func (this *Msf) Listening(address string) {
	tcpListen, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("failed to create socket listener")
		panic(err)
	}
	go this.SessionMaster.HeartBeat(40) //hold heartbeat 40s
	fd := uint32(0)
	for {
		conn, err := tcpListen.Accept()
		if err != nil {
			logger.Error("tcp server accept connection fail")
			continue
		}
		if this.EventPool.OnHand(fd, conn) == false {
			return
		}
		this.SessionMaster.SetSession(fd, conn)
		go this.SocketType.ConnHandle(this, this.SessionMaster.GetSessionById(fd))
		fd++
	}
}

func (this *Msf) Hook(fd uint32, requestData []byte) {
	this.EventPool.OnMessage(fd, requestData)
	if len(this.EventPool.methods) == 0 {
		return
	}

	var result []byte
	result = this.EventPool.methods[BEFOREACTION](fd, requestData)

	result = this.EventPool.methods[DEFAULTACTION](fd, requestData)
	if result != nil {
		this.des <- result
	}

	result = this.EventPool.methods[AFTERACTION](fd, requestData)
	return
}

/*--------------------------------------------------*/
