function changeTab(w) {
    if (activeBurger) {
        toggleBurger();
    }

    let dashboard = document.getElementById("dashboard");
    let server = document.getElementById("server");
    let user = document.getElementById("user");

    let dashtab = document.getElementById("dashtab");
    let servertab = document.getElementById("servertab");
    let usertab = document.getElementById("usertab");

    let dashicon = document.getElementById("dashicon");
    let servericon = document.getElementById("servericon");
    let usericon = document.getElementById("usericon");
    if (w === 0) {
        dashboard.style.display = "flex";
        server.style.display = "none";
        user.style.display = "none";

        dashtab.classList.remove("disabled");
        servertab.classList.add("disabled");
        usertab.classList.add("disabled");

        dashicon.setAttribute("fill", "#00873e")
        servericon.setAttribute("fill", "#FFFFFF73")
        usericon.setAttribute("fill", "#FFFFFF73")
    } else if (w === 1) {
        dashboard.style.display = "none";
        server.style.display = "flex";
        user.style.display = "none";

        dashtab.classList.add("disabled");
        servertab.classList.remove("disabled");
        usertab.classList.add("disabled");

        dashicon.setAttribute("fill", "#FFFFFF73")
        servericon.setAttribute("fill", "#00873e")
        usericon.setAttribute("fill", "#FFFFFF73")
    } else if (w === 2) {
        dashboard.style.display = "none";
        server.style.display = "none";
        user.style.display = "flex";

        dashtab.classList.add("disabled");
        servertab.classList.add("disabled");
        usertab.classList.remove("disabled");

        dashicon.setAttribute("fill", "#FFFFFF73")
        servericon.setAttribute("fill", "#FFFFFF73")
        usericon.setAttribute("fill", "#00873e")
    }
}