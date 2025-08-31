//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// getPlatformComponents شناسه‌های سخت‌افزاری را در لینوکس از فایل سیستم و دستورات شل جمع‌آوری می‌کند.
func getPlatformComponents() ([]string, error) {
	var components []string

	// تابع کمکی برای خواندن یک فایل و افزودن محتوای آن به لیست شناسه‌ها
	readFile := func(path, prefix string) {
		data, err := os.ReadFile(path)
		if err == nil {
			val := strings.TrimSpace(string(data))
			if val != "" {
				components = append(components, prefix+val)
			}
		}
	}

	// تابع کمکی برای اجرای یک دستور شل و افزودن خروجی آن به لیست شناسه‌ها
	runCmd := func(prefix string, name string, args ...string) {
		out, err := exec.Command(name, args...).Output()
		if err == nil {
			val := strings.TrimSpace(string(out))
			if val != "" {
				components = append(components, prefix+val)
			}
		}
	}

	// UUID from DMI
	readFile("/sys/class/dmi/id/product_uuid", "UUID:")
	// Board Serial from DMI
	readFile("/sys/class/dmi/id/board_serial", "BOARD:")
	// Disk Serial for /dev/sda
	runCmd("DISK:", "lsblk", "-d", "-n", "-o", "SERIAL", "/dev/sda")

	// MAC address of the default network interface
	// این دستور اینترفیس شبکه پیش‌فرض را پیدا کرده و آدرس MAC آن را می‌خواند
	out, err := exec.Command("sh", "-c", "ip route get 1.1.1.1 | awk '/src/ {print $5}'").Output()
	if err == nil {
		iface := strings.TrimSpace(string(out))
		if iface != "" {
			readFile(fmt.Sprintf("/sys/class/net/%s/address", iface), "MAC:")
		}
	}

	return components, nil
}
