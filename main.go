package main

import (
	"fmt"
	"log"
)

func main() {
	// تابع اصلی را برای گرفتن اثر انگشت فراخوانی می‌کند
	fingerprint, err := getSystemFingerprint()
	if err != nil {
		log.Fatalf("خطا در تولید اثر انگشت: %v", err)
	}

	// نتیجه را در خروجی چاپ می‌کند
	fmt.Println("--- اثر انگشت سیستم شما ---")
	fmt.Println(fingerprint)
	fmt.Println("---------------------------")
}
