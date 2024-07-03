package preharvest

import "github.com/Copped-Inc/aio-types/captcha/preharvest"

type Tasks []preharvest.Task

func (t Tasks) Len() int {
	return len(t)
}
func (s Tasks) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Tasks) Less(i, j int) bool {
	return s[i].Date.After(s[j].Date)
}
