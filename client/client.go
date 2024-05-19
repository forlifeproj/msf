package flcli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/forlifeproj/msf/consul"
	cclient "github.com/rpcxio/rpcx-consul/client"
	rclient "github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/share"
)

// CallDesc RPC参数
type CallDesc struct {
	LocalServiceName string        // <非必填>本次请求主调服务名
	ServiceName      string        // <必填>本次请求被调服务名, 对应toml配置文件中的一段
	Timeout          time.Duration // <非必填>RPC超时时间
	CodecType        protocol.SerializeType
}

type ServiceInfo struct {
	SvrBasePath   string
	SvrName       string
	InterfaceName string
}

type FlClient struct {
	RpcCli      rclient.XClient
	SvrInfo     ServiceInfo
	ReqMetaData map[string]string
	ResMetaData map[string]string
}

func NewClient(callDesc CallDesc) *FlClient {
	// get consul addr to-do
	consulAddr := consul.GetConsulAddr()

	// parse svr_addr
	flC := &FlClient{}
	flC.ReqMetaData = make(map[string]string)
	flC.ResMetaData = make(map[string]string)

	flC.ParseSvrInfo(callDesc.ServiceName)
	svrDiscovery, _ := cclient.NewConsulDiscovery(
		flC.SvrInfo.SvrBasePath,
		flC.SvrInfo.SvrName,
		[]string{consulAddr},
		nil)
	option := rclient.DefaultOption
	if callDesc.CodecType > protocol.SerializeNone && callDesc.CodecType <= protocol.Thrift {
		option.SerializeType = callDesc.CodecType
	}

	flC.RpcCli = rclient.NewXClient(
		flC.SvrInfo.SvrName,
		rclient.Failtry,
		rclient.RandomSelect,
		svrDiscovery,
		option)
	return flC
}

func (f *FlClient) Close() {
	f.RpcCli.Close()
}

func (f *FlClient) AddReqMetaData(k, v string) {
	if f.ReqMetaData != nil {
		f.ReqMetaData[k] = v
	}
}

func (f *FlClient) GetResMetaData(key string) string {
	if f.ResMetaData != nil {
		return f.ResMetaData[key]
	}
	return ""
}

func (f *FlClient) DoRequest(ctx context.Context, req interface{}, rsp interface{}) (err error) {

	ctx = context.WithValue(ctx, share.ReqMetaDataKey, f.ReqMetaData)
	ctx = context.WithValue(ctx, share.ResMetaDataKey, make(map[string]string))

	err = f.RpcCli.Call(ctx, f.SvrInfo.InterfaceName, req, rsp)
	resMeta := ctx.Value(share.ResMetaDataKey).(map[string]string)
	for k, v := range resMeta {
		if f.ResMetaData != nil {
			f.ResMetaData[k] = v
		}
	}
	return err
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
