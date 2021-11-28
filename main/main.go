package main

import (
	"fmt"
	"ginDemo/main/blog"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var server *gin.Engine
var session map[string]string = map[string]string{}
var blogConfig blog.Blog

func main() {
	content, err:=ioutil.ReadFile("./blog.yaml")
	if err != nil{
		fmt.Println(err)
	}
	err2 := yaml.Unmarshal(content, &blogConfig)
	if err2!=nil{
		fmt.Println(err2)
	}
	fmt.Println(blogConfig.Path)
	server = gin.Default()
	server.LoadHTMLGlob("resource/static/*")
	server.MaxMultipartMemory = 1024 << 20
	server.GET("/index", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})
	server.GET("/login", func(context *gin.Context) {
		context.HTML(http.StatusOK, "login.html", nil)
	})
	server.POST("/webhook", func(context *gin.Context) {
		blog.Webhook(context, blogConfig)
		context.Done()
	})
	server.POST("/doLogin", doLogin)

	server.Use(Validate())
	server.GET("/listfile", listFiles)
	server.GET("/download/:fileName", download)

	server.POST("/uploadfile", upload)
	//server.GET("/setCookie", func(context *gin.Context) {
	//	context.SetCookie("name", "gxf", 10, "/", "localhost", false, true)
	//})
	server.Run(":9999")
}

//鉴权中间件
func Validate() gin.HandlerFunc {
	return func(context *gin.Context) {
		log.Printf("middleware invoke, url: %v", context.Request.URL.Path)
		//if strings.Contains(context.Request.URL.Path, "/login") || strings.Contains(context.Request.URL.Path, "/doLogin") {
		//	context.Next()
		//}
		if token, err := context.Cookie("my_token"); err == nil {
			if userName, ok := session[token]; ok {
				log.Printf("user[%v] request url: [%v]", userName, context.Request.URL.Path)
				context.Next()
				return
			}
		}
		log.Println("redirect to /login")
		context.Abort()
		context.Redirect(http.StatusTemporaryRedirect, "/login")

		//redirectInternal(context, "/login", "GET")

	}
}

func doLogin(c *gin.Context) {
	userName := c.PostForm("userName")
	password := c.PostForm("password")
	if userName == "admin" && password == "admin" {
		token := uuid.NewV4().String()
		log.Printf("toke is %v", token)
		log.Println(uuid.NewV4().String())
		session[token] = userName
		c.SetCookie("my_token", token, 60*60, "/", ".", false, true)
		//c.Request.URL.Path = "/listfile"
		c.Redirect(http.StatusFound,"/listfile")
	} else {
		c.JSON(http.StatusBadRequest, "login failed, userName or password validate failed.")
	}
}

func redirectInternal(c *gin.Context, url string, method string) {
	c.Request.URL.Path = url
	c.Request.Method = method
	server.HandleContext(c)
}

func upload(context *gin.Context) {
	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	} else {
		context.SaveUploadedFile(file, getUpLoadPath(file.Filename))
		context.JSON(http.StatusOK, gin.H{
			"message": "upload file success.",
		})
	}
}
func download(context *gin.Context) {
	fileName := context.Param("fileName")
	filePath := getDownLoadDir() + string(filepath.Separator) + fileName
	context.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	context.Writer.Header().Add("Content-Type", "application/octet-stream")
	context.File(filePath)
}

func listFiles(c *gin.Context) {
	downLoadDir := getDownLoadDir()
	log.Printf("download dir is: %v", downLoadDir)
	if files, err := ioutil.ReadDir(downLoadDir); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
	} else {
		log.Printf("files is : %v", files)
		c.HTML(200, "index.html", gin.H{
			"fileList": files,
			"descript": "this is a template",
			"title":    "this is a title",
		})
	}

}

func getDownLoadDir() string {

	downLoadDir := getResourcePath() + string(filepath.Separator) + "download"
	return downLoadDir
}

func getUpLoadPath(fileName string) string {
	upLoadDir := getResourcePath() + string(filepath.Separator) + "upload" + string(filepath.Separator) + fileName
	log.Printf("upload file: %v", upLoadDir)
	return upLoadDir
}

func getProjectPath() string {
	dir, _ := os.Getwd()
	projectPath := filepath.Dir(dir)
	return projectPath
}

func getResourcePath() string {
	return getProjectPath() + string(filepath.Separator) + "resource"

}
