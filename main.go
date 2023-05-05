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
	arg18 := "ultrafast"
	arg19 := "-bufsize"
	arg20 := "6M"
	arg21 := "-c:a"
	arg22 := "aac"
	arg23 := "-filter:v:0"
	arg24 := "scale=w=480:h=270:force_original_aspect_ratio=decrease"
	arg25 := "-maxrate:v:0"
	arg26 := "1000k"
	arg27 := "-b:a:0"
	arg28 := "500k"
	arg29 := "-filter:v:1"
	arg30 := "scale=w=640:h=360:force_original_aspect_ratio=decrease, drawtext=textfile=activetext.txt:fontfile=./fonts/OpenSans.ttf:fontsize=12:fontcolor=white:reload=10:y=h-line_h-32:x=w-(mod(2*n\\,w+tw)-tw/20):box=1:boxcolor=black@0.5:boxborderw=2"
	arg31 := "-maxrate:v:1"
	arg32 := "2000k"
	arg33 := "-b:a:1"
	arg34 := "1000k"
	arg35 := "-filter:v:2"
	arg36 := "scale=w=1280:h=720:force_original_aspect_ratio=decrease, drawtext=textfile=activetext.txt:fontfile=./fonts/OpenSans.ttf:fontsize=24:fontcolor=white:reload=10:y=h-line_h-52:x=w-(mod(4*n\\,w+tw)-tw/40):box=1:boxcolor=black@0.5:boxborderw=2"
	arg37 := "-maxrate:v:2"
	arg38 := "5000k"
	arg39 := "-b:a:2"
	arg40 := "2000k"
	arg41 := "-var_stream_map"
	arg42 := "v:0,a:0,name:480p v:1,a:1,name:720p v:2,a:2,name:1080p"
	arg43 := "-f"
	arg44 := "hls"
	arg45 := "-hls_list_size"
	arg46 := "0"
	arg47 := "-hls_time"
	arg48 := "10"
	arg49 := "-hls_flags"
	arg50 := "delete_segments"
	arg51 := "-master_pl_name"
	arg52 := "stream.m3u8"
	arg53 := "-y"
	arg54 := "./media/stream-%v.m3u8"

	cmd := exec.Command( app, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10, arg11, arg12, arg13, arg14, arg15, arg16, arg17, arg18, arg19, arg20, arg21, arg22, arg23, arg24, arg25, arg26, arg27, arg28, arg29, arg30, arg31, arg32, arg33, arg34, arg35, arg36, arg37, arg38, arg39, arg40, arg41, arg42, arg43, arg44, arg45, arg46, arg47, arg48, arg49, arg50, arg51, arg52, arg53, arg54 )

	go func() {

		if err := cmd.Start(); err != nil {
			_ = os.RemoveAll(fmt.Sprintf("%s", "./media"))
			fmt.Println(err)
		}
		if err := cmd.Wait(); err != nil {
			fmt.Println(err)
			_ = os.RemoveAll(fmt.Sprintf("%s", "./media"))
			// panic(err)
		}
	} ()
}

