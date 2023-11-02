package main

import (
	"log"
	"math"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

func enabledComp(compEnabledInput string) bool {
	if strings.ToLower(compEnabledInput) == "true" {
		return true
	} else {
		return false
	}
}

func checkCPUUsage() {
	comp := "cpu"
	if enabledComp(cpuEnabled) {
		var count int
		log.Println(monitorTag, comp, "MONITORING STARTED.")
		for {
			rawPerc, _ := cpu.Percent(time.Second, false) //Get CPU usage percent
			cpuPerc := math.Round(rawPerc[0]*100) / 100

			if cpuPerc >= cpuThreshold {
				count = AmountHit(count, cpuAlertTime, 1, comp)
			} else {
				count = AmountHit(count, cpuAlertTime, 0, comp)
			}
		}
	} else {
		log.Println(monitorTag, comp, "MONITORING DISABLED")
	}
}

func checkRAMUsage() {
	comp := "ram"
	if enabledComp(ramEnabled) {
		var count int
		log.Println(monitorTag, comp, "MONITORING STARTED.")
		for {
			rawPerc, _ := mem.VirtualMemory()
			ramPerc := math.Round(rawPerc.UsedPercent*100) / 100

			if ramPerc >= ramThreshold {
				count = AmountHit(count, ramAlertTime, 1, comp)
			} else {
				count = AmountHit(count, ramAlertTime, 0, comp)
			}
			time.Sleep(time.Second)
		}
	} else {
		log.Println(monitorTag, comp, "MONITORING DISABLED")
	}
}

func AmountHit(count, alertTime, command int, comp string) int {
	if count == (alertTime - 1) {
		log.Println(warningTag, "REPORTING AMOUNT HIT!")
		switch comp {
		case "cpu":
			go sendReport("cpu")
		case "ram":
			log.Println("RAM TRIGGERED")
		}
		return 0
	} else if command == 1 {
		count++
		log.Println(warningTag, "CURRENT COUNT:", count, "Second(s) FOR", comp)
		return count
	} else if command == 0 {
		return 0
	} else {
		return 0
	}
}
