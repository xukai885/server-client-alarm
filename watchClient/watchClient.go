package watchClient

import (
	"fmt"
	"server-client-alarm/modules"
	"server-client-alarm/settings"
	"server-client-alarm/tosend"
	"time"
)

var ClientListSum []*modules.Client
var LastClientListSum []modules.Client

var noChange map[string]int64

func Init() {
	// 启动定时任务

	ticker := time.NewTicker(settings.Conf.Watch.Wait * time.Second)
	defer ticker.Stop()

	noChange := make(map[string]int64)
	InitClientListSum()
	for {
		select {
		case <-ticker.C:
			// 去发送短信
			hasTimeChanged(noChange)
			if noChange != nil {
				for id, v := range noChange {
					if v != 2 {
						for i := range ClientListSum {
							if ClientListSum[i].Id == id {
								noChange[ClientListSum[i].Id] = 2
								tosend.ToWebHook(fmt.Sprintf("IP:%s,名字:%s掉线了", ClientListSum[i].Ip, ClientListSum[i].Name), "Alarm")
							}
						}
					}
				}
			}
		}
	}
}
func hasTimeChanged(noChange map[string]int64) {
outerLoop:
	// 遍历每个对象，比较 Time 是否有变化
	for i := range ClientListSum {
		if ClientListSum[i].Time != LastClientListSum[i].Time {
			LastClientListSum[i].Time = ClientListSum[i].Time
			old_len := len(noChange)
			delete(noChange, ClientListSum[i].Id)
			new_len := len(noChange)
			if new_len < old_len {
				tosend.ToWebHook(fmt.Sprintf("IP:%s,名字:%s恢复了", ClientListSum[i].Ip, ClientListSum[i].Name), "Recover")
			}
		} else {
			if time.Since(ClientListSum[i].Time) > settings.Conf.Watch.Threshold*time.Minute {
				if noChange[ClientListSum[i].Id] == 3 || noChange[ClientListSum[i].Id] == 2 {
				} else if noChange[ClientListSum[i].Id] == 1 {
					noChange[ClientListSum[i].Id] = 2
				} else {
					noChange[ClientListSum[i].Id] = 1
					break outerLoop
				}
				if time.Since(ClientListSum[i].Time) > settings.Conf.Watch.Threshold*time.Minute && time.Since(ClientListSum[i].Time) <= settings.Conf.Watch.RepeatInterval*time.Minute+settings.Conf.Watch.Threshold*time.Minute {
					noChange[ClientListSum[i].Id] = 2
				} else {
					noChange[ClientListSum[i].Id] = 3
					ClientListSum[i].Time = time.Now()
					LastClientListSum[i].Time = ClientListSum[i].Time
				}
			}
		}
	}
}

func InitClientListSum() {
	for i := range ClientListSum {
		ClientListSum[i].Time = time.Now()
		LastClientListSum = append(LastClientListSum, *ClientListSum[i])
	}
}
