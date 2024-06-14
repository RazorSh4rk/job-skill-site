package browser

import (
	"fmt"
	"log"
	"os"

	"github.com/tebeka/selenium"
)

type Browser struct {
	driverPath string
	port       int
	opts       []selenium.ServiceOption
	caps       selenium.Capabilities
	service    *selenium.Service
	remoteUrl  string
}

type IBrowser interface {
	FromFirefox()
	FromChrome()
	Destroy()
	GetPage(url string)
}

func (b *Browser) FromFirefox() {
	b.driverPath = fmt.Sprintf("%s/geckodriver", getPath())
	b.port = 6789
	b.remoteUrl = "http://localhost:%d"
	b.opts = []selenium.ServiceOption{
		selenium.StartFrameBuffer(),
		selenium.GeckoDriver(b.driverPath),
		selenium.Output(nil),
	}
	b.caps = selenium.Capabilities{"browserName": "firefox"}
	service, err := selenium.NewGeckoDriverService(b.driverPath, b.port, b.opts...)
	if err != nil {
		log.Fatal("Error creating firefox service: ", err.Error())
	}
	b.service = service
}

func (b *Browser) FromChrome() {
	b.driverPath = fmt.Sprintf("%s/chromedriver", getPath())
	b.port = 6789
	b.remoteUrl = "http://localhost:%d/wd/hub"
	b.opts = []selenium.ServiceOption{
		selenium.ChromeDriver(b.driverPath),
		selenium.Output(nil),
	}
	b.caps = selenium.Capabilities{"browserName": "chrome"}
	// b.caps.AddChrome(chrome.Capabilities{
	// 	Args: []string{
	// 		"--headless",
	// 	},
	// })
	service, err := selenium.NewChromeDriverService(b.driverPath, b.port, b.opts...)
	if err != nil {
		log.Fatal("Error creating chrome service: ", err.Error())
	}
	b.service = service
}

func (b *Browser) Destroy() {
	b.service.Stop()
}

func (b *Browser) GetPage(url string) string {
	if b == &(Browser{}) {
		log.Fatal("Initialize the browser with .FromFirefox() or .FromChrome()")
	}

	driver, err := selenium.NewRemote(b.caps, fmt.Sprintf(b.remoteUrl, b.port))
	if err != nil {
		fmt.Errorf("Error creating chrome driver: %s", err.Error())
	}
	defer driver.Quit()

	err = driver.Get(url)
	if err != nil {
		log.Fatal("Could not get page: ", err)
	}

	html, err := driver.PageSource()
	if err != nil {
		log.Fatal("Could not get page source:", err)
	}
	return html
}

func getPath() string {
	p, err := os.Getwd()
	if err != nil {
		return "."
	}
	return p
}
