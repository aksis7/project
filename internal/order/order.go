package order

type Order struct {
	OrderUID          string   `json:"order_uid" gorm:"primaryKey"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery" gorm:"foreignKey:OrderUID;references:OrderUID"`
	Payment           Payment  `json:"payment" gorm:"foreignKey:OrderUID;references:OrderUID"`
	Items             []Item   `json:"items" gorm:"foreignKey:OrderUID;references:OrderUID"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerID        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	ShardKey          string   `json:"shardkey"`
	SmID              int      `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

// Указываем явное имя таблицы для модели `Order`
func (Order) TableName() string {
	return "orders"
}

type Delivery struct {
	OrderUID string `json:"order_uid" gorm:"primaryKey"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Zip      string `json:"zip"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Region   string `json:"region"`
	Email    string `json:"email"`
}

// Указываем явное имя таблицы для модели `Delivery`
func (Delivery) TableName() string {
	return "delivery"
}

type Payment struct {
	OrderUID     string `json:"order_uid" gorm:"primaryKey"`
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

// Указываем явное имя таблицы для модели `Payment`
func (Payment) TableName() string {
	return "payment"
}

type Item struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
	OrderUID    string `json:"order_uid"`
}

// Указываем явное имя таблицы для модели `Item`
func (Item) TableName() string {
	return "items"
}
