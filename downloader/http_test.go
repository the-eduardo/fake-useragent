package downloader

import (
	"fmt"
	"github.com/walkmiao/fake-useragent/setting"
	"testing"
)

func TestDownload_Get(t *testing.T) {
	downloader := Download{
		Delay:   setting.HTTP_DELAY,
		Timeout: setting.HTTP_TIMEOUT,
	}

	resp, err := downloader.Get("https://developers.whatismybrowser.com")
	if err != nil {
		t.Errorf("downloader.Get err: %v", err)
	}
	fmt.Println("Status\n", resp.Status, "\nResp:\n", resp)
}
