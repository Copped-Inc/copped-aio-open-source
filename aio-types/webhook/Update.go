package webhook

import (
	"github.com/Copped-Inc/aio-types/branding"
	"github.com/Copped-Inc/aio-types/helper"
	"os"
	"time"

	"github.com/infinitare/disgo"
)

func (b *Body) AddEmbed(a webhook, args ...string) *Body {
	embed := disgo.Embed{
		Timestamp: disgo.Timestamp(time.Now().UTC().Format("2006-01-02T15:04:05.999Z")),
		Color:     branding.Grey,
		Footer: &disgo.EmbedFooter{
			Text:     "Copped AIO",
			Icon_URL: "https://cdn.discordapp.com/attachments/901123672809566258/1007591468485066792/discord.png",
		},
	}

	switch webhook(a) {
	case Test:
		embed.Title = "Test Webhook"
		embed.Description = "Your Webhook has been successfully created!"

	case DataRequest:
		embed.Title = "Data Requested"
		embed.Description = "Instance `" + args[0] + "` needs to retrieve your data. This may occur because the instance has restarted or updated. To perform this action, please just visit the [Dashboard](https://aio.copped-inc.com) once."

	case UpdateClient:
		version, err := os.ReadFile("version-client-local-v3")
		if err != nil {
			version = []byte("Unknown")
		}

		embed.Title = "Client Update available"
		embed.Description = "Version `" + string(version) + "` is now pushed. All running instances will be restarted and you will receive a notification that you need to log in to your dashboard once."

	case UpdatePayments:
		version, err := os.ReadFile("version-payments")
		if err != nil {
			version = []byte("Unknown")
		}

		embed.Title = "Payments Update available"
		embed.Description = "Version `" + string(version) + "` is now pushed. The update for your payment instance will be performed after the next restart."

	case DataReceived:
		embed.Title = "Data Received"
		embed.Description = "Your instances have received all data."

	case NewProduct:
		embed.Title = args[0]
		embed.Description = "This product is classified as profitable. Your Copped AIO instance will run automatically when enabled.\n\nPlease solve captchas manually at **https://aio.copped-inc.com/captcha** to ensure that your instance is working."
		embed.Fields = []disgo.EmbedField{
			{
				Name:   "Product",
				Value:  args[1],
				Inline: true,
			},
			{
				Name:   "StockX",
				Value:  args[2],
				Inline: true,
			},
			{
				Name:   "Price",
				Value:  args[3],
				Inline: true,
			},
		}
		embed.Thumbnail = &disgo.EmbedThumbnail{URL: args[4]}

	case Log:
		embed.Title = helper.Server
		embed.Description = "```" + args[0] + "```"

	case NewCheckout:
		embed.Title = args[0]
		embed.Description = "A Copped AIO client checked out successfully."
		embed.Thumbnail = &disgo.EmbedThumbnail{URL: args[2]}
		embed.URL = args[1]
		embed.Fields = []disgo.EmbedField{
			{
				Name:   "Store",
				Value:  args[3],
				Inline: true,
			},
			{
				Name:   "Size",
				Value:  args[4],
				Inline: true,
			},
			{
				Name:   "Price",
				Value:  args[5],
				Inline: true,
			},
		}

	case NewCheckoutLink:
		embed.Title = "Open Checkout Link"
		embed.Description = "Click link to open the checkout page.\n[Mobile](" + args[0] + ")"
		embed.URL = args[0]

	case PingFailed:
		embed.Title = args[0]
		embed.Description = "The API is not responding. Please check!\n\n```" + args[1] + "```"
		embed.URL = "https://" + args[0]
		embed.Color = branding.Red

	case PingSuccess:
		embed.Title = args[0]
		embed.Description = "The API is up and running again."
		embed.URL = "https://" + args[0]

	case Whitelist:
		embed.Title = "Global Whitelist"
		embed.Description = "These are the products that are globally whitelisted. You can add or remove these products from your dashboard.\n\n"

	case MonitorDisabled:
		embed.Title = "Monitor " + args[0] + " Disabled"
		embed.Description = "One Monitor has been disabled. for security / spam reasons. Please check the [Developer Dashboard](https://aio.copped-inc.com) for more information.\nTo restart click [here](https://database.copped-inc.com/monitor/restart/" + args[0] + ").\n\n```" + args[1] + "```"

	case MonitorRestarted:
		embed.Title = "Monitor Restarted"
		embed.Description = "`" + args[0] + "` has been restarted successfully."

	case ErrorLog:
		embed.Title = "Error Log"
		embed.Description = "Instance `" + args[0] + "` logged an error. To download the logs, visit the [Developer Dashboard](https://aio.copped-inc.com)\n```json\n" + args[1] + "\n```"

	case InstanceLogout:
		embed.Title = "Instance Logout"
		embed.Description = "Your Instance `" + args[0] + "` has been logged out. Please check the [Dashboard](https://aio.copped-inc.com) for more information. If you can't solve the problem, please contact our support via Discord."

	case AiPredictIntra:
		embed.Title = args[0] + " Prediction"

		embed.Fields = []disgo.EmbedField{
			{
				Name:   "Einstieg",
				Value:  "`" + args[1] + "€`",
				Inline: true,
			},
			{
				Name:   "Prognose",
				Value:  args[2],
				Inline: true,
			},
			{
				Name:   "\u200b",
				Value:  "\u200b",
				Inline: true,
			},
			{
				Name:  "Hinweis",
				Value: "Dies ist keine Anlageberatung, sondern nur ein KI-Modell, welches auf Basis des bisherigen Kursverlaufes die wahrscheinliche Kursentwicklung für den Tag angibt.",
			},
			{
				Name:  "Legende",
				Value: "📈 — Der Kurs wird zum Tagesende über `" + args[1] + "€` gestiegen sein.\n📉 — Der Kurs wird zum Tagesende unter `" + args[1] + "€` gefallen sein.",
			},
		}

	case AiPredictDiff:
		embed.Title = args[0] + " Prediction"

		embed.Fields = []disgo.EmbedField{
			{
				Name:   "vor. Close",
				Value:  "`" + args[1] + "€`",
				Inline: true,
			},
			{
				Name:   "Prognose",
				Value:  args[2],
				Inline: true,
			},
			{
				Name:   "\u200b",
				Value:  "\u200b",
				Inline: true,
			},
			{
				Name:  "Hinweis",
				Value: "Dies ist keine Anlageberatung, sondern nur ein KI-Modell, welches auf Basis des bisherigen Kursverlaufes die wahrscheinliche Kursentwicklung für den Tag angibt.",
			},
			{
				Name:  "Legende",
				Value: "📈 — Der Kurs wird zum Tagesbeginn über `" + args[1] + "€` gestiegen sein.\n📉 — Der Kurs wird zum Tagesbeginn unter `" + args[1] + "€` gefallen sein.",
			},
		}
	case ISINList:
		embed.Title = "Tesla Faktor Zertifikate"

		embed.Fields = []disgo.EmbedField{
			{
				Name:   "Long",
				Value:  "**2x** [DE000GQ8HHB4](https://www.boerse.de/wertpapier/DE000GQ8HHB4) [[SC](https://de.scalable.capital/broker/security?isin=DE000GQ8HHB4)]\n**5x** [DE000GQ7TR82](https://www.boerse.de/wertpapier/DE000GQ7TR82) [[SC](https://de.scalable.capital/broker/security?isin=DE000GQ7TR82)]\n**10x** [DE000GQ85Q46](https://www.boerse.de/wertpapier/DE000GQ85Q46) [[SC](https://de.scalable.capital/broker/security?isin=DE000GQ85Q46)]",
				Inline: true,
			},
			{
				Name:   "Short",
				Value:  "**2x** [DE000GQ8HMF5](https://www.boerse.de/wertpapier/DE000GQ8HMF5) [[SC](https://de.scalable.capital/broker/security?isin=DE000GQ8HMF5)]\n**5x** [DE000GQ7GHU8](https://www.boerse.de/wertpapier/DE000GQ7GHU8) [[SC](https://de.scalable.capital/broker/security?isin=DE000GQ7GHU8)]\n**10x** [DE000GQ768W1](https://www.boerse.de/wertpapier/DE000GQ768W1) [[SC](https://de.scalable.capital/broker/security?isin=DE000GQ768W1)]",
				Inline: true,
			},
		}
	}

	b.Embeds = append(b.Embeds, embed)
	return b
}

func (b *Body) SetFields(a ...disgo.EmbedField) *Body {
	b.Embeds[len(b.Embeds)-1].Fields = a
	return b
}
