function setInstances() {
    let filter = 0;

    if (document.getElementById("allinstance").classList.contains("switchactive")) {
        filter = 1;
    } else if (document.getElementById("stoppedinstance").classList.contains("switchactive")) {
        filter = 2;
    } else if (document.getElementById("runninginstance").classList.contains("switchactive")) {
        filter = 3;
    }

    data.instances = data.instances == null ? [] : data.instances;

    let instanceslist = document.getElementById("instancelist");
    instanceslist.innerHTML = "";

    for (let i = 0; i < data.instances.length; i++) {
        if (filter === 1 || filter === 2 && data.instances[i].status === "Stopped" || filter === 3 && data.instances[i].status === "Running") {
            let instance = data.instances[i];
            let instancediv = document.createElement("div");

            let pausebutton = data.instances[i].status === "Running" ? `<svg onclick="pause('${instance.id}')" xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="24" height="24" viewBox="0 0 172 172"><g fill="none" fill-rule="nonzero" stroke="none" stroke-width="1" stroke-linecap="butt" stroke-linejoin="miter" stroke-miterlimit="10" stroke-dasharray="" stroke-dashoffset="0" font-family="none" font-size="none" style="mix-blend-mode: normal"><path d="M0,172v-172h172v172z" fill="none"></path><g fill="#FFFFFF73"><path d="M57.33333,35.83333c-7.91917,0 -14.33333,6.41417 -14.33333,14.33333v71.66667c0,7.91917 6.41417,14.33333 14.33333,14.33333c7.91917,0 14.33333,-6.41417 14.33333,-14.33333v-71.66667c0,-7.91917 -6.41417,-14.33333 -14.33333,-14.33333zM114.66667,35.83333c-7.91917,0 -14.33333,6.41417 -14.33333,14.33333v71.66667c0,7.91917 6.41417,14.33333 14.33333,14.33333c7.91917,0 14.33333,-6.41417 14.33333,-14.33333v-71.66667c0,-7.91917 -6.41417,-14.33333 -14.33333,-14.33333z"></path></g></g></svg>`
                : `<svg onclick="pause('${instance.id}')" xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="24" height="24" viewBox="0 0 172 172" style=" fill:#000000;"><g fill="none" fill-rule="nonzero" stroke="none" stroke-width="1" stroke-linecap="butt" stroke-linejoin="miter" stroke-miterlimit="10" stroke-dasharray="" stroke-dashoffset="0" font-family="none" font-size="none" style="mix-blend-mode: normal"><path d="M0,172v-172h172v172z" fill="none"></path><g fill="#FFFFFF73"><path d="M57.33333,48.891v74.22517c0,5.6545 6.24217,9.08017 11.01517,6.04867l58.31517,-37.109c4.429,-2.8165 4.429,-9.27367 0,-12.09017l-58.31517,-37.12333c-4.773,-3.03867 -11.01517,0.39417 -11.01517,6.04867z"></path></g></g></svg>`;

            instancediv.classList.add("instance");
            instancediv.innerHTML = `
                <div class="actions" id="actions">
                        ${pausebutton}
                        <svg onclick="remove('${instance.id}')" xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="24" height="24" viewBox="0 0 172 172"><g fill="none" fill-rule="nonzero" stroke="none" stroke-width="1" stroke-linecap="butt" stroke-linejoin="miter" stroke-miterlimit="10" stroke-dasharray="" stroke-dashoffset="0" font-family="none" font-size="none" style="mix-blend-mode: normal"><path d="M0,172v-172h172v172z" fill="none"></path><g fill="#FFFFFF73"><path d="M71.66667,14.33333l-7.16667,7.16667h-28.66667c-4.3,0 -7.16667,2.86667 -7.16667,7.16667c0,4.3 2.86667,7.16667 7.16667,7.16667h14.33333h71.66667h14.33333c4.3,0 7.16667,-2.86667 7.16667,-7.16667c0,-4.3 -2.86667,-7.16667 -7.16667,-7.16667h-28.66667l-7.16667,-7.16667zM35.83333,50.16667v93.16667c0,7.88333 6.45,14.33333 14.33333,14.33333h71.66667c7.88333,0 14.33333,-6.45 14.33333,-14.33333v-93.16667zM64.5,64.5c4.3,0 7.16667,2.86667 7.16667,7.16667v64.5c0,4.3 -2.86667,7.16667 -7.16667,7.16667c-4.3,0 -7.16667,-2.86667 -7.16667,-7.16667v-64.5c0,-4.3 2.86667,-7.16667 7.16667,-7.16667zM107.5,64.5c4.3,0 7.16667,2.86667 7.16667,7.16667v64.5c0,4.3 -2.86667,7.16667 -7.16667,7.16667c-4.3,0 -7.16667,-2.86667 -7.16667,-7.16667v-64.5c0,-4.3 2.86667,-7.16667 7.16667,-7.16667z"></path></g></g></svg>
                    </div>
                    <table class="table">
                        <tr>
                            <td>Server €/h</td>
                            <td>${instance.price}</td>
                        </tr>
                        <tr>
                            <td>Provider</td>
                            <td>${instance.provider}</td>
                        </tr>
                        <tr>
                            <td>ID</td>
                            <td>${instance.id}</td>
                        </tr>
                        <tr>
                            <td>Status</td>
                            <td>${instance.status}</td>
                        </tr>
                        <tr>
                            <td>Task Max</td>
                            <td>${instance.task_max}</td>
                        </tr>
                        <tr>
                            <td>Region</td>
                            <td>${instance.region}</td>
                        </tr>
                    </table>
            `
            instanceslist.appendChild(instancediv)
        }
    }

    if (instanceslist.children.length === 0) {
        instanceslist.innerHTML = `
            <div class="nothing-list">
                <svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="24" height="24" viewBox="0 0 172 172"><g fill="none" fill-rule="nonzero" stroke="none" stroke-width="1" stroke-linecap="butt" stroke-linejoin="miter" stroke-miterlimit="10" stroke-dasharray="" stroke-dashoffset="0" font-family="none" font-size="none" style="mix-blend-mode: normal"><path d="M0,172v-172h172v172z" fill="none"></path><g fill="#FFFFFF73"><path d="M135.45,0c-6.97406,0.16125 -13.82719,2.33813 -19.6725,6.45l-28.4875,20.1025l7.8475,11.395l28.595,-20.21c4.79719,-3.35937 10.57531,-4.66281 16.34,-3.655c5.76469,1.00781 10.83063,4.23281 14.19,9.03c3.37281,4.79719 4.66281,10.57531 3.655,16.34c-1.00781,5.76469 -4.23281,10.81719 -9.03,14.19c-2.09625,1.47813 -23.08562,16.20563 -24.51,17.2c-5.46906,3.84313 -13.07469,5.76469 -22.4675,1.6125l-12.255,8.4925c13.06125,10.13188 30.46281,9.70188 42.57,1.1825c1.42438,-0.99437 22.54813,-15.84281 24.6175,-17.3075c7.78031,-5.4825 12.98063,-13.61219 14.62,-23.005c1.63938,-9.37937 -0.44344,-18.86625 -5.9125,-26.66c-5.4825,-7.79375 -13.71969,-12.98062 -23.1125,-14.62c-2.33812,-0.40312 -4.66281,-0.59125 -6.9875,-0.5375zM9.5675,3.3325c-2.67406,0.25531 -4.95844,2.05594 -5.83187,4.59563c-0.88688,2.55312 -0.20156,5.375 1.74687,7.22937l37.84,37.84c1.65281,2.05594 4.32688,2.98313 6.89344,2.39188c2.58,-0.59125 4.58219,-2.59344 5.17344,-5.17344c0.59125,-2.56656 -0.33594,-5.24062 -2.39188,-6.89344l-37.84,-37.84c-1.29,-1.35719 -3.07719,-2.13656 -4.945,-2.15c-0.215,-0.01344 -0.43,-0.01344 -0.645,0zM26.5525,87.29l-20.1025,28.595c-5.4825,7.79375 -7.55188,17.17313 -5.9125,26.5525c1.63938,9.39281 6.82625,17.63 14.62,23.1125c7.79375,5.4825 17.28063,7.55188 26.66,5.9125c9.37938,-1.63937 17.53594,-6.83969 23.005,-14.62c1.46469,-2.06937 16.31313,-23.19312 17.3075,-24.6175c8.51938,-12.10719 8.94938,-29.40125 -1.1825,-42.4625l-8.4925,12.1475c4.15219,9.39281 2.23063,16.99844 -1.6125,22.4675c-0.99437,1.42438 -15.72187,22.52125 -17.2,24.6175c-3.37281,4.79719 -8.42531,7.91469 -14.19,8.9225c-5.76469,1.00781 -11.54281,-0.28219 -16.34,-3.655c-4.79719,-3.35937 -8.02219,-8.42531 -9.03,-14.19c-1.00781,-5.76469 0.29563,-11.54281 3.655,-16.34l20.21,-28.595z"></path></g></g></svg>
                Nothing here
            </div>
        `;
    }

    if (data.user.instance_limit > data.instances.length) {
        document.getElementById("addinstance").style.display = "block";
    } else {
        document.getElementById("addinstance").style.display = "none";
        document.getElementById("addinstancefield").style.display = "none";
        document.getElementById("addinstance").innerHTML = "+";
    }
}

function pause(id) {
    for (let i = 0; i < data.instances.length; i++) {
        if (data.instances[i].id === id) {
            if (data.instances[i].status === "Running") {
                data.instances[i].status = "Stopped";
            } else if (data.instances[i].status === "Stopped") {
                data.instances[i].status = "Running";
            } else {
                return;
            }
            setInstances()

            fetch(activedata + "data/instance", {
                method: "PATCH",
                credentials: "include",
                headers: {
                    "Password": password
                },
                body: JSON.stringify({
                    id: requestId(),
                    instance: data.instances[i]
                })
            }).then(resp => {
                if (resp.status !== 200) {
                    throw new Error("Error: " + resp.status);
                }
            }).catch(err => {
                error(err);
            })
            return
        }
    }
}

function remove(id) {
    for (let i = 0; i < data.instances.length; i++) {
        if (data.instances[i].id === id) {
            let insta = data.instances[i];
            data.instances.splice(i, 1);
            setInstances()

            fetch(activedata + "data/instance", {
                method: "DELETE",
                credentials: "include",
                headers: {
                    "Password": password
                },
                body: JSON.stringify({
                    id: requestId(),
                    instance: insta
                })
            }).then(resp => {
                if (resp.status !== 200) {
                    throw new Error("Error: " + resp.status);
                }
            }).catch(err => {
                error(err);
            })
            return
        }
    }
}

function changeInstance(name) {
    document.getElementById("allinstance").classList.remove("switchactive");
    document.getElementById("stoppedinstance").classList.remove("switchactive");
    document.getElementById("runninginstance").classList.remove("switchactive");
    document.getElementById(name + "instance").classList.add("switchactive");
    setInstances();
}

let addinstance = false

function addInstance() {
    if (addinstance) {
        document.getElementById("addinstancefield").style.display = "none";
        document.getElementById("addinstance").innerHTML = "<svg class=\"svgbutton\" onclick=\"addInstance()\" version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\" width=\"24px\" height=\"24px\" viewBox=\"0,0,256,256\"><g fill-opacity=\"0.45098\" fill=\"#FFFFFFE5\" fill-rule=\"nonzero\" stroke=\"none\" stroke-width=\"1\" stroke-linecap=\"butt\" stroke-linejoin=\"miter\" stroke-miterlimit=\"10\" stroke-dasharray=\"\" stroke-dashoffset=\"0\" font-family=\"none\" font-size=\"none\" style=\"mix-blend-mode: normal\"><g transform=\"scale(10.66667,10.66667)\"><path d=\"M20,11h-7v-7c0,-0.552 -0.448,-1 -1,-1c-0.552,0 -1,0.448 -1,1v7h-7c-0.552,0 -1,0.448 -1,1c0,0.552 0.448,1 1,1h7v7c0,0.552 0.448,1 1,1c0.552,0 1,-0.448 1,-1v-7h7c0.552,0 1,-0.448 1,-1c0,-0.552 -0.448,-1 -1,-1z\"></path></g></g></svg>";
    } else {
        document.getElementById("addinstancefield").style.display = "block";
        document.getElementById("addinstance").innerHTML = "<svg class=\"svgbutton\" onclick=\"addInstance()\" version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\" width=\"24px\" height=\"24px\" viewBox=\"0,0,256,256\"><g fill-opacity=\"0.45098\" fill=\"#ffffff\" fill-rule=\"nonzero\" stroke=\"none\" stroke-width=\"1\" stroke-linecap=\"butt\" stroke-linejoin=\"miter\" stroke-miterlimit=\"10\" stroke-dasharray=\"\" stroke-dashoffset=\"0\" font-family=\"none\" font-size=\"none\" style=\"mix-blend-mode: normal\"><g transform=\"scale(10.66667,10.66667)\"><path d=\"M5.99023,4.99023c-0.40692,0.00011 -0.77321,0.24676 -0.92633,0.62377c-0.15312,0.37701 -0.06255,0.80921 0.22907,1.09303l5.29297,5.29297l-5.29297,5.29297c-0.26124,0.25082 -0.36647,0.62327 -0.27511,0.97371c0.09136,0.35044 0.36503,0.62411 0.71547,0.71547c0.35044,0.09136 0.72289,-0.01388 0.97371,-0.27511l5.29297,-5.29297l5.29297,5.29297c0.25082,0.26124 0.62327,0.36648 0.97371,0.27512c0.35044,-0.09136 0.62411,-0.36503 0.71547,-0.71547c0.09136,-0.35044 -0.01388,-0.72289 -0.27512,-0.97371l-5.29297,-5.29297l5.29297,-5.29297c0.29576,-0.28749 0.38469,-0.72707 0.22393,-1.10691c-0.16075,-0.37985 -0.53821,-0.62204 -0.9505,-0.60988c-0.2598,0.00774 -0.50638,0.11632 -0.6875,0.30273l-5.29297,5.29297l-5.29297,-5.29297c-0.18827,-0.19353 -0.4468,-0.30272 -0.7168,-0.30273z\"></path></g></g></svg>";
        getInstanceCode();
    }
    addinstance = !addinstance;
}

function getInstanceCode() {
    let downloadload = document.getElementById("downloadload");
    let addinstancebutton = document.getElementById("addinstancebutton");
    let instancecode = document.getElementById("instancecode");

    downloadload.style.display = "flex";
    addinstancebutton.style.display = "none";
    instancecode.style.display = "none";

    fetch(activedata + "instance", {
        method: "GET",
        credentials: "include",
        headers: {
            "Password": password
        }
    }).then(resp => {
        if (resp.status === 200) {
            resp.json().then(data => {
                instancecode.innerHTML = data.code;
                downloadload.style.display = "";
                addinstancebutton.style.display = "";
                instancecode.style.display = "";
            })
        } else {
            throw new Error("Error: " + resp.status);
        }
    }).catch(err => {
        instancecode.innerHTML = "Error";
        downloadload.style.display = "";
        addinstancebutton.style.display = "";
        instancecode.style.display = "";
        error(err);
    })
}

function copyCode() {
    let instancecode = document.getElementById("instancecode");
    navigator.clipboard.writeText(instancecode.innerHTML).then(() => {
        let instancecode = document.getElementById("addinstancebutton");
        instancecode.innerHTML = "Copied!";
        sleep(1000).then(() => {
            instancecode.innerHTML = "Copy Code";
        })
    });
}

function createCloud() {
    toggleLoading();

    fetch(activedata + "data/instance", {
        method: "POST",
        credentials: "include",
    }).then(resp => {
        if (resp.status !== 200) {
            throw new Error("Error: " + resp.status);
        }

        toggleLoading();
        addInstance();
    }).catch(err => {
        toggleLoading();
        error(err);
    })
}

function download(file) {
    toggleLoading();

    fetch(activedata + "instance/download/" + file, {
        method: "GET",
        credentials: "include",
        headers: {
            "Password": password
        }
    }).then(resp => {
        if (resp.status === 200) {
            return resp.blob();
        } else {
            throw new Error("Error: " + resp.status);
        }
    }).then(data => {
        let a = document.createElement("a");
        a.href = window.URL.createObjectURL(data);
        if (file === "payments") a.download = "payments.exe";
        else a.download = "copped-aio.exe";

        a.click();

        sleep(1000).then(() => toggleLoading());
    }).catch(err => {
        toggleLoading();
        error(err);
    })
}

function toggleLoading() {
    let downloadload = document.getElementById("downloadload");
    let createcloud = document.getElementById("create-cloud");
    let downloadbuttonclient = document.getElementById("downloadbuttonclient");

    downloadload.style.display = downloadload.style.display === "flex" ? "" : "flex";
    createcloud.style.display = createcloud.style.display === "none" ? "" : "none";
    downloadbuttonclient.style.display = downloadbuttonclient.style.display === "none" ? "" : "none";
}