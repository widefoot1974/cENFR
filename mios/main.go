package main

import (
	"log"
	"os"
)

func main() {

	proc_name := os.Args[0]
	log.Println("############################################################")
	log.Printf(" [%v] Started.\n", proc_name)
	log.Println("############################################################")

	log.Println("############################################################")
	log.Printf(" [%v] Ended.\n", proc_name)
	log.Println("############################################################")

}
