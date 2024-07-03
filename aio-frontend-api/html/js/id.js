let requestId = () => {
    let s4 = () => {
        return Math.floor((1 + Math.random()) * 0x10000)
            .toString(16)
            .substring(1);
    }

    let id = (s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + '-' + s4() + s4() + s4()).toString()
    addId(id);

    return id;
}

function checkId(d) {
    if (d?.body !== undefined) {
        return requestIds.includes(d.body.id);
    }
    return false;
}

function addId(id) {
    requestIds.push(id);
}

function removeId(id) {
    requestIds.splice(requestIds.indexOf(id), 1);
}

let requestIds = [];