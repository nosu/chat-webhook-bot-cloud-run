package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func createMessage() string {
	msg := "この内容がチャットルームに投稿されます"

	return msg
}

func sendMessageToChat(msg string, webhookUrl string) {
	payload, err := json.Marshal(struct {
		Text string `json:"text"`
	}{
		Text: msg,
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := http.Post(webhookUrl, "application/json; charset=UTF-8", bytes.NewReader(payload))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP: %v\n", resp.StatusCode)
	}
}

func main() {
	url, ret := os.LookupEnv("WEBHOOK_URL")
	if !ret {
		log.Fatal("Environment Variable 'WEBHOOK_URL' is required. Set it in Cloud Run Console before run the app.")
	}

	rootHandler := func(w http.ResponseWriter, req *http.Request) {
		msg := createMessage()
		fmt.Println(msg)
		sendMessageToChat(msg, url)
	}

	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
