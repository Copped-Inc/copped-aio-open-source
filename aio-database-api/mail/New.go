package mail

func New() *Mail {
	return &Mail{}
}

func (m Mail) SetTitle(title string) *Mail {
	m.Title = title
	return &m
}

func (m Mail) SetSubtitle(subtitle string) *Mail {
	m.Subtitle = subtitle
	return &m
}

func (m Mail) SetText(text string) *Mail {
	m.Text = text
	return &m
}

func (m Mail) SetButton(text, url string) *Mail {
	m.ButtonUrl = url
	m.ButtonText = text
	m.Button = true
	return &m
}

func (m Mail) SetBelowButton(text string) *Mail {
	m.BelowButton = text
	return &m
}
