package main

import (
	"sort"
	"strings"

	"github.com/StackExchange/wmi"
)

func getPlatformComponents() ([]string, error) {
	var components []string

	// System UUID
	var s []struct{ UUID string }
	if err := wmi.Query("SELECT UUID FROM Win32_ComputerSystemProduct", &s); err == nil && len(s) > 0 && s[0].UUID != "" {
		components = append(components, "UUID:"+s[0].UUID)
	}

	// Motherboard Serial Number
	var board []struct{ SerialNumber string }
	if err := wmi.Query("SELECT SerialNumber FROM Win32_BaseBoard", &board); err == nil && len(board) > 0 && board[0].SerialNumber != "" {
		components = append(components, "BOARD:"+board[0].SerialNumber)
	}

	// First Disk Drive Serial Number (only if valid)
	var disk []struct {
		SerialNumber string
		Index        int
	}
	if err := wmi.Query("SELECT SerialNumber, Index FROM Win32_DiskDrive", &disk); err == nil {
		for _, d := range disk {
			if d.Index == 0 {
				serial := strings.TrimSpace(d.SerialNumber)
				// فقط اگر serial معتبر باشد اضافه کن
				if serial != "" && !strings.HasPrefix(serial, "0000") && !strings.Contains(serial, "0000_") {
					components = append(components, "DISK:"+serial)
				}
				break
			}
		}
	}

	// MAC Address of the first physical adapter
	var mac []struct {
		MACAddress      string
		PhysicalAdapter bool
	}
	if err := wmi.Query("SELECT MACAddress, PhysicalAdapter FROM Win32_NetworkAdapter WHERE PhysicalAdapter = True", &mac); err == nil {
		for _, m := range mac {
			if m.MACAddress != "" {
				components = append(components, "MAC:"+m.MACAddress)
				break
			}
		}
	}

	// CPU Processor ID
	var cpu []struct{ ProcessorId string }
	if err := wmi.Query("SELECT ProcessorId FROM Win32_Processor", &cpu); err == nil && len(cpu) > 0 && cpu[0].ProcessorId != "" {
		components = append(components, "CPU:"+cpu[0].ProcessorId)
	}

	// 🔥 مرتب‌سازی الفبایی
	sort.Strings(components)

	return components, nil
}
