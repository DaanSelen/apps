package main

type Alert struct {
	ID       int    `json:"id"`
	Hostname string `json:"hostname"`
	IpAddr   string `json:"ipaddress"`
	Comp     string `json:"comp"`
	Time     string `json:"time"`
}

func main() {
	go initDBConnection()
	initHTTP()
}

func monitorUptimeAlert(alert Alert) {
	insertEntry(alert)
}

func monitorCPUAlert(alert Alert) {
	insertEntry(alert)
}

func monitorRAMAlert(alert Alert) {
	insertEntry(alert)
}

func checkComp(component string) []Alert {
	return dataCheck(component)
}
