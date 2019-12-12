package smartSock

import (
	"fmt"
	"strconv"
	"time"
)

const CONSTMLENGTH = 4

type CommSocket struct {
}

func (this *CommSocket) ConnHandle(msf *Msf, sess *Session) {
	defer func() {
		msf.SessionMaster.DelSessionById(sess.id)
		logger.Warnf("connection id: %d, is closed", sess.id)
		msf.EventPool.OnClose(sess.id)
	}()
	var errs error
	var tempBuff []byte
	var data []byte
	for {
		readBuff := make([]byte, 1000)
		n, err := sess.conn.Read(readBuff)
		if err != nil {
			logger.Errorln("socket read failed")
			return
		}
		sess.UpdateTime()
		tempBuff = append(tempBuff,readBuff[:n]...)
		for {
			tempBuff, data, errs = this.Depack(tempBuff)
			if errs != nil {
				logger.Errorln("depack package fail")
				return
			}
			if data == nil {                        //message incomplete
				break
			}
			if len(data) == 0{                      //heart-beat message
				continue
			}
			go msf.Hook(sess.id, data)                           //request message
			continue
		}
	}
}

func (this *CommSocket) Pack(message []byte) []byte {
	return append(this.Int2Bytes(len(message)), message...)
}

func (this *CommSocket) Depack(buff []byte)([]byte, []byte, error) {
	length := len(buff)
	if length < CONSTMLENGTH {
		return buff,nil, nil
	}
	msgLength,err := this.Bytes2Int(buff[:CONSTMLENGTH])
	if err != nil {
		return buff, nil, err
	}
	if msgLength == 0{
		return buff[CONSTMLENGTH:],[]byte(""),nil
	}
	if length < CONSTMLENGTH + msgLength {
		return buff, nil, nil
	}
	data := buff[CONSTMLENGTH:CONSTMLENGTH + msgLength]
	buffs := buff[CONSTMLENGTH + msgLength:]
	return buffs, data, nil
}

func (this *CommSocket) Int2Bytes(n int) []byte{
	return []byte(fmt.Sprintf("%04d", n))
}

func (this *CommSocket) Bytes2Int(b []byte) (int, error) {
	return strconv.Atoi(string(b))
}
/*----------------------------------------------------*/

type CliSocket struct{
}

func (this *CliSocket) ConnHandle(mcf *Mcf, sess *Session) {
	//defer mcf.SessionMaster.sessPtr.close()
	for{
		if mcf.SessionMaster.sessPtr == nil || mcf.SessionMaster.isAvailable == false {
			time.Sleep(100 * time.Millisecond)
			logger.Infoln("current connection is unavailable")
		} else {
			select {
			case v := <- mcf.resource:
				v = mcf.SocketType.Pack(v)
				err := mcf.SessionMaster.sessPtr.write(v)
				if err != nil {
					logger.Errorln("socket write failed, response connect will close")
					mcf.SessionMaster.SetFlagFalse()
					continue
				}
			}
		}
	}
}

func (this *CliSocket) Pack(message []byte) []byte{
	return append(this.IntToBytes(len(message)), message...)
}

func (this *CliSocket)IntToBytes(n int) []byte {
	return []byte(fmt.Sprintf("%04d",n))
}
