package services

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/leiwenxuan/crontab/infra"

	"github.com/leiwenxuan/crontab/common"
	log "github.com/sirupsen/logrus"

	"github.com/leiwenxuan/crontab/infra/base"
	"go.etcd.io/etcd/clientv3"
)

type EtcdRegisterStarter struct {
	infra.BaseStarter
}

func (s *EtcdRegisterStarter) Start(ctx infra.StarterContext) {
	var (
		kv      clientv3.KV
		lease   clientv3.Lease
		localIp string
		err     error
	)
	// 获取etcd clientMongo
	client := base.EtcdClient()
	// 得到kv 和 lease 租约id
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	// 获取本机IP
	if localIp, err = GetLocalIP(); err != nil {
		return
	}
	GRegister = &RegisterEtcd{
		client:  client,
		kv:      kv,
		lease:   lease,
		localIP: localIp,
	}
	// 服务注册
	go GRegister.keepOnline()
}

// 注册节点到etcd: /cron/workers/IP地址
type RegisterEtcd struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	localIP string // 本机IP
}

var (
	GRegister *RegisterEtcd
)

// 服务注册， 自动续租
func InitRegister() (err error) {
	var (
		kv      clientv3.KV
		lease   clientv3.Lease
		localIp string
	)
	// 获取etcd client
	client := base.EtcdClient()
	// 得到kv 和 lease 租约id
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)
	// 获取本机IP
	localIp = GetOutboundIP()

	GRegister = &RegisterEtcd{
		client:  client,
		kv:      kv,
		lease:   lease,
		localIP: localIp,
	}
	// 服务注册
	go GRegister.keepOnline()
	return err
}

// 服务注册
func (register *RegisterEtcd) keepOnline() {
	var (
		regKey         string
		leaseGrantResp *clientv3.LeaseGrantResponse
		err            error
		keepAliveChan  <-chan *clientv3.LeaseKeepAliveResponse
		keepAliveReap  *clientv3.LeaseKeepAliveResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
	)
	for {
		// 注册路径
		regKey = common.JOB_WORKER_DIR + register.localIP
		cancelFunc = nil

		// 创建租约， 10秒
		if leaseGrantResp, err = register.lease.Grant(context.TODO(), 10); err != nil {
			goto RETRY
		}
		// 自动续租
		if keepAliveChan, err = register.lease.KeepAlive(context.TODO(), leaseGrantResp.ID); err != nil {
			log.Debug("err", err)
			goto RETRY
		}

		cancelCtx, cancelFunc = context.WithCancel(context.TODO())
		// 注册到， etcd， 自动续租
		if _, err = register.kv.Put(cancelCtx, regKey, "", clientv3.WithLease(leaseGrantResp.ID)); err != nil {
			goto RETRY
		}
		// 处理应答租约
		for {
			select {
			case keepAliveReap = <-keepAliveChan:
				if keepAliveReap == nil {
					// 续租失败
					goto RETRY
				}
			}
		}
	RETRY:
		time.Sleep(1 * time.Second)
		if cancelFunc != nil {
			// 当cancelfunc 创建成功后， 重新续租失败， 取消
			log.Error("register 续租失败", time.Now().Format("2006-01-02 15:04:05"))
			cancelFunc()
		}
	}
}

// 获取本机网卡IP
func GetLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIpNet bool
	)

	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}

	fmt.Println(addrs)
	// 取第一个非lo的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址: ipv4,ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String() // 192.168.1.1
				return
			}
		}
	}
	err = errors.New("没有找到网卡IP")
	return
}

// Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	return localAddr.IP.String()
}
