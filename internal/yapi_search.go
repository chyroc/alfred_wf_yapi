package internal

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/chyroc/gorequests"
)

type YapiInterface struct {
	Host      string        `json:"host"`
	EditUID   int           `json:"edit_uid"`
	Status    string        `json:"status"`
	APIOpened bool          `json:"api_opened"`
	Tag       []interface{} `json:"tag"`
	ID        int           `json:"_id"`
	Method    string        `json:"method"`
	Catid     int           `json:"catid"`
	Title     string        `json:"title"`
	Path      string        `json:"path"`
	ProjectID int           `json:"project_id"`
	UID       int           `json:"uid"`
	AddTime   int           `json:"add_time"`
	UpTime    int           `json:"up_time"`
}

type searchResp struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Data    struct {
		Count int              `json:"count"`
		Total int              `json:"total"`
		List  []*YapiInterface `json:"list"`
	} `json:"data"`
}

func (r *YapiConfig) Search(s string) ([]*YapiInterface, error) {
	s = strings.ToLower(s)
	uri := fmt.Sprintf("%s/api/interface/list?token=%s&limit=2000", r.Host, r.Token)
	resp := new(searchResp)
	if err := gorequests.New(http.MethodGet, uri).Unmarshal(resp); err != nil {
		return nil, err
	} else if resp.Errcode != 0 {
		return nil, fmt.Errorf("%d: %s", resp.Errcode, resp.Errmsg)
	}
	log.Printf("[yapi] 搜索到接口: %d, %s: %d\n", r.ID, r.Name, len(resp.Data.List))

	var result = make([]*YapiInterface, 0)
	for _, v := range resp.Data.List {
		if strings.Contains(strings.ToLower(v.Title), s) || strings.Contains(strings.ToLower(v.Path), s) {
			v.Host = r.Host
			result = append(result, v)
		}
	}

	return result, nil
}
