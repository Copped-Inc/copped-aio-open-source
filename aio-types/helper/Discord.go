package helper

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/Copped-Inc/aio-types/discord"
)

func GetDiscordResp(code, redirect string) (DiscordMeResp, string, error) {

	tResp, err := GetAccessToken(code, redirect)
	if err != nil {
		return DiscordMeResp{}, "", err
	}

	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, "https://discord.com/api/v"+discord.API_Version+"/users/@me", nil)
	if err != nil {
		return DiscordMeResp{}, "", err
	}

	req.Header.Add("Authorization", "Bearer "+tResp.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return DiscordMeResp{}, "", err
	}

	var mResp DiscordMeResp
	err = json.NewDecoder(resp.Body).Decode(&mResp)
	return mResp, tResp.AccessToken, err

}

func JoinServer(id, accessToken string) error {

	client := http.Client{}
	addGuildReq := AddGuildReq{AccessToken: accessToken}

	b, err := json.Marshal(addGuildReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, "https://discord.com/api/v"+discord.API_Version+"/guilds/SERVERID/members/"+id, bytes.NewReader(b)) // Insert Server ID here
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", discord.Bearer)
	_, err = client.Do(req)
	return err

}

func GetAccessToken(code, redirect string) (DiscordTokenResp, error) {

	client := http.Client{}

	form := url.Values{}
	form.Add("code", code)
	form.Add("grant_type", "authorization_code")
	form.Add("redirect_uri", redirect+"/")
	form.Add("scope", "identify%20email%20guilds.join")

	req, err := http.NewRequest(http.MethodPost, "https://discord.com/api/v"+discord.API_Version+"/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return DiscordTokenResp{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(discord.Application_ID, discord.Oauth2_Secret)

	resp, err := client.Do(req)
	if err != nil {
		return DiscordTokenResp{}, err
	}

	var tResp DiscordTokenResp
	err = json.NewDecoder(resp.Body).Decode(&tResp)
	return tResp, err

}

type DiscordTokenResp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type DiscordMeResp struct {
	Id            string      `json:"id"`
	Username      string      `json:"username"`
	Avatar        string      `json:"avatar"`
	Discriminator string      `json:"discriminator"`
	PublicFlags   int         `json:"public_flags"`
	Flags         int         `json:"flags"`
	Banner        string      `json:"banner"`
	BannerColor   interface{} `json:"banner_color"`
	AccentColor   interface{} `json:"accent_color"`
	Locale        string      `json:"locale"`
	MfaEnabled    bool        `json:"mfa_enabled"`
	PremiumType   int         `json:"premium_type"`
	Email         string      `json:"email"`
	Verified      bool        `json:"verified"`
}

type AddGuildReq struct {
	AccessToken string `json:"access_token"`
}
