package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/chyroc/alfred_wf_yapi/internal"
)

func run() error {
	scope, arg := internal.Default.InputWithPre()
	switch scope {
	case "config":
		if arg == "" {
			configs, err := internal.Default.GetYapiConfigList()
			if err != nil {
				return err
			}
			internal.Default.NewItem("[提示] 添加配置: yapi.config <host> <token>")
			if len(configs) == 0 {
				internal.Default.NewItem("[警告] 尚未配置任何配置")
				internal.Default.SendFeedback()
				return nil
			}

			for i, v := range configs {
				internal.Default.NewItem(fmt.Sprintf("[%d] %s", i+1, v.Name)).Subtitle(v.Host)
			}
			if len(configs) == 0 {
				internal.Default.NewItem("[警告] 尚未配置任何配置")
			}
			internal.Default.SendFeedback()
			return nil
		}

		args := strings.Split(arg, " ")
		if len(args) != 2 {
			internal.Default.Warn("配置格式 yapi.config <host> <token>", "")
			return nil
		}
		host := strings.TrimSpace(args[0])
		uri, err := url.Parse(host)
		if err != nil || uri == nil || uri.Host == "" {
			internal.Default.Warn(fmt.Sprintf("%s 不是合法 host，请重试(注意添加 http 头)", arg), "")
			return nil
		}
		host = fmt.Sprintf("%s://%s", uri.Scheme, uri.Host)

		token := strings.TrimSpace(args[1])
		if token == "" {
			internal.Default.Warn("token 为空", "")
			return nil
		}

		configBytes, _ := json.Marshal(internal.YapiConfig{
			Host:  host,
			Token: token,
		})
		internal.Default.NewItem(fmt.Sprintf("回车以录入 %s", host)).Valid(true).Arg(string(configBytes))
		internal.Default.SendFeedback()
	case "config.select":
		config := new(internal.YapiConfig)
		if err := json.Unmarshal([]byte(arg), config); err != nil {
			return err
		}
		return internal.Default.AddYapiConfig(config.Host, config.Token)
	case "search":
		result, err := internal.Default.Search(arg)
		if err != nil {
			return err
		}
		if len(result) == 0 {
			internal.Default.Warn(fmt.Sprintf("%s 没有搜索到接口", arg), "")
			return nil
		}

		for _, v := range result {
			internal.Default.NewItem(v.Title).Subtitle(v.Path).Valid(true).Arg(fmt.Sprintf("%s/project/%d/interface/api/%d", v.Host, v.ProjectID, v.ID))
		}

		internal.Default.SendFeedback()
	case "search.select":
		internal.Open(arg)
	}

	return nil
}

func main() {
	internal.Init()

	log.Printf("run\n")

	internal.Default.Run(run)
}
