package consul

import (
	"fmt"
	"sync"

	"github.com/forlifeproj/msf/config"
)

type ConsulConfig struct {
	ConsulAddr string
}

type ConsulUtils struct {
	ConsulCfg ConsulConfig
}

var onceConusl sync.Once
var ConsulInstance *ConsulUtils

type SvrCfg struct {
	Server struct {
		ConsulAddr string `default:""`
	}
}

func Init(cfg string) error {
	svrCfg := SvrCfg{}

	if err := config.ParseConfigWithPath(&svrCfg, cfg); err != nil {
		fmt.Printf("load svrCfg failed. err:%+v cfg:%s", err, cfg)
		return err
	}

	SetConsulAddr(svrCfg.Server.ConsulAddr)
	return nil
}

func NewConsulUtils() *ConsulUtils {
	onceConusl.Do(func() {
		ConsulInstance = &ConsulUtils{}
	})
	return ConsulInstance
}

func SetConsulAddr(addr string) {
	utils := NewConsulUtils()
	if utils != nil {
		utils.ConsulCfg.ConsulAddr = addr
		fmt.Printf("set  consulAddr=%s suuc \n", addr)
	}
}

func GetConsulAddr() string {
	utils := NewConsulUtils()
	if utils != nil {
		fmt.Printf("get  consulAddr=%s suuc \n", utils.ConsulCfg.ConsulAddr)
		return utils.ConsulCfg.ConsulAddr
	}
	fmt.Printf("utils:%+v", utils)
	return ""
}
