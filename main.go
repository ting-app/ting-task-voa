package main

import (
	"github.com/ting-app/ting-task-voa/ting"
	"log"
)

func main() {
	err := ting.RunTask("https://learningenglish.voanews.com/api/zpyp_e-rm_")

	if err != nil {
		log.Printf("Run task error %v\n", err)
	}
}
