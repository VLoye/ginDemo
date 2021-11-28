package blog

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"os/exec"
)

/**
  blog github webhook
  todo:
	1. 增加配置读取；
	2. 添加handler
	3. 执行git pull
	4. hexo g
 */

type Blog struct {
	Path string `yaml:"path"`
}


func Webhook(ctx *gin.Context, blogConfig Blog) {
	body,_:= ioutil.ReadAll(ctx.Request.Body)
	log.Println(string(body))


	var tempMap map[string]interface{}
	json.Unmarshal(body,&tempMap)
	if tempMap["pusher"]!= nil {
		log.Printf("pusher is %s\n", tempMap["pusher"])

		cmd1 := exec.Command("cd", blogConfig.Path)
		err1:= cmd1.Run()
		if err1!=nil{
			panic(err1)
		}
		cmd2 := exec.Command("git", "pull")
		err2:= cmd2.Run()
		if err2!=nil{
			panic(err2)
		}

		cmd3 := exec.Command("hexo","g")
		err3:= cmd3.Run()
		if err3!=nil{
			panic(err3)
		}
		log.Println("execute webhook success.")
	}

}

