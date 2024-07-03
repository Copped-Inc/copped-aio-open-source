function connect(reconnect = false) {
    let ws = new WebSocket(activedata.replaceAll("http", "ws") + "websocket");
    ws.onopen = () => {
        if (reconnect) signal("Reconnected");
        sendPassword(ws);
    }

    ws.onmessage = function (event) {
        let d = JSON.parse(event.data);
        switch (d.op) {
            case 1:
                ping(ws).then();
                break;
            case 2:
                updateData(d.data).then();
                break;
            case 3:
                sendPassword(ws);
                break;
            case 8:
                notifications.unshift(d.data);
                setNotifications();
                break;
            case 9:
                let i = notifications.indexOf(notifications.find(n => n.id === d.data.id));
                notifications[i].title = d.data.title;
                notifications[i].text = d.data.text;
                notifications[i].created_at = d.data.created_at;
                setNotifications();
                break;
            case 10:
                if (d?.data === undefined) {
                    fetchNotifications()
                } else {
                    removeNotification(d.data.id);
                }
                break;
        }
    };

    ws.onclose = function () {
        error("Connection lost");
        connect(true);
    };
}

async function ping(ws) {
    await sleep(20000);
    if (ws.readyState !== 1) {
        return;
    }

    ws.send(JSON.stringify({
        op: 1
    }));
}

async function updateData(d) {
    if (checkId(d)) {
        removeId(d);
        return;
    }

    switch (d.action) {
        case AddWebhook:
            data.settings.webhooks = data.settings.webhooks == null ? [] : data.settings.webhooks;
            data.settings.webhooks.push(d.body.webhook);

            setWebhooks();
            break;
        case DeleteWebhook:
            data.settings.webhooks = data.settings.webhooks == null ? [] : data.settings.webhooks;
            data.settings.webhooks.splice(data.settings.webhooks.indexOf(d.body.webhook), 1);

            setWebhooks();
            break;
        case UpdateStores:
            switch (d.body.store) {
                case "kith_eu":
                    data.settings.stores.kith_eu = d.body.value;
                    break
            }

            setStores();
            break;
        case UpdateInstances:
            data.instances = d.body;

            setInstances();
            break;
        case UpdateSession:
            data.session = d.body.session;

            setStart();
            break;
        case AddCheckout:
            d.body.checkout.date = new Date(d.body.checkout.date);
            data.checkouts.unshift(d.body.checkout);

            setChart();
            setCheckouts();
            break;
        case UpdateBilling:
            data.billing = d.body.billing;
            break;
        case UpdateShipping:
            data.shipping = d.body.shipping;
            break;
        case CreateNotification:
            notifications.unshift(d.body);
            setNotifications();
            break;
        case UpdateNotification:
            notifications[notifications.indexOf(notifications.find(n => n.id === d.body.id))] = d.body;
            if (d.body.global) {
                fetchGlobalNotification(d.body.id);
            } else {
                setNotifications();
            }
            break;
        case DeleteNotification:
            if (d?.body === undefined) {
                notifications = []
                setNotifications()
            } else {
                removeNotification(d.body.id);
            }
            break;
        case AddWhitelist:
            data.whitelist = data.whitelist == null ? [] : data.whitelist;
            data.whitelist.push(d.body.product);

            setWhitelist();
            break;
        case RemoveWhitelist:
            data.whitelist = data.whitelist == null ? [] : data.whitelist;
            data.whitelist.splice(data.whitelist.indexOf(d.body.product), 1);

            setWhitelist();
            break;
        case UpdateNotificationReadstate:
            notifications.forEach(n => n.read = d.body)
            setNotifications();
            break
    }
}

const AddWebhook = 1;
const DeleteWebhook = 2;
const UpdateStores = 3;
const UpdateInstances = 4;
const UpdateSession = 5;
const AddCheckout = 6;
const UpdateBilling = 7;
const UpdateShipping = 8;
const CreateNotification = 9;
const UpdateNotification = 10;
const DeleteNotification = 11;
const AddWhitelist = 12;
const RemoveWhitelist = 13;
const UpdateNotificationReadstate = 14

function sendPassword(ws) {
    ws.send(JSON.stringify({
        op: 3,
        data: password
    }));
}

function parseDate() {
    data.checkouts = data.checkouts == null ? [] : data.checkouts;
    for (let i = 0; i < data.checkouts.length; i++) {
        data.checkouts[i].date = new Date(data.checkouts[i].date);
    }
}