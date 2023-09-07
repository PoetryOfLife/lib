package etcd

import (
	"context"
	"github.com/coreos/etcd/storage/storagepb"
	"go.etcd.io/etcd/clientv3"
	"log"
	"sync"
	"time"
)

type ServiceDiscovery struct {
	cli        *clientv3.Client
	serverList map[string]string
	lock       sync.Mutex
}

func NewServiceDiscovery(endpoints []string) *ServiceDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	return &ServiceDiscovery{
		cli:        cli,
		serverList: make(map[string]string),
	}
}

func (s *ServiceDiscovery) WatchService(prefix string) error {
	resp, err := s.cli.Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range resp.Kvs {
		s.SetServiceList(string(kv.Key), string(kv.Value))
	}
	go s.watcher(prefix)
	return nil
}

func (s *ServiceDiscovery) watcher(prefix string) {
	rch := s.cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	log.Printf("watching prefix:%s now...", prefix)
	for wresp := range rch {
		for _, event := range wresp.Events {
			switch event.Type {
			case storagepb.PUT:
				s.SetServiceList(string(event.Kv.Key), string(event.Kv.Value))
			case storagepb.DELETE:
				s.DelServiceList(string(event.Kv.Key))
			}
		}
	}
}

func (s *ServiceDiscovery) SetServiceList(key, val string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.serverList[key] = val
	log.Println("del key:", key)
}

func (s *ServiceDiscovery) DelServiceList(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.serverList, key)
	log.Println("del key:", key)
}

func (s *ServiceDiscovery) GetServices() []string {
	s.lock.Lock()
	defer s.lock.Unlock()
	addrs := make([]string, 0)

	for _, v := range s.serverList {
		addrs = append(addrs, v)
	}
	return addrs
}

func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}
