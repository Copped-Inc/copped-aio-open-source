import './main.css';

import logo from './assets/logo.png';
import load_button from './assets/load_button.svg';
import {Loggedin, Login, Check} from "../wailsjs/go/main/App";

document.getElementById("logo").src = logo;
document.getElementById("load_button").src = load_button;
document.getElementById("instanceinput").addEventListener("keydown", function (e) {
    if (e.key === "Enter") {
        window.login();
    }
});

window.onload = function () {
    Loggedin().then((result) => {
            console.log(result);
            if (!result) toggleFields();
            lastCheck = result;
            check();
        }).catch((err) => {
            console.error(err);
        });
};

window.login = function () {
    let code = document.getElementById("instanceinput").value;
    if (code === "") {
        return;
    }

    document.getElementById("passwordmessage").style.visibility = "none";
    toggleFields();
    Login(code).then((result) => {
            if (!result) {
                document.getElementById("passwordmessage").style.visibility = "visible";
                document.getElementById("passwordmessage").innerHTML = "Invalid instance code or instance already running";

                toggleFields();
            }

        }).catch((err) => {
            console.error(err);
        });
}

function toggleFields() {
    console.log("toogleFields");
    document.getElementById("loginload").style.display = document.getElementById("loginload").style.display === "none" ? "flex" : "none";
    document.getElementById("loginbutton").style.display = document.getElementById("loginbutton").style.display === "none" ? "flex" : "none";
    document.getElementById("instanceinput").style.display = document.getElementById("instanceinput").style.display === "none" ? "flex" : "none";
}

let lastCheck;

function check() {
    Check().then((result) => {
            if (result && !lastCheck) {
                document.getElementById("passwordmessage").style.visibility = "";
                toggleFields();
            } else if (!result && lastCheck) {
                document.getElementById("passwordmessage").style.visibility = "visible";
                document.getElementById("passwordmessage").innerHTML = "Disconnected. Reconnecting...";
                toggleFields();
            }
            lastCheck = result;
            setTimeout(check, 1000);
        }).catch((err) => {
            console.error(err);
        });
}
