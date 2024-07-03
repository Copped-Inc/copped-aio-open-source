function setUser() {
    document.getElementById("userlogo").src = "https://cdn.discordapp.com/avatars/" + data.user.id + "/" + data.user.picture;
    document.getElementById("name").innerHTML = data.user.name;
    document.getElementById("email").innerHTML = data.user.email;
    document.getElementById("plan").innerHTML = getPlan();
    document.getElementById("instanceLimit").innerHTML = data.user.instance_limit;
    document.getElementById("subscription-manage").onclick = () => {
        if (data.user.subscription.plan === 3) signal("You are a developer, you don't have a subscription")
        else document.location = activedata + "stripe/customer-portal"
    }
}

function logout() {
    sessionStorage.removeItem("password");
    removeCookie("authorization");
    location.reload();
}

function getPlan() {
    switch (data.user.subscription.plan) {
        case 1:
            return "Friends & Family";
        case 2:
            return "Basic";
        case 3:
            return "Developer";
    }
}
