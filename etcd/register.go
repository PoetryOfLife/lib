package etcd

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"log"
	"time"
)

type ServiceRegister struct {
	cli           *clientv3.Client
	leaseID       clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	val           string
}

func NewServiceRegister(endpoint []string, key, val string, lease int64) (*ServiceRegister, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoint,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	ser := &ServiceRegister{
		cli: client,
		key: key,
		val: val,
	}

	err = ser.putKeyWithLease(lease)
	if err != nil {
		return nil, err
	}
	return ser, nil
}

func (s *ServiceRegister) putKeyWithLease(lease int64) error {
	resp, err := s.cli.Lease.Create(context.TODO(), lease)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	_, err = s.cli.Put(context.Background(), s.key, s.val, clientv3.WithLease(clientv3.LeaseID(resp.ID)))
	if err != nil {
		return err
	}
	leaseRespChan, err := s.cli.KeepAlive(context.Background(), clientv3.LeaseID(resp.ID))
	if err != nil {
		return err
	}
	s.leaseID = clientv3.LeaseID(resp.ID)
	log.Println(s.leaseID)
	s.keepAliveChan = leaseRespChan
	log.Printf("Put key:%s  val:%s  success!", s.key, s.val)
	return nil
}

func (s *ServiceRegister) ListenLeaseRespChan() {
	for leaseKeepResp := range s.keepAliveChan {
		log.Println("续约成功", leaseKeepResp)
	}
	log.Println("关闭续租")
}

func (s *ServiceRegister) Close() error {
	_, err := s.cli.Revoke(context.Background(), s.leaseID)
	if err != nil {
		return err
	}
	log.Println("撤销租约")
	return s.cli.Close()
}
