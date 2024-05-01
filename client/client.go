package flcli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/forlifeproj/msf/consul"
	cclient "github.com/rpcxio/rpcx-consul/client"
	rclient "github.com/smallnest/rpcx/client"
)

// CallDesc RPC参数
type CallDesc struct {
	LocalServiceName string        // <非必填>本次请求主调服务名
	ServiceName      string        // <必填>本次请求被调服务名, 对应toml配置文件中的一段
	Timeout          time.Duration // <非必填>RPC超时时间
}

type ServiceInfo struct {
	SvrBasePath   string
	SvrName       string
	InterfaceName string
}

type FlClient struct {
	RpcCli rclient.XClient

	SvrInfo ServiceInfo
}

func NewClient(callDesc CallDesc) *FlClient {
	// get consul addr to-do
	consulAddr := consul.GetConsulAddr()

	// parse svr_addr
	flC := &FlClient{}
	flC.ParseSvrInfo(callDesc.ServiceName)
	svrDiscovery, _ := cclient.NewConsulDiscovery(
		flC.SvrInfo.SvrBasePath,
		flC.SvrInfo.SvrName,
		[]string{consulAddr},
		nil)
	flC.RpcCli = rclient.NewXClient(
		flC.SvrInfo.SvrName,
		rclient.Failtry,
		rclient.RandomSelect,
		svrDiscovery,
		rclient.DefaultOption)
	return flC
}

func (f *FlClient) Close() {
	f.RpcCli.Close()
}

func (f *FlClient) DoRequest(ctx context.Context, req interface{}, rsp interface{}) error {
	return f.RpcCli.Call(ctx, f.SvrInfo.InterfaceName, req, rsp)
}

func (f *FlClient) ParseSvrInfo(serviceName string) {
	vecSplit := strings.Split(serviceName, ".")
	if len(vecSplit) >= 0 {
		f.SvrInfo.SvrBasePath = fmt.Sprintf("/%s_%s", consul.GetConsulEnvironment(), vecSplit[0])
	}
	if len(vecSplit) >= 1 {
		f.SvrInfo.SvrName = vecSplit[1]
	}
	if len(vecSplit) >= 2 {
		f.SvrInfo.InterfaceName = vecSplit[2]
	}
	fmt.Printf("svrInfo:%+v \n", f.SvrInfo)
}
