package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/glog"
)

const (
	LOCAL_IP = "LOCAL_IP"
	OMPI     = "ompi"
	MPICH    = "mpich"
)

func GetLocalIP() string {
	return os.Getenv(LOCAL_IP)
}

func WriteMPIHostFile(MPIImplementation, hostFilePath string, ips []string, replicasPerNode int) error {
	//mpi需要在本机跑进程，需要LOCAL_IP环境变量获取本机ip
	//kai有不把进程下发到本机的需求，不需要本机ip
	localIP := GetLocalIP()
	if len(localIP) != 0 {
		ips = append([]string{localIP}, ips...)
	}

	if len(ips) == 0 {
		return fmt.Errorf("no ips get, cannot write %s", hostFilePath)
	}

	ipsWithReplicas := []string{}

	switch MPIImplementation {
	case MPICH:
		for _, ip := range ips {
			ipsWithReplicas = append(ipsWithReplicas, fmt.Sprintf("%v:%v", ip, replicasPerNode))
		}
	case OMPI:
		for _, ip := range ips {
			ipsWithReplicas = append(ipsWithReplicas, fmt.Sprintf("%v slots=%v", ip, replicasPerNode))
		}
	default:
		//保留原来的hostfile设置，这种设置对mpich和open-mpi通用
		//NOTE: 这种指定ip的方式在open-mpi 4.0.0中，如果replicasPerNode是1，mpirun在每个节点启动的进程数与cpu核数相等
		//TODO: open-mpi和mpich使用专有的hostfile稳定后可以去掉这个设置
		for _, ip := range ips {
			for i := 0; i < replicasPerNode; i++ {
				ipsWithReplicas = append(ipsWithReplicas, ip)
			}
		}
	}

	paths := strings.Split(hostFilePath, ",")
	var content string
	for _, path := range paths {
		content = strings.Join(ipsWithReplicas, "\n") + "\n"
		if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
			glog.Errorf("failed to write hostfile %s: %v", path, err)
			return err
		}
	}
	glog.Infof("succeeded to write hostfile %s: content=\n%s", hostFilePath, content)
	return nil
}
