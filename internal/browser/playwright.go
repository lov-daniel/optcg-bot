package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mxschmitt/playwright-go"
)

func ClickButton(page playwright.Page, className string) {
	// Find visible sign-in buttons
	fmt.Println((page.URL()))
	buttons := page.Locator(className)

	count, err := buttons.Count()
	check(err)

	fmt.Println("Visible sign-in buttons found:", count)

	if count == 0 {
		log.Fatal("Could not find Sign In button")
	}

	button := buttons.First()

	text, err := button.TextContent()
	check(err)

	fmt.Println("Clicking button:\n", text)

	err = button.Click(playwright.LocatorClickOptions{
		Timeout: playwright.Float(10000),
	})
	check(err)

	fmt.Printf("Clicked: %s\n", className)
}

func main() {
	pw, err := playwright.Run()
	check(err)
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	check(err)
	defer browser.Close()

	context, err := browser.NewContext()
	check(err)
	defer context.Close()

	page, err := context.NewPage()
	check(err)

	// Increase timeout for dynamic pages
	page.SetDefaultTimeout(10000)

	fmt.Println("Opening site...")

	_, err = page.Goto(
		"https://cardkaizoku.com",
		playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
		},
	)
	check(err)

	title, err := page.Title()
	check(err)

	fmt.Println("Page title:", title)

	ClickButton(page, "button.sign-in-button:visible")
	ClickButton(page, "button.auth-patreon-button:visible")
	time.Sleep(2 * time.Second)
	googleButton := page.
		FrameLocator("iframe[title='Sign in with Google Button']").
		Locator("div.nsm7Bb-HzV7m-LgbsSe-bN97Pc-sM5MNb")

	popup, err := page.ExpectPopup(func() error {
		err := googleButton.Click()
		return err
	})
	check(err)

	fmt.Println("Google popup opened:")
	fmt.Println(popup.URL())

	// Print URL changes (useful for OAuth debugging)
	page.On("framenavigated", func(frame playwright.Frame) {
		fmt.Println("Navigated:", frame.URL())
	})

	// Keep browser alive
	select {}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
