package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Andreifx02/forum/internal/domain"
	"github.com/google/uuid"
)

func (b *Bot) CreatePost(nickname string, topic string, text string, date time.Time) error {
	resp, err := http.Get(fmt.Sprintf("%s/user/%s", b.serverAddress, nickname))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(resp.Status)
	}
	
	
	var user domain.User
	json.NewDecoder(resp.Body).Decode(&user)

	type request struct {
		AuthorID uuid.UUID `json:"author_id"`
		Topic    string    `json:"topic"`
		Text     string    `json:"text"`
	}

	bts, _ := json.Marshal(request{
		AuthorID: user.ID,
		Topic: topic,
		Text: text,
	})
	
	resp, err = http.Post(fmt.Sprintf("%s/post/create",b.serverAddress),"application/json", bytes.NewBuffer(bts))
	if err != nil {
		return err
	}
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(resp.Status)
	}
	return nil
}