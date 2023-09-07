package etcd

import (
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	var endpoints = []string{"localhost:2379"}
	ser, err := NewServiceRegister(endpoints, "/web/node1", "localhost:8000", 5)
	if err != nil {
		t.Error(err.Error())
		return
	}
	go ser.ListenLeaseRespChan()
	select {
	case <-time.After(20 * time.Second):
		ser.Close()
	}
}
