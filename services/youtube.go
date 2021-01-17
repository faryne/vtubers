package services

import (
	"context"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	"os"
	"time"
)

type (
	YoutubeStruct struct {
		Client     *youtube.Service
		SearchList *youtube.SearchListCall
		VideosList *youtube.VideosListCall
	}
)

func InitYoutubeClient(filename string) (*YoutubeStruct, error) {
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
		Client:     client,
		SearchList: client.Search.List([]string{"snippet"}),
		VideosList: client.Videos.List([]string{"liveStreamingDetails"}),
	}, nil
}
func (s *YoutubeStruct) getLives(channelId string, eventType string) (*youtube.SearchListResponse, error) {
	response, err := s.SearchList.ChannelId(channelId).EventType(eventType).Type("video").Do()
	return response, err
}
func (s *YoutubeStruct) GetUpcomingLive(channelId string) (*youtube.SearchListResponse, error) {
	response, err := s.getLives(channelId, "upcoming")
	return response, err
}
func (s *YoutubeStruct) GetCompletedLive(channelId string) (*youtube.SearchListResponse, error) {
	response, err := s.getLives(channelId, "completed")
	return response, err
}
func (s *YoutubeStruct) GetNowLive(channelId string) (*youtube.SearchListResponse, error) {
	response, err := s.getLives(channelId, "live")
	return response, err
}

func (s *YoutubeStruct) GetVideo(videoId string) (*youtube.Video, error) {
	video, err := s.VideosList.Id(videoId).Do()

	if err != nil {
		return nil, err
	}
	if len(video.Items) == 0 {
		return nil, nil
	}
	//video.Items[0].LiveStreamingDetails.ActualEndTime
	return video.Items[0], nil
}

func (s *YoutubeStruct) GetLiveMessages(livechatId string, callback func(*youtube.LiveChatMessageListResponse) error) error {
	req := s.Client.LiveChatMessages.List(livechatId, []string{"id", "snippet", "authorDetails"})

	var pageToken = ""
	for {
		resp, err := req.PageToken(pageToken).MaxResults(100).Do()
		if err != nil {
			return err
		}
		if len(resp.Items) == 0 {
			break
		}
		if resp.NextPageToken != "" {
			pageToken = resp.NextPageToken
		}
		callback(resp)
		time.Sleep(150 * time.Millisecond)
	}
	return nil
}
