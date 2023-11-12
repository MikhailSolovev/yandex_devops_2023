package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"syscall"

	"github.com/shirou/gopsutil/mem"
)

func main() {
	p := NewProcess()

	if err := p.getPathToExecutable(); err != nil {
		log.Fatalf("failed to get path to exec: %v", err)
	}
	if err := p.getMemStats(); err != nil {
		log.Fatalf("failed to get mem stats: %v", err)
	}
	if err := p.getNumberOfUsedDescriptors(); err != nil {
		log.Fatalf("failed to get descriptors: %v", err)
	}

	s := NewSystem()

	if err := s.getMemInfo(); err != nil {
		log.Fatalf("failed to get mem info: %v", err)
	}

	if err := s.getProcessorInfo(); err != nil {
		log.Fatalf("failed to get processor info: %v", err)
	}

	result := Result{
		Process: p,
		System:  s,
	}

	data, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Fatalf("failed to marshall result: %v", err)
	}

	fmt.Println(string(data))
}

type Result struct {
	Process *Process `json:"process"`
	System  *System  `json:"system"`
}

type Process struct {
	Path       string `json:"path"`
	Descriptor struct {
		SoftMax uint64 `json:"soft_max"`
		HardMax uint64 `json:"hard_max"`
	} `json:"descriptor"`
	Mem struct {
		Heap struct {
			// Heap memory allocated and still in use
			Alloc uint64 `json:"alloc"`
			// Heap memory obtained from the OS
			Sys uint64 `json:"sys"`
			// Heap memory in us
			Inuse uint64 `json:"inuse"`
			// Heap memory that is idle
			Idle uint64 `json:"idle"`
			// Heap memory released to the OS
			Released uint64 `json:"released"`
		} `json:"heap"`
		Stack struct {
			// Stack memory obtained from the OS
			Sys uint64 `json:"sys"`
			// Stack memory in use
			Inuse uint64 `json:"inuse"`
		} `json:"stack"`
		// Total memory obtained from the OS
		Sys uint64 `json:"sys"`
		// Bytes allocated and still in use
		Alloc uint64 `json:"alloc"`
		// Total bytes allocated (even if freed)
		TotalAlloc uint64 `json:"total_alloc"`
	} `json:"mem"`
}

func NewProcess() *Process {
	return &Process{}
}

func (p *Process) getPathToExecutable() error {
	path, err := os.Executable()
	if err != nil {
		return err
	}

	p.Path = path

	return nil
}

func (p *Process) getNumberOfUsedDescriptors() error {
	var rlimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		return err
	}

	p.Descriptor.SoftMax = rlimit.Cur
	p.Descriptor.HardMax = rlimit.Max

	return nil
}

func (p *Process) getMemStats() error {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	p.Mem.Heap.Alloc = memStats.HeapAlloc
	p.Mem.Heap.Sys = memStats.HeapSys
	p.Mem.Heap.Inuse = memStats.HeapInuse
	p.Mem.Heap.Idle = memStats.HeapIdle
	p.Mem.Heap.Released = memStats.HeapReleased

	p.Mem.Stack.Sys = memStats.StackSys
	p.Mem.Stack.Inuse = memStats.StackInuse

	p.Mem.Sys = memStats.Sys
	p.Mem.Alloc = memStats.Alloc
	p.Mem.TotalAlloc = memStats.TotalAlloc

	return nil
}

type System struct {
	Processor struct {
		NumCPU int `json:"num_cpu"`
	} `json:"processor"`
	Mem struct {
		Active uint64 `json:"active"`
		Used   uint64 `json:"used"`
		Total  uint64 `json:"total"`
	} `json:"mem"`
}

func (s *System) getMemInfo() error {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	s.Mem.Active = memInfo.Active
	s.Mem.Used = memInfo.Used
	s.Mem.Total = memInfo.Total

	return nil
}

func (s *System) getProcessorInfo() error {
	numCPU := runtime.NumCPU()

	s.Processor.NumCPU = numCPU

	return nil
}

func NewSystem() *System {
	return &System{}
}
