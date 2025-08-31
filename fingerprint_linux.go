//go:build linux
// +build linux

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func getPlatformComponents() ([]string, error) {
	var components []string

	// تابع کمکی: خواندن فایل و تمیز کردن محتوا
	readFile := func(path string) (string, error) {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(data)), nil
	}

	// 1. System UUID (از DMI)
	if uuid, err := readFile("/sys/class/dmi/id/product_uuid"); err == nil && uuid != "" {
		components = append(components, "UUID:"+uuid)
	}

	// 2. شماره سریال مادربرد
	if board, err := readFile("/sys/class/dmi/id/board_serial"); err == nil && board != "" {
		components = append(components, "BOARD:"+board)
	}

	// 3. شماره سریال دیسک اول (sda)
	if serial, err := getDiskSerial("/dev/sda"); err == nil && serial != "" {
		components = append(components, "DISK:"+serial)
	}

	// 4. آدرس MAC اولین آداپتور فیزیکی (روی مسیر پیش‌فرض)
	if mac, err := getPrimaryMAC(); err == nil && mac != "" {
		components = append(components, "MAC:"+mac)
	}

	// 🔥 مرتب‌سازی برای ثبات (مثل پایتون)
	sort.Strings(components)

	return components, nil
}

func getDiskSerial(device string) (string, error) {
	cmd := exec.Command("lsblk", "-d", "-n", "-o", "SERIAL", device)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = io.Discard // خطاها نادیده گرفته می‌شوند

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	serial := strings.TrimSpace(out.String())
	if serial == "" {
		return "", fmt.Errorf("no serial found")
	}

	// فیلتر مقادیر پیش‌فرض/نامعتبر (مثل 0000...)
	if strings.HasPrefix(serial, "0000") || strings.Contains(serial, "0000_") {
		return "", fmt.Errorf("invalid serial")
	}

	return serial, nil
}

func getPrimaryMAC() (string, error) {
	// پیدا کردن اینترفیس پیش‌فرض
	cmd := exec.Command("ip", "route")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "default") {
			parts := strings.Fields(line)
			for i, field := range parts {
				if field == "dev" && i+1 < len(parts) {
					iface := parts[i+1]
					if iface == "" {
						continue
					}

					// خواندن MAC از /sys/class/net/...
					macPath := fmt.Sprintf("/sys/class/net/%s/address", iface)
					data, err := os.ReadFile(macPath)
					if err != nil {
						continue
					}

					mac := strings.TrimSpace(string(data))
					if mac != "" && mac != "00:00:00:00:00:00" && !strings.HasPrefix(mac, "ff:ff") {
						return mac, nil
					}
				}
			}
			break
		}
	}

	return "", fmt.Errorf("primary MAC not found")
}
