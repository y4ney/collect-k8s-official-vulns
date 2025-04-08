package cmd

/*
Copyright © 2025 Yaney yangli.yaney@foxmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/y4ney/collect-k8s-official-vulns/pkg"
	"golang.org/x/term"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	fetch     bool
	cacheDir  string
	update    bool
	translate bool
	apikey    string
)
var rootCmd = &cobra.Command{
	Use:   "collect-k8s-official-vulns",
	Short: "K8S 官方漏洞",
	Long: `
 ██ ▄█▀  ██████  ▒█████    █████▒ █████▒██▓ ▄████▄   ██▓ ▄▄▄       ██▓  ██▒   █▓ █    ██  ██▓     ███▄    █ 
 ██▄█▒ ▒██    ▒ ▒██▒  ██▒▓██   ▒▓██   ▒▓██▒▒██▀ ▀█  ▓██▒▒████▄    ▓██▒ ▓██░   █▒ ██  ▓██▒▓██▒     ██ ▀█   █ 
▓███▄░ ░ ▓██▄   ▒██░  ██▒▒████ ░▒████ ░▒██▒▒▓█    ▄ ▒██▒▒██  ▀█▄  ▒██░  ▓██  █▒░▓██  ▒██░▒██░    ▓██  ▀█ ██▒
▓██ █▄   ▒   ██▒▒██   ██░░▓█▒  ░░▓█▒  ░░██░▒▓▓▄ ▄██▒░██░░██▄▄▄▄██ ▒██░   ▒██ █░░▓▓█  ░██░▒██░    ▓██▒  ▐▌██▒
▒██▒ █▄▒██████▒▒░ ████▓▒░░▒█░   ░▒█░   ░██░▒ ▓███▀ ░░██░ ▓█   ▓██▒░██████▒▒▀█░  ▒▒█████▓ ░██████▒▒██░   ▓██░
▒ ▒▒ ▓▒▒ ▒▓▒ ▒ ░░ ▒░▒░▒░  ▒ ░    ▒ ░   ░▓  ░ ░▒ ▒  ░░▓   ▒▒   ▓▒█░░ ▒░▓  ░░ ▐░  ░▒▓▒ ▒ ▒ ░ ▒░▓  ░░ ▒░   ▒ ▒ 
░ ░▒ ▒░░ ░▒  ░ ░  ░ ▒ ▒░  ░      ░      ▒ ░  ░  ▒    ▒ ░  ▒   ▒▒ ░░ ░ ▒  ░░ ░░  ░░▒░ ░ ░ ░ ░ ▒  ░░ ░░   ░ ▒░
░ ░░ ░ ░  ░  ░  ░ ░ ░ ▒   ░ ░    ░ ░    ▒ ░░         ▒ ░  ░   ▒     ░ ░     ░░   ░░░ ░ ░   ░ ░      ░   ░ ░ 
░  ░         ░      ░ ░                 ░  ░ ░       ░        ░  ░    ░  ░   ░     ░         ░  ░         ░ 
                                           ░                                ░                               
⎈ 可爬取、更新并将 K8S 官方漏洞翻译（通过 DeepSeek）成简体中文，欢迎使用～ ⎈`,
	Run: run,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	InitLogger()
	rootCmd.Flags().BoolVarP(&fetch, "fetch", "f", false, "爬取 K8S 官方漏洞")
	rootCmd.Flags().BoolVarP(&update, "update", "u", false, "更新 K8S 官方漏洞，建议额外指定 -t 以翻译更新的内容")
	rootCmd.Flags().BoolVarP(&translate, "translate", "t", false, "翻译 K8S 官方漏洞")
	rootCmd.Flags().StringVarP(&cacheDir, "cache-dir", "c", CacheDir(), "指定缓存目录")
	rootCmd.Flags().StringVarP(&apikey, "apikey", "a", "", "指定 DeepSeek 的 API Key，可设置为环境变量 API_KEY")
}

func InitLogger() {
	var (
		defaultLogger = zerolog.New(os.Stderr)
		logLevel      = zerolog.TraceLevel
	)

	zerolog.SetGlobalLevel(logLevel)
	if term.IsTerminal(int(os.Stdout.Fd())) {
		defaultLogger = zerolog.New(zerolog.NewConsoleWriter())
	}
	log.Logger = defaultLogger.With().Timestamp().Stack().Logger()
}

func run(_ *cobra.Command, _ []string) {
	if fetch {
		fetchK8sVuln()
	}
	if translate && !update {
		translateK8sVulns()
	}
	if update {
		updateK8sVuln()
	}
}

func fetchK8sVuln() {
	data, err := pkg.Fetch()
	if err != nil {
		log.Fatal().Err(err).Msg("无法爬取 K8S 官方漏洞")
	}
	for _, item := range data.Items {
		if item.ContentText == "" {
			log.Debug().Str("CVE 编号", item.ID).Msg("漏洞信息为空，已跳过")
			continue
		}
		dir := filepath.Join(cacheDir, item.DatePublished.Format("2006"),
			item.DatePublished.Format("01"))
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			log.Fatal().Err(err).Str("目录", dir).Msg("无法创建目录")
		}
		err = item.Save(dir)
		if err != nil {
			log.Fatal().Err(err).Str("目录", dir).Str("CVE 编号", item.ID).Msg("无法保存漏洞")
		}
		log.Info().Str("CVE 编号", item.ID).Msg("已成功保存漏洞")
	}
	log.Info().Str("缓存目录", cacheDir).Int("漏洞总数", len(data.Items)).Msg("已成功保存所有的 K8S 官方漏洞")
}

func translateK8sVulns() {
	checkApiKey()
	err := filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal().Err(err).Str("路径", path).Msg("无法遍历路径")
		}
		if !info.IsDir() && filepath.Ext(path) == ".md" && !strings.HasSuffix(path, "_zh.md") {
			translateK8sVuln(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Str("缓存目录", cacheDir).Msg("无法遍历目录")
	}
	log.Info().Msg("已成功翻译并保存好所有漏洞信息")
}

func updateK8sVuln() {
	// 获取最新的漏洞信息
	data, err := pkg.Fetch()
	if err != nil {
		log.Fatal().Err(err).Msg("无法爬取 K8S 官方漏洞")
	}

	// 若最新时间早于当前，则证明当前的数据已为最新数据，暂不更新
	if data.KubernetesIo.UpdatedAt.Before(time.Now()) {
		log.Info().Msg("当前漏洞数据已为最新数据，暂不更新")
		return
	}

	var (
		oldVulns    = make(map[string]string) // 旧的漏洞信息
		updateVulns []string                  // 更新的漏洞信息
	)

	// 提取历史的漏洞信息
	err = filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal().Err(err).Str("路径", path).Msg("无法遍历路径")
		}
		if !info.IsDir() && filepath.Ext(path) == ".md" && !strings.HasSuffix(path, "_zh.md") {
			fileContent, err := os.ReadFile(path)
			if err != nil {
				log.Fatal().Err(err).Str("文件", path).Msg("无法读取文件内容")
			}
			oldVulns[path] = string(fileContent)
		}
		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Str("缓存目录", cacheDir).Msg("无法遍历目录")
	}

	// 对比并更新漏洞信息
	for _, newVuln := range data.Items {
		path := filepath.Join(cacheDir, newVuln.DatePublished.Format("2006"),
			newVuln.DatePublished.Format("01"), newVuln.ID+".md")
		// 若文件内容相同，则 continue
		if oldVulns[path] == newVuln.ContentText {
			continue
		}

		_, exists := oldVulns[path]
		// 若新增漏洞，则可能需要创建目录
		if !exists {
			err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
			if err != nil && !os.IsExist(err) {
				log.Fatal().Err(err).Str("目录", filepath.Dir(path)).Msg("无法创建目录")
			}
			log.Info().Str("CVE 编号", newVuln.ID).Msg("新增 1 个漏洞")
		}
		err = os.WriteFile(path, []byte(newVuln.ContentText), 0644)
		if err != nil {
			log.Fatal().Err(err).Str("文件", path).Msg("无法更新漏洞")
		}
		updateVulns = append(updateVulns, path)
		log.Info().Str("文件", path).Msg("已成功保存更新后的漏洞信息")
	}

	// 若需要，则翻译更新后的漏洞信息
	if translate {
		checkApiKey()
		for _, path := range updateVulns {
			translateK8sVuln(path)
		}
	}
}

func CacheDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir = os.TempDir()
	}
	return filepath.Join(dir, "k8s-official-vulns")
}

// checkApiKey 若未指定 apikey，则读取环境变量 API_KEY
func checkApiKey() {
	if apikey == "" {
		log.Info().Msg("未指定 api key，尝试从环境变量中获取。")
		apikey = os.Getenv("API_KEY")
		if apikey == "" {
			log.Fatal().Msg("api key 未空")
		}
		log.Info().Msg("已从环境变量中读取到 api key")
	}
}

// translateK8sVuln 翻译单个 k8s 漏洞
func translateK8sVuln(path string) {
	// 读取文件
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal().Err(err).Str("文件", path).Msg("无法读取文件内容")
	}

	// 翻译文件
	res, err := pkg.Request(string(data), apikey)
	if err != nil {
		log.Fatal().Err(err).Str("文件", path).Msg("无法翻译漏洞信息")
	}
	log.Info().Str("文件", path).Msg("已成功翻译漏洞信息")

	// 写入翻译内容
	newPath := strings.Replace(path, ".md", "_zh.md", -1)
	err = os.WriteFile(newPath, []byte(res.Choices[0].Message.Content), 0644)
	if err != nil {
		log.Fatal().Err(err).Str("新文件", newPath).Msg("无法写入新文件")
	}
	log.Info().Str("新文件", newPath).Msg("已成功保存翻译后的漏洞信息")
}
