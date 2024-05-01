package flsvr

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	cserver "github.com/rpcxio/rpcx-consul/serverplugin"
	rpcx_svr "github.com/smallnest/rpcx/server"

	"github.com/forlifeproj/msf/config"
	consul "github.com/forlifeproj/msf/consul"
	fllog "github.com/forlifeproj/msf/log"
)

type SvrCfg struct {
	Server struct {
		Name        string `default:""`
		Address     string `default:""`
		ConsulAddr  string `default:""`
		Environment string `default:"test"`
	}
}

type ServerInfo struct {
	SvrAddr     string
	ConsulAddr  string
	BasePath    string
	App         string
	SvrName     string
	Environment string
}

// v0.1.2
type FLSvr struct {
	s       *rpcx_svr.Server
	svrInfo *ServerInfo
}

func NewFLServer(cfg string) *FLSvr {
	if len(cfg) == 0 {
		panic("cfg empty")
	}
	flSvr := &FLSvr{}
	var err error
	flSvr.svrInfo, err = loadSvrCfgInfo(cfg)
	if err != nil || flSvr.svrInfo == nil {
		panic("load svrcfg failed")
	}
	flSvr.s = rpcx_svr.NewServer()
	registerConuslPlugin(flSvr.s, flSvr.svrInfo.SvrAddr, flSvr.svrInfo.ConsulAddr, flSvr.svrInfo.BasePath)
	return flSvr
}

func (f *FLSvr) getSvrName() string {
	return f.svrInfo.SvrName
}

func (f *FLSvr) getSvrAddr() string {
	return f.svrInfo.SvrAddr
}

func (f *FLSvr) RegisterHandler(svrHandle interface{}) error {
	f.s.RegisterName(f.getSvrName(), svrHandle, "")
	fllog.Log().Debug("consulAddr:%s", consul.GetConsulAddr())

	return nil
}

func (f *FLSvr) RegisterFunc(fn interface{}) {
	f.s.RegisterFunction(f.getSvrName(), fn, "")
}

func (f *FLSvr) StartServer() error {
	if err := f.s.Serve("tcp", f.getSvrAddr()); err != nil {
		fllog.Log().Error("serve failed. err:", err)
		return err
	}
	fllog.Log().Error("start server succ")
	return nil
}

func loadSvrCfgInfo(cfg string) (*ServerInfo, error) {
	svrCfg := SvrCfg{}
	if err := config.ParseConfigWithPath(&svrCfg, cfg); err != nil {
		fllog.Log().Error("load svr logcfg failed.", err, cfg)
		return nil, err
	}
	fllog.Log().Debug(fmt.Sprintf("svrCfg:%+v", svrCfg))
	svrCfg.Server.Environment = strings.ToLower(svrCfg.Server.Environment)
	if len(svrCfg.Server.Name) == 0 || len(svrCfg.Server.Address) == 0 ||
		(svrCfg.Server.Environment != "test" && svrCfg.Server.Environment != "proc") {
		fllog.Log().Error(fmt.Sprintf("SvrAddr:[%s] or SvrName:[%s] or Environment:[%s] invalid",
			svrCfg.Server.Address, svrCfg.Server.Name, svrCfg.Server.Environment))
		return nil, errors.New("svrcfg invalid")
	}

	strApp, svrName := parseAppSvrName(svrCfg.Server.Name)
	if len(strApp) == 0 || len(svrName) == 0 {
		fllog.Log().Error(fmt.Sprintf("app:[%s] or svrName:[%s] empty", strApp, svrName))
		return nil, errors.New("parse server name failed")
	}

	if len(svrCfg.Server.ConsulAddr) == 0 {
		localIP := getLocalIp()
		if len(localIP) == 0 {
			fllog.Log().Error("localIP empty")
			return nil, errors.New("consulAddr empty")
		}
		svrCfg.Server.ConsulAddr = localIP + ":8500"
	}

	fllog.Log().Debug(fmt.Sprintf("SvrAddr:%s ConsulAddr:%s App:%s SvrName:%s Environment:%s",
		svrCfg.Server.Address, svrCfg.Server.ConsulAddr, strApp, svrName, svrCfg.Server.Environment))
	if len(svrCfg.Server.ConsulAddr) > 0 {
		consul.SetConsulAddr(svrCfg.Server.ConsulAddr)
	}
	if len(svrCfg.Server.Environment) > 0 {
		consul.SetConsulEnvironment(svrCfg.Server.Environment)
	}
	svrInfo := &ServerInfo{
		SvrAddr:     svrCfg.Server.Address,
		ConsulAddr:  svrCfg.Server.ConsulAddr,
		BasePath:    fmt.Sprintf("/%s_%s", svrCfg.Server.Environment, strApp),
		App:         strApp,
		SvrName:     svrName,
		Environment: svrCfg.Server.Environment,
	}
	fllog.Log().Debug(fmt.Sprintf("svrInfo:%+v", svrInfo))
	return svrInfo, nil
}

func registerConuslPlugin(s *rpcx_svr.Server, svrAddr, conuslAddr, basePath string) {
	r := &cserver.ConsulRegisterPlugin{
		ServiceAddress: "tcp@" + svrAddr,
		ConsulServers:  []string{conuslAddr},
		BasePath:       basePath,
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		fllog.Log().Error("register consul failed. err=", err)
	}

	s.Plugins.Add(r)
	fllog.Log().Debug("register consul succ!")
}

func parseAppSvrName(src string) (string, string) {
	vecSplit := strings.Split(src, ".")
	if len(vecSplit) != 2 {
		return "", ""
	}
	return vecSplit[0], vecSplit[1]
}

func getLocalIp() string {
	iface, err := net.InterfaceByName("eth0")
	if err != nil {
		fllog.Log().Error("Error:", err)
		return ""
	}

	addrs, err := iface.Addrs()
	if err != nil {
		fllog.Log().Error("Error:", err)
		return ""
	}

	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			fllog.Log().Error("Error:", err)
			continue
		}
		if ip.To4() != nil {
			fllog.Log().Debug("IPv4:", ip)
			return ip.String()
		}
		//  else {
		// 	fmt.Println("IPv6:", ip)
		// 	return
		// }
	}
	return ""
}
