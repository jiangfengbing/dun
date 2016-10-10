// Copyright (c) 2016 Jiang Fengbing
// Released under the MIT license -- see LICENSE for details

package dun

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os/exec"
	"sort"
	"time"
)

// DunChecker 易盾文本检查器
type DunChecker struct {
	secretID   string // 产品密钥ID，产品标识
	secretKey  string // 产品私有密钥，服务端生成签名信息使用，请严格保管，避免泄露
	businessID string // 业务ID，易盾根据产品业务特点分配
}

const (
	version = "v2"
	apiURL  = "https://api.aq.163.com/v2/text/check"
)

// NewChecker 创建DunChecker
func NewChecker(secretID, secretKey, businessID string) *DunChecker {
	return &DunChecker{
		secretID:   secretID,
		secretKey:  secretKey,
		businessID: businessID,
	}
}

func (d *DunChecker) genSignature(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	buf := new(bytes.Buffer)
	for _, key := range keys {
		buf.WriteString(key)
		buf.WriteString(params[key])
	}
	buf.WriteString(d.secretKey)
	hash := md5.Sum(buf.Bytes())
	return hex.EncodeToString(hash[:])
}

func (d *DunChecker) genDataID() string {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		out, _ := exec.Command("uuidgen").Output()
		return string(out)
	}
	return hex.EncodeToString(u)
}

// Check 检查文本内容
func (d *DunChecker) Check(text string) (action int, err error) {
	action = 3
	defer func() {
		if err2 := recover(); err2 != nil {
			err = fmt.Errorf("panic: %v", err2)
		}
	}()
	params := make(map[string]string)
	params["secretId"] = d.secretID
	params["businessId"] = d.businessID
	params["version"] = version
	params["timestamp"] = fmt.Sprintf("%d", time.Now().UnixNano()/(1000*1000))
	params["nonce"] = fmt.Sprintf("%d", rand.Int31())
	params["dataId"] = d.genDataID()
	params["dataOpType"] = "1"
	params["content"] = text
	params["signature"] = d.genSignature(params)
	client := &http.Client{}
	var data url.Values
	for k, v := range params {
		data.Set(k, v)
	}
	var resp *http.Response
	for i := 0; i < 3; i++ {
		body := bytes.NewBufferString(data.Encode())
		req, err := http.NewRequest("POST", apiURL, body)
		req.Header.Add("Content-type", "application/x-www-form-urlencoded")
		resp, err = client.Do(req)
		if err == nil {
			break
		}
	}
	if err != nil {
		return
	}
	var respData []byte
	if respData, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	var ret map[string]interface{}
	if err = json.Unmarshal(respData, &ret); err != nil {
		return
	}
	code := int(ret["code"].(float64))
	msg := ret["msg"].(string)
	if code != 200 {
		err = fmt.Errorf("ERROR: code=%d, msg=%s", code, msg)
		return
	}
	action = int(ret["action"].(float64))
	return
}
