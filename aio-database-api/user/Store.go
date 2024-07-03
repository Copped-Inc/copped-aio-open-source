package user

import (
	"errors"
	"math/big"

	"github.com/Copped-Inc/aio-types/modules"
)

type Store string

func (s Store) Parse() (*big.Int, error) {
	if s == "" {
		return nil, nil
	}

	stores, ok := new(big.Int).SetString(string(s), 0)
	if !ok {
		return nil, errors.New("conversion from string to big int failed")
	}

	return stores, nil
}

func (s Store) IsEnabled(store modules.Site) bool {
	bint, err := s.Parse()
	if err != nil || bint == nil {
		return false
	}

	return store.Parse().Cmp((new(big.Int)).And(bint, store.Parse())) == 0
}

func StoreFromBigint(b *big.Int) Store {
	if b == new(big.Int) {
		return ""
	} else {
		return Store(b.String())
	}
}

func (s *Store) Remove(store modules.Site) error {
	currentStores, err := s.Parse()
	if err != nil {
		return err
	} else if currentStores == nil {
		return nil
	}

	if store.Parse().Cmp(currentStores) == 0 {
		*s = ""
	} else if store.Parse().Cmp((new(big.Int)).And(currentStores, store.Parse())) == 0 {
		newStores := new(big.Int)
		for _, site := range modules.Sites {
			if site == store {
				continue
			} else if site := site.Parse(); site.Cmp((new(big.Int)).And(currentStores, site)) == 0 {
				newStores.Or(newStores, site)
			}
		}

		*s = StoreFromBigint(newStores)
	}

	return nil
}

func (s *Store) Add(store modules.Site) error {
	currentStores, err := s.Parse()
	if err != nil {
		return err
	}

	if currentStores == nil {
		*s = StoreFromBigint(store.Parse())
	} else {
		*s = StoreFromBigint(new(big.Int).Or(currentStores, store.Parse()))
	}

	return nil
}
