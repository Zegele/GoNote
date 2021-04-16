// Package storage is part of a hypothetical cloud storage server.
//!+main
package storage

import (
	"fmt"
	"log"
	"net/smtp"
)

var usage = make(map[string]int64)

func bytesInUse(username string) int64 { return usage[username] }

// Emain sender configuration.
// NOTE: never put passwords in source code!
const sender = "notification@example.com"
const password = "correcthorsebatterystaple"
const hostname = "smtp.example.com"

const template = `Warning: you are using %d bytes of storage,
%d%% of your quota.`

// %d%%是啥意思？？？

func CheckQuota(username string) {
	used := bytesInUse(username)
	const quota = 100000000 //1GB
	percent := 100 * used / quota
	if percent < 90 {
		return // OK
	}
	msg := fmt.Sprintf(template, used, percent)
	auth := smtp.PlainAuth("", sender, password, hostname)
	err := smtp.SendMail(hostname+":587", auth, sender,
		[]string{username}, [](msg))
	if err != nil {
		log.Printf("smtp.SendMail(%s) failded: %s", username, err)
	}
}

//!-main
