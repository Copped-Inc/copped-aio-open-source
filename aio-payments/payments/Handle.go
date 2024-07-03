package payments

import "strings"

func (p *Payments) Handle(auth string) {

	p.auth = auth
	println("Handling payment", p.Id)

	p.State = Accepted
	err := p.post()
	if err != nil {
		println(err.Error())
		return
	}

	p.data = make(map[string]string)
	split := strings.Split(p.Data, ";")
	for _, s := range split {
		if s == "" {
			continue
		}

		split := strings.Split(s, ":")
		p.data[split[0]] = split[1]
	}

	switch p.Store {
	case "aboutyou":
		err = p.aboutyou()
		break
	}

	if err != nil {
		println(err.Error())
		p.State = Declined

		err = p.post()
		if err != nil {
			println(err.Error())
		}

		return
	}

	p.State = Finalized
	err = p.post()

}
