function setStores() {
    data.settings.stores = data.settings.stores || new Map();
    for (let [storename, enabled] of new Map(Object.entries(data.settings.stores))) {
        let store = document.getElementById(storename);
        if (store === null) {
            continue;
        }
        store.children.item(+!enabled).classList.add("switchactive");
        store.children.item(+enabled).classList.remove("switchactive");
    }
}

function updateStores(store, value) {
    data.settings.stores[store] = value;
    setStores();

    fetch(activedata + "data/stores", {
        method: "PATCH",
        credentials: "include",
        headers: {
            "Password": password
        },
        body: JSON.stringify({
            "Id": requestId(),
            "Store": store,
            "Value": value
        })
    }).then(resp => {
        if (resp.status !== 200) {
            throw new Error("Error: " + resp.status);
        }
    }).catch(err => {
        error(err);
    });
}