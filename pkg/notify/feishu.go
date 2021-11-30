// Package notify @Description  TODO
// @Author  	 jiangyang
// @Created  	 2021/11/30 11:50 下午
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func PostFieShu(info string) error {
	msg := FeiShuMsg{
		MsgType: "text",
	}
	msg.Content.Text = fmt.Sprintf("working [ %s ] : %s", os.Getenv("APP_ENV"), info)
	marshal, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", "https://open.feishu.cn/open-apis/bot/v2/hook/67c37caa-a7c2-44b3-8726-b784081c2102", bytes.NewReader(marshal))
	if err != nil {
		return err
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return nil
}

type FeiShuMsg struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}
