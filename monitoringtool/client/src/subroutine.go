package main

import (
	"log"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func subroutine() {
	cpu := checkCPU()
	ram := checkRAM()
	log.Println("USAGES, CPU:", strconv.Itoa(cpu)+"%", "RAM:", strconv.Itoa(ram)+"%")
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
