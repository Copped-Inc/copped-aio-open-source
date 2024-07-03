const localdata = "http://localhost:91/"
const webdata = "https://database.copped-inc.com/"
let activedata = webdata

let sitekey = ""

window.onload = function () {
    if (getCookie("localhost") != null) {
        activedata = localdata;
    }

    fetch(activedata + "captcha/challenge", {
        method: "GET",
        credentials: "include",
    }).then(resp => {
        if (resp.status !== 200) {
            throw new Error("Error: " + resp.status);
        }
        return resp.json()
    }).then(resp => {
        sitekey = resp.sitekey;
        document.getElementById("captcha").setAttribute("data-sitekey", sitekey);

        let s = document.createElement( 'script' );
        s.setAttribute( 'src', "js/captcha.min.js" );
        s.setAttribute( 'async', "" );

        document.body.appendChild( s );
    }).catch(err => {
        error(err);
    })
}

function OnCaptchaSolved() {
    let captcha = document.getElementById("captcha");
    let captcharesp = captcha.children.item(0).getAttribute("data-hcaptcha-response");

    fetch(activedata + "captcha/" + sitekey, {
        method: "POST",
        credentials: "include",
        headers: {
            "Browser": "true"
        },
        body: JSON.stringify({
            errorId: 0,
            errorCode: "",
            errorDescription: "",
            solution: {
                gRecaptchaResponse: captcharesp
            },
            status: "ready"
        })
    }).then(resp => {
        if (resp.status !== 200) {
            throw new Error("Error: " + resp.status);
        }
        location.reload();
    }).catch(err => {
        error(err);
    })
}
