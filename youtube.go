package vtubers

import (
	"context"
	"fmt"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"os"
	"time"
)

type (
	YoutubeStruct struct {
		ChannelId   string
		Client      *youtube.Service
		SearchList  *youtube.SearchListCall
		VideosList  *youtube.VideosListCall
		ChannelList *youtube.ChannelsListCall
	}
)

func New(filename string, channelId string) (*YoutubeStruct, error) {
	var err error
	var client *youtube.Service

	_, err = os.Stat(filename)
	if err != nil {
		return &YoutubeStruct{}, err
	}
	client, err = youtube.NewService(context.Background(), option.WithCredentialsFile(filename))
	if err != nil {
		return &YoutubeStruct{}, err
	}
	return &YoutubeStruct{
		ChannelId:   channelId,
		Client:      client,
		SearchList:  client.Search.List([]string{"snippet"}),
		VideosList:  client.Videos.List([]string{"liveStreamingDetails", "snippet"}),
		ChannelList: client.Channels.List([]string{"snippet", "statistics", "brandingSettings"}),
	}, nil
}

func (s *YoutubeStruct) GetChannelInfo() (*youtube.ChannelListResponse, error) {
	response, err := s.ChannelList.Id(s.ChannelId).Do()
	//response.Items[0].Statistics.SubscriberCount
	return response, err
}

func (s *YoutubeStruct) getLives(eventType string) (*youtube.SearchListResponse, error) {
	response, err := s.SearchList.ChannelId(s.ChannelId).EventType(eventType).Type("video").Do()
	return response, err
}
func (s *YoutubeStruct) GetUpcomingLive() (*youtube.SearchListResponse, error) {
	response, err := s.getLives("upcoming")
	return response, err
}
func (s *YoutubeStruct) GetCompletedLive() (*youtube.SearchListResponse, error) {
	response, err := s.getLives("completed")
	return response, err
}
func (s *YoutubeStruct) GetNowLive() (*youtube.SearchListResponse, error) {
	response, err := s.getLives("live")
	return response, err
}

func (s *YoutubeStruct) GetVideo(videoId string) (*youtube.VideoListResponse, error) {
	video, err := s.VideosList.Id(videoId).Do()

	if err != nil {
		return nil, err
	}
	return video, nil
}

func (s *YoutubeStruct) GetLiveMessages(livechatId string, callback func(*youtube.LiveChatMessageListResponse) error) error {
	req := s.Client.LiveChatMessages.List(livechatId, []string{"id", "snippet", "authorDetails"}).MaxResults(100)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
				time.Sleep(150 * time.Millisecond)
				return
			default:
				if ctx.Err() != nil {
					fmt.Println(ctx.Err())
					return
				}
				time.Sleep(150 * time.Millisecond)
			}
		}
	}()

	err := req.Pages(ctx, func(response *youtube.LiveChatMessageListResponse) error {
		time.Sleep(150 * time.Millisecond)

		return callback(response)
	})
	cancel()
	return err
}
