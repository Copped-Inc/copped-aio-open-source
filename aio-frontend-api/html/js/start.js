let running = false;

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

function setStart() {
    running = (data?.session?.status || "")  === "Running";

    document.getElementById("running").innerHTML = (data?.session?.tasks || 0).toString();
    document.getElementById("instances").innerHTML = (data?.session?.instances || []).length.toString();
    document.getElementById("checkout").innerHTML = (data?.checkouts || []).length.toString();
    document.getElementById("status").innerHTML = data?.session?.status || "not Running";
    document.getElementById("startbutton").innerHTML = running ?
        "<svg class=\"svgbutton\" onclick=\"changeStart()\" xmlns=\"http://www.w3.org/2000/svg\" x=\"0px\" y=\"0px\" width=\"60\" height=\"60\" viewBox=\"0 0 172 172\"><g fill=\"none\" fill-rule=\"nonzero\" stroke=\"none\" stroke-width=\"1\" stroke-linecap=\"butt\" stroke-linejoin=\"miter\" stroke-miterlimit=\"10\" stroke-dasharray=\"\" stroke-dashoffset=\"0\" font-family=\"none\" font-size=\"none\" style=\"mix-blend-mode: normal\"><path d=\"M0,172v-172h172v172z\" fill=\"none\"></path><g fill=\"#FFFFFF73\"><path d=\"M57.33333,35.83333c-7.91917,0 -14.33333,6.41417 -14.33333,14.33333v71.66667c0,7.91917 6.41417,14.33333 14.33333,14.33333c7.91917,0 14.33333,-6.41417 14.33333,-14.33333v-71.66667c0,-7.91917 -6.41417,-14.33333 -14.33333,-14.33333zM114.66667,35.83333c-7.91917,0 -14.33333,6.41417 -14.33333,14.33333v71.66667c0,7.91917 6.41417,14.33333 14.33333,14.33333c7.91917,0 14.33333,-6.41417 14.33333,-14.33333v-71.66667c0,-7.91917 -6.41417,-14.33333 -14.33333,-14.33333z\"></path></g></g></svg>" :
        "<svg class=\"svgbutton\" onclick=\"changeStart()\" xmlns=\"http://www.w3.org/2000/svg\" x=\"0px\" y=\"0px\" width=\"80\" height=\"80\" viewBox=\"0 0 172 172\"><g fill=\"none\" fill-rule=\"nonzero\" stroke=\"none\" stroke-width=\"1\" stroke-linecap=\"butt\" stroke-linejoin=\"miter\" stroke-miterlimit=\"10\" stroke-dasharray=\"\" stroke-dashoffset=\"0\" font-family=\"none\" font-size=\"none\" style=\"mix-blend-mode: normal\"><path d=\"M0,172v-172h172v172z\" fill=\"none\"></path><g id=\"original-icon\" fill=\"#FFFFFF73\"><path d=\"M57.33333,48.891v74.22517c0,5.6545 6.24217,9.08017 11.01517,6.04867l58.31517,-37.109c4.429,-2.8165 4.429,-9.27367 0,-12.09017l-58.31517,-37.12333c-4.773,-3.03867 -11.01517,0.39417 -11.01517,6.04867z\"></path></g></g></svg>";
}

function changeStart() {
    if (data?.session?.status === "Stopped by Admin") {
        error("Session was forced stopped by Admin");
        return;
    }

    if (running) {
        data?.session === undefined ? data.session = { status: "not Running" } : data.session.status = "not Running";
        data.session.instances = [];
        data.session.checkouts = [];
        data.session.tasks = 0;
    } else {
        if (data?.shipping === undefined || data.shipping.length === 0 || data.billing === undefined || data.billing.length === 0) {
            error("Please add at least one shipping and billing profile");
            return;
        }

        data?.session === undefined ? data.session = { status: "Running" } : data.session.status = "Running";

        let runningInstances = [];
        let runningTasks = 0;
        for (let i = 0; i < data.instances.length; i++) {
            if (data.instances[i].status === "Running") {
                runningInstances.push(data.instances[i]);
                runningTasks += parseInt(data.instances[i].task_max);
            }
        }
        data.session.instances = runningInstances;
        data.session.tasks = runningTasks;
    }
    setStart();
    updateSession();
}

function updateSession() {
    fetch(activedata + "data/session", {
        method: "PATCH",
        credentials: "include",
        headers: {
            "Password": password
        },
        body: JSON.stringify({
            "Id": requestId(),
            "Session": data.session
        })
    }).then(resp => {
        if (resp.status !== 200) {
            throw new Error("Error: " + resp.status);
        }
    }).catch(err => {
        error(err);
    })
}
