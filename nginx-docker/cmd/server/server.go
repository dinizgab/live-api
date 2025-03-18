package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

    room := r.URL.Query().Get("room")
    if room == "" {
        http.Error(w, "Room parameter is required", http.StatusBadRequest)
        return
    }

	ffmpeg := exec.Command("ffmpeg",
		"-i", "pipe:0",
		"-c:v", "libx264", "-preset", "ultrafast", "-tune", "zerolatency",
		"-c:a", "aac",
		"-f", "flv", fmt.Sprintf("rtmp://localhost:1935/live/%s", room),
	)

    fmt.Println(room)

	ffmpegIn, err := ffmpeg.StdinPipe()
	if err != nil {
		log.Fatalf("Error creating FFmpeg pipe: %v+", err)
	}

	if err := ffmpeg.Start(); err != nil {
		log.Fatalf("Error starting FFmpeg: %v+", err)
	}

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message: ", err)
			break
		}

		_, err = ffmpegIn.Write(data)
		if err != nil {
			fmt.Println("Error writing message: ", err)
			break
		}
	}

    ffmpegIn.Close()
    ffmpeg.Wait()
}

func main() {
	http.HandleFunc("/stream", streamHandler)

    log.Println("WebSocket server running on ws://localhost:3000/stream")

    log.Fatal(http.ListenAndServe(":3000", nil))
}
