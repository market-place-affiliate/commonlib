package shopee

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"resty.dev/v3"
)

type ShopeeRepository interface {
	GetProductOfferListV2(shopId, itemId string) (ShopeeGetProductOfferList, error)
}

type shopeeRepository struct {
	appId       string
	appSecret   string
	restyClient *resty.Client
}

func NewShopeeRepository(appId, appSecret string) ShopeeRepository {
	return &shopeeRepository{
		appId:       appId,
		appSecret:   appSecret,
		restyClient: resty.New().SetBaseURL("https://open-api.affiliate.shopee.co.th/graphql"),
	}
}

func (s *shopeeRepository) GetProductOfferListV2(shopId, itemId string) (ShopeeGetProductOfferList, error) {
	rawQuery := `
	{
		productOfferV2(shopId: %s,itemId: %s) {
			nodes {
				productName
				itemId
				commissionRate
				commission
				price
				sales
				imageUrl
				shopName
				productLink
				offerLink
				periodStartTime
				periodEndTime
				priceMin
				priceMax
				productCatIds
				ratingStar
				priceDiscountRate
				shopId
				shopType
				sellerCommissionRate
				shopeeCommissionRate
			}
			pageInfo {
				page
				limit
				hasNextPage
				scrollId
			}
		}
	}`
	query := fmt.Sprintf(rawQuery, shopId, itemId)

	factor := fmt.Sprint(s.appId, fmt.Sprintf("%d", time.Now().Unix()), query, s.appSecret)
	var message bytes.Buffer
	message.WriteString(factor)
	hash := hmac.New(sha256.New, []byte(s.appSecret))
	hash.Write(message.Bytes())
	sign := strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
	request := s.restyClient.R()
	request.SetHeader("Authorization", sign)
	request.SetHeader("Credential", s.appId)
	request.SetHeader("Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	request.SetBody(map[string]string{
		"query": query,
	})

	var response ShopeeGetProductOfferList
	request.SetResult(&response)
	_, err := request.Post("")
	if err != nil {
		return ShopeeGetProductOfferList{}, err
	}

	return response, nil
}
