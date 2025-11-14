package lazada

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"resty.dev/v3"
)

type LazadaApiGateway string

const (
	Version = "lazada-sdk-go-20251113"

	ApiGatewaySG LazadaApiGateway = "https://api.lazada.sg/rest"
	ApiGatewayMY LazadaApiGateway = "https://api.lazada.com.my/rest"
	ApiGatewayVN LazadaApiGateway = "https://api.lazada.vn/rest"
	ApiGatewayTH LazadaApiGateway = "https://api.lazada.co.th/rest"
	ApiGatewayPH LazadaApiGateway = "https://api.lazada.com.ph/rest"
	ApiGatewayID LazadaApiGateway = "https://api.lazada.co.id/rest"
)

type LazadaRepository interface {
	GetProductFeed(cred LazadaCredentials, productId string, page, limit int) (LazadaResponse[[]ProductFeedResponse], error)
	GetBatchPromoteLink(cred LazadaCredentials, inputType, inputValue string, sub [6]string) (LazadaResponse[BatchPromoteLinkResponse], error)
}

type LazadaCredentials struct {
	AppKey     string
	AppSecret  string
	SignMethod string
	UserToken  string
}

type lazadaRepository struct {
	apiGateway  LazadaApiGateway
	restyClient *resty.Client
}

func NewLazadaRepository(apiGateway LazadaApiGateway) LazadaRepository {
	return &lazadaRepository{
		apiGateway:  apiGateway,
		restyClient: resty.New(),
	}
}

func (l *lazadaRepository) GetProductFeed(cred LazadaCredentials, productId string, page, limit int) (LazadaResponse[[]ProductFeedResponse], error) {
	request := l.restyClient.R()
	apiPath := "/marketing/product/feed"
	sysParams := map[string]string{
		"app_key":     cred.AppKey,
		"sign_method": "sha256",
		"timestamp":   fmt.Sprintf("%d000", time.Now().Unix()),
	}
	apiParams := map[string]string{
		"offerType": "1",
		"userToken": cred.UserToken,
		"productId": productId,
		"page":      fmt.Sprintf("%d", page),
		"limit":     fmt.Sprintf("%d", limit),
	}

	keys := []string{}
	union := map[string]string{}
	for key, val := range sysParams {
		union[key] = val
		keys = append(keys, key)
	}
	for key, val := range apiParams {
		union[key] = val
		keys = append(keys, key)
	}

	// sort sys params and api params by key
	sort.Strings(keys)

	var message bytes.Buffer
	message.WriteString(apiPath)
	for _, key := range keys {
		message.WriteString(fmt.Sprintf("%s%s", key, union[key]))
	}

	hash := hmac.New(sha256.New, []byte(cred.AppSecret))
	hash.Write(message.Bytes())

	sign := strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
	request.SetQueryParams(union)
	request.SetQueryParam("sign", sign)

	var lazadaResp LazadaResponse[[]ProductFeedResponse]
	request.SetResult(&lazadaResp)
	_, err := request.Get(string(l.apiGateway) + apiPath)
	if err != nil {
		return LazadaResponse[[]ProductFeedResponse]{}, err
	}

	return lazadaResp, err
}

func (l *lazadaRepository) GetBatchPromoteLink(cred LazadaCredentials, inputType, inputValue string, sub [6]string) (LazadaResponse[BatchPromoteLinkResponse], error) {
	request := l.restyClient.R()
	apiPath := "/marketing/getlink"
	sysParams := map[string]string{
		"app_key":     cred.AppKey,
		"sign_method": "sha256",
		"timestamp":   fmt.Sprintf("%d000", time.Now().Unix()),
	}
	apiParams := map[string]string{
		"userToken":  cred.UserToken,
		"inputType":  inputType,
		"inputValue": inputValue,
	}
	for i, v := range sub {
		if v == "" {
			continue
		}
		apiParams[fmt.Sprintf("sub%d", i+1)] = v
	}

	keys := []string{}
	union := map[string]string{}
	for key, val := range sysParams {
		union[key] = val
		keys = append(keys, key)
	}
	for key, val := range apiParams {
		union[key] = val
		keys = append(keys, key)
	}

	// sort sys params and api params by key
	sort.Strings(keys)

	var message bytes.Buffer
	message.WriteString(apiPath)
	for _, key := range keys {
		message.WriteString(fmt.Sprintf("%s%s", key, union[key]))
	}

	hash := hmac.New(sha256.New, []byte(cred.AppSecret))
	hash.Write(message.Bytes())

	sign := strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
	request.SetQueryParams(union)
	request.SetQueryParam("sign", sign)

	var lazadaResp LazadaResponse[BatchPromoteLinkResponse]
	request.SetResult(&lazadaResp)
	_, err := request.Get(string(l.apiGateway) + apiPath)
	if err != nil {
		return LazadaResponse[BatchPromoteLinkResponse]{}, err
	}
	// log.Println(string(resp.Bytes()))
	return lazadaResp, err
}
