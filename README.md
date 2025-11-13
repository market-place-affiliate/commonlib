# Golang unofficial sdk affiliate Lazada, Shopee 

# Lazada
## Get credential Lazada from https://adsense.lazada.co.th
```
  appKey := "123456"
	appSecret := "xxx"
	userToken := "xxx"

	lazadaRepo := lazada.NewLazadaRepository(lazada.ApiGatewayTH, appKey, appSecret, "sha256", userToken)

  feedProduct, err := lazadaRepo.GetProductFeed(1,10)
	if err != nil {
	  log.Fatal(err)
	}
	log.Println(feedProduct)
```

# Shopee
## Get credential Shopee from https://affiliate.shopee.co.th
```
  appId := "xxx"
  appSecret := "xxx"

  shopeeRepo := shopee.NewShopeeRepository(appId,appSecret)


  productUrl := https://shopee.co.th/example-i.70998059.24612422412
  shopId,itemId := shopee.ExtractShopIdAndItemIdFromLink(productUrl) // 70998059,24612422412
  
  offerResp,err := shopeeRepo.GetProductOfferListV2(shopId,itemId)
  if err != nil {
	  log.Fatal(err)
	}
	log.Println(offerResp)
```