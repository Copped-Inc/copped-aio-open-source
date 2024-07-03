let charts = {}
let chartData = {}
let leaved = false;
let activelog = "database";

let showLogs = function (name) {};
let openRequest = function (id) {};
let openUser = function (id) {};
let updateStoresUser = function (store, value) {};
let changeStartUser = function () {};
let verifyChange = function () {};
let userSwitch = function (enable) {};
let onInputCheckout = function (name) {};
let onInputInstances = function (name) {};
let openCheckout = function (date) {};
let openInstance = function (date) {};

function setDevChart(name, labels, newAmounts, oldAmounts) {
    if (charts[name] != null) {
        charts[name].destroy();
    }

    charts[name] = new Chart(document.getElementById(name).getContext('2d'), {
        type: 'line',
        color: "#DE3E12",
        data: {
            labels: labels,
            datasets: [
                {
                    label: "new",
                    data: newAmounts,
                    fill: false,
                    borderWidth: 2,
                    borderColor: "rgb(0,135,62)",
                    pointRadius: 0,
                    pointBorderColor: "#FFFFFF73",
                },
                {
                    label: "old",
                    data: oldAmounts,
                    fill: false,
                    borderWidth: 2,
                    borderColor: "#25252a",
                    pointRadius: 0,
                    pointBorderColor: "#FFFFFF73",
                }
            ]
        },
        options: {
            animation: false,
            interaction: {
                mode: 'index',
                intersect: false,
            },
            maintainAspectRatio: false,
            responsive: true,
            scales: {
                y: {
                    grid: {
                        color: "transparent",
                        borderColor: "transparent",
                    },
                    ticks: {
                        display: false,
                    }
                },
                x: {
                    grid: {
                        color: "transparent",
                        borderColor: "#25252a",
                    },
                    ticks: {
                        display: false,
                    }
                }
            },
            plugins: {
                legend: {
                    display: false
                },
                tooltip: {
                    backgroundColor: "transparent",
                    titleColor: "transparent",
                    callbacks: {
                        label: function (context) {
                            if (leaved) return;
                            let ext = name.includes("latency") ? " ms" : " r/m";

                            let label = context.label.split(";")
                            document.getElementById("new" + name + "label").innerHTML = label[0];
                            document.getElementById("old" + name + "label").innerHTML = label[1];

                            document.getElementById("old" + name).innerHTML = oldAmounts[context.dataIndex].toFixed(2) + ext;
                            if (context.dataIndex < newAmounts.length) {
                                document.getElementById("new" + name).innerHTML = newAmounts[context.dataIndex].toFixed(2) + ext;
                                setRatio(name, newAmounts[context.dataIndex], oldAmounts[context.dataIndex]);
                            } else {
                                document.getElementById("new" + name).innerHTML = "-" + ext;
                                setRatio(name, null, oldAmounts[context.dataIndex]);
                            }
                        }
                    },
                }
            }
        }
    });
}

function mouseEnter(name) {
    leaved = false;
    document.getElementById("new" + name + "label").style.color = "var(--darkgrey)";
    document.getElementById("old" + name + "label").style.color = "var(--darkgrey)";
}

function mouseLeave(name) {
    leaved = true;
    let ext = name.includes("latency") ? " ms" : " r/m";

    document.getElementById("new" + name + "label").style.color = "transparent";
    document.getElementById("old" + name + "label").style.color = "transparent";
    document.getElementById("new" + name).innerHTML = chartData[name] + ext;
    document.getElementById("old" + name).innerHTML = chartData["old" + name] + ext;
    setRatio(name, chartData[name], chartData["old" + name]);
}

function setRatio(name, reqRpm, oldReqRpm) {
    let ratioElement = document.getElementById(name + "ratio");
    ratioElement.classList.remove("positive", "negative", "neutral");

    let ratio = function () {
        if (reqRpm === null) {
            ratioElement.innerHTML = "";
            return NaN;
        }

        if (oldReqRpm == 0 && reqRpm == 0) {
            return 0;
        }

        if (oldReqRpm == 0) {
            return 200;
        }

        if (reqRpm == 0) {
            return 0;
        }

        if (oldReqRpm == reqRpm) {
            return 0;
        }

        return isNaN((reqRpm / oldReqRpm) * 100) ? 0 : ((reqRpm / oldReqRpm) * 100).toFixed(0);
    }();

    if (isNaN(ratio)) {
        return;
    }

    if (parseFloat(reqRpm) > parseFloat(oldReqRpm)) {
        ratioElement.classList.add("negative");
        ratioElement.innerHTML = "+" + (ratio - 100) + "%";
    } else if (parseFloat(reqRpm) < parseFloat(oldReqRpm)) {
        ratioElement.classList.add("positive");
        ratioElement.innerHTML = "-" + (100 - ratio) + "%";
    } else {
        ratioElement.classList.add("neutral");
        ratioElement.innerHTML = ratio + "%";
    }
}

function changeLogs(name) {
    document.getElementById("databaselogs").classList.remove("switchactive");
    document.getElementById("monitorlogs").classList.remove("switchactive");
    document.getElementById("servicelogs").classList.remove("switchactive");
    document.getElementById("instanceslogs").classList.remove("switchactive");
    activelog = name;
    document.getElementById(name + "logs").classList.add("switchactive");
    showLogs(name);
}
