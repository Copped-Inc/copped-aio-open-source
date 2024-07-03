package linked_roles

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/Copped-Inc/aio-types/console"

	"github.com/Copped-Inc/aio-types/discord"
	"github.com/Copped-Inc/aio-types/subscriptions"
)

func RegisterMetadata() {

	var (
		expected_metadata, current_metadata []metadata
		client                              = &http.Client{}
	)

	req, err := http.NewRequest(http.MethodGet, "https://discord.com/api/v"+discord.API_Version+"/applications/"+discord.Application_ID+"/role-connections/metadata", nil)
	if err != nil {
		console.Log(err)
		return
	}

	req.Header.Set("Content-type", "application/json")
	req.Header.Add("Authorization", discord.Bearer)

	res, err := client.Do(req)
	if err != nil {
		console.Log(err)
		return
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		console.Log(err)
		return
	}

	err = json.Unmarshal(data, &current_metadata)
	if err != nil {
		console.Log(err)
		return
	}

	for _, plan := range subscriptions.Plans {
		expected_metadata = append(expected_metadata, metadata{Type: 7, Name: plan.GetData().Name, Description: "User has selected " + plan.GetData().Name + " plan.", Key: strconv.Itoa(int(plan))})
	}

	if len(current_metadata) != len(expected_metadata) {

		data, err = json.Marshal(expected_metadata)
		if err != nil {
			console.Log(err)
			return
		}

		req, err = http.NewRequest(http.MethodPut, "https://discord.com/api/v"+discord.API_Version+"/applications/"+discord.Application_ID+"/role-connections/metadata", bytes.NewBuffer(data))
		if err != nil {
			console.Log(err)
			return
		}

		req.Header.Set("Content-type", "application/json")
		req.Header.Add("Authorization", discord.Bearer)

		_, err = client.Do(req)
		if err != nil {
			console.Log(err)
			return
		}

		console.Log("Linked role metadata was updated successfully.")

	} else {
		console.Log("Currently used linked role metadata is up to date.")
	}
}

type metadata struct {
	Type        int    `json:"type"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
