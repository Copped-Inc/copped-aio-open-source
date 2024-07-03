let chart;

function setChart() {
    let date = new Date();
    let labels = [];
    let amounts = [];

    data.checkouts = data.checkouts == null ? [] : data.checkouts;

    if (document.getElementById("daychart").classList.contains("switchactive")) {
        for (let i = 0; i < 12; i++) {
            labels.push(date.getHours())

            let amount = 0;
            for (let j = 0; j < data.checkouts.length; j++) {
                if (data.checkouts[j].date.getHours() === date.getHours() && data.checkouts[j].date.getDate() === date.getDate() && data.checkouts[j].date.getMonth() === date.getMonth() && data.checkouts[j].date.getFullYear() === date.getFullYear())
                    amount += data.checkouts[j].price;
            }

            amounts.push(amount);
            date.setHours(date.getHours() - 1);
        }
        labels[0] = "Now";
    } else if (document.getElementById("weekchart").classList.contains("switchactive")) {
        for (let i = 0; i < 7; i++) {
            labels.push(date.getDate() + "." + (date.getMonth() + 1) + ".");

            let amount = 0;
            for (let j = 0; j < data.checkouts.length; j++) {
                if (data.checkouts[j].date.getDate() === date.getDate() && data.checkouts[j].date.getMonth() === date.getMonth()) {
                    amount += data.checkouts[j].price;
                }
            }

            amounts.push(amount);
            date.setDate(date.getDate() - 1);
        }
        labels[0] = "Today";
    } else if (document.getElementById("monthchart").classList.contains("switchactive")) {
        labels[0] = "This week";
        labels[1] = "Last week";
        labels[2] = "2 weeks ago";
        labels[3] = "3 weeks ago";

        amounts[0] = 0;
        amounts[1] = 0;
        amounts[2] = 0;
        amounts[3] = 0;

        for (let i = 0; i < 4; i++) {
            date.setDate(date.getDate() - 7);
            let amount = 0;
            for (let j = 0; j < data.checkouts.length; j++) {

                if (data.checkouts[j].date.getTime() >= date.getTime() && data.checkouts[j].date.getTime() <= date.getTime() + 8.64e+7*6) {
                    amount += data.checkouts[j].price;
                }
            }
            amounts[i] = amount;
        }
    }
    labels.reverse();
    amounts.reverse();

    if (chart !== undefined) chart.destroy();

    chart = new Chart(document.getElementById('checkouts').getContext('2d'), {
        type: 'line',
        color: "#DE3E12",
        data: {
            labels: labels,
            datasets: [{
                label: "Money spent",
                data: amounts,
                fill: true,
                borderColor: "#00873E",
                backgroundColor: "rgba(0,135,62,0.05)",
                pointRadius: 0,
                cubicInterpolationMode: 'monotone'
            }]
        },
        options: {
            scales: {
                y: {
                    grid: {
                        color: "#1f1e1e",
                        borderColor: "#181A19",
                    },
                    ticks: {
                        color: "#FFFFFF73",
                        font: {
                            size: 14,
                        }
                    }
                },
                x: {
                    grid: {
                        color: "#181A19",
                        borderColor: "#181A19",
                    },
                    ticks: {
                        color: "#FFFFFF73",
                        font: {
                            size: 14,
                        }
                    }
                }
            },
            plugins: {
                legend: {
                    display: false
                }
            }
        }
    });
}

function changeChart(name) {
    document.getElementById("daychart").classList.remove("switchactive");
    document.getElementById("weekchart").classList.remove("switchactive");
    document.getElementById("monthchart").classList.remove("switchactive");
    document.getElementById(name + "chart").classList.add("switchactive");
    setChart();
}