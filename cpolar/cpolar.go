package cpolar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ruokeqx/cpolarPorter/env"
)

const ApiLogin = "/api/v1/user/login"
const ApiTunnels = "/api/v1/tunnels"

// Environment variables names.
const (
	envNamespace = "CPOLAR_"

	EnvUrl      = envNamespace + "URL"
	EnvUserName = envNamespace + "USERNAME"
	EnvPassWord = envNamespace + "PASSWORD"
)

type CpolarCredential struct {
	UserName string `json:"email"`
	PassWord string `json:"password"`
}

type CpolarConnector struct {
	Url   string
	token string
	CpolarCredential
}

func NewCpolarConnector() *CpolarConnector {
	return &CpolarConnector{
		Url: env.GetOrFile(EnvUrl),
		CpolarCredential: CpolarCredential{
			UserName: env.GetOrFile(EnvUserName),
			PassWord: env.GetOrFile(EnvPassWord),
		},
	}
}

func (cc *CpolarConnector) Login() error {
	credentialBytes, err := json.Marshal(cc.CpolarCredential)
	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf("%s%s", cc.Url, ApiLogin), "application/json", bytes.NewReader(credentialBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var r Response
	err = json.Unmarshal(content, &r)
	if err != nil {
		return err
	}
	cc.token = r.Data.Token
	return nil
}

func (cc *CpolarConnector) Tunnels() ([]Tunnel, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", cc.Url, ApiTunnels), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+cc.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var r Response
	err = json.Unmarshal(content, &r)
	if err != nil {
		return nil, err
	}
	var ret []Tunnel
	for _, item := range r.Data.Items {
		if item.Status == "active" {
			ret = append(ret, item.PublishTunnels...)
		}
	}
	return ret, nil
}
