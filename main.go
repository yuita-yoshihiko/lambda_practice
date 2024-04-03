package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

func main() {
	// AWS Session の初期化
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1"), // 例: us-east-1
	})
	if err != nil {
		fmt.Println("AWS session の作成に失敗:", err)
		return
	}

	// SSM Service のクライアントを作成
	ssmSvc := ssm.New(sess)

	// Parameter Store から Slack Webhook URL を取得
	paramName := "lambda_practice" // 保存したパラメータの名前
	param, err := ssmSvc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(paramName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		fmt.Println("Parameter の取得に失敗:", err)
		return
	}

	slackWebhookURL := *param.Parameter.Value

	// Slack に送信するメッセージの構造体
	message := struct {
		Text string `json:"text"`
	}{
		Text: "こんにちは、Slack!",
	}

	// JSON 形式にエンコード
	messageBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("JSON エンコードに失敗:", err)
		return
	}

	// Slack の Incoming Webhook URL に POST リクエストを送信
	resp, err := http.Post(slackWebhookURL, "application/json", bytes.NewBuffer(messageBytes))
	if err != nil {
		fmt.Println("Slack への POST リクエストに失敗:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("メッセージを Slack に送信しました。ステータスコード:", resp.StatusCode)
}
