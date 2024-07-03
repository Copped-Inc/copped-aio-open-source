package product

type Product struct {
	Name      string      `json:"name" firestore:"name"`
	SKU       string      `json:"sku" firestore:"sku"`
	StockX    string      `json:"stockx,omitempty" firestore:"stockx,omitempty"`
	Image     string      `json:"image" firestore:"image"`
	State     State       `json:"state" firestore:"state"`
	Handles   []string    `json:"handles" firestore:"handles"`
	Price     float64     `json:"price" firestore:"price"`
	UserState []UserState `json:"user_state,omitempty" firestore:"user_state,omitempty"`
}

type UserState struct {
	ID    string `json:"id" firestore:"id"`
	SKU   string `json:"sku" firestore:"sku"`
	State State  `json:"state" firestore:"state"`
}

type State int

const (
	None State = iota - 1
	Whitelisted
	Blacklisted
)
