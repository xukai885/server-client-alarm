package watchClient

import (
	"fmt"
	"log"
	"server-client-alarm/modules"
	"server-client-alarm/settings"
	"server-client-alarm/tosend"
	"time"
)

var ClientListSum []*modules.Client
var lastClientListSum []modules.Client

var noChange map[string]string

func Init() {
	// 启动定时任务

	ticker := time.NewTicker(settings.Conf.Watch.Wait * time.Second)
	defer ticker.Stop()

	noChange := make(map[string]string)
	InitClientListSum()
	for {
		select {
		case <-ticker.C:
			log.Println("监听中")
			// 去发送短信
			hasTimeChanged(noChange)
			if noChange != nil {
				for id, _ := range noChange {
					log.Printf("%s已掉线", id)
					for i := range ClientListSum {
						if ClientListSum[i].Id == id {
							// 调用webhook
							tosend.ToWebHook(fmt.Sprintf("IP:%s,名字:%s掉线了", ClientListSum[i].Ip, ClientListSum[i].Name), "Alarm")
						}
					}
				}
			}
		}
	}
}
func hasTimeChanged(noChange map[string]string) {

	// 遍历每个对象，比较 Time 是否有变化
	for i := range ClientListSum {
		if ClientListSum[i].Time != lastClientListSum[i].Time {
			lastClientListSum[i].Time = ClientListSum[i].Time
			delete(noChange, ClientListSum[i].Id)
			tosend.ToWebHook(fmt.Sprintf("IP:%s,名字:%s恢复了", ClientListSum[i].Ip, ClientListSum[i].Name), "Recover")
		} else if time.Since(ClientListSum[i].Time) > settings.Conf.Watch.Threshold*time.Minute {
			noChange[ClientListSum[i].Id] = "true"
		}
	}
}

func InitClientListSum() {
	for i := range ClientListSum {
		ClientListSum[i].Time = time.Now()
		lastClientListSum = append(lastClientListSum, *ClientListSum[i])
	}
}
