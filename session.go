package smartSock

import (
	"net"
	"sync"
	"time"
)

type Session struct {
	id uint32
	conn net.Conn
	times int64
	lock sync.Mutex
}

func NewSession(id uint32, conn net.Conn) *Session {
	return &Session{
		id: id,
		conn: conn,
		times: time.Now().Unix(),
	}
}

func (this *Session) write(msg []byte) error {
	_, err := this.conn.Write(msg)
	return err
}

func (this *Session) close() {
	_ = this.conn.Close()
}

func (this *Session) UpdateTime() {
	this.times = time.Now().Unix()
}
/*---------------------session server manager-----------------------*/
type SessionSM struct {
	ser *Msf
	sessions sync.Map
}

func NewSessionSM(msf *Msf) *SessionSM {
	if msf == nil {
		return nil
	}
	return &SessionSM{
		ser: msf,
	}
}

func (this *SessionSM) GetSessionById(id uint32) *Session {
	tem, exist := this.sessions.Load(id)
	if exist {
		if sess, ok := tem.(*Session); ok {
			return sess
		}
	}
	return nil
}

func (this *SessionSM) SetSession(fd uint32, conn net.Conn) {
	sess := NewSession(fd, conn)
	this.sessions.Store(fd, sess)
}

func (this *SessionSM) DelSessionById(id uint32) {
	tem, exist := this.sessions.Load(id)
	if exist {
		if sess, ok := tem.(*Session); ok {
			sess.close()
		}
	}
	this.sessions.Delete(id)
}

func (this *SessionSM) WriteById(id uint32, msg []byte) bool {
	msg = this.ser.SocketType.Pack(msg)
	tem, exist := this.sessions.Load(id)
	if exist {
		if sess, ok := tem.(*Session); ok {
			if err := sess.write(msg); err == nil {
				return true
			}
		}
	}
	this.DelSessionById(id)
	return false
}

func (this *SessionSM) WriteToAll(msg []byte) {
	msg = this.ser.SocketType.Pack(msg)
	this.sessions.Range(func(key, val interface{}) bool {
		if val, ok := val.(*Session); ok {
			if err := val.write(msg); err != nil {
				this.DelSessionById(key.(uint32))
			}
		}
		return true
	})
}

func (this *SessionSM) HeartBeat(num int64) {
	for {
		time.Sleep(time.Second)
		this.sessions.Range(func(key, val interface{}) bool {
			tem, ok := val.(*Session)
			if !ok {
				return true
			}
			if time.Now().Unix() - tem.times > num {
				this.DelSessionById(key.(uint32))
			}
			return true
		})
	}
}

/*------------------------------session client manager----------------------*/
type SessionCM struct {
	cli      *Mcf
	sessPtr  *Session
	isAvailable bool
	lock sync.Mutex
}
func NewSessionCM(mcf *Mcf) *SessionCM {
	if mcf == nil {
		return nil
	}
	return &SessionCM{
		cli:    mcf,
		isAvailable: false,
	}
}
func (this *SessionCM) SetFlagTrue() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.isAvailable = true
}
func (this *SessionCM) SetFlagFalse() {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.isAvailable = false
}
func (this *SessionCM) SetSession(fd uint32, conn net.Conn) {
	sess := NewSession(fd, conn)
	this.sessPtr = sess
	this.SetFlagTrue()
	logger.Infoln("create socket connection to tip")
}
func (this *SessionCM) WriteToSev(msg []byte) bool {
	msg = this.cli.SocketType.Pack(msg)
	if this.isAvailable == false {
		return false
	}
	if err := this.sessPtr.write(msg); err == nil {
		return true
	}
	return false
}
func (this *SessionCM) TimedSend(num int64) {
	msg := this.cli.SocketType.Pack([]byte(""))
	for {
		if this.isAvailable{
			err := this.sessPtr.write(msg)
			if err != nil {
				this.SetFlagFalse()
				logger.Errorln("failed to send heartbeat package")
			} else {
				logger.Infoln("send heartbeat package to tip")
			}
		}
		time.Sleep(20 * time.Second)
	}
}
func (this *SessionCM) ReConnect(address string) {
	fd := uint32(0)
	for {
		if this.isAvailable == false {
			if tcpconn, err := net.Dial("tcp", address); err == nil {
				logger.Infoln("response socket create")
				this.SetSession(fd, tcpconn)
				fd++
			}
		}
		logger.Infoln("current response connection is available")
		time.Sleep(20 * time.Second)
	}
}



