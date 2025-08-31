package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// getSystemFingerprint با هش کردن شناسه‌های سخت‌افزاری، یک اثر انگشت منحصر به فرد تولید می‌کند.
// این تابع یک تابع مخصوص پلتفرم را برای دریافت اجزای سخت‌افزاری فراخوانی می‌کند.
func getSystemFingerprint() (string, error) {
	// تابع مخصوص پلتفرم برای جمع‌آوری شناسه‌های سخت‌افزاری فراخوانی می‌شود.
	components, err := getPlatformComponents()
	if err != nil {
		return "", fmt.Errorf("جمع آوری اطلاعات سخت‌افزار با مشکل مواجه شد: %w", err)
	}

	// اگر هیچ شناسه‌ای یافت نشد، خطا برگردانده می‌شود.
	if len(components) == 0 {
		return "", fmt.Errorf("هیچ مشخصه‌ای برای تولید اثر انگشت یافت نشد")
	}

	// شناسه‌ها مرتب‌سازی می‌شوند تا اطمینان حاصل شود که اثر انگشت همیشه برای یک ماشین یکسان است.
	sort.Strings(components)

	// شناسه‌ها با یک جداکننده به یک رشته واحد متصل می‌شوند.
	joined := strings.Join(components, "||")

	// هش SHA-256 از رشته ترکیبی محاسبه می‌شود.
	hash := sha256.Sum256([]byte(joined))
	fingerprint := hex.EncodeToString(hash[:])

	return fingerprint, nil
}
