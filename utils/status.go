package utils

import (
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"

	"github.com/sofianinho/vnf-api-golang/vnf/types"
)

var cpu_tick = float64(100)

//SysStatus returns the system status in terms of load average, cpu, memory, and network
func SysStatus()(*types.Status, error) {
	a, err := ldStatus()
	if err!=nil{
		return nil, err
	}
	c, err := cpuStatus()
	if err!=nil{
		return nil, err
	}
	m, err := memStatus()
	if err!=nil{
		return nil, err
	}
	n, err := netStatus()
	if err!=nil{
		return nil, err
	}
	ret := &types.Status{
		Load:	a,
		CPU:	c,
		Memory:	m,
		Net:	n,
	}
	return ret, nil
}

func netStatus()(*types.Net,error){
	n, err := net.IOCounters(false)
	if err != nil {
		return nil, fmt.Errorf("Cannot get network stats on system: %s", err)
	} 

	ret := &types.Net{
		BSent: ConvertBytes(n[0].BytesSent),
		BRecv: ConvertBytes(n[0].BytesRecv),
		PSent: ConvertBytes(n[0].PacketsSent),
		PRecv: ConvertBytes(n[0].PacketsRecv),
	}

	return ret, nil
}

func cpuStatus()(*types.Cpu, error){
	c, err := cpu.Times(false)
	if err != nil {
		return nil, fmt.Errorf("Cannot get system cpu stats: %s", err)
	}
	ret := &types.Cpu{
		User:	c[0].User/cpu_tick,
		System:	c[0].System/cpu_tick,
		Idle:	c[0].Idle/cpu_tick,
	}
	return ret, nil
}

func ldStatus()(*types.Load, error){
	a, err := load.Avg()
	if err != nil {
		return nil, fmt.Errorf("Cannot get system load average: %s", err)
	}
	ret := &types.Load{
		Load1:	a.Load1,
		Load5:	a.Load5,
		Load15:	a.Load15,
	}
	return ret, nil
}

func memStatus()(*types.Mem, error){
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("Cannot get memory stats: %s", err)
	}

	ret := &types.Mem{
		Total:	ConvertBytes(v.Total),
		Used:	ConvertBytes(v.Used),
		Free:	ConvertBytes(v.Free),
	}
	return ret, nil
}
