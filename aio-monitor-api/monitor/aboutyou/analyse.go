package aboutyou

import "github.com/Copped-Inc/aio-types/console"

func analyse(res []product) {

	for _, re := range res {
		if re.Id == lastProduct.Id {
			break
		}

		if lastProduct.Id == 0 {
			continue
		}

		console.Log("Aboutyou Found", re.Attributes.Name.Values.Label, re.Attributes.Brand.Values.Value, re.Id)

		if re.Attributes.Brand.Values.Value == "nike" {
			go re.send()
		}
	}

	lastProduct = res[0]

}
