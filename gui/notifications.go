package gui

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)

// ShowNotification displays a Windows toast notification
func ShowNotification(title, message string) {
	// Use PowerShell to show Windows 10/11 toast notification
	script := fmt.Sprintf(`
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
$template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02)
$toastXml = [xml] $template.GetXml()
$toastXml.GetElementsByTagName("text")[0].AppendChild($toastXml.CreateTextNode("%s")) > $null
$toastXml.GetElementsByTagName("text")[1].AppendChild($toastXml.CreateTextNode("%s")) > $null
$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($toastXml.OuterXml)
$toast = [Windows.UI.Notifications.ToastNotification]::new($xml)
$toast.ExpirationTime = [DateTimeOffset]::Now.AddSeconds(3)
$notifier = [Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("Text Expander")
$notifier.Show($toast)
`, escapeForPowerShell(title), escapeForPowerShell(message))

	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to show notification: %v", err)
	}
}

// ShowExpansionNotification shows a notification when text is expanded
func ShowExpansionNotification(trigger, replacement string) {
	// Truncate long replacements for notification
	displayText := replacement
	if len(displayText) > 50 {
		displayText = displayText[:47] + "..."
	}

	message := fmt.Sprintf("%s → %s", trigger, displayText)
	ShowNotification("Text Expanded", message)
}

// escapeForPowerShell escapes special characters for PowerShell
func escapeForPowerShell(s string) string {
	s = strings.ReplaceAll(s, `"`, `'`)
	s = strings.ReplaceAll(s, "`", "")
	s = strings.ReplaceAll(s, "$", "")
	return s
}
