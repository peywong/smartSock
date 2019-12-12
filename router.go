package smartSock

import (
	"net"
	"reflect"
)

const (
	DEFAULTACTION = "Default"
	BEFOREACTION = "BeforeRequest"
	AFTERACTION = "AfterRequest"
)

type module interface {
	Default(fd uint32, data []byte) []byte
	BeforeRequest(fd uint32, data []byte) []byte
	AfterRequest(fd uint32, data []byte) []byte
}

type eventer interface {
	OnHand(fd uint32, conn net.Conn) bool
	OnClose(fd uint32)
	OnMessage(fd uint32, data []byte) bool
}

type RoutersMap struct {
	//map[methodname]handlefunc
	methods map[string]func(uint32, []byte) []byte
	//module[methodname]
	mods module
	events eventer
}

func NewRoutersMap() *RoutersMap {
	return &RoutersMap{
		methods: make(map[string]func(uint32, []byte) []byte),
	}
}

func (this *RoutersMap) RegisterEvent(events eventer) {
	this.events = events
}

func (this *RoutersMap) RegisterStructFun(moduleName string, mod module) bool {
	if this.mods != nil {
		return false
	}
	this.mods = mod

	temType := reflect.TypeOf(mod)
	temValue := reflect.ValueOf(mod)
	for i := 0; i < temType.NumMethod(); i++ {
		tem := temValue.Method(i).Interface()
		if temFunc, ok := tem.(func(uint32, []byte) []byte); ok {
			this.methods[temType.Method(i).Name] = temFunc
		}
	}
	return true
}

func (this *RoutersMap) OnClose(fd uint32) {
	if this.events != nil {
		this.events.OnClose(fd)
	}
}

func (this *RoutersMap) OnHand(fd uint32, conn net.Conn) bool {
	if this.events != nil {
		return	this.events.OnHand(fd, conn)
	}
	return true
}

func (this *RoutersMap) OnMessage(fd uint32, data []byte) bool {
	if this.events != nil {
		return this.events.OnMessage(fd, data)
	}
	return true
}