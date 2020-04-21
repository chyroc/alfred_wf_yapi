package internal

import (
	"fmt"
	"net/http"

	"github.com/chyroc/gorequests"
)

type YapiProject struct {
	SwitchNotice bool   `json:"switch_notice"`
	IsMockOpen   bool   `json:"is_mock_open"`
	Strice       bool   `json:"strice"`
	IsJSON5      bool   `json:"is_json5"`
	ID           int    `json:"_id"`
	Name         string `json:"name"`
	Basepath     string `json:"basepath"`
	ProjectType  string `json:"project_type"`
	UID          int    `json:"uid"`
	GroupID      int    `json:"group_id"`
	Icon         string `json:"icon"`
	Color        string `json:"color"`
	Psm          string `json:"psm"`
	AddTime      int    `json:"add_time"`
	UpTime       int    `json:"up_time"`
	Env          []struct {
		Header []interface{} `json:"header"`
		Global []interface{} `json:"global"`
		ID     string        `json:"_id"`
		Name   string        `json:"name"`
		Domain string        `json:"domain"`
	} `json:"env"`
	Tag  []interface{} `json:"tag"`
	Cat  []interface{} `json:"cat"`
	Role bool          `json:"role"`
}

type projectResp struct {
	Errcode int          `json:"errcode"`
	Errmsg  string       `json:"errmsg"`
	Data    *YapiProject `json:"data"`
}

func (r *YapiConfig) Project() (*YapiProject, error) {
	uri := fmt.Sprintf("%s/api/project/get?token=%s", r.Host, r.Token)
	resp := new(projectResp)
	if err := gorequests.New(http.MethodGet, uri).Unmarshal(resp); err != nil {
		return nil, err
	} else if resp.Errcode != 0 {
		return nil, fmt.Errorf("%d: %s", resp.Errcode, resp.Errmsg)
	}

	return resp.Data, nil
}
