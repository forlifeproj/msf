package flsvr

import (
	"errors"
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
		Name       string `default:""`
		Address    string `default:""`
		ConsulAddr string `default:""`
	}
}

// v0.1.1
type FLSvr struct {
	s          *rpcx_svr.Server
	svrAddr    string
	consulAddr string
	basePath   string
	svrName    string
}

func NewFLServer(cfg string) *FLSvr {
	if len(cfg) == 0 {
		panic("cfg empty")
	}
	flSvr := &FLSvr{}
	svrAddr, consulAddr, basePath, svrName, err := loadSvrCfgInfo(cfg)
	if err != nil {
		panic("load svrcfg failed")
	}
	flSvr.svrAddr = svrAddr
	flSvr.consulAddr = consulAddr
	flSvr.basePath = basePath
	flSvr.svrName = svrName
	flSvr.s = rpcx_svr.NewServer()
	registerConuslPlugin(flSvr.s, svrAddr, consulAddr, basePath)
	return flSvr
}

func (f *FLSvr) RegisterHandler(svrHandle interface{}) error {
	f.s.RegisterName(f.svrName, svrHandle, "")
	fllog.Log().Debug("consulAddr:%s", consul.GetConsulAddr())

	return nil
}

func (f *FLSvr) RegisterFunc(fn interface{}) {
	f.s.RegisterFunction(f.svrName, fn, "")
}

func (f *FLSvr) StartServer() error {
	if err := f.s.Serve("tcp", f.svrAddr); err != nil {
		fllog.Log().Error("serve failed. err:", err)
		return err
	}
	fllog.Log().Error("start server succ")
	return nil
}

func loadSvrCfgInfo(cfg string) (string, string, string, string, error) {
	svrCfg := SvrCfg{}
	if err := config.ParseConfigWithPath(&svrCfg, cfg); err != nil {
		fllog.Log().Error("load svr logcfg failed.", err, cfg)
		return "", "", "", "", err
	}
	fllog.Log().Debug("svrCfg:%+v", svrCfg)
	if len(svrCfg.Server.Name) == 0 || len(svrCfg.Server.Address) == 0 {
		return "", "", "", "", errors.New("svrcfg invalid")
	}

	basePath, svrName := parseSvrName(svrCfg.Server.Name)
	if len(basePath) == 0 || len(svrName) == 0 {
		fllog.Log().Error("basePath or svrName empty", basePath, svrName)
		return "", "", "", "", errors.New("parse server name failed")
	}

	if len(svrCfg.Server.ConsulAddr) == 0 {
		localIP := getLocalIp()
		if len(localIP) == 0 {
			fllog.Log().Error("localIP empty")
			return "", "", "", "", errors.New("consulAddr empty")
		}
		svrCfg.Server.ConsulAddr = localIP + ":8500"
	}
	fllog.Log().Debug(svrCfg.Server.Address, svrCfg.Server.ConsulAddr, basePath, svrName)
	if len(svrCfg.Server.ConsulAddr) > 0 {
		consul.SetConsulAddr(svrCfg.Server.ConsulAddr)
	}
	fllog.Log().Debug("svrCfg=", svrCfg)
	return svrCfg.Server.Address, svrCfg.Server.ConsulAddr, basePath, svrName, nil
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

func parseSvrName(name string) (string, string) {
	vecSplit := strings.Split(name, ".")
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

// v0.1.0
// func Server(cfg string, svrHandle interface{}) error {
// 	svrAddr, consulAddr, basePath, svrName, err := loadSvrCfgInfo(cfg)
// 	if err != nil {
// 		return err
// 	}
// 	s := rpcx_svr.NewServer()
// 	registerConuslPlugin(s, svrAddr, consulAddr, basePath)
// 	s.RegisterName(svrName, svrHandle, "")
// 	fllog.Log().Debug("consulAddr:%s", consul.GetConsulAddr())
// 	if err := s.Serve("tcp", svrAddr); err != nil {
// 		fllog.Log().Error("serve failed. err:", err)
// 		return err
// 	}

// 	fllog.Log().Error("start server success!")
// 	return nil
// }
