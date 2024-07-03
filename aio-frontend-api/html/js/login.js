function AddHandler() {
    canLogin = true;
    document.getElementById("passwordinput").addEventListener("keypress", function (e) {
        if (e.key === "Enter") {
            login();
        }
    })
    document.getElementById("passwordconfirminput").addEventListener("keypress", function (e) {
        if (e.key === "Enter") {
            login();
        }
    })
    document.getElementById("passwordinput").addEventListener("keydown", function (e) {
        if (e.key !== "Enter") {
            document.getElementById("passwordmessage").style.visibility = "hidden";
        }
    })
    document.getElementById("webhookinput").addEventListener("keydown", function (e) {
        if (e.key === "Enter") {
            addWebhook();
        }
    })
    document.getElementById("keywords-input").addEventListener("keydown", function (e) {
        if (e.key === "Enter") {
            addWhitelist();
        }
    })
}

let data
let confirm = false
let password = ""
let canLogin = false

function login() {
    if (!canLogin) {
        return;
    }

    let loginload = document.getElementById("loginload");
    let loginbutton = document.getElementById("loginbutton");
    let passwordinput = document.getElementById("passwordinput");
    let passwordconfirminput = document.getElementById("passwordconfirminput");

    if (passwordinput.value === "") {
        return;
    }

    if (passwordinput.value !== passwordconfirminput.value && confirm) {
        document.getElementById("passwordmessage").style.visibility = "visible";
        return;
    }

    loginload.style.display = "flex";
    loginbutton.style.display = "none";
    passwordinput.style.display = "none";
    passwordconfirminput.style.display = "none";

    fetch(activedata + "data", {
        method: "GET",
        credentials: "include",
        headers: {
            "password": passwordinput.value,
            "confirm": confirm ? "true" : ""
        }
    }).then(resp => {
        if (resp.status === 403) {
            loginload.style.display = "none";
            loginbutton.style.display = "";
            passwordinput.style.display = "";

            document.getElementById("passwordmessage").style.visibility = "visible";
        } else if (resp.status === 200) {
            resp.json().then(d => {
                data = d;
                data.settings = data.settings || {};
                password = passwordinput.value

                sessionStorage.setItem("password", password);
                load();

                document.getElementById("dashboard").style.display = "flex";
                document.getElementById("logo").src = "https://cdn.discordapp.com/avatars/" + data.user.id + "/" + data.user.picture;
                sleep(200).then(() => {
                    document.getElementById("passwordwindow").style.transform = "translate(0, -100%)";
                    sleep(200).then(() => document.getElementById("passwordwindow").style.display = "none")
                });
            })
        } else if (resp.status === 203) {
            loginload.style.display = "none";
            loginbutton.style.display = "";
            passwordinput.style.display = "";
            passwordconfirminput.style.display = "";

            confirm = true
            document.getElementById("passwordconfirminput").style.display = "block";
        } else {
            window.location.href = activedata + "login";
        }
    })
}

function load() {
    parseDate();
    setChart();
    setCheckouts();
    setStores();
    setUser();
    setWebhooks();
    setStart();
    setInstances();
    connect();
    setDrag();
    fetchNotifications();
    setWhitelist();
}
