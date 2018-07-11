package main

import "log"

// for now mock and print to stdout
// better would be to use something like logrus

func LogIt(q interface{}) {
	log.Println(q)
}
