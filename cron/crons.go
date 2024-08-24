package cron

import (
	"fmt"
	"gian/db"

	cron "gopkg.in/robfig/cron.v2"
)

func RunCron() {
	c := cron.New()

	//@every 00h00m00s
	c.AddFunc("@every 01h00m10s", backupUpload)

	c.Start()
}

func backupUpload() {
	fmt.Println("CRON FUNCtion calling in golang")
	db.Backup()
}
