package shopee

type ShopeeGetProductOfferList struct {
	Data struct {
		ProductOfferV2 struct {
			Nodes []struct {
				ProductName          string `json:"productName"`
				ItemID               int64  `json:"itemId"`
				CommissionRate       string `json:"commissionRate"`
				Commission           string `json:"commission"`
				Price                string `json:"price"`
				Sales                int    `json:"sales"`
				ImageURL             string `json:"imageUrl"`
				ShopName             string `json:"shopName"`
				ProductLink          string `json:"productLink"`
				OfferLink            string `json:"offerLink"`
				PeriodStartTime      int    `json:"periodStartTime"`
				PeriodEndTime        int64  `json:"periodEndTime"`
				PriceMin             string `json:"priceMin"`
				PriceMax             string `json:"priceMax"`
				ProductCatIds        []int  `json:"productCatIds"`
				RatingStar           string `json:"ratingStar"`
				PriceDiscountRate    int    `json:"priceDiscountRate"`
				ShopID               int    `json:"shopId"`
				ShopType             []int  `json:"shopType"`
				SellerCommissionRate string `json:"sellerCommissionRate"`
				ShopeeCommissionRate string `json:"shopeeCommissionRate"`
			} `json:"nodes"`
			PageInfo struct {
				Page        int         `json:"page"`
				Limit       int         `json:"limit"`
				HasNextPage bool        `json:"hasNextPage"`
				ScrollID    interface{} `json:"scrollId"`
			} `json:"pageInfo"`
		} `json:"productOfferV2"`
	} `json:"data"`
}

type ShopeeGetShortLink struct {
	Data struct {
		GenerateShortLink struct {
			ShortLink string `json:"shortLink"`
		} `json:"generateShortLink"`
	} `json:"data"`
	Errors []struct {
		Message    string `json:"message"`
		Extensions struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"extensions"`
	} `json:"errors"`
}
