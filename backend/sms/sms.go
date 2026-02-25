package sms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var (
	Username string
	ApiKey   string
	SenderID string
)

func Init() {
	Username = os.Getenv("SMS_USERNAME")
	ApiKey = os.Getenv("SMS_API_KEY")
	SenderID = os.Getenv("SMS_SENDER_ID")
}

func SendSMS(to string, message string) error {
	if Username == "" || ApiKey == "" {
		fmt.Println("SMS simulation (no config):", to, "->", message)
		return nil
	}

	url := "https://api.sandbox.africastalking.com/version1/messaging"

	payload := map[string]string{
		"username": Username,
		"to":       to,
		"message":  message,
		"from":     SenderID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", ApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fmt.Println("SMS sent! Status:", resp.Status)
	return nil
}
