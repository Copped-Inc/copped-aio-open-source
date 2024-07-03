const localdata = "http://localhost:91/"
const webdata = "https://database.copped-inc.com/"
let activedata = webdata

window.addEventListener("load", (event) => {
    if (document.getElementById("generate-keys-window") == null) onLoad();
});

function onLoad() {
    document.getElementById("loginload").style.display = "";
    document.getElementById("passwordinput").style.display = "";
    document.getElementById("loginbutton").style.display = "";

    if (getCookie("localhost") != null) {
        activedata = localdata;
    }

    fetch(activedata + "data", {
        method: "GET",
        credentials: "include",
        cache: "no-cache"
    }).then(resp => {
        if (resp.redirected) {
            window.location.href = resp.url

        } else if (resp.status !== 203) {
            throw new Error("Error: " + resp.status)

        } else if (sessionStorage.getItem("password")) {
            document.getElementById("passwordinput").value = sessionStorage.getItem("password");
            canLogin = true;
            login()
            
        }
        AddHandler();
    }).catch(() => {
        window.location.href = activedata + "login";
    })
}
