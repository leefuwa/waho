package main

import (
	"appservice/services"
	"github.com/streadway/amqp"
	"waho/register"
	"waho/register/base"
	"waho/static"
)


func main() {
	register.New().Start()
	forever := make(chan bool)
	// 启动先更新一轮缓存
	go func() {
		// 更新缓存
		//go services.UpdateCache()
	}()
	// 启动集合缓存更新器
	go base.Consumer(static.UpdWholeQueueName, UpdCacheQueue)
	// 发送短信任务
	go base.Consumer(static.SendPhoneCodeQueueName, services.SendPhoneCode)
	<-forever
}

// 更新缓存
// @author Fuwa
// @date 2020-12-02
func UpdCacheQueue(body string, d *amqp.Delivery)  {
	switch body {
		case static.ClassifyMusterKey:
			go func() {
				//services.UpdateCache()
				d.Ack(false)
			}()
			break;
	}
}