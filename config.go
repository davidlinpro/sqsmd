package main

import (
	"log"
	"os"
)

type StatsConfig struct {
	SQSBASEURL   string
	SQSPREFIX    string
	AWSREGION    string
	AWSACCESSKEY string
	AWSSECRETKEY string
}

var Config StatsConfig

func (c *StatsConfig) Init() {
	c.SQSBASEURL = os.Getenv("SQSBASEURL")
	if c.SQSBASEURL == "" {
		log.Fatal("Required SQSBASEURL")
	}
	c.SQSPREFIX = os.Getenv("SQSPREFIX")
	if c.SQSPREFIX == "" {
		log.Println("Empty SQSPREFIX -- Getting all values")
	}
	c.AWSREGION = os.Getenv("AWSREGION")
	if c.AWSREGION == "" {
		log.Fatal("Required AWSREGION")
	}
	c.AWSACCESSKEY = os.Getenv("AWSACCESSKEY")
	if c.AWSACCESSKEY == "" {
		log.Fatal("Required AWSACCESSKEY")
	}
	c.AWSSECRETKEY = os.Getenv("AWSSECRETKEY")
	if c.AWSSECRETKEY == "" {
		log.Fatal("Required AWSSECRETKEY")
	}
}
