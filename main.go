package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type ChangeText struct {
	Text    string    `json:"text"`
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
	arg1 := "-i"
	arg2 := url
	arg3 := "-map"
	arg4 := "0:v:0"
	arg5 := "-map"
	arg6 := "0:a:0"
	arg7 := "-map"
	arg8 := "0:v:0"
	arg9 := "-map"
	arg10 := "0:a:0"
	arg11 := "-map"
	arg12 := "0:v:0"
	arg13 := "-map"
	arg14 := "0:a:0"
	arg15 := "-c:v"
	arg16 := "libx264"
	arg17 := "-preset:v"
	arg18 := "superfast"
	arg19 := "-c:a"
	arg20 := "aac"
	arg21 := "-x264opts"
	arg22 := "keyint=123:min-keyint=20:no-scenecut"
	arg23 := "-sc_threshold"
	arg24 := "0"
	arg25 := "-filter:v:0"
	arg26 := "scale=w=480:h=270:force_original_aspect_ratio=decrease, drawtext=textfile=activetext.txt:fontfile=./fonts/OpenSans.ttf:fontsize=8:fontcolor=white:reload=10:y=h-line_h-22:x=w-(mod(2*n\\,w+tw)-tw/20):box=1:boxcolor=black@0.5:boxborderw=1"
	arg27 := "-maxrate:v:0"
	arg28 := "600k"
	arg29 := "-b:a:0"
	arg30 := "500k"
	arg31 := "-filter:v:1"
	arg32 := "scale=w=640:h=360:force_original_aspect_ratio=decrease, drawtext=textfile=activetext.txt:fontfile=./fonts/OpenSans.ttf:fontsize=12:fontcolor=white:reload=10:y=h-line_h-32:x=w-(mod(2*n\\,w+tw)-tw/20):box=1:boxcolor=black@0.5:boxborderw=2"
	arg33 := "-maxrate:v:1"
	arg34 := "1500k"
	arg35 := "-b:a:1"
	arg36 := "1000k"
	arg37 := "-filter:v:2"
	arg38 := "scale=w=1280:h=720:force_original_aspect_ratio=decrease, drawtext=textfile=activetext.txt:fontfile=./fonts/OpenSans.ttf:fontsize=24:fontcolor=white:reload=10:y=h-line_h-52:x=w-(mod(4*n\\,w+tw)-tw/40):box=1:boxcolor=black@0.5:boxborderw=2"
	arg39 := "-maxrate:v:2"
	arg40 := "3000k"
	arg41 := "-b:a:2"
	arg42 := "2000k"
	arg43 := "-var_stream_map"
	arg44 := "v:0,a:0,name:480p v:1,a:1,name:720p v:2,a:2,name:1080p"
	arg45 := "-f"
	arg46 := "hls"
	arg47 := "-hls_list_size"
	arg48 := "6"
	arg49 := "-threads"
	arg50 := "0"
	arg51 := "-hls_time"
	arg52 := "10"
	arg53 := "-hls_flags"
	arg54 := "delete_segments"
	arg55 := "-master_pl_name"
	arg56 := "stream.m3u8"
	arg57 := "-y"
	arg58 := "./media/stream-%v.m3u8"

	cmd := exec.Command( app, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10, arg11, arg12, arg13, arg14, arg15, arg16, arg17, arg18, arg19, arg20, arg21, arg22, arg23, arg24, arg25, arg26, arg27, arg28, arg29, arg30, arg31, arg32, arg33, arg34, arg35, arg36, arg37, arg38, arg39, arg40, arg41, arg42, arg43, arg44, arg45, arg46, arg47, arg48, arg49, arg50, arg51, arg52, arg53, arg54, arg55, arg56, arg57, arg58 )

	go func() {
		if err := cmd.Start(); err != nil {
			_ = os.RemoveAll(fmt.Sprintf("%s", "./media"))
			fmt.Println(err)
		}
		if err := cmd.Wait(); err != nil {
			_ = os.RemoveAll(fmt.Sprintf("%s", "./media"))
			panic(err)
		}
	} ()
}

