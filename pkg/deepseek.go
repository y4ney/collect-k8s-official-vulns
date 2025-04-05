package pkg

// This file was generated from JSON Schema using quick type, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//	deepSeekRes, err := UnmarshalDeepSeekRes(bytes)
//	bytes, err = deepSeekRes.Marshal()

import (
	"bytes"
	"encoding/json"
	"golang.org/x/xerrors"
	"io"
	"net/http"
)

const (
	deepseekUrl  = "https://api.deepseek.com/chat/completions"
	systemPrompt = `你是一位经验丰富的 Kubernetes 安全专家，精通 Kubernetes 安全公告的翻译工作。
					你能精准地将英文内容翻译为地道、专业的简体中文，采用中国 Kubernetes 安全社区常用的术语。
					你会完整保留原文中的 Markdown 格式，包括标题、代码块、列表、链接等内容。不翻译代码块、命令、路径、配置字段，
					仅翻译说明性文字。不要添加任何解释或额外信息。`
	userPrompt = `请将以下 Kubernetes 安全公告内容翻译为简体中文，保持 Markdown 格式不变，不翻译代码块、命令和配置字段：

				---

`
)

type DeepSeekRes struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
}

type Choice struct {
	Index        int64       `json:"index"`
	Message      Message     `json:"message"`
	Logprobs     interface{} `json:"logprobs"`
	FinishReason string      `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens          int64               `json:"prompt_tokens"`
	CompletionTokens      int64               `json:"completion_tokens"`
	TotalTokens           int64               `json:"total_tokens"`
	PromptTokensDetails   PromptTokensDetails `json:"prompt_tokens_details"`
	PromptCacheHitTokens  int64               `json:"prompt_cache_hit_tokens"`
	PromptCacheMissTokens int64               `json:"prompt_cache_miss_tokens"`
}

type PromptTokensDetails struct {
	CachedTokens int64 `json:"cached_tokens"`
}

type DeepSeekReq struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func UnmarshalDeepSeekRes(data []byte) (DeepSeekRes, error) {
	var r DeepSeekRes
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *DeepSeekRes) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalDeepSeekReq(data []byte) (DeepSeekReq, error) {
	var r DeepSeekReq
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *DeepSeekReq) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func Request(content, apiKey string) (*DeepSeekRes, error) {
	// 设置 AI 的身份和提示词
	reqBody := DeepSeekReq{
		Model: "deepseek-chat",
		Messages: []ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt + content},
		},
		Stream: false,
	}
	jsonData, err := reqBody.Marshal()
	if err != nil {
		return nil, xerrors.Errorf("faild to marshal request body: %v", err)
	}
	req, err := http.NewRequest("POST", deepseekUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, xerrors.Errorf("faild to create request: %v", err)
	}

	// 设置 api key
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("faild to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("faild to read %s:%v", deepseekUrl, err)
	}
	resData, err := UnmarshalDeepSeekRes(body)
	if err != nil {
		return nil, xerrors.Errorf("faild to unmarshal %s:%v", deepseekUrl, err)
	}

	return &resData, nil
}
