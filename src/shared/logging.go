package shared

import "log"

func Info(v ...any) {
	log.Println(v...)
}

func Debug(v ...any) {
	log.Println(v...)
}
