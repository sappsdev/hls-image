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

	gin.SetMode("debug")
	// gin.SetMode("release")

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
	args := []string{"-i", url}
	args = append(args, "-map", "0:v:0", "-map", "0:a:0")
	args = append(args, "-map", "0:v:0", "-map", "0:a:0")
	args = append(args, "-map", "0:v:0", "-map", "0:a:0")
	args = append(args, "-c:v", "libx264", "-preset:v", "superfast")
	args = append(args, "-c:a", "aac")
	args = append(args, "-x264opts", "keyint=123:min-keyint=20:no-scenecut")
	args = append(args, "-sc_threshold", "0")
	args = append(args, "-filter:v:0")
	args = append(args, "scale=w=480:h=270:force_original_aspect_ratio=decrease, drawtext=textfile=activetext.txt:fontsize=8:fontcolor=white:reload=10:y=h-line_h-22:x=w-(mod(2*n\\,w+tw)-tw/20):box=1:boxcolor=black@0.5:boxborderw=1")
	args = append(args, "-maxrate:v:0", "600k", "-b:a:0", "500k")
	args = append(args, "-filter:v:1")
	args = append(args, "scale=w=640:h=360:force_original_aspect_ratio=decrease, drawtext=textfile=activetext.txt:fontsize=12:fontcolor=white:reload=10:y=h-line_h-32:x=w-(mod(2*n\\,w+tw)-tw/20):box=1:boxcolor=black@0.5:boxborderw=2")
	args = append(args, "-maxrate:v:1", "1500k", "-b:a:1", "1000k")
	args = append(args, "-filter:v:2")
	args = append(args, "scale=w=1280:h=720:force_original_aspect_ratio=decrease, drawtext=textfile=activetext.txt:fontsize=24:fontcolor=white:reload=10:y=h-line_h-52:x=w-(mod(4*n\\,w+tw)-tw/40):box=1:boxcolor=black@0.5:boxborderw=5")
	args = append(args, "-maxrate:v:2", "3000k", "-b:a:2", "2000k")
	args = append(args, "-var_stream_map", "v:0,a:0,name:480p v:1,a:1,name:720p v:2,a:2,name:1080p")
	args = append(args, "-f", "hls")
	args = append(args, "-hls_list_size", "6", "-threads", "0", "-hls_time", "10",)
	args = append(args, "-hls_flags", "delete_segments")
	args = append(args, "-master_pl_name", "stream.m3u8", "-y", fmt.Sprintf("./media/%s", "stream-%v.m3u8"))

	go func() {
		cmd := exec.Command("ffmpeg", args...)
		if err := cmd.Start(); err != nil {
			panic(err.Error())
		}
	}()
}

