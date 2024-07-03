package subscriptions

type (
	Plan int

	planData struct {
		Name  string
		Price string
	}
)

const (
	None Plan = iota
	Friends_and_Family
	Basic
	Developer
)

var (
	plans = map[Plan]planData{
		Friends_and_Family: {
			Price: "price_1MtEnmETdQPiNAdkf8F4b3WZ",
			Name:  "F&F",
		},
		Basic: {
			Price: "price_1NLmKVETdQPiNAdkwLmi7lk0",
			Name:  "Basic",
		},
		Developer: {
			Name: "Developer",
		},
	}

	Plans = []Plan{Friends_and_Family, Basic, Developer}
)

func (p Plan) GetData() planData {
	return plans[p]
}
