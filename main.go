package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
)

type GlobalConfig struct {
	PoolSize int
}

func DefaultGC() *GlobalConfig {
	gc := new(GlobalConfig)
	gc.PoolSize = 1000
	return gc
}

func main() {
	c := flag.String("c", "gofluent.conf", "config filepath")
	p := flag.String("p", "", "write cpu profile to file")
	v := flag.String("v", "error.log", "log file path")
	flag.Parse()

	f, err := os.OpenFile(*v, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("os.Open failed, err:", err)
	}
	defer f.Close()

	log.SetOutput(f)

	if *p != "" {
		f, err := os.Create(*p)
		if err != nil {
			log.Fatalln(err)
		}
		pprof.StartCPUProfile(f)
		defer func() {
			pprof.StopCPUProfile()
			f.Close()
		}()
	}

	gc := DefaultGC()
	config := NewPipeLineConfig(gc)
	config.LoadConfig(*c)

	Run(config)
}
