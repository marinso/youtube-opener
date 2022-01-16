package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var APIKey = os.Getenv("YT_API_KEY")
var channelID = flag.String("chid", "UCQpBmjL9kJ768fty4tbqa6Q", "Channel ID")

func main() {
	flag.Parse()

	ctx := context.Background()
	youtubeService, err := youtube.NewService(ctx, option.WithAPIKey(APIKey))

	if err != nil {
		log.Fatalf("Error creating YouTube client: %v", err)
	}

	var lastVideoId string

	for {
		currentViedoID := checkForViedoID(youtubeService, *channelID)

		if lastVideoId != currentViedoID {
			lastVideoId = currentViedoID
			err := openBrowser("https://www.youtube.com/watch?v=" + lastVideoId)

			if err != nil {
				fmt.Println(err)
			}
		}

		time.Sleep(time.Second * 60)
	}
}

func checkForViedoID(service *youtube.Service, channelID string) string {
	call := service.Search.List([]string{"id"}).
		Order("date").
		ChannelId(channelID).
		MaxResults(1)

	response, err := call.Do()

	if err != nil {
		log.Fatal(err)
	}

	return response.Items[0].Id.VideoId
}

func openBrowser(url string) error {
	var args []string 

	switch runtime.GOOS {
	case "darwin":
		args = []string{"open"}
	case "windows":
		args = []string{"cmd", "/c", "start"}
	default: 
		args = []string{"xdg-open"}
	}

	cmd := exec.Command(args[0], append(args[1:], url)...)
	return cmd.Start()
}