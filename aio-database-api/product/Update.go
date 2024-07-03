package product

import "strings"

func (p *Product) UpdateName(name string) *Product {
	p.Name = name
	return p
}

func (p *Product) UpdateStockX(stockx string) *Product {
	p.StockX = stockx
	return p
}

func (p *Product) UpdateState(state State) *Product {
	p.State = state
	return p
}

func (p *Product) UpdateImage(image string) *Product {
	p.Image = image
	return p
}

func (p *Product) UpdateHandles(handles []string) *Product {
	p.Handles = handles
	return p
}

func (p *Product) AddHandle(handle string) *Product {
	p.Handles = append(p.Handles, strings.ToLower(handle))
	return p
}

func (p *Product) UpdatePrice(price float64) *Product {
	p.Price = price
	return p
}

func (p *Product) UpdateUserState(userState []UserState) *Product {
	p.UserState = userState
	return p
}

func (p *Product) AddUserState(userState UserState) *Product {
	for _, state := range p.UserState {
		if state.ID == userState.ID {
			state.State = userState.State
			return p
		}
	}

	p.UserState = append(p.UserState, userState)
	return p
}
