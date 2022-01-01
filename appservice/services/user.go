package services

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"time"
	"waho/comm"
	"waho/core"
	"waho/register/base"
	"waho/static"
)

func SendPhoneCode(body string, d *amqp.Delivery)  {
	var cache static.SendPhoneCodeCurrentDayLogCache
	err := json.Unmarshal([]byte(body), &cache)
	if err != nil {
		log.Info("SendPhoneCode Err: ", err)
		d.Ack(false)
		return
	}
	err = base.ValidateStruct(&cache)
	if err != nil {
		log.Info("SendPhoneCode Err: ", err)
		d.Ack(false)
		return
	}
	cacheList := core.GetCore().GetSendCodeCache(cache.UserId)
	//TODO SEND SMS
	cache.Status = 1
	cache.SendTime = time.Now().Unix()
	cacheList = append(cacheList, cache)

	base.RdbHSet(static.SendPhoneCodeCurrentDayLogKey,cache.UserId , comm.ToJsonNotErr(cacheList))

	// set 24:00 expire
	t := time.Now()
	timeAddDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, 1).UTC()
	endtime := t.UTC()
	base.RdbExpire(static.SendPhoneCodeCurrentDayLogKey, timeAddDay.Sub(endtime))
	log.Info("SendPhoneCode Success: ", comm.ToJsonNotErr(cache))
	d.Ack(true)
}