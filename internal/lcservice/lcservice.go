package lcservice

import (
	"context"
	"errors"

	"github.com/Lunarisnia/socrates/chatlogger"
	"google.golang.org/api/youtube/v3"
)

type LiveChatService interface {
	Poll() (chatlogger.ChatLogger, error)
}

type liveChatImpl struct {
	youtubeService *youtube.Service

	latestChatId  string
	latestTokenId string
}

func NewService(ctx context.Context) (LiveChatService, error) {
	youtubeService, err := youtube.NewService(ctx)
	liveChat := liveChatImpl{
		youtubeService: youtubeService,
	}
	if err != nil {
		return nil, err
	}
	lives, err := youtubeService.LiveBroadcasts.
		List([]string{"snippet"}).Mine(true).Do()
	if err != nil {
		return nil, err
	}
	if len(lives.Items) < 1 {
		return nil, errors.New("no live broadcast found")
	}
	liveChat.latestChatId = lives.Items[0].Snippet.LiveChatId

	return &liveChat, nil
}

func (l *liveChatImpl) Poll() (chatlogger.ChatLogger, error) {
	chatContainer := chatlogger.NewChatLogger()
	chatResp, err := l.youtubeService.LiveChatMessages.
		List(l.latestChatId, []string{"snippet", "author_details"}).
		PageToken(l.latestTokenId).Do()
	if err != nil {
		return nil, err
	}
	l.latestTokenId = chatResp.NextPageToken

	for _, chat := range chatResp.Items {
		chatContainer.Log(chatlogger.Chat{
			Username: chat.AuthorDetails.DisplayName,
			Content:  chat.Snippet.TextMessageDetails.MessageText,
		})
	}

	return chatContainer, nil
}
