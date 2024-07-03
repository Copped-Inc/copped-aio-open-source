import {initializeApp} from "https://www.gstatic.com/firebasejs/9.15.0/firebase-app.js"
import {
    child,
    get,
    getDatabase,
    onValue,
    ref,
    remove,
    set,
    limitToLast,
    orderByChild,
    query
} from "https://www.gstatic.com/firebasejs/9.15.0/firebase-database.js"

const firebaseConfig = {
    apiKey: "",
    authDomain: "",
    databaseURL: "",
    projectId: "",
    storageBucket: "",
    messagingSenderId: "",
    appId: "",
    measurementId: ""
};  // Insert your firebase config here

const app = initializeApp(firebaseConfig);
const database = getDatabase(app, );
const dbRef = ref(database);
const months = ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"];
const states = ["Ok", "Error", "Timeout"]

let today = new Date();
let timeframeToday = today.getFullYear() + "-" + today.toLocaleDateString('en-US', { month: 'long'}) + "-" + today.getDate();
let timeframeFrom = "";
let timeframeTo = "";
let stats = {};
let displayedRequests = {};
let allLogs = {};
let requestLimit = 50;
let logLimit = 50;

let openedUser = {};
let globalCheckouts = [];
let globalInstances = [];
let checkedUrls = {};

function init() {
    if (getCookie("localhost") != null) {
        activedata = localdata;
    }

    document.getElementById("switchview").addEventListener("click", switchView);
    document.getElementById("close-generate-keys").addEventListener("click", toggleGenerateWindow);
    document.getElementById("generate-keys").addEventListener("click", toggleGenerateWindow);
    document.getElementById("logs-download").addEventListener("click", logsDownload);
    document.getElementById("delete-logs").addEventListener("click", toggleDeleteLogs);
    document.getElementById("close-delete-logs").addEventListener("click", toggleDeleteLogs);
    document.getElementById("submit-delete").addEventListener("click", submitDeleteLogs);

    document.getElementById("monitor-restart").addEventListener("click", toggleRestartWindow);
    document.getElementById("close-monitor-restart").addEventListener("click", toggleRestartWindow);

    document.getElementById("copy-request-id").addEventListener("click", copy("copy-request-id"));
    document.getElementById("copy-request-user-id").addEventListener("click", copy("copy-request-user-id"));
    document.getElementById("copy-user-id").addEventListener("click", copy("copy-user-id"));
    document.getElementById("copy-link").addEventListener("click", copy("copy-link"));
    document.getElementById("copy-checkout-link").addEventListener("click", copy("copy-checkout-link"));

    document.getElementById("user-email").addEventListener("change", onInputUser);
    document.getElementById("user-instance-limit").addEventListener("change", onInputUser);
    document.getElementById("user-picture").addEventListener("change", onInputUser);

    document.getElementById("checkout-open-user").addEventListener("click", () => {
        document.getElementById("checkout-window").style.display = "";
        openUser(openedUser.id)
    });

    document.getElementById("instance-open-user").addEventListener("click", () => {
        document.getElementById("instance-window").style.display = "";
        openUser(openedUser.id)
    });

    document.getElementById("submit-generate").addEventListener("click", submitGenerate);
    document.getElementById("submit-restart").addEventListener("click", submitRestart);
    document.getElementById("submit-update").addEventListener("click", submitUpdate);
    window.addEventListener("resize", verifySize);

    parseOldData();
    showGlobalStatistics();
    fetchAllUser();

    openRequest = openRequestFunc;
    openUser = openUserFunc;
    showLogs = showLogsFunc;
    updateStoresUser = updateStoresUserFunc;
    changeStartUser = changeStartUserFunc;
    verifyChange = verifyChangeFunc;
    userSwitch = userSwitchFunc;
    onInputCheckout = onInputCheckoutFunc;
    onInputInstances = onInputInstancesFunc;
    openCheckout = openCheckoutFunc;
    openInstance = openInstanceFunc;
}

function parseOldData() {
    get(child(dbRef, "serverstats/" + activedata.split("/")[2].replaceAll(".", "-"))).then((snapshot) => {
        if (snapshot.exists()) {
            let keys = Object.keys(snapshot.val());
            let vals = Object.values(snapshot.val());
            for (const key in keys) {
                if (keys[key] === timeframeToday || vals[key].requests === undefined || vals[key].stats !== undefined) {
                    continue;
                }
                let data = parseData(vals[key].requests);
                set(ref(database, "serverstats/" + activedata.split("/")[2].replaceAll(".", "-") + "/" + keys[key] + "/stats"), data);
                remove(ref(database, "serverstats/" + activedata.split("/")[2].replaceAll(".", "-") + "/" + keys[key] + "/requests"));
                remove(ref(database, "serverstats/" + activedata.split("/")[2].replaceAll(".", "-") + "/" + keys[key] + "/logs"));
                set(ref(database, "old-serverstats/" + activedata.split("/")[2].replaceAll(".", "-") + "/" + keys[key]), vals[key]);
            }
        } else {
            error("no notparsed data available");
        }

        fetchOldStats();
    });
}

function fetchOldStats() {
    get(child(dbRef, "serverstats/" + activedata.split("/")[2].replaceAll(".", "-"))).then((snapshot) => {
        if (snapshot.exists()) {
            stats = snapshot.val();
        } else {
            error("no stats data available");
        }

        const r = ref(database, "serverstats/" + activedata.split("/")[2].replaceAll(".", "-") + "/" + timeframeToday + "/requests");
        onValue(r, (snapshot) => {
            let snapchotData = snapshot.val()
            let downtime = stats[timeframeToday].stats === undefined ? {} : stats[timeframeToday].stats.downtime;
            stats[timeframeToday].stats = parseData(snapchotData);
            stats[timeframeToday].stats.downtime = downtime

            showRequests(snapchotData).then();
            showData();
        });

        const ld = ref(database, "serverstats/" + activedata.split("/")[2].replaceAll(".", "-") + "/" + timeframeToday + "/logs");
        onValue(ld, (snapshot) => {
            animate("database");
            allLogs.database = snapshot.val();
            showLogsFunc("database").then();
        });

        const lm = ref(database, "serverstats/monitor-copped-inc-com/" + timeframeToday + "/logs");
        onValue(lm, (snapshot) => {
            animate("monitor");
            allLogs.monitor = snapshot.val();
            showLogsFunc("monitor").then();
        });

        const la = ref(database, "serverstats/service-copped-inc-com/" + timeframeToday + "/logs");
        onValue(la, (snapshot) => {
            animate("service");
            allLogs.service = snapshot.val();
            showLogsFunc("service").then();
        });

        const li = ref(database, "serverstats/instances-copped-inc-com/" + timeframeToday + "/logs");
        onValue(li, (snapshot) => {
            animate("instances");
            allLogs.instances = snapshot.val();
            showLogsFunc("instances").then();
        });

        fetchUptime("database");
        fetchUptime("monitor");
        fetchUptime("aio");
        fetchUptime("service");
        fetchUptime("instances");
    });
}

function fetchUptime(name) {
    let server = name + "-copped-inc-com"
    get(child(dbRef, "uptime/" + server + "/")).then((snapshot) => {
        if (snapshot.exists()) {
            let keys = Object.keys(snapshot.val());
            let vals = Object.values(snapshot.val());
            for (const key in keys) {
                if (stats[keys[key]] === undefined) stats[keys[key]] = {};
                if (stats[keys[key]].stats === undefined) stats[keys[key]].stats = {};
                if (stats[keys[key]].stats.downtime === undefined) stats[keys[key]].stats.downtime = {};
                stats[keys[key]].stats.downtime[name] = vals[key].downtime_duration;
            }
        } else {
            error(name + " has no uptime data");
        }
        checkUp(server, name);
    });
}

function checkUp(server, name) {
    const r = ref(database, "uptime/" + server + "/" + timeframeToday);
    onValue(r, (snapshot) => {
        if (!snapshot.exists()) return;
        let data = snapshot.val();
        if (data.status === 0) {
            document.getElementById(name + "icon").setAttribute("fill", "#00873E");
            document.getElementById(name + "status").innerHTML = "Online";
        } else if (data.status === 1) {
            document.getElementById(name + "icon").setAttribute("fill", "#c22315");
            document.getElementById(name + "status").innerHTML = "Offline";
        }
        if (stats[timeframeToday].stats.downtime === undefined) {
            stats[timeframeToday].stats.downtime = {};
        }
        stats[timeframeToday].stats.downtime[name] = data.downtime_duration;
        showData();
    });
}

async function showRequests(data) {
    while (requestLimit === 0) {
        await sleep(10)
    }
    let l = requestLimit;
    requestLimit = 0;

    let list = document.getElementById("requestlist");
    list.innerHTML = "";

    let requests = Object.values(data).sort((a, b) => new Date(b.start) - new Date(a.start));
    let input = document.getElementById("requestinput").value.toLowerCase();
    let search = input.length > 0;

    displayedRequests = {};

    for (const r in requests) {
        if (requests[r].query === undefined) requests[r].query = (requests[r].user + requests[r].err + requests[r].id + requests[r].method + requests[r].path + requests[r].status_code + states[requests[r].closing_state]).toLowerCase();
        if (list.childNodes.length >= l) {
            list.innerHTML += `
                <div class="more">
                    <h6 class="more" id="more-requests">load more</h6>
                </div>
            `

            document.getElementById("more-requests").addEventListener("click", () => {
                requestLimit += 50;
                list.innerHTML = ""
                showRequests(data);
            });
            break;
        }

        let request = requests[r];
        if (search && !request.query.includes(input)) continue;
        displayedRequests[request.id] = request;

        let add = request.status_code >= 500 ? " negative" : "";
        let time = getTimeString(request.start)

        list.innerHTML += `
            <div class="request" onclick="openRequest('${request.id}')">
                <div class="left">
                    <h6 class="statuscode${add}">${request.status_code}</h6>
                    <h4>${request.method}</h4>
                    <h4>${request.path}</h4>
                </div>
                <h6 class="time">${time}</h6>
            </div>
        `
    }

    if (!search) {
        document.getElementById("requestinput").removeEventListener("input", () => {});
        document.getElementById("requestinput").addEventListener("input", () => {
            requestLimit = 50;
            showRequests(data);
        });
    }

    if (list.innerHTML === "") {
        list.innerHTML = `
            <div class="nothing-list">
                <svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="24" height="24" viewBox="0 0 172 172"><g fill="none" fill-rule="nonzero" stroke="none" stroke-width="1" stroke-linecap="butt" stroke-linejoin="miter" stroke-miterlimit="10" stroke-dasharray="" stroke-dashoffset="0" font-family="none" font-size="none" style="mix-blend-mode: normal"><path d="M0,172v-172h172v172z" fill="none"></path><g fill="#FFFFFF73"><path d="M135.45,0c-6.97406,0.16125 -13.82719,2.33813 -19.6725,6.45l-28.4875,20.1025l7.8475,11.395l28.595,-20.21c4.79719,-3.35937 10.57531,-4.66281 16.34,-3.655c5.76469,1.00781 10.83063,4.23281 14.19,9.03c3.37281,4.79719 4.66281,10.57531 3.655,16.34c-1.00781,5.76469 -4.23281,10.81719 -9.03,14.19c-2.09625,1.47813 -23.08562,16.20563 -24.51,17.2c-5.46906,3.84313 -13.07469,5.76469 -22.4675,1.6125l-12.255,8.4925c13.06125,10.13188 30.46281,9.70188 42.57,1.1825c1.42438,-0.99437 22.54813,-15.84281 24.6175,-17.3075c7.78031,-5.4825 12.98063,-13.61219 14.62,-23.005c1.63938,-9.37937 -0.44344,-18.86625 -5.9125,-26.66c-5.4825,-7.79375 -13.71969,-12.98062 -23.1125,-14.62c-2.33812,-0.40312 -4.66281,-0.59125 -6.9875,-0.5375zM9.5675,3.3325c-2.67406,0.25531 -4.95844,2.05594 -5.83187,4.59563c-0.88688,2.55312 -0.20156,5.375 1.74687,7.22937l37.84,37.84c1.65281,2.05594 4.32688,2.98313 6.89344,2.39188c2.58,-0.59125 4.58219,-2.59344 5.17344,-5.17344c0.59125,-2.56656 -0.33594,-5.24062 -2.39188,-6.89344l-37.84,-37.84c-1.29,-1.35719 -3.07719,-2.13656 -4.945,-2.15c-0.215,-0.01344 -0.43,-0.01344 -0.645,0zM26.5525,87.29l-20.1025,28.595c-5.4825,7.79375 -7.55188,17.17313 -5.9125,26.5525c1.63938,9.39281 6.82625,17.63 14.62,23.1125c7.79375,5.4825 17.28063,7.55188 26.66,5.9125c9.37938,-1.63937 17.53594,-6.83969 23.005,-14.62c1.46469,-2.06937 16.31313,-23.19312 17.3075,-24.6175c8.51938,-12.10719 8.94938,-29.40125 -1.1825,-42.4625l-8.4925,12.1475c4.15219,9.39281 2.23063,16.99844 -1.6125,22.4675c-0.99437,1.42438 -15.72187,22.52125 -17.2,24.6175c-3.37281,4.79719 -8.42531,7.91469 -14.19,8.9225c-5.76469,1.00781 -11.54281,-0.28219 -16.34,-3.655c-4.79719,-3.35937 -8.02219,-8.42531 -9.03,-14.19c-1.00781,-5.76469 0.29563,-11.54281 3.655,-16.34l20.21,-28.595z"></path></g></g></svg>
                Nothing here
            </div>
        `
    }

    requestLimit = l;
}

async function showLogsFunc(name) {
    if (name !== activelog) return;
    let data = allLogs[name];
    if (data === null) data = {};

    while (logLimit === 0) {
        await sleep(10)
    }
    let l = logLimit;
    logLimit = 0;

    let list = document.getElementById("loglist");
    list.innerHTML = "";

    let logs = Object.values(data).sort((a, b) => new Date(b.time) - new Date(a.time));
    let input = document.getElementById("loginput").value.toLowerCase();
    let search = input.length > 0;

    for (const r in logs) {
        if (logs[r].content === undefined) continue;
        if (logs[r].query === undefined) logs[r].query = (logs[r].id + logs[r].content.join("") + states[logs[r].state]).toLowerCase();
        if (list.childNodes.length >= l) {
            list.innerHTML += `
                <div class="more">
                    <h6 class="more" id="more-logs">load more</h6>
                </div>
            `

            document.getElementById("more-logs").addEventListener("click", () => {
                logLimit += 50;
                list.innerHTML = ""
                showLogsFunc(name);
            });
            break;
        }

        let log = logs[r];
        if (search) {
            if (search && !log.query.includes(input)) continue;
        }

        let add = log.state === 1 ? " negative" : "";
        let time = getTimeString(log.time);

        let click = function () {
            if (log.ref === "") {
                return `class="log notclickable"`;
            } else {
                return `class="log" onclick="openRequest('${log.ref}')"`;
            }
        }();

        list.innerHTML += `
            <div ${click}>
                <div class="left">
                    <h6 class="statuscode${add}">${states[log.state]}</h6>
                    <h4>${log.content.join(" ")}</h4>
                </div>
                <h6 class="time">${time}</h6>
            </div>
        `
    }

    if (!search) {
        document.getElementById("loginput").removeEventListener("input", () => {});
        document.getElementById("loginput").addEventListener("input", () => {
            logLimit = 50;
            showLogsFunc(name);
        });
    }

    if (list.innerHTML === "") {
        list.innerHTML = `
            <div class="nothing-list">
                <svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="24" height="24" viewBox="0 0 172 172"><g fill="none" fill-rule="nonzero" stroke="none" stroke-width="1" stroke-linecap="butt" stroke-linejoin="miter" stroke-miterlimit="10" stroke-dasharray="" stroke-dashoffset="0" font-family="none" font-size="none" style="mix-blend-mode: normal"><path d="M0,172v-172h172v172z" fill="none"></path><g fill="#FFFFFF73"><path d="M135.45,0c-6.97406,0.16125 -13.82719,2.33813 -19.6725,6.45l-28.4875,20.1025l7.8475,11.395l28.595,-20.21c4.79719,-3.35937 10.57531,-4.66281 16.34,-3.655c5.76469,1.00781 10.83063,4.23281 14.19,9.03c3.37281,4.79719 4.66281,10.57531 3.655,16.34c-1.00781,5.76469 -4.23281,10.81719 -9.03,14.19c-2.09625,1.47813 -23.08562,16.20563 -24.51,17.2c-5.46906,3.84313 -13.07469,5.76469 -22.4675,1.6125l-12.255,8.4925c13.06125,10.13188 30.46281,9.70188 42.57,1.1825c1.42438,-0.99437 22.54813,-15.84281 24.6175,-17.3075c7.78031,-5.4825 12.98063,-13.61219 14.62,-23.005c1.63938,-9.37937 -0.44344,-18.86625 -5.9125,-26.66c-5.4825,-7.79375 -13.71969,-12.98062 -23.1125,-14.62c-2.33812,-0.40312 -4.66281,-0.59125 -6.9875,-0.5375zM9.5675,3.3325c-2.67406,0.25531 -4.95844,2.05594 -5.83187,4.59563c-0.88688,2.55312 -0.20156,5.375 1.74687,7.22937l37.84,37.84c1.65281,2.05594 4.32688,2.98313 6.89344,2.39188c2.58,-0.59125 4.58219,-2.59344 5.17344,-5.17344c0.59125,-2.56656 -0.33594,-5.24062 -2.39188,-6.89344l-37.84,-37.84c-1.29,-1.35719 -3.07719,-2.13656 -4.945,-2.15c-0.215,-0.01344 -0.43,-0.01344 -0.645,0zM26.5525,87.29l-20.1025,28.595c-5.4825,7.79375 -7.55188,17.17313 -5.9125,26.5525c1.63938,9.39281 6.82625,17.63 14.62,23.1125c7.79375,5.4825 17.28063,7.55188 26.66,5.9125c9.37938,-1.63937 17.53594,-6.83969 23.005,-14.62c1.46469,-2.06937 16.31313,-23.19312 17.3075,-24.6175c8.51938,-12.10719 8.94938,-29.40125 -1.1825,-42.4625l-8.4925,12.1475c4.15219,9.39281 2.23063,16.99844 -1.6125,22.4675c-0.99437,1.42438 -15.72187,22.52125 -17.2,24.6175c-3.37281,4.79719 -8.42531,7.91469 -14.19,8.9225c-5.76469,1.00781 -11.54281,-0.28219 -16.34,-3.655c-4.79719,-3.35937 -8.02219,-8.42531 -9.03,-14.19c-1.00781,-5.76469 0.29563,-11.54281 3.655,-16.34l20.21,-28.595z"></path></g></g></svg>
                Nothing here
            </div>
        `
    }

    logLimit = l;
}

function getTimeString(start) {
    let time = new Date(start);
    let hours = time.getHours().toString().length === 1 ? "0" + time.getHours() : time.getHours();
    let minutes = time.getMinutes().toString().length === 1 ? "0" + time.getMinutes() : time.getMinutes();
    let seconds = time.getSeconds().toString().length === 1 ? "0" + time.getSeconds() : time.getSeconds();
    return `${hours}:${minutes}:${seconds}`;
}

function openRequestFunc(id) {
    let request = displayedRequests[id];
    if (request === undefined || request.logs.length === 0) {
        get(child(dbRef, "serverstats/" + activedata.split("/")[2].replaceAll(".", "-") + "/" + timeframeToday + "/requests/" + id)).then((snapshot) => {
            if (!snapshot.exists()) return;
            displayedRequests[id].logs = snapshot.val().logs;
            openRequestFunc(id);
        });
        return
    }

    try {
        request.request_body = atob(request.request_body);
        request.request_body = JSON.stringify(JSON.parse(request.request_body), null, 4);
    } catch (e) {
        request.request_body = request.request_body === undefined ? "" : request.request_body;
    }

    try {
        request.response_body = atob(request.response_body);
        request.response_body = JSON.stringify(JSON.parse(request.response_body), null, 4);
    } catch (e) {
        request.response_body = request.response_body === undefined ? "" : request.response_body;
    }

    document.getElementById("details-path-short").innerHTML = "/" + request.path.toString().split("/")[1];
    document.getElementById("details-path").innerHTML = request.path;
    document.getElementById("details-id").innerHTML = request.id;
    document.getElementById("details-user").innerHTML = request.user;
    document.getElementById("details-method").innerHTML = request.method;
    document.getElementById("details-status").innerHTML = request.status_code;
    document.getElementById("details-time").innerHTML = new Date(request.start).toLocaleString();
    document.getElementById("details-duration").innerHTML = (request.duration * 1e-6).toFixed(2) + "ms";
    document.getElementById("details-state").innerHTML = states[request.closing_state];
    document.getElementById("details-err").innerHTML = request.err;
    document.getElementById("details-requestbody").innerHTML = request.request_body;
    document.getElementById("details-responsebody").innerHTML = request.response_body;

    document.getElementById("details-load").style.display = "flex";
    document.getElementById("details").style.display = "flex";
    document.getElementById("details-area").style.display = "none";

    document.getElementById("error-attribute").style.display = request.err === "" ? "none" : "flex";
    document.getElementById("requestbody-attribute").style.display = request.request_body === "" ? "none" : "";
    document.getElementById("responsebody-attribute").style.display = request.response_body === "" ? "none" : "";

    document.getElementById("details-area").innerHTML = "";
    fetchLogs(request.logs);
}

function copy(id) {
    return function () {
        copyAttribute(id)
    }
}

function copyAttribute(name) {
    let element = document.getElementById(name);
    let text = element.children[1].innerHTML;
    navigator.clipboard.writeText(element.children[1].innerHTML);

    element.children[1].innerHTML = "Copied!";
    sleep(1000).then(() => {
        element.children[1].innerHTML = text;
    })
}

function switchView() {
    let details = document.getElementById("details");
    let window = details.children[0];
    let state = window.children[1].style.display === "none";
    for (let e of window.children) {
        if (e.classList.contains("attributes")) {
            if (e.classList.contains("hasarea")) {
                e.style.display = state ? "" : "flex";
            } else {
                if (e.id === "error-attribute" && document.getElementById("details-err").innerHTML === "") {
                    continue;
                }

                e.style.display = state ? "" : "none";
            }
        }
    }
}

function verifySize() {
    if (document.getElementById("details").children[0].children[1].style.display === "none" && window.innerHeight > 1000) {
        switchView();
    }
}

function fetchLogs(logs) {
    get(child(dbRef, "serverstats/" + activedata.split("/")[2].replaceAll(".", "-") + "/" + timeframeToday + "/logs/" + logs[0])).then((snapshot) => {
        if (snapshot.exists()) {
            let log = snapshot.val();
            document.getElementById("details-area").innerHTML += `[${states[log.state]}] [${new Date(log.time).toLocaleTimeString()}] ${log.content.join(" ")}\n\n`;
            logs.shift();
            if (logs.length > 0) fetchLogs(logs);
            else {
                document.getElementById("details-load").style.display = "none";
                document.getElementById("details-area").style.display = "block";
                document.getElementById("details-area").innerHTML = document.getElementById("details-area").innerHTML.slice(0, -2);
            }
        } else {
            error("No data available");
        }
    });
}

function showData() {
    let from = timeframeFrom;
    let to = timeframeTo;

    let date = new Date(new Date().setFullYear(parseInt(from.split("-")[0]), months.indexOf(from.split("-")[1]), from.split("-")[2]));
    let toDate = new Date(new Date().setFullYear(parseInt(to.split("-")[0]), months.indexOf(to.split("-")[1]), to.split("-")[2]));

    if (date.toString() === toDate.toString()) document.getElementById("newtimeframe").innerHTML = date.getDate() + " " + months[date.getMonth()] + " " + date.getFullYear();
    else document.getElementById("newtimeframe").innerHTML = date.getDate() + " " + months[date.getMonth()] + " " + date.getFullYear() + " - " + toDate.getDate() + " " + months[toDate.getMonth()] + " " + toDate.getFullYear();

    let labels = [];

    let allData = {
        capLatHour: [],
        capRpmHour: [],
        errorRpmHour: [],
        reqLatHour: [],
        reqRpmHour: [],
        captchalatency: 0.0,
        captchaload: 0.0,
        errorload: 0.0,
        requestlatency: 0.0,
        requestload: 0.0,
        solvedCaptcha: 0,
        downtime: {
            database: 0,
            monitor: 0,
            aio: 0,
            service: 0,
        },
    };

    let oldAllData = {
        capLatHour: [],
        capRpmHour: [],
        errorRpmHour: [],
        reqLatHour: [],
        reqRpmHour: [],
        captchalatency: 0.0,
        captchaload: 0.0,
        errorload: 0.0,
        requestlatency: 0.0,
        requestload: 0.0,
        solvedCaptcha: 0,
        downtime: {
            database: 0,
            monitor: 0,
            aio: 0,
            service: 0,
        },
    }

    let today = false;
    let days = 1;
    while (date.getFullYear() + "-" + date.toLocaleDateString('en-US', { month: 'long'}) + "-" + date.getDate() !== to) {
        if (!today) {
            if (date.getFullYear() + "-" + date.toLocaleDateString('en-US', { month: 'long'}) + "-" + date.getDate() === timeframeToday) today = true;
            allData = appendData(allData, date, today);
        }
        date.setDate(date.getDate() + 1);
        days++;
    }

    if (!today) {
        if (date.getFullYear() + "-" + date.toLocaleDateString('en-US', { month: 'long'}) + "-" + date.getDate() === timeframeToday) today = true;
        allData = appendData(allData, date, today);
    }

    if (allData.captchalatency !== 0) allData.captchalatency = allData.capLatHour.reduce((a, b) => a + b,  0) / allData.capLatHour.filter((x) => x !== 0).length;
    allData.captchaload /= days;
    allData.errorload /= days;
    if (allData.requestlatency !== 0) allData.requestlatency = allData.reqLatHour.reduce((a, b) => a + b,  0) / allData.reqLatHour.filter((x) => x !== 0).length;
    allData.requestload /= days;

    date.setDate(date.getDate() - (days * 2) + 1);
    let oldTimeframe = date.getDate() + " " + months[date.getMonth()] + " " + date.getFullYear();
    while (date.getFullYear() + "-" + date.toLocaleDateString('en-US', { month: 'long'}) + "-" + date.getDate() !== from) {
        labels = appendLabels(labels, date, days);
        oldAllData = appendData(oldAllData, date, false);
        date.setDate(date.getDate() + 1);
    }
    date.setDate(date.getDate() - 1);

    let oldToTimeframe = date.getDate() + " " + months[date.getMonth()] + " " + date.getFullYear();
    if (oldTimeframe !== oldToTimeframe) oldTimeframe += " - " + oldToTimeframe;
    document.getElementById("oldtimeframe").innerHTML = oldTimeframe;

    if (oldAllData.captchalatency !== 0) oldAllData.captchalatency = oldAllData.capLatHour.reduce((a, b) => a + b,  0) / oldAllData.capLatHour.filter((x) => x !== 0).length;
    oldAllData.captchaload /= days;
    oldAllData.errorload /= days;
    if (oldAllData.requestlatency !== 0) oldAllData.requestlatency = oldAllData.reqLatHour.reduce((a, b) => a + b,  0) / oldAllData.reqLatHour.filter((x) => x !== 0).length;
    oldAllData.requestload /= days;

    chartData.oldrequestlatency = oldAllData.requestlatency.toFixed(2);
    chartData.oldrequestload = oldAllData.requestload.toFixed(2);
    chartData.olderrorload = oldAllData.errorload.toFixed(2);
    chartData.oldSolvedCaptcha = oldAllData.solvedCaptcha;
    chartData.oldcaptchaload = oldAllData.captchaload.toFixed(2);
    chartData.oldcaptchalatency = oldAllData.captchalatency.toFixed(2);

    chartData.requestlatency = allData.requestlatency.toFixed(2);
    chartData.requestload = allData.requestload.toFixed(2);
    chartData.errorload = allData.errorload.toFixed(2);
    chartData.solvedCaptcha = allData.solvedCaptcha;
    chartData.captchaload = allData.captchaload.toFixed(2);
    chartData.captchalatency = allData.captchalatency.toFixed(2);

    mouseLeave("requestload");
    mouseLeave("errorload");
    mouseLeave("requestlatency");
    mouseLeave("captchaload");
    mouseLeave("captchalatency");

    setDevChart("requestload", labels, allData.reqRpmHour, oldAllData.reqRpmHour)
    setDevChart("errorload", labels, allData.errorRpmHour, oldAllData.errorRpmHour)
    setDevChart("requestlatency", labels, allData.reqLatHour, oldAllData.reqLatHour)
    setDevChart("captchaload", labels, allData.capRpmHour, oldAllData.capRpmHour)
    setDevChart("captchalatency", labels, allData.capLatHour, oldAllData.capLatHour)

    document.getElementById("newcaptchatotal").innerHTML = allData.solvedCaptcha.toString();
    document.getElementById("oldcaptchatotal").innerHTML = oldAllData.solvedCaptcha.toString();
    document.getElementById("newtotalcaptchacost").innerHTML = (allData.solvedCaptcha / 1000 * 1.0).toFixed(2) + "€";
    document.getElementById("oldtotalcaptchacost").innerHTML = (oldAllData.solvedCaptcha / 1000 * 1.0).toFixed(2) + "€";

    setRatio("totalcaptchasolved", allData.solvedCaptcha, oldAllData.solvedCaptcha);
    setRatio("totalcaptchacosts", allData.solvedCaptcha / 1000 * 1.0, oldAllData.solvedCaptcha / 1000 * 1.0);

    let millis = allData.reqLatHour.length * 3.6e6;
    if (allData.downtime.database !== 0) document.getElementById("databaseratio").innerHTML = 100 - (allData.downtime.database / millis * 100).toFixed(2) + "%";
    else document.getElementById("databaseratio").innerHTML = "100.00%";
    if (allData.downtime.monitor !== 0) document.getElementById("monitorratio").innerHTML = 100 - (allData.downtime.monitor / millis * 100).toFixed(2) + "%";
    else document.getElementById("monitorratio").innerHTML = "100.00%";
    if (allData.downtime.aio !== 0) document.getElementById("aioratio").innerHTML = 100 - (allData.downtime.aio / millis * 100).toFixed(2) + "%";
    else document.getElementById("aioratio").innerHTML = "100.00%";
    if (allData.downtime.aio !== 0) document.getElementById("serviceratio").innerHTML = 100 - (allData.downtime.aio / millis * 100).toFixed(2) + "%";
    else document.getElementById("serviceratio").innerHTML = "100.00%";
}

function appendLabels(labels, date, days) {
    let dateLabel = date.getDate() + ". " + date.toLocaleDateString('en-US', { month: 'long'}) + " " + date.getFullYear();
    date = new Date(new Date().setDate(date.getDate() + days));
    let dateLabel2 = date.getDate() + ". " + date.toLocaleDateString('en-US', { month: 'long'}) + " " + date.getFullYear();

    for (let i = 0; i < 24; i++) {
        labels.push(dateLabel2 + " " + i + ":00" + ";" + dateLabel + " " + i + ":00");
    }
    return labels;
}

function appendData(allData, date, today) {
    let timeframe = date.getFullYear() + "-" + date.toLocaleDateString('en-US', { month: 'long'}) + "-" + date.getDate();
    let data = stats[timeframe] === undefined ? {} : stats[timeframe].stats;

    if (data.capLatHour == null) {
        data.capLatHour = new Array(24).fill(0);
        data.capRpmHour = new Array(24).fill(0);
        data.errorRpmHour = new Array(24).fill(0);
        data.reqLatHour = new Array(24).fill(0);
        data.reqRpmHour = new Array(24).fill(0);
        data.captchalatency = 0.0;
        data.captchaload = 0.0;
        data.errorload = 0.0;
        data.requestlatency = 0.0;
        data.requestload = 0.0;
        data.solvedCaptcha = 0.0;
    }

    if (today) {
        data.capLatHour = data.capLatHour.slice(0, new Date().getHours() + 1);
        data.capRpmHour = data.capRpmHour.slice(0, new Date().getHours() + 1);
        data.errorRpmHour = data.errorRpmHour.slice(0, new Date().getHours() + 1);
        data.reqLatHour = data.reqLatHour.slice(0, new Date().getHours() + 1);
        data.reqRpmHour = data.reqRpmHour.slice(0, new Date().getHours() + 1);
    }

    allData.capLatHour = allData.capLatHour.concat(data.capLatHour);
    allData.capRpmHour = allData.capRpmHour.concat(data.capRpmHour);
    allData.errorRpmHour = allData.errorRpmHour.concat(data.errorRpmHour);
    allData.reqLatHour = allData.reqLatHour.concat(data.reqLatHour);
    allData.reqRpmHour = allData.reqRpmHour.concat(data.reqRpmHour);

    allData.captchalatency += parseFloat(data.captchalatency);
    allData.captchaload += parseFloat(data.captchaload);
    allData.errorload += parseFloat(data.errorload);
    allData.requestlatency += parseFloat(data.requestlatency);
    allData.requestload += parseFloat(data.requestload);
    allData.solvedCaptcha += parseInt(data.solvedCaptcha);
    if (data.downtime != null) {
        if (data.downtime.database != null) allData.downtime.database += data.downtime.database;
        if (data.downtime.monitor != null) allData.downtime.monitor += data.downtime.monitor;
        if (data.downtime.aio != null) allData.downtime.aio += data.downtime.aio;
        if (data.downtime.service != null) allData.downtime.aio += data.downtime.service;
    }

    return allData;
}

function parseData(data) {
    let dataArr = Object.values(data)
    let nonCaptchaReq = dataArr.filter((x) => !(x.path.includes("captcha") && !x.path.includes("challenge")));
    let captchaReq = dataArr.filter((x) => (x.path.includes("captcha") && !x.path.includes("challenge") && !x.path.includes("preharvest")));

    let parsedData = {
        reqLatHour: new Array(24).fill(0),
        reqRpmHour: new Array(24).fill(0),
        errorRpmHour: new Array(24).fill(0),
        capLatHour: new Array(24).fill(0),
        capRpmHour: new Array(24).fill(0),
    };

    let errRpm = 0;
    let lat = 0;
    let reqCap = 0;
    let solvedCap = 0;
    let capLat = 0;

    for (const key of nonCaptchaReq) {
        let hour = new Date(key.start).getHours();
        parsedData.reqLatHour[hour] += key.duration * 1e-6;
        parsedData.reqRpmHour[hour] += 1;
        if (key.closing_state !== 0) parsedData.errorRpmHour[hour] += 1;
    }

    for (const key of captchaReq) {
        if (key.method === "GET" || key.method === "OPTIONS") {
            let hour = new Date(key.start).getHours();
            parsedData.capLatHour[hour] += key.duration * 1e-6;
            parsedData.capRpmHour[hour] += 1;
            reqCap++;
        } else solvedCap++;
    }

    for (const hour in parsedData.reqRpmHour) {
        lat += parsedData.reqLatHour[hour];
        errRpm += parsedData.errorRpmHour[hour];

        let hlat = parsedData.reqLatHour[hour] / parsedData.reqRpmHour[hour];
        if (!isNaN(hlat)) parsedData.reqLatHour[hour] = hlat;

        parsedData.reqRpmHour[hour] = parsedData.reqRpmHour[hour] / 60;
        parsedData.errorRpmHour[hour] = parsedData.errorRpmHour[hour] / 60;

        capLat += parsedData.capLatHour[hour];

        let clat = parsedData.capLatHour[hour] / parsedData.capRpmHour[hour];
        if (!isNaN(clat)) parsedData.capLatHour[hour] = clat;

        parsedData.capRpmHour[hour] = parsedData.capRpmHour[hour] / 60;
    }

    parsedData.requestlatency = (lat / nonCaptchaReq.length).toFixed(2);
    parsedData.requestload = (nonCaptchaReq.length / (24 * 60)).toFixed(2);
    parsedData.errorload = (errRpm / (24 * 60)).toFixed(2);
    parsedData.solvedCaptcha = solvedCap;
    parsedData.captchaload = (reqCap / (24 * 60)).toFixed(2);
    parsedData.captchalatency = isNaN(capLat / reqCap) ? 0 : (capLat / reqCap).toFixed(2);

    return parsedData;
}

function animate(server) {
    let element = document.getElementById(server + "icon");
    element.style.opacity = "0";
    sleep(50).then(() => {
        element.style.opacity = "1";
    });
}

function toggleGenerateWindow() {
    let element = document.getElementById("generate-keys-window");
    element.style.display = element.style.display === "" ? "flex" : "";
}

function submitGenerate() {
    let field = document.getElementById("copy-link");
    let loading = document.getElementById("copy-link-loading");
    let link = document.getElementById("purchase-link");

    field.style.display = "flex";
    loading.style.display = "block";
    link.style.display = "none";

    fetch(activedata + "purchase", {
        method: "POST",
        credentials: "include",
        body: JSON.stringify({
            "plan": parseInt(document.getElementById('plan-dropdown-input').getAttribute('data-selected')),
            "stock": parseInt(document.getElementById("usage-limit").value),
            "instance_limit": parseInt(document.getElementById("instance-limit").value),
        })
    }).then(resp => {
        if (resp.status !== 201) {
            throw new Error("Statuscode " + resp.status);
        }

        return resp.json();
    }).then(data => {
        link.innerHTML = data.link;
        loading.style.display = "none";
        link.style.display = "block";
    }).catch(err => {
        field.style.display = "";
        error(err);
    });
}

function toggleRestartWindow() {
    let element = document.getElementById("monitor-restart-window");
    element.style.display = element.style.display === "" ? "flex" : "";
}

function submitRestart() {
    let loading = document.getElementById("restart-load");
    let field = document.getElementById("restart-field");

    loading.style.display = "flex";
    field.style.display = "none";

    fetch(activedata + "monitor/restart/" + document.getElementById("restart-dropdown-input").getAttribute("data-selected"), {
        method: "GET",
        credentials: "include",
    }).then(resp => {
        if (resp.status !== 200) {
            throw new Error("Statuscode " + resp.status);
        }

        field.style.display = "";
        loading.style.display = "";
        toggleRestartWindow();
    }).catch(err => {
        field.style.display = "";
        loading.style.display = "";
        error(err);
    });
}

function showGlobalStatistics() {
    const lm = ref(database, "userstats/global");
    onValue(lm, (snapshot) => {
        let global = snapshot.val();
        document.getElementById("monitor-pings").innerHTML = global.monitor.pings + " Pings";
        document.getElementById("monitor-products").innerHTML = global.monitor.products + " Products";
    });
}

function fetchAllUser() {
    const r = ref(database, "userstats/global/user/updates");
    onValue(r, () => {
        fetch(activedata + "user", {
            method: "GET",
            credentials: "include",
        }).then(resp => {
            if (resp.status === 401) {
                removeCookie("authorization");
                window.location.href = activedata + "login";
                return
            } if (resp.status !== 200) {
                throw new Error("Statuscode " + resp.status);
            }
            return resp.json();
        }).then(data => {
            displayUser(data);
            if (password === "") {
                console.log("Password is empty");
                onLoad();
            }
        }).catch(err => error(err));
    });
}

function displayUser(data) {
    let list = document.getElementById("user-stats-list");
    list.innerHTML = "";
    if (data.length === 0) {
        list.innerHTML = `
            <div class="nothing-list flex-parent">
                <svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="24" height="24" viewBox="0 0 172 172"><g fill="none" fill-rule="nonzero" stroke="none" stroke-width="1" stroke-linecap="butt" stroke-linejoin="miter" stroke-miterlimit="10" stroke-dasharray="" stroke-dashoffset="0" font-family="none" font-size="none" style="mix-blend-mode: normal"><path d="M0,172v-172h172v172z" fill="none"></path><g fill="#FFFFFF73"><path d="M135.45,0c-6.97406,0.16125 -13.82719,2.33813 -19.6725,6.45l-28.4875,20.1025l7.8475,11.395l28.595,-20.21c4.79719,-3.35937 10.57531,-4.66281 16.34,-3.655c5.76469,1.00781 10.83063,4.23281 14.19,9.03c3.37281,4.79719 4.66281,10.57531 3.655,16.34c-1.00781,5.76469 -4.23281,10.81719 -9.03,14.19c-2.09625,1.47813 -23.08562,16.20563 -24.51,17.2c-5.46906,3.84313 -13.07469,5.76469 -22.4675,1.6125l-12.255,8.4925c13.06125,10.13188 30.46281,9.70188 42.57,1.1825c1.42438,-0.99437 22.54813,-15.84281 24.6175,-17.3075c7.78031,-5.4825 12.98063,-13.61219 14.62,-23.005c1.63938,-9.37937 -0.44344,-18.86625 -5.9125,-26.66c-5.4825,-7.79375 -13.71969,-12.98062 -23.1125,-14.62c-2.33812,-0.40312 -4.66281,-0.59125 -6.9875,-0.5375zM9.5675,3.3325c-2.67406,0.25531 -4.95844,2.05594 -5.83187,4.59563c-0.88688,2.55312 -0.20156,5.375 1.74687,7.22937l37.84,37.84c1.65281,2.05594 4.32688,2.98313 6.89344,2.39188c2.58,-0.59125 4.58219,-2.59344 5.17344,-5.17344c0.59125,-2.56656 -0.33594,-5.24062 -2.39188,-6.89344l-37.84,-37.84c-1.29,-1.35719 -3.07719,-2.13656 -4.945,-2.15c-0.215,-0.01344 -0.43,-0.01344 -0.645,0zM26.5525,87.29l-20.1025,28.595c-5.4825,7.79375 -7.55188,17.17313 -5.9125,26.5525c1.63938,9.39281 6.82625,17.63 14.62,23.1125c7.79375,5.4825 17.28063,7.55188 26.66,5.9125c9.37938,-1.63937 17.53594,-6.83969 23.005,-14.62c1.46469,-2.06937 16.31313,-23.19312 17.3075,-24.6175c8.51938,-12.10719 8.94938,-29.40125 -1.1825,-42.4625l-8.4925,12.1475c4.15219,9.39281 2.23063,16.99844 -1.6125,22.4675c-0.99437,1.42438 -15.72187,22.52125 -17.2,24.6175c-3.37281,4.79719 -8.42531,7.91469 -14.19,8.9225c-5.76469,1.00781 -11.54281,-0.28219 -16.34,-3.655c-4.79719,-3.35937 -8.02219,-8.42531 -9.03,-14.19c-1.00781,-5.76469 0.29563,-11.54281 3.655,-16.34l20.21,-28.595z"></path></g></g></svg>
                Nothing here
            </div>
        `
        return
    }

    for (let user of data) {
        let url = `https://cdn.discordapp.com/avatars/${user.id}/${user.picture}`;
        if (checkedUrls[url] === undefined) checkedUrls[url] = urlExists(url) ? url : "https://cdn.discordapp.com/embed/avatars/0.png";
        url = checkedUrls[url]

        if (user?.checkouts === undefined) user.checkouts = [];

        list.innerHTML += `
            <div class="user-element" onclick="openUser('${user.id}')">
                <img src="${url}" alt="User">
                <div class="right">
                    <h4>${user.name}</h4>
                    <div class="attributes">
                        <h6 class="${user.state!=="Running"?"":"online"}">${user.state || "Offline"}</h6>
                        <h6>${user.checkouts?.length || 0} Checkouts</h6>
                        <h6>${user.instances?.length || 0} Clients</h6>
                    </div>
                </div>
            </div>
        `
    }

    document.getElementById("stats-all-user").innerHTML = data.length + " User";
    document.getElementById("stats-online-user").innerHTML = data.filter(x => x.state === "Running").length + " User Online";
    document.getElementById("overall-friends").innerHTML = data.filter(x => x.plan === 1).length + " Friends";
    document.getElementById("overall-basic").innerHTML = data.filter(x => x.plan === 2).length + " Basic";
    document.getElementById("overall-developer").innerHTML = data.filter(x => x.plan === 3).length + " Developer";
    document.getElementById("online-user").innerHTML = data.filter(x => x.state === "Running").length + " User";
    document.getElementById("online-clients").innerHTML = data.filter(x => x.state === "Running").reduce((a, b) => {
        return a + (b?.instances || []).filter(x => x.status === "Running").length
    }, 0) + " Clients";

    globalCheckouts = data.reduce((a, b) => a.concat(function (){
        b.checkouts.forEach(x => {
            x.user = b.id;
            x.username = b.name;
            x.query = (b.id + b.name + x.date + x.name + x.link + x.image + x.store + x.size + x.price.toFixed(2)).replaceAll(" ", "").toLowerCase();
        });
        return b.checkouts;
    }()), []);

    globalInstances = data.reduce((a, b) => a.concat(function (){
        if (b?.instances === undefined) return [];
        b.instances.forEach(x => {
            x.user = b.id;
            x.username = b.name;
            x.query = (b.id + b.name + x.id + x.price.toFixed(2) + x.provider + x.region + x.status + x.task_max).replaceAll(" ", "").toLowerCase();
        });
        return b.instances;
    }()), []);

    document.getElementById("global-checkouts").innerHTML = globalCheckouts.length + " Checkouts";
    document.getElementById("product-checkouts").innerHTML = Object.keys(function () {
        let products = {};
        globalCheckouts.forEach(x => {
            products[x.name] = products[x.name]
        });

        return products;
    }()).length + " Products";

    parseCheckouts(globalCheckouts, "global");
    parseInstances(globalInstances, "global");
}

function parseCheckouts(checkouts, name) {
    let checkoutlist = document.getElementById(name + "-checkout-list");
    checkoutlist.innerHTML = "";

    if (checkouts.length === 0) {
        checkoutlist.innerHTML = `
            <div class="nothing-list">
                <svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="24" height="24" viewBox="0 0 172 172"><g fill="none" fill-rule="nonzero" stroke="none" stroke-width="1" stroke-linecap="butt" stroke-linejoin="miter" stroke-miterlimit="10" stroke-dasharray="" stroke-dashoffset="0" font-family="none" font-size="none" style="mix-blend-mode: normal"><path d="M0,172v-172h172v172z" fill="none"></path><g fill="#FFFFFF73"><path d="M135.45,0c-6.97406,0.16125 -13.82719,2.33813 -19.6725,6.45l-28.4875,20.1025l7.8475,11.395l28.595,-20.21c4.79719,-3.35937 10.57531,-4.66281 16.34,-3.655c5.76469,1.00781 10.83063,4.23281 14.19,9.03c3.37281,4.79719 4.66281,10.57531 3.655,16.34c-1.00781,5.76469 -4.23281,10.81719 -9.03,14.19c-2.09625,1.47813 -23.08562,16.20563 -24.51,17.2c-5.46906,3.84313 -13.07469,5.76469 -22.4675,1.6125l-12.255,8.4925c13.06125,10.13188 30.46281,9.70188 42.57,1.1825c1.42438,-0.99437 22.54813,-15.84281 24.6175,-17.3075c7.78031,-5.4825 12.98063,-13.61219 14.62,-23.005c1.63938,-9.37937 -0.44344,-18.86625 -5.9125,-26.66c-5.4825,-7.79375 -13.71969,-12.98062 -23.1125,-14.62c-2.33812,-0.40312 -4.66281,-0.59125 -6.9875,-0.5375zM9.5675,3.3325c-2.67406,0.25531 -4.95844,2.05594 -5.83187,4.59563c-0.88688,2.55312 -0.20156,5.375 1.74687,7.22937l37.84,37.84c1.65281,2.05594 4.32688,2.98313 6.89344,2.39188c2.58,-0.59125 4.58219,-2.59344 5.17344,-5.17344c0.59125,-2.56656 -0.33594,-5.24062 -2.39188,-6.89344l-37.84,-37.84c-1.29,-1.35719 -3.07719,-2.13656 -4.945,-2.15c-0.215,-0.01344 -0.43,-0.01344 -0.645,0zM26.5525,87.29l-20.1025,28.595c-5.4825,7.79375 -7.55188,17.17313 -5.9125,26.5525c1.63938,9.39281 6.82625,17.63 14.62,23.1125c7.79375,5.4825 17.28063,7.55188 26.66,5.9125c9.37938,-1.63937 17.53594,-6.83969 23.005,-14.62c1.46469,-2.06937 16.31313,-23.19312 17.3075,-24.6175c8.51938,-12.10719 8.94938,-29.40125 -1.1825,-42.4625l-8.4925,12.1475c4.15219,9.39281 2.23063,16.99844 -1.6125,22.4675c-0.99437,1.42438 -15.72187,22.52125 -17.2,24.6175c-3.37281,4.79719 -8.42531,7.91469 -14.19,8.9225c-5.76469,1.00781 -11.54281,-0.28219 -16.34,-3.655c-4.79719,-3.35937 -8.02219,-8.42531 -9.03,-14.19c-1.00781,-5.76469 0.29563,-11.54281 3.655,-16.34l20.21,-28.595z"></path></g></g></svg>
                Nothing here
            </div>
        `;
    }

    checkouts.sort(function(a, b) {
        let c = new Date(a.date);
        let d = new Date(b.date);
        return c-d;
    });

    for (let checkout of checkouts) {
        checkoutlist.innerHTML += `
            <div class="checkout clickable" onclick="openCheckout('${checkout.date}', '${name}')">
                <img src="${checkout.image}" alt="">
                <div class="description">
                    <div class="split">
                        <h4>${checkout.name}</h4>
                        <h5 class="date">${new Date(checkout.date).toLocaleDateString()}</h5>
                    </div>
                    <div class="attributes">
                        <h6>${checkout.username}</h6>
                        <h6>Size: ${checkout.size}</h6>
                        <h6>Paid: ${checkout.price}€</h6>
                    </div>
                </div>
            </div>
        </div>
        `;
    }
}

function openCheckoutFunc(date, name) {
    document.getElementById("checkout-close").onclick = () => {
        document.getElementById("checkout-window").style.display = "";
        if (name === "user") openUser(openedUser.id);
    }

    document.getElementById("checkout-window").style.display = "flex";
    document.getElementById("user-window").style.display = "";
    for (const globalCheckout of globalCheckouts) {
        if (globalCheckout.date === date) {
            document.getElementById("checkout-top-name").innerHTML = globalCheckout.name;
            document.getElementById("checkout-picture").src = globalCheckout.image;
            document.getElementById("checkout-link").innerHTML = globalCheckout.link;
            document.getElementById("checkout-date").innerHTML = new Date(globalCheckout.date).toLocaleString();
            document.getElementById("checkout-price").innerHTML = globalCheckout.price.toString() + "€";
            document.getElementById("checkout-size").innerHTML = globalCheckout.size;
            document.getElementById("checkout-store").innerHTML = globalCheckout.store;
            document.getElementById("checkout-open-user").innerHTML = "Open " + globalCheckout.username;
            openedUser.id = globalCheckout.user;
            break;
        }
    }
}

function onInputCheckoutFunc(name) {
    let input = document.getElementById(name + "-checkout-input").value.toLowerCase().replaceAll(" ", "");
    let checkouts = function () {
        if (name === "global") {
            return globalCheckouts;
        } else {
            return openedUser.checkouts;
        }
    }()

    checkouts = checkouts.filter(checkout => {
        return checkout.query.includes(input);
    })

    parseCheckouts(checkouts, name);
}

function parseInstances(data, name) {
    let instancelist = document.getElementById(name + "-instances-list");
    instancelist.innerHTML = "";

    if (data.length === 0) {
        instancelist.innerHTML = `
            <div class="nothing-list">
                <svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="24" height="24" viewBox="0 0 172 172"><g fill="none" fill-rule="nonzero" stroke="none" stroke-width="1" stroke-linecap="butt" stroke-linejoin="miter" stroke-miterlimit="10" stroke-dasharray="" stroke-dashoffset="0" font-family="none" font-size="none" style="mix-blend-mode: normal"><path d="M0,172v-172h172v172z" fill="none"></path><g fill="#FFFFFF73"><path d="M135.45,0c-6.97406,0.16125 -13.82719,2.33813 -19.6725,6.45l-28.4875,20.1025l7.8475,11.395l28.595,-20.21c4.79719,-3.35937 10.57531,-4.66281 16.34,-3.655c5.76469,1.00781 10.83063,4.23281 14.19,9.03c3.37281,4.79719 4.66281,10.57531 3.655,16.34c-1.00781,5.76469 -4.23281,10.81719 -9.03,14.19c-2.09625,1.47813 -23.08562,16.20563 -24.51,17.2c-5.46906,3.84313 -13.07469,5.76469 -22.4675,1.6125l-12.255,8.4925c13.06125,10.13188 30.46281,9.70188 42.57,1.1825c1.42438,-0.99437 22.54813,-15.84281 24.6175,-17.3075c7.78031,-5.4825 12.98063,-13.61219 14.62,-23.005c1.63938,-9.37937 -0.44344,-18.86625 -5.9125,-26.66c-5.4825,-7.79375 -13.71969,-12.98062 -23.1125,-14.62c-2.33812,-0.40312 -4.66281,-0.59125 -6.9875,-0.5375zM9.5675,3.3325c-2.67406,0.25531 -4.95844,2.05594 -5.83187,4.59563c-0.88688,2.55312 -0.20156,5.375 1.74687,7.22937l37.84,37.84c1.65281,2.05594 4.32688,2.98313 6.89344,2.39188c2.58,-0.59125 4.58219,-2.59344 5.17344,-5.17344c0.59125,-2.56656 -0.33594,-5.24062 -2.39188,-6.89344l-37.84,-37.84c-1.29,-1.35719 -3.07719,-2.13656 -4.945,-2.15c-0.215,-0.01344 -0.43,-0.01344 -0.645,0zM26.5525,87.29l-20.1025,28.595c-5.4825,7.79375 -7.55188,17.17313 -5.9125,26.5525c1.63938,9.39281 6.82625,17.63 14.62,23.1125c7.79375,5.4825 17.28063,7.55188 26.66,5.9125c9.37938,-1.63937 17.53594,-6.83969 23.005,-14.62c1.46469,-2.06937 16.31313,-23.19312 17.3075,-24.6175c8.51938,-12.10719 8.94938,-29.40125 -1.1825,-42.4625l-8.4925,12.1475c4.15219,9.39281 2.23063,16.99844 -1.6125,22.4675c-0.99437,1.42438 -15.72187,22.52125 -17.2,24.6175c-3.37281,4.79719 -8.42531,7.91469 -14.19,8.9225c-5.76469,1.00781 -11.54281,-0.28219 -16.34,-3.655c-4.79719,-3.35937 -8.02219,-8.42531 -9.03,-14.19c-1.00781,-5.76469 0.29563,-11.54281 3.655,-16.34l20.21,-28.595z"></path></g></g></svg>
                Nothing here
            </div>
        `;
    }

    for (const instance of data) {
        instancelist.innerHTML += `
            <div class="log" onclick="openInstance('${instance.id}', '${name}')">
                <div class="left">
                    <h6 class="statuscode">${instance.status}</h6>
                    <h4>ID ${instance.id}</h4>
                    <h4>${instance.provider}</h4>
                </div>
                <h6 class="time">${instance.username}</h6>
            </div>
        `;
    }
}

function openInstanceFunc(id, name) {
    document.getElementById("instance-close").onclick = () => {
        document.getElementById("instance-window").style.display = "";
        if (name === "user") openUser(openedUser.id);
    }

    document.getElementById("instance-window").style.display = "flex";
    document.getElementById("user-window").style.display = "";
    for (const globalInstance of globalInstances) {
        if (globalInstance.id === id) {
            document.getElementById("instance-top-id").innerHTML = globalInstance.id;
            document.getElementById("instance-provider").innerHTML = globalInstance.provider;
            document.getElementById("instance-region").innerHTML = globalInstance.region;
            document.getElementById("instance-status").innerHTML = globalInstance.status;
            document.getElementById("instance-tasks").innerHTML = globalInstance.task_max;
            document.getElementById("instance-price").innerHTML = globalInstance.price + "€";
            document.getElementById("instance-open-user").innerHTML = "Open " + globalInstance.username;
            openedUser.id = globalInstance.user;
            break;
        }
    }
}

function onInputInstancesFunc(name) {
    let input = document.getElementById(name + "-instances-input").value.toLowerCase().replaceAll(" ", "");
    let instances = function () {
        if (name === "global") {
            return globalInstances;
        } else {
            return openedUser.instances;
        }
    }()

    instances = instances.filter(checkout => {
        return checkout.query.includes(input);
    })

    parseInstances(instances, name);
}

function urlExists(url) {
    let http = new XMLHttpRequest();
    http.open('HEAD', url, false);
    http.send();
    return http.status !== 404;
}

async function openUserFunc(id) {
    let loading = document.getElementById("user-load");
    let window = document.getElementById("user-window");
    let field = document.getElementById("user-field");
    loading.style.display = "flex";
    window.style.display = "flex";
    field.style.display = "none";
    userSwitch('info');

    let res = await fetch(activedata + "user/" + id, {
        method: "GET",
        credentials: "include",
    }).catch(err => error(err));

    if (res.status !== 200) {
        error("Error: " + res.status);
        window.style.display = "";
    } else {
        openedUser = await res.json().catch(err => error(err));
        document.getElementById("submit-update").style.display = "none";
        document.getElementById("user-top-name").innerHTML = openedUser.user.name;
        document.getElementById("user-id").innerHTML = openedUser.user.id.toString();
        document.getElementById("user-data").innerHTML = new Blob([openedUser.data]).size + " Bytes";
        document.getElementById("user-tasks").innerHTML = (openedUser.session?.tasks || 0) + " Tasks";
        document.getElementById("user-instances").innerHTML = (openedUser.session?.instances?.length.toString() || "0") + " Instances";
        document.getElementById("user-checkouts").innerHTML = (openedUser.session?.checkouts?.length.toString() || "0") + " Checkouts";
        document.getElementById("user-status").innerHTML = openedUser.session?.status || "Offline";
        document.getElementById("user-startbutton").innerHTML = openedUser.session?.status === "Running" ?
            "<svg class=\"svgbutton\" onclick=\"changeStartUser()\" xmlns=\"http://www.w3.org/2000/svg\" x=\"0px\" y=\"0px\" width=\"60\" height=\"60\" viewBox=\"0 0 172 172\"><g fill=\"none\" fill-rule=\"nonzero\" stroke=\"none\" stroke-width=\"1\" stroke-linecap=\"butt\" stroke-linejoin=\"miter\" stroke-miterlimit=\"10\" stroke-dasharray=\"\" stroke-dashoffset=\"0\" font-family=\"none\" font-size=\"none\" style=\"mix-blend-mode: normal\"><path d=\"M0,172v-172h172v172z\" fill=\"none\"></path><g fill=\"#FFFFFF73\"><path d=\"M57.33333,35.83333c-7.91917,0 -14.33333,6.41417 -14.33333,14.33333v71.66667c0,7.91917 6.41417,14.33333 14.33333,14.33333c7.91917,0 14.33333,-6.41417 14.33333,-14.33333v-71.66667c0,-7.91917 -6.41417,-14.33333 -14.33333,-14.33333zM114.66667,35.83333c-7.91917,0 -14.33333,6.41417 -14.33333,14.33333v71.66667c0,7.91917 6.41417,14.33333 14.33333,14.33333c7.91917,0 14.33333,-6.41417 14.33333,-14.33333v-71.66667c0,-7.91917 -6.41417,-14.33333 -14.33333,-14.33333z\"></path></g></g></svg>":
            "<svg class=\"svgbutton\" onclick=\"changeStartUser()\" xmlns=\"http://www.w3.org/2000/svg\" x=\"0px\" y=\"0px\" width=\"80\" height=\"80\" viewBox=\"0 0 172 172\"><g fill=\"none\" fill-rule=\"nonzero\" stroke=\"none\" stroke-width=\"1\" stroke-linecap=\"butt\" stroke-linejoin=\"miter\" stroke-miterlimit=\"10\" stroke-dasharray=\"\" stroke-dashoffset=\"0\" font-family=\"none\" font-size=\"none\" style=\"mix-blend-mode: normal\"><path d=\"M0,172v-172h172v172z\" fill=\"none\"></path><g id=\"original-icon\" fill=\"#FFFFFF73\"><path d=\"M57.33333,48.891v74.22517c0,5.6545 6.24217,9.08017 11.01517,6.04867l58.31517,-37.109c4.429,-2.8165 4.429,-9.27367 0,-12.09017l-58.31517,-37.12333c-4.773,-3.03867 -11.01517,0.39417 -11.01517,6.04867z\"></path></g></g></svg>";

        document.getElementById("user-plan-dropdown").children[openedUser.user.subscription.plan].click();
        document.getElementById("user-email").value = openedUser.user.email;
        document.getElementById("user-instance-limit").value = openedUser.user.instance_limit;

        let url = "https://cdn.discordapp.com/avatars/" + openedUser.user.id + "/" + openedUser.user.picture;
        let exists = urlExists(url);
        document.getElementById("user-picture").value = openedUser.user.picture;
        document.getElementById("user-picture-exists").innerHTML = exists ? "✓" : "✗";
        document.getElementById("user-picture-preview").src = exists ? url : "https://cdn.discordapp.com/embed/avatars/0.png";
        setStoresUser();

        openedUser.checkouts.forEach(checkout => {
            checkout.username = openedUser.user.name;
            checkout.query = (openedUser.id + openedUser.name + checkout.date + checkout.name + checkout.link + checkout.image + checkout.store + checkout.size + checkout.price.toFixed(2)).replaceAll(" ", "").toLowerCase();
        })

        openedUser.instances = openedUser.instances || []
        openedUser.instances.forEach(instance => {
            instance.user = openedUser.user.id;
            instance.username = openedUser.user.name;
            instance.query = (openedUser.id + openedUser.name + instance.id + instance.price.toFixed(2) + instance.provider + instance.region + instance.status + instance.task_max).replaceAll(" ", "").toLowerCase();
        })

        parseCheckouts(openedUser.checkouts, "user");
        parseInstances(openedUser.instances, "user");
    }

    field.style.display = "";
    loading.style.display = "";
}

function setStoresUser() {
    openedUser.settings = openedUser.settings || {};
    openedUser.settings.stores = openedUser.settings?.stores || new Map();
    let stores = new Map(Object.entries(openedUser.settings.stores))
    for (const child of document.getElementById("user-stores").children) {
        if (!child.id.includes("user-")) continue
        let enabled = stores.get(child.id.replace("user-", ""));
        if (enabled === undefined) enabled = false;
        child.children.item(+!enabled).classList.add("switchactive");
        child.children.item(+enabled).classList.remove("switchactive");
    }
}

function updateStoresUserFunc(store, value) {
    openedUser.settings.stores[store] = value;
    setStoresUser();
    verifyChange().then(() => {});
}

function changeStartUserFunc() {
    if (!openedUser.session || openedUser.session.status === "Offline") {
        return
    }

    openedUser.session.status = openedUser.session.status === "Running" ? "Stopped by Admin" : "Running";

    document.getElementById("user-status").innerHTML = openedUser.session.status;
    document.getElementById("user-startbutton").innerHTML = openedUser.session.status === "Running" ?
        "<svg class=\"svgbutton\" onclick=\"changeStartUser()\" xmlns=\"http://www.w3.org/2000/svg\" x=\"0px\" y=\"0px\" width=\"60\" height=\"60\" viewBox=\"0 0 172 172\"><g fill=\"none\" fill-rule=\"nonzero\" stroke=\"none\" stroke-width=\"1\" stroke-linecap=\"butt\" stroke-linejoin=\"miter\" stroke-miterlimit=\"10\" stroke-dasharray=\"\" stroke-dashoffset=\"0\" font-family=\"none\" font-size=\"none\" style=\"mix-blend-mode: normal\"><path d=\"M0,172v-172h172v172z\" fill=\"none\"></path><g fill=\"#FFFFFF73\"><path d=\"M57.33333,35.83333c-7.91917,0 -14.33333,6.41417 -14.33333,14.33333v71.66667c0,7.91917 6.41417,14.33333 14.33333,14.33333c7.91917,0 14.33333,-6.41417 14.33333,-14.33333v-71.66667c0,-7.91917 -6.41417,-14.33333 -14.33333,-14.33333zM114.66667,35.83333c-7.91917,0 -14.33333,6.41417 -14.33333,14.33333v71.66667c0,7.91917 6.41417,14.33333 14.33333,14.33333c7.91917,0 14.33333,-6.41417 14.33333,-14.33333v-71.66667c0,-7.91917 -6.41417,-14.33333 -14.33333,-14.33333z\"></path></g></g></svg>":
        "<svg class=\"svgbutton\" onclick=\"changeStartUser()\" xmlns=\"http://www.w3.org/2000/svg\" x=\"0px\" y=\"0px\" width=\"80\" height=\"80\" viewBox=\"0 0 172 172\"><g fill=\"none\" fill-rule=\"nonzero\" stroke=\"none\" stroke-width=\"1\" stroke-linecap=\"butt\" stroke-linejoin=\"miter\" stroke-miterlimit=\"10\" stroke-dasharray=\"\" stroke-dashoffset=\"0\" font-family=\"none\" font-size=\"none\" style=\"mix-blend-mode: normal\"><path d=\"M0,172v-172h172v172z\" fill=\"none\"></path><g id=\"original-icon\" fill=\"#FFFFFF73\"><path d=\"M57.33333,48.891v74.22517c0,5.6545 6.24217,9.08017 11.01517,6.04867l58.31517,-37.109c4.429,-2.8165 4.429,-9.27367 0,-12.09017l-58.31517,-37.12333c-4.773,-3.03867 -11.01517,0.39417 -11.01517,6.04867z\"></path></g></g></svg>";

    verifyChange().then(() => {});
}

async function onInputUser() {
    openedUser.user.email = document.getElementById("user-email").value
    openedUser.user.instance_limit = parseInt(document.getElementById("user-instance-limit").value)
    let newPicture = document.getElementById("user-picture").value
    if (newPicture !== openedUser.user.picture) {
        let url = "https://cdn.discordapp.com/avatars/" + openedUser.user.id + "/" + newPicture;
        let exists = urlExists(url);
        document.getElementById("user-picture-exists").innerHTML = exists ? "✓" : "✗";
        document.getElementById("user-picture-preview").src = exists ? url : "https://cdn.discordapp.com/embed/avatars/0.png";
    }

    openedUser.user.picture = newPicture;
    verifyChange().then(() => {});
}

async function verifyChangeFunc() {
    if (document.getElementById("user-field").style.display === "none") return
    let res = await fetch(activedata + "user/" + openedUser.user.id, {
        method: "GET",
        credentials: "include",
    }).catch(err => error(err));

    if (res.status !== 200) {
        error("Error: " + res.status);
        return
    }

    let defaultUser = await res.json().catch(err => error(err));

    let storeChanged = false;
    for (const [k, v] of Object.entries(openedUser.settings?.stores)) {
        let oldStore = defaultUser.settings?.stores?.[k] || false
        if (oldStore !== v) {
            storeChanged = true;
            break;
        }
    }

    if (defaultUser.session && (defaultUser.session.status !== openedUser.session?.status) ||
        document.getElementById('user-plan-dropdown-input').getAttribute('data-selected') !== (openedUser.user.subscription.plan).toString() ||
        defaultUser.user.email !== openedUser.user.email ||
        defaultUser.user.instance_limit !== openedUser.user.instance_limit ||
        defaultUser.user.picture !== openedUser.user.picture ||
        storeChanged) {

        document.getElementById("submit-update").style.display = "";

    } else {
        document.getElementById("submit-update").style.display = "none";
    }
}

function submitUpdate() {
    document.getElementById("user-update-load").style.display = "flex";
    document.getElementById("submit-update").style.display = "none";

    openedUser.user.subscription.plan = parseInt(document.getElementById('user-plan-dropdown-input').getAttribute('data-selected'));
    fetch(activedata + "user/" + openedUser.user.id, {
        method: "PATCH",
        credentials: "include",
        body: JSON.stringify(openedUser),
    }).then(resp => {
        document.getElementById("user-update-load").style.display = "";
        if (resp.status !== 200) {
            document.getElementById("submit-update").style.display = "";
            throw new Error("Statuscode " + resp.status);
        }

        verifyChange().then(() => {});
    }).catch(err => error(err));
}

function userSwitchFunc(enable) {
    document.getElementById("user-info-window").style.display = "none";
    document.getElementById("user-checkouts-window").style.display = "none";
    document.getElementById("user-instances-window").style.display = "none";
    document.getElementById("user-"+enable+"-window").style.display = "";

    document.getElementById("user-info-switch").classList.remove("switchactive");
    document.getElementById("user-checkouts-switch").classList.remove("switchactive");
    document.getElementById("user-instances-switch").classList.remove("switchactive");
    document.getElementById("user-"+enable+"-switch").classList.add("switchactive");
}

async function logsDownload() {
    document.getElementById("logs-download-loading").style.display = "flex";
    document.getElementById("logs-download").style.display = "none";

    let instance = document.getElementById("instance-top-id").innerHTML
    let amount = parseInt(document.getElementById("instance-log-amount").value)

    let q = query(ref(database, "userstats/user/" + openedUser.id + "/logs/" + instance), orderByChild("date"), limitToLast(amount));
    let snapshot = await get(q);
    if (snapshot.exists()) {
        let values = snapshot.val();
        let logs = Object.values(values)
        let csv = "Date,State,Message\n"
        for (const log of logs) {
            csv += log.date + "," + log.state + ",\"" + log.message.replaceAll("\"", "'") + "\"\n";
        }

        save(openedUser.id + "-" + instance + ".csv", csv);
    } else {
        error("No logs found");
    }

    document.getElementById("logs-download-loading").style.display = "";
    document.getElementById("logs-download").style.display = "";
}

function toggleDeleteLogs() {
    document.getElementById("delete-date").value = timeframeToday;
    document.getElementById("delete-logs-window").style.display = document.getElementById("delete-logs-window").style.display === "" ? "flex" : "";
}

function submitDeleteLogs() {
    let time = document.getElementById("delete-date").value;
    let server = document.getElementById("delete-dropdown-input").getAttribute("data-selected");
    remove(ref(database, "old-serverstats/" + server + "/" + time + "/logs"));
    remove(ref(database, "serverstats/" + server + "/" + time + "/logs"));
    remove(ref(database, "old-serverstats/" + server + "/" + time + "/requests"));
    remove(ref(database, "serverstats/" + server + "/" + time + "/requests"));
    toggleDeleteLogs();
}

$(function() {
    let start = moment();
    let end = moment();

    function cb(start, end) {
        timeframeFrom = start.format('YYYY-MMMM-D');
        timeframeTo = end.format('YYYY-MMMM-D');
        showData();

        let s = start.format('MMMM D, YYYY');
        let e = end.format('MMMM D, YYYY');
        if (s === e) e = "";
        else e = " - " + e;

        $('#reportrange span').html(s + e);
    }

    $('#reportrange').daterangepicker({
        opens: 'left',
        startDate: start,
        endDate: end,
        ranges: {
            'Today': [moment(), moment()],
            'Yesterday': [moment().subtract(1, 'days'), moment().subtract(1, 'days')],
            'Last 7 Days': [moment().subtract(6, 'days'), moment()],
            'Last 30 Days': [moment().subtract(29, 'days'), moment()],
            'This Month': [moment().startOf('month'), moment().endOf('month')],
            'Last Month': [moment().subtract(1, 'month').startOf('month'), moment().subtract(1, 'month').endOf('month')]
        }
    }, cb);

    cb(start, end);
});

init();

window.addEventListener("readystatechange", (event) => {
    if (event.readyState === "complete") init();
});