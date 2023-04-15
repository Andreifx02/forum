package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func (b *Bot) SignUp(nickname string) error {
	type request struct {
		Nickname string `json:"nickname"`
	}
	bts, _ := json.Marshal(request{
		Nickname: nickname,
	})
	resp, err := http.Post(fmt.Sprintf("%s/user/create", b.serverAddress), "application/json", bytes.NewBuffer(bts))	
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(resp.Status)
	}  
	return nil
}