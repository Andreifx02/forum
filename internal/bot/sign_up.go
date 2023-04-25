package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf(string(b))
	}  
	return nil
}