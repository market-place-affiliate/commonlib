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
	GetProductOfferListV2(cred ShopeeCredentials, shopId, itemId string) (ShopeeGetProductOfferList, error)
	GetShortLink(cred ShopeeCredentials, originalUrl string, sub [5]string) (ShopeeGetShortLink, error)
}

type ShopeeCredentials struct {
	AppId     string
	AppSecret string
}

type shopeeRepository struct {
	restyClient *resty.Client
}

func NewShopeeRepository() ShopeeRepository {
	return &shopeeRepository{
		restyClient: resty.New().SetBaseURL("https://open-api.affiliate.shopee.co.th/graphql"),
	}
}

func ExtractShopIdAndItemIdFromLink(link string) (string, string, error) {
	parts := strings.Split(link, "-i.")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid link format")
	}

	idParts := strings.Split(parts[1], ".")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("invalid link format")
	}

	shopId := idParts[0]
	itemId := strings.Split(idParts[1], "?")[0]

	return shopId, itemId, nil
}

func (s *shopeeRepository) GetProductOfferListV2(cred ShopeeCredentials, shopId, itemId string) (ShopeeGetProductOfferList, error) {
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

	factor := fmt.Sprint(cred.AppId, fmt.Sprintf("%d", time.Now().Unix()), query, cred.AppSecret)
	var message bytes.Buffer
	message.WriteString(factor)
	hash := hmac.New(sha256.New, []byte(cred.AppSecret))
	hash.Write(message.Bytes())
	sign := strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
	request := s.restyClient.R()
	request.SetHeader("Authorization", sign)
	request.SetHeader("Credential", cred.AppId)
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

func (s *shopeeRepository) GetShortLink(cred ShopeeCredentials, originalUrl string, sub [5]string) (ShopeeGetShortLink, error) {
	query := `
	mutation {
		generateShortLink(input:{originUrl:"` + originalUrl + `",subIds:["` + strings.Join(sub[:], `","`) + `"]}){
			shortLink
		}
	}
	`
	factor := fmt.Sprint(cred.AppId, fmt.Sprintf("%d", time.Now().Unix()), query, cred.AppSecret)
	var message bytes.Buffer
	message.WriteString(factor)
	hash := hmac.New(sha256.New, []byte(cred.AppSecret))
	hash.Write(message.Bytes())
	sign := strings.ToUpper(hex.EncodeToString(hash.Sum(nil)))
	request := s.restyClient.R()
	request.SetHeader("Authorization", sign)
	request.SetHeader("Credential", cred.AppId)
	request.SetHeader("Timestamp", fmt.Sprintf("%d", time.Now().Unix()))
	request.SetBody(map[string]string{
		"query": query,
	})

	var response ShopeeGetShortLink
	request.SetResult(&response)
	_, err := request.Post("")
	if err != nil {
		return ShopeeGetShortLink{}, err
	}

	return response, nil
}
