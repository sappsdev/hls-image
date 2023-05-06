package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type ChangeText struct {
	Text string `json:"text"`
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, token")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	_ = os.RemoveAll(fmt.Sprintf("%s", "./media"))
	_ = os.Mkdir(fmt.Sprintf("%s", "./media"), os.ModePerm)

	file, err := os.OpenFile("activetext.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	start()

	// gin.SetMode("debug")
	gin.SetMode("release")

	r := gin.New()
	r.Use(CORSMiddleware())
	r.Static("/media", "./media")
	r.POST("/text", HandlerChangeText)
	r.Run(":4000")

}

func HandlerChangeText(c *gin.Context) {
	changeText := new(ChangeText)

	if err := c.BindJSON(&changeText); err != nil {
		c.JSON(400, gin.H{
			"message": "text data error",
		})
		return
	}

	file, err := os.OpenFile("activetext.txt", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(changeText.Text)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	c.JSON(200, gin.H{
		"message": "text data change",
	})
}

func start() {
	fmt.Println("Starting ffmpeg")
	url := os.Getenv("URL")
	fmt.Println(url)

	app := "ffmpeg"
	arg1 := "-re"
	arg2 := "-i"
	arg3 := url
	arg4 := "-reconnect_at_eof"
	arg5 := "-reconnect_streamed"
	arg6 := "-reconnect_on_network_error"
	arg7 := "-reconnect_on_http_error"
	arg8 := "-err_detect"
	arg9 := "ignore_err"
	arg10 := "-map"
	arg11 := "0:v:0"
	arg12 := "-map"
	arg13 := "0:a:0"
	arg14 := "-map"
	arg15 := "0:v:0"
	arg16 := "-map"
	arg17 := "0:a:0"
	arg18 := "-map"
	arg19 := "0:v:0"
	arg20 := "-map"
	arg21 := "0:a:0"
	arg22 := "-c:v"
	arg23 := "libx264"
	arg24 := "-preset:v"
	arg25 := "ultrafast"
	arg26 := "-c:a"
	arg27 := "aac"
	arg46 := "-b:a:0"
	arg47 := "64k"
	arg48 := "-bufsize:v:0"
	arg49 := "1000k"
	arg50 := "-maxrate:v:0"
	arg51 := "3000k"
	arg38 := "-b:a:1"
	arg39 := "64k"
	arg40 := "-bufsize:v:1"
	arg41 := "2000k"
	arg42 := "-maxrate:v:1"
	arg43 := "5000k"
	arg30 := "-b:a:2"
	arg31 := "64k"
	arg32 := "-bufsize:v:2"
	arg33 := "4000k"
	arg34 := "-maxrate:v:2"
	arg35 := "8000k"
	arg52 := "-var_stream_map"
	arg53 := "v:0,a:0,name:480p v:1,a:1,name:720p v:2,a:2,name:1080p"
	arg54 := "-f"
	arg55 := "hls"
	arg56 := "-hls_list_size"
	arg57 := "10"
	arg58 := "-hls_time"
	arg59 := "5"
	arg60 := "-hls_flags"
	arg61 := "delete_segments"
	arg62 := "-master_pl_name"
	arg63 := "stream.m3u8"
	arg64 := "-y"
	arg65 := "./media/stream-%v.m3u8"

	go func() {
		for {
			cmd := exec.Command(app, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10, arg11, arg12, arg13, arg14, arg15, arg16, arg17, arg18, arg19, arg20, arg21, arg22, arg23, arg24, arg25, arg26, arg27, arg30, arg31, arg32, arg33, arg34, arg35, arg38, arg39, arg40, arg41, arg42, arg43, arg46, arg47, arg48, arg49, arg50, arg51, arg52, arg53, arg54, arg55, arg56, arg57, arg58, arg59, arg60, arg61, arg62, arg63, arg64, arg65)
			if err := cmd.Start(); err != nil {
				_ = os.RemoveAll(fmt.Sprintf("%s", "./media"))
				fmt.Println(err)
			}
			if err := cmd.Wait(); err != nil {
				fmt.Println(err)
				_ = os.RemoveAll(fmt.Sprintf("%s", "./media"))
				// panic(err)
			}
		}
	}()
}
