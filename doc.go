package smartSock

/*
import (
	"log"
	"net"
	"time"
)

type event struct {
}

func (this event) OnHand(fd uint32, conn net.Conn) bool {
	log.Println(fd, "链接成功类")
	return true
}

func (this event) OnClose(fd uint32) {
	log.Println(fd, "链接断开类")
}

func (this event) OnMessage(fd uint32, msg []byte) bool {
	log.Println("这个是接受消息事件",string(msg))
	return true
}

type Test struct{
}


func (this Test) Default(fd uint32,data []byte) []byte {
	log.Println("default")
	return nil
}

func (this Test) BeforeRequest(fd uint32,data []byte) []byte {
	log.Println("before")
	return nil
}

func (this Test) AfterRequest(fd uint32,data []byte) []byte{
	log.Println("after")
	return nil
}

func main() {
	var buf = make(chan []byte, 10000)

	var ser = NewMsf(&CommSocket{}, buf)
	ser.EventPool.RegisterEvent(&event{})
	ser.EventPool.RegisterStructFun("test", &Test{})

	var cli = NewMcf(&CliSocket{}, buf)

	go ser.Listening(":9123")
	time.Sleep(time.Second * 2)
	go cli.Dial("10.193.0.45:8234")
	select {}

}
 */
