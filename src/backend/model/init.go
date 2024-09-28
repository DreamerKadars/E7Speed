package model

import "time"

func InitModel() {
	go func() {
		for {
			time.Sleep(time.Second)
			model.Run()
		}
	}()
}
