package cloopen

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"
)

type SMS struct {
	c *Client
}

func (c *Client) SMS() *SMS {
	return &SMS{c}
}

type SendRequest struct {
	To         string   `json:"to" xml:"to"`
	TemplateId string   `json:"templateId" xml:"templateId"`
	Template string     `json:"template" xml:"template"`
	Datas      map[string]string `json:"datas" xml:"datas>key>value"`
	international bool
}

type responseData struct {
	SmsMessageSid string `xml:"smsMessageSid"`
	DateCreated   string `xml:"dateCreated"`
}

type SendResponse struct {
	StatusCode  string `xml:"statusCode"`
	StatusMsg   string `xml:"statusMsg"`
	TemplateSMS responseData
}

func (sms *SMS) Send(input *SendRequest) (*SendResponse, error) {
	if input == nil {
		input = &SendRequest{}
	}
	input.international = false
	tos := strings.Split(input.To, ",")
	if strings.HasPrefix(tos[0], "00") {
		to := string(tos[0][2:])
		if !strings.HasPrefix(to, "86") {
			input.international = true
		}
	}
// 	input.international = true
	err := input.Verify()
	if err != nil {
		return nil, err
	}

	uri := strings.Join([]string{"/2013-12-26/Accounts/", sms.c.config.APIAccount, "/SMS/TemplateSMS"}, "")

	if input.international {
		uri = strings.Join([]string{"/v2/account/", sms.c.config.APIAccount, "/international/send"}, "")
	}

	r := sms.c.newRequest(HTTP_POST, sms.c.config.SmsHost, uri)
	ct := getHeaderContentType(sms.c.config.ContentType)
	r.header.Set(HEADER_CONTENT_TYPE, ct)
	r.header.Set(HEADER_ACCEPT, ct)

	auth, sig := buildSign(sms.c.config.APIAccount, sms.c.config.APIToken)
	r.header.Set(HEADER_AUTH, auth)
	r.params.Set(URL_PARAM_SIG, sig)


	sms.buildBody(r, input)

	resp, err := sms.c.doRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SendResponse
	if err = sms.c.handleResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (req *SendRequest) Verify() error {
	if len(req.To) == 0 {
		return errors.New("Miss param: to")
	}
	if !req.international && len(req.TemplateId) == 0 {
		return errors.New("Miss param:templateId")
	}
	return nil
}

func buildSign(account, token string) (auth, sig string) {
	timeStr := time.Now().Format("20060102150405")
	sigValue := Md5(strings.Join([]string{account, token, timeStr}, ""))
	authValue := Base64URL(strings.Join([]string{account, timeStr}, ":"))
	return authValue, sigValue
}

func getHeaderContentType(contentType string) string {
	if contentType == CONTENT_JSON {
		return HEADER_CONTENT_JSON
	} else {
		return HEADER_CONTENT_XML
	}
}

func (sms *SMS)buildBody(request *request, input *SendRequest) {
	buf := new(bytes.Buffer)
	var data map[string]interface{}
	if input.international {
		template := input.Template
		for key, value := range input.Datas {
			template = strings.Replace(template, fmt.Sprintf("{{%s}}", key), value, 1)
		}
		data = map[string]interface{}{
			"appId" : sms.c.config.AppId,
			"content": template,
			"mobile": input.To,
		}
	} else {
		var arguments []string
		for _, value := range input.Datas {
			arguments = append(arguments, value)
		}
		data = map[string]interface{}{
			"appId" : sms.c.config.AppId,
			"templateId": input.TemplateId,
			"datas": arguments,
			"to": input.To,
		}
	}
	if sms.c.config.ContentType == CONTENT_JSON {
		_ = json.NewEncoder(buf).Encode(data)
	} else {
		_ = xml.NewEncoder(buf).Encode(data)
	}
	request.body = buf
}
