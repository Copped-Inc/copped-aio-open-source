function downloadShipping() {
    data.shipping = data.shipping == null ? [] : data.shipping;
    let shipping = "LASTNAME,ADDRESS1,ADDRESS2,CITY,STATE,ZIP,COUNTRY,EMAIL";
    for (let i = 0; i < data.shipping.length; i++) {
        shipping += "\n" + data.shipping[i].last + "," + data.shipping[i].address1 + "," + data.shipping[i].address2 + "," + data.shipping[i].city + "," + data.shipping[i].state + "," + data.shipping[i].zip + "," + data.shipping[i].country + "," + data.shipping[i].email;
    }

    save("shipping.csv", shipping);
    toggleDataDropdown();
}

function downloadBilling() {
    data.billing = data.billing == null ? [] : data.billing;
    let billing = "CCNUMBER,MONTH,YEAR,CVV";
    for (let i = 0; i < data.billing.length; i++) {
        billing += "\n" + data.billing[i].ccnumber + "," + data.billing[i].month + "," + data.billing[i].year + "," + data.billing[i].cvv;
    }

    save("billing.csv", billing);
    toggleDataDropdown();
}

function save(filename, data) {
    const blob = new Blob([data], {type: 'text/csv'});
    if(window.navigator.msSaveOrOpenBlob) {
        window.navigator.msSaveBlob(blob, filename);
    }
    else{
        const elem = window.document.createElement('a');
        elem.href = window.URL.createObjectURL(blob);
        elem.download = filename;
        document.body.appendChild(elem);
        elem.click();
        document.body.removeChild(elem);
    }
}

async function readFile(text) {
    text = text.replaceAll("ß", "ss");
    text = text.replaceAll("ä", "ae");
    text = text.replaceAll("ö", "oe");
    text = text.replaceAll("ü", "ue");
    text = text.replaceAll("Ä", "Ae");
    text = text.replaceAll("Ö", "Oe");
    text = text.replaceAll("Ü", "Ue");
    let tmap = {};

    let lines = text.split("\n");

    lines[0] = lines[0].toLowerCase().trim().replaceAll(" ", "");
    let keys = lines[0].split(",");
    for (let i = 1; i < lines.length; i++) {
        let line = lines[i];
        let values = line.split(",");
        for (let j = 0; j < keys.length; j++) {
            tmap[keys[j]] = tmap[keys[j]] == null ? [] : tmap[keys[j]];
            tmap[keys[j]][i - 1] = values[j].trim();
        }
    }

    let billing = [];
    let shipping = [];
    for (let i = 0; i < keys.length; i++) {
        let key = keys[i];
        if (key.includes("ccnumber") || key.includes("cardnumber")) {
            for (let j = 0; j < tmap[key].length; j++) {
                billing[j] = billing[j] == null ? {} : billing[j];
                billing[j].ccnumber = tmap[key][j];
                billing[j].ccnumber = billing[j].ccnumber.replaceAll(" ", "");
                billing[j].ccnumber = billing[j].ccnumber.replace(/\d{4}/g, "$& ");
                billing[j].ccnumber = billing[j].ccnumber.substring(0, billing[j].ccnumber.length - 1);
            }
        } else if (key.includes("month")) {
            for (let j = 0; j < tmap[key].length; j++) {
                billing[j] = billing[j] == null ? {} : billing[j];
                billing[j].month = tmap[key][j];
                if (billing[j].month.length === 2 && billing[j].month.charAt(0) === "0") {
                    billing[j].month = billing[j].month.charAt(1);
                }
            }
        } else if (key.includes("year")) {
            for (let j = 0; j < tmap[key].length; j++) {
                billing[j] = billing[j] == null ? {} : billing[j];
                if (tmap[key][j].length < 4) {
                    tmap[key][j] = "20" + tmap[key][j];
                }
                billing[j].year = tmap[key][j];
            }
        } else if (key.includes("cvv") || key.includes("csv") || key.includes("cvc")) {
            for (let j = 0; j < tmap[key].length; j++) {
                billing[j] = billing[j] == null ? {} : billing[j];
                billing[j].cvv = tmap[key][j];
            }
        } else if (key.includes("last")) {
            for (let j = 0; j < tmap[key].length; j++) {
                shipping[j] = shipping[j] == null ? {} : shipping[j];
                shipping[j].last = tmap[key][j];
            }
        } else if (key.includes("address") && !key.includes("2") && !key.includes("3") && !key.includes("two") && !key.includes("three")) {
            for (let j = 0; j < tmap[key].length; j++) {
                shipping[j] = shipping[j] == null ? {} : shipping[j];
                shipping[j].address1 = tmap[key][j];
            }
        } else if (key.includes("2") || key.includes("two")) {
            for (let j = 0; j < tmap[key].length; j++) {
                shipping[j] = shipping[j] == null ? {} : shipping[j];
                shipping[j].address2 = tmap[key][j];
            }
        } else if (key.includes("city")) {
            for (let j = 0; j < tmap[key].length; j++) {
                shipping[j] = shipping[j] == null ? {} : shipping[j];
                shipping[j].city = tmap[key][j];
            }
        } else if (key.includes("state")) {
            for (let j = 0; j < tmap[key].length; j++) {
                shipping[j] = shipping[j] == null ? {} : shipping[j];
                shipping[j].state = tmap[key][j];
            }
        } else if (key.includes("zip") || key.includes("postal")) {
            for (let j = 0; j < tmap[key].length; j++) {
                shipping[j] = shipping[j] == null ? {} : shipping[j];
                shipping[j].zip = tmap[key][j];
            }
        } else if (key.includes("country")) {
            for (let j = 0; j < tmap[key].length; j++) {
                shipping[j] = shipping[j] == null ? {} : shipping[j];
                shipping[j].country = await checkCountry(tmap[key][j]);
                if (shipping[j].country === "") {
                    error("Country not found, please check")
                    return;
                }
            }
        } else if (key.includes("mail")) {
            for (let j = 0; j < tmap[key].length; j++) {
                shipping[j] = shipping[j] == null ? {} : shipping[j];
                shipping[j].email = tmap[key][j];
                shipping[j].email = shipping[j].email.toLowerCase().replaceAll("random", "")
                if (shipping[j].email[0] === "@") {
                    shipping[j].email = shipping[j].email.split("@")[1];
                }
            }
        }
    }

    if (billing.length > 0) {
        data.billing = billing;
        updateBilling();
    }

    if (shipping.length > 0) {
        data.shipping = shipping;
        updateShipping();
    }

    if (billing.length === 0 && shipping.length === 0) error("No data found");
    else {
        signal("Data updated " +
            function () {
                if (billing.length > 0 && shipping.length > 0) return "shipping & billing";
                else if (billing.length > 0) return "billing";
                else if (shipping.length > 0) return "shipping";
            }()
        )
    }
}

async function checkCountry(country) {
    let path = (country.length <= 3 ? "https://restcountries.com/v3.1/alpha/" : "https://restcountries.com/v3.1/name/") + country;
    let resp = await fetch(path);
    let json = await resp.json();
    return json[0].name.common || "";
}

function setDrag() {
    let dropzone = document.getElementById('dropzone');

    function showDropZone() {
        dropzone.style.visibility = "visible";
    }
    function hideDropZone() {
        dropzone.style.visibility = "hidden";
    }

    function allowDrag(e) {
        e.dataTransfer.dropEffect = 'copy';
        e.preventDefault();
    }

    function handleDrop(e) {
        e.preventDefault();
        hideDropZone();

        if (e.dataTransfer.items) {
            for (let i = 0; i < e.dataTransfer.items.length; i++) {
                if (e.dataTransfer.items[i].kind === 'file') {
                    const file = e.dataTransfer.items[i].getAsFile();
                    file.text().then(t => readFile(t));
                }
            }
        } else {
            for (let i = 0; i < e.dataTransfer.files.length; i++) {
                let file = e.dataTransfer.files[i];
                file.text().then(t => readFile(t));
            }
        }
    }

    window.addEventListener('dragenter', function(e) {
        showDropZone();
    });

    dropzone.addEventListener('dragenter', allowDrag);
    dropzone.addEventListener('dragover', allowDrag);

    dropzone.addEventListener('dragleave', function(e) {
        hideDropZone();
    });

    dropzone.addEventListener('drop', handleDrop);
}

function updateBilling() {
    fetch(activedata + "data/billing", {
        method: "PATCH",
        credentials: "include",
        headers: {
            "Password": password
        },
        body: JSON.stringify({
            "Id": requestId(),
            "Billing": data.billing
        })
    }).then(resp => {
        if (resp.status !== 200) {
            throw new Error("Error: " + resp.status);
        }
    }).catch(err => {
        error(err);
    })
}

function updateShipping() {
    fetch(activedata + "data/shipping", {
        method: "PATCH",
        credentials: "include",
        headers: {
            "Password": password
        },
        body: JSON.stringify({
            "Id": requestId(),
            "Shipping": data.shipping
        })
    }).then(resp => {
        if (resp.status !== 200) {
            throw new Error("Error: " + resp.status);
        }
    }).catch(err => {
        error(err);
    })
}

function toggleDataDropdown() {
    let dropdown = document.getElementById("data-dropdown");
    dropdown.style.display = dropdown.style.display === "flex" ? "none" : "flex";
}
