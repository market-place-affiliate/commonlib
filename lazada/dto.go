package lazada

type LazadaResponse[T []ProductFeedResponse | BatchPromoteLinkResponse] struct {
	Result struct {
		Data    T    `json:"data"`
		Success bool `json:"success"`
	} `json:"result"`
	Code      string `json:"code"`
	RequestID string `json:"request_id"`
	TraceID   string `json:"_trace_id_"`
}

type ProductFeedResponse struct {
	TotalCommissionAmount float64  `json:"totalCommissionAmount"`
	SellerName            string   `json:"sellerName,omitempty"`
	DiscountPrice         float64  `json:"discountPrice"`
	CpsCommissionAmount   float64  `json:"cpsCommissionAmount"`
	Sales7D               int      `json:"sales7d"`
	Pictures              []string `json:"pictures"`
	ProductName           string   `json:"productName"`
	SellerID              int64    `json:"sellerId"`
	OutOfStock            bool     `json:"outOfStock"`
	Currency              string   `json:"currency"`
	Stock                 int      `json:"stock"`
	BrandName             string   `json:"brandName"`
	ProductID             int64    `json:"productId"`
	CategoryL1            int      `json:"categoryL1"`
	BonusOfferFlag        bool     `json:"bonusOfferFlag"`
	CpsCommissionRate     float64  `json:"cpsCommissionRate"`
	BrandID               int      `json:"brandId"`
	TotalCommissionRate   float64  `json:"totalCommissionRate"`
}

type BatchPromoteLinkResponse struct {
	URLBatchGetLinkInfoList []struct {
		RegularCommission    string `json:"regularCommission"`
		ProductID            string `json:"productId"`
		OriginalURL          string `json:"originalUrl"`
		RegularPromotionLink string `json:"regularPromotionLink"`
		ProductName          string `json:"productName"`
		Class                string `json:"class"`
	} `json:"urlBatchGetLinkInfoList"`
	ErrorInfoList []struct {
		InputValue string `json:"inputValue"`
		ErrorCode  string `json:"errorCode"`
		Class      string `json:"class"`
		ErrorMsg   string `json:"errorMsg"`
	} `json:"errorInfoList"`
	Class      string `json:"class"`
	ErrorCount int    `json:"errorCount"`
}
