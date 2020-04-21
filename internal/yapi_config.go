package internal

import (
	"encoding/json"
	"log"
	"sync"
)

type YapiConfig struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Host  string `json:"host"`
	Token string `json:"token"`
}

const yapiConfigListKey = "yapi_list.alfred_wf_yapi.com"

func (r *Workflow) AddYapiConfig(host, token string) error {
	// config
	config := &YapiConfig{
		Host:  urlWithScheme(host),
		Token: token,
	}
	project, err := config.Project()
	if err != nil {
		return err
	}
	config.ID = project.ID
	config.Name = project.Name

	// check exist
	configs, err := r.GetYapiConfigList()
	if err != nil {
		return err
	}

	for _, v := range configs {
		if v.ID == config.ID {
			return nil
		}
	}

	configs = append(configs, config)
	bs, _ := json.Marshal(configs)
	_ = r.Keychain.Set(yapiConfigListKey, string(bs))
	return nil
}

func (r *Workflow) GetYapiConfigList() ([]*YapiConfig, error) {
	s, err := r.Keychain.Get(yapiConfigListKey)
	if err != nil {
		if err.Error() == "password not found" {
			return nil, nil
		}
		return nil, err
	}
	var config = make([]*YapiConfig, 0)
	if err = json.Unmarshal([]byte(s), &config); err != nil {
		_ = r.Keychain.Delete(yapiConfigListKey)
		log.Printf("[keychain] get yapi config failed: %s", err)
		return nil, nil
	}

	return config, nil
}

func (r *Workflow) Search(s string) ([]*YapiInterface, error) {
	configs, err := r.GetYapiConfigList()
	if err != nil {
		return nil, err
	}

	var result = make([]*YapiInterface, 0)
	var wait = sync.WaitGroup{}
	var lock = sync.Mutex{}
	var finalErr error
	for _, v := range configs {
		wait.Add(1)
		go func(config YapiConfig) {
			defer wait.Done()
			list, err := config.Search(s)
			if err != nil {
				log.Printf("[search] host=%s, failed: %s", v.Host, err)
				finalErr = err
			} else {
				lock.Lock()
				result = append(result, list...)
				lock.Unlock()
			}
		}(*v)
	}
	wait.Wait()

	if finalErr != nil {
		return nil, finalErr
	}

	return result, nil
}
