package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mxschmitt/playwright-go"
)

var pw *playwright.Playwright
var page playwright.Page
var err error
var browser playwright.Browser
var context playwright.BrowserContext
var DEBUG = (os.Getenv("DEBUG") == "true")

func New(username string, password string) {
	// generates pw client
	pw, err = playwright.Run()
	check(err)

	// creating a new browser
	browser, err = pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(DEBUG),
	})
	check(err)

	// create browser context
	context, err = browser.NewContext()
	check(err)

	// create a page that we can use
	page, err = context.NewPage()
	check(err)

	page.SetDefaultTimeout(10000)

	fmt.Println("Opening site...")

	_, err = page.Goto(
		"https://cardkaizoku.com",
		playwright.PageGotoOptions{
			WaitUntil: playwright.WaitUntilStateNetworkidle,
		},
	)
	check(err)

	auth(username, password)
}

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

func Close() {
	pw.Stop()
	browser.Close()
	context.Close()
	page.Close()
}

func auth(username string, password string) {
	ClickButton(page, "button.sign-in-button:visible")
	page.GetByPlaceholder("Email").Fill(username)
	page.GetByPlaceholder("Password").Fill(password)
	ClickButton(page, "button.auth-submit-button:visible")

	page.Locator("button:has-text('Next'):visible").Click()
}

func main() {
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
