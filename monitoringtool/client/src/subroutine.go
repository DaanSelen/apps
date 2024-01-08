package main

import (
	"log"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func startSubroutine() { //CONFVAL DIRECTORY
	uptime_on := confVal[1]
	log.Println("Uptime subroutine set to:", uptime_on)
	mon_Interval, err := strconv.Atoi(confVal[3])
	if err != nil {
		log.Println("No or incorrect monitoring interval set, defaulting to 30 seconds.")
		mon_Interval = 30
	}
	log.Println("Monitoring interval set to", mon_Interval, "seconds")

	subroutine(uptime_on, mon_Interval)
}

func subroutine(uptime_on string, mon_Interval int) {
	for {
		var uptime *time.Duration
		cpu := checkCPU()
		ram := checkRAM()
		if uptime_bool := (strings.ToLower(uptime_on) == "true"); uptime_bool {
			uptime = checkUptime()
			log.Println("USAGES, CPU:", strconv.Itoa(cpu)+"%", "RAM:", strconv.Itoa(ram)+"%", "UPTIME:", uptime)
			conn.Write([]byte("cpuutil:" + strconv.Itoa(cpu)))
			conn.Write([]byte("ramutil:" + strconv.Itoa(ram)))
			conn.Write([]byte("uptime:" + uptime.String()))
			log.Println("DATA SEND TO SERVER")
		} else {
			log.Println("USAGES, CPU:", strconv.Itoa(cpu)+"%", "RAM:", strconv.Itoa(ram)+"%")
			conn.Write([]byte("cpuutil:" + strconv.Itoa(cpu)))
			conn.Write([]byte("ramutil:" + strconv.Itoa(ram)))
			log.Println("DATA SEND TO SERVER")
		}

		time.Sleep(time.Duration(mon_Interval-1) * time.Second)
	}
}

func checkCPU() int {
	percentages, err := cpu.Percent(time.Second, false)
	if err == nil && len(percentages) > 0 {
		return int(percentages[0])
	} else {
		return 0
	}
}

func checkRAM() int {
	memory, err := mem.VirtualMemory()
	if err == nil {
		return int((float64(memory.Used) / float64(memory.Total)) * 100.0)
	} else {
		return 0
	}
}

func checkUptime() *time.Duration {
	var info syscall.Sysinfo_t
	syscall.Sysinfo(&info)
	uptime := time.Duration(info.Uptime) * time.Second
	return &uptime
}
