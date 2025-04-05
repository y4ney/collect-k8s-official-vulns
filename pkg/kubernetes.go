// This file was generated from JSON Schema using quick type, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    kubernetes, err := UnmarshalKubernetes(bytes)
//    bytes, err = kubernetes.Marshal()

package pkg

import (
	"encoding/json"
	"fmt"
	"golang.org/x/xerrors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Status string

const (
	k8sUrl = "https://kubernetes.io/docs/reference/issues-security/official-cve-feed/index.json"

	Fixed Status = "fixed"
	Open  Status = "open"
)

type Kubernetes struct {
	KubernetesIo KubernetesKubernetesIo `json:"_kubernetes_io"`
	Authors      []Author               `json:"authors"`
	Description  string                 `json:"description"`
	FeedURL      string                 `json:"feed_url"`
	HomePageURL  string                 `json:"home_page_url"`
	Items        []Item                 `json:"items"`
	Title        string                 `json:"title"`
	Version      string                 `json:"version"`
}

type Author struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Item struct {
	KubernetesIo  ItemKubernetesIo `json:"_kubernetes_io"`
	ContentText   string           `json:"content_text"`   // 内容描述，为 markdown 格式
	DatePublished time.Time        `json:"date_published"` // 公布时间
	ExternalURL   string           `json:"external_url"`   // 外部 URL，一般指的是 http://www.cve.org
	ID            string           `json:"id"`             // CVE 编号
	Status        Status           `json:"status"`         // 状态，共有 open 和 fixed 两种格式
	Summary       string           `json:"summary"`        // 总结标题
	URL           string           `json:"url"`            // github 上所对应的 issue 连接
}

type ItemKubernetesIo struct {
	GoogleGroupURL string `json:"google_group_url"` // Google 讨论组 url
	IssueNumber    int64  `json:"issue_number"`     // Github 上对应的 issue 号
}

type KubernetesKubernetesIo struct {
	FeedRefreshJob string    `json:"feed_refresh_job"`
	UpdatedAt      time.Time `json:"updated_at"` // 更新时间
}

func Fetch() (*Kubernetes, error) {
	resp, err := http.Get(k8sUrl)
	if err != nil {
		return nil, xerrors.Errorf("faild to fetch %s:%v", k8sUrl, err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("faild to read %s:%v", k8sUrl, err)
	}
	kubeData, err := UnmarshalKubernetes(body)
	if err != nil {
		return nil, xerrors.Errorf("faild to unmarshal %s:%v", k8sUrl, err)
	}

	return &kubeData, nil
}

func UnmarshalKubernetes(data []byte) (Kubernetes, error) {
	var r Kubernetes
	err := json.Unmarshal(data, &r)
	return r, err
}

func (k *Kubernetes) Marshal() ([]byte, error) {
	return json.Marshal(k)
}

func (i *Item) Save(dir string) (err error) {
	file := filepath.Join(dir, fmt.Sprintf("%s.md", i.ID))
	err = os.WriteFile(file, []byte(i.ContentText), 0644)
	if err != nil {
		return xerrors.Errorf("faild to write %s:%v", file, err)
	}
	return nil
}
