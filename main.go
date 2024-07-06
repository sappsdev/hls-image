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
	url := "http://191.97.14.38:25461/casa/aptv2023/216.ts"
	fmt.Println(url)

	app := "ffmpeg"
	args := []string{
		"-re",
		"-i", url,
		"-reconnect_at_eof", "1",
		"-reconnect_streamed", "1",
		"-reconnect_on_network_error", "1",
		"-reconnect_on_http_error", "404,502",
		"-err_detect", "ignore_err",
		"-c:v", "libx264",
		"-c:a", "aac",
		"-f", "hls",
		"-hls_list_size", "10",
		"-hls_time", "5",
		"-hls_flags", "delete_segments",
		"-y", "./media/stream.m3u8",
	}

	go func() {
		for {
			cmd := exec.Command(app, args...)
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
