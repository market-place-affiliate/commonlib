package shopee

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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
	debug       bool
	restyClient *resty.Client
}

func NewShopeeRepository(debug bool) ShopeeRepository {
	return &shopeeRepository{
		debug:       debug,
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
	rawQuery := `{"query":"{ productOfferV2(shopId: %s, itemId: %s) { nodes { productName itemId commissionRate commission price sales imageUrl shopName productLink offerLink periodStartTime periodEndTime priceMin priceMax productCatIds ratingStar priceDiscountRate shopId shopType sellerCommissionRate shopeeCommissionRate } pageInfo { page limit hasNextPage scrollId } } }"}`
	query := fmt.Sprintf(rawQuery, shopId, itemId)
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	factor := fmt.Sprint(cred.AppId, timestamp, query, cred.AppSecret)
	hash := sha256.Sum256([]byte(factor))
	sign := hex.EncodeToString(hash[:])

	authHeader := fmt.Sprintf(
		"SHA256 Credential=%s, Timestamp=%s, Signature=%s",
		cred.AppId,
		timestamp,
		sign,
	)

	request := s.restyClient.R()
	request.SetHeader("Authorization", authHeader)
	request.SetContentType("application/json")
	request.SetBody(query)

	var response ShopeeGetProductOfferList
	request.SetResult(&response)
	resp, err := request.Post("")
	if err != nil {
		return ShopeeGetProductOfferList{}, err
	}
	if s.debug {
		fmt.Printf("Shopee GetProductOfferListV2 Response: %+v\n", string(resp.Bytes()))
	}
	return response, nil
}

func (s *shopeeRepository) GetShortLink(cred ShopeeCredentials, originalUrl string, subid [5]string) (ShopeeGetShortLink, error) {
	sub := []string{}
	for _, v := range subid {
		if v != "" {
			sub = append(sub, v)
		}
	}
	query := fmt.Sprintf(`
mutation {
	generateShortLink(input:{originUrl:"` + originalUrl + `",subIds:["` + strings.Join(sub[:], `","`) + `"]}){
		shortLink
	}
}`, originalUrl, strings.Join(sub, `","`))

	payloadMap := map[string]string{
		"query": query,
	}

	payloadBytes, _ := json.Marshal(payloadMap)
	payload := string(payloadBytes)

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	factor := fmt.Sprint(cred.AppId, timestamp, query, cred.AppSecret)
	hash := sha256.Sum256([]byte(factor))
	sign := hex.EncodeToString(hash[:])

	authHeader := fmt.Sprintf(
		"SHA256 Credential=%s, Timestamp=%s, Signature=%s",
		cred.AppId,
		timestamp,
		sign,
	)

	request := s.restyClient.R()
	request.SetHeader("Authorization", authHeader)
	request.SetContentType("application/json")
	request.SetBody(payload)

	var response ShopeeGetShortLink
	request.SetResult(&response)
	resp, err := request.Post("")
	if err != nil {
		return ShopeeGetShortLink{}, err
	}
	if s.debug {
		fmt.Printf("Shopee GetShortLink Response: %+v\n", string(resp.Bytes()))
	}
	return response, nil
}
