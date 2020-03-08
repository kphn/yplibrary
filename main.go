package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/kphn/yplibrary/data"
)

type tomlConfig struct {
	Port            string `toml:"port"`
	IntervalSeconds int    `toml:"intervalSeconds"`
	Wechat          Wechat `toml:"wechat"`
}

type Wechat struct {
	OpenID    string `toml:"openId"`
	SessionID string `toml:"sessionId"`
}

func main() {

	cfgFile := flag.String("cfg", "config/test.toml", "config file")
	flag.Parse()

	var config tomlConfig
	if _, err := toml.DecodeFile(*cfgFile, &config); err != nil {
		panic(err)
	}

	ticker := time.Tick(time.Duration(config.IntervalSeconds) * time.Second)

	go func() {
		for {
			select {
			case <-ticker:
				data.GetYpLibVisitNum(config.Wechat.OpenID, config.Wechat.SessionID)
			default:
			}
		}
	}()

	r := gin.Default()

	r.GET("/visitors", func(c *gin.Context) {
		file := fmt.Sprintf("%s.html", time.Now().Format("20060102"))
		c.File(file)
		return
	})

	r.GET("/visitors/:day", func(c *gin.Context) {
		day := c.Param("day")

		file := fmt.Sprintf("%s.html", day)
		c.File(file)
		return
	})

	r.Run(config.Port)
}
