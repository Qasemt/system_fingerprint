//go:build windows

package main

import (
	"github.com/StackExchange/wmi"
)

// getPlatformComponents شناسه‌های سخت‌افزاری را در ویندوز با استفاده از WMI جمع‌آوری می‌کند.
func getPlatformComponents() ([]string, error) {
	var components []string

	// System UUID
	var s []struct{ UUID string }
	if err := wmi.Query("SELECT UUID FROM Win32_ComputerSystemProduct", &s); err == nil && len(s) > 0 {
		components = append(components, "UUID:"+s[0].UUID)
	}

	// Motherboard Serial Number
	var board []struct{ SerialNumber string }
	if err := wmi.Query("SELECT SerialNumber FROM Win32_BaseBoard", &board); err == nil && len(board) > 0 {
		components = append(components, "BOARD:"+board[0].SerialNumber)
	}

	// First Disk Drive Serial Number
	var disk []struct{ SerialNumber string }
	if err := wmi.Query("SELECT SerialNumber FROM Win32_DiskDrive WHERE Index = 0", &disk); err == nil && len(disk) > 0 {
		components = append(components, "DISK:"+disk[0].SerialNumber)
	}

	// MAC Address of the first physical adapter
	var mac []struct{ MACAddress string }
	if err := wmi.Query("SELECT MACAddress FROM Win32_NetworkAdapter WHERE PhysicalAdapter = True AND MACAddress IS NOT NULL", &mac); err == nil && len(mac) > 0 {
		components = append(components, "MAC:"+mac[0].MACAddress)
	}

	// CPU Processor ID
	var cpu []struct{ ProcessorId string }
	if err := wmi.Query("SELECT ProcessorId FROM Win32_Processor", &cpu); err == nil && len(cpu) > 0 {
		components = append(components, "CPU:"+cpu[0].ProcessorId)
	}

	return components, nil
}
