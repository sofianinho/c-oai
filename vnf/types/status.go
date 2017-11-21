package types

//Load is for the load average
type Load struct{
	Load1	float64	`json:"load_avg_1"`
	Load5	float64 `json:"load_avg_5"`
	Load15	float64	`json:"load_avg_15"`
}

//Cpu is for the CPU stats
type Cpu struct{
	User	float64	`json:"user_per"`
	System	float64	`json:"sys_per"`
	Idle	float64	`json:"idle_per"`
}

//Mem is for the memory stats
type Mem struct{
	Total	string	`json:"total"`
	Used	string	`json:"user"`
	Free	string	`json:"free"`
}

//Net is for the network stats
type Net struct{
	BSent	string	`json:"bytes_sent"`
	BRecv	string	`json:"bytes_recv"`
	PSent	string	`json:"packets_sent"`
	PRecv	string	`json:"packets_recv"`
}

//Status is the global status for the VNF and instances
type Status struct{
	Load	*Load	`json:"load_avg"`
	CPU		*Cpu	`json:"cpu_stats"`
	Memory	*Mem	`json:"memory"`
	Net		*Net	`json:"network"`
}