package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	// 首页
	c.HTML(http.StatusOK, "index/index.html", gin.H{
		"tt": "下载页面",
		// "zo": "zo.apk",
		// "yg": "yg.apk",
	})
}

func DownApp(c *gin.Context) {
	// 下载资料 ，下载的函数
	// 从文件系统读取打开这个要下载的文件，
	// 然后，循环读取1M字节的内容，直到读完所有的内容，分块读取能提高效率，不会堵塞， 输出的时候用到了 Header()
	filen := c.Query("file")
	fmt.Println("=======>")
	fmt.Println(filen)
	fmt.Println("=======>")
	filep := path.Join("static", filen)
	// filep := "static/zoom.apk"
	// filen := "zoom.apk"
	file, err := os.Open(filep)
	if err != nil {
		c.AbortWithError(404, err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		c.AbortWithError(404, err)
		return
	}
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+filen)

	c.Writer.Header().Set("Content-Type", "application/octet-stream")

	c.Writer.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))
	c.Writer.Flush()
	var offset int64 = 0
	var bufsize int64 = 1024 * 1024
	buf := make([]byte, bufsize)

	for {
		n, err := file.ReadAt(buf, offset)
		if err != nil && err != io.EOF {
			log.Println("read file error!", err)
			break
		}
		if n == 0 {
			break
		}
		_, err = c.Writer.Write(buf[:n])
		if err != nil {
			log.Println("write file error ", err)
			break
		}
		offset += int64(n)
	}
	c.Writer.Flush()
}
func List(c *gin.Context) {
	// 显示列表页面
	c.HTML(http.StatusOK, "index/list.html", gin.H{})
}
func Play(c *gin.Context) {
	// fn := c.Query("fn")
	// filename := path.Join("./static", fn)

	file, err := os.Open("./static/a11.mp4")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	c.File(file.Name())
	return
	// file, err := os.Open("static/")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// defer file.Close()

	//
}
func main() {
	fmt.Println("start service!")
	r := gin.Default()
	r.StaticFS("/static", http.Dir("/static"))
	r.LoadHTMLGlob("templates/**/*")
	r.GET("/", Index)
	r.GET("/list", List)
	r.GET("/play", Play)
	r.GET("/download/", DownApp)
	r.Run(":8900")
}
