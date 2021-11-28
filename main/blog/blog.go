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

		//cmd1 := exec.Command("cd", blogConfig.Path)
		//b1, err1 :=cmd1.Output()
		////err1:= cmd1.Run()
		//if err1!=nil{
		//	panic(err1)
		//}
		//log.Printf("%s\n", string(b1))
		cmd2 := exec.Command("git", "pull")
		cmd2.Path=blogConfig.Path
		b2,err2:= cmd2.Output()
		if err2!=nil{
			panic(err2)
		}
		log.Printf("%s\n", string(b2))

		cmd3 := exec.Command("hexo","g")
		cmd3.Path=blogConfig.Path
		b3,err3:= cmd3.Output()
		if err3!=nil{
			panic(err3)
		}
		log.Printf("%s\n", string(b3))
		log.Println("execute webhook success.")
	}

}

