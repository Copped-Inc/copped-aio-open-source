let notifications = [];
let activeNotification = false;

function fetchNotifications() {

    fetch(activedata + "notifications/users/@me", {
        method: "GET",
        credentials: "include",
        headers: {
            "password": password,
        }
    }).then(resp => { if (resp.status != 204) return resp.json() }).then(data => {

        notifications = data || [];
        setNotifications()
        notifications.forEach(message => {
            if (message.global) {
                fetchGlobalNotification(message.id);
            }
        });
    }).catch(err => {
        error(err);
    });

}

function fetchGlobalNotification(id) {
    fetch(activedata + "notifications/" + id, {
        method: "GET",
        credentials: "include",
        headers: {
            "password": password,
        }
    }).then(resp => {
        if (resp.status === 404) {
            removeNotification(id)
            fetch(activedata + "notifications/" + id + "/users/@me", {
                method: "DELETE",
                credentials: "include",
                headers: {
                    "password": password,
                }
            }).catch(err => {
                error(err);
            })
        } else {
            return resp.json();
        }
    }).then(message => {
        let i = notifications.indexOf(notifications.find(n => n.id === message.id))
        notifications[i].title = message.title;
        notifications[i].text = message.text;
        notifications[i].created_at = message.created_at;
        setNotifications();
    });
}

function setNotifications() {

    let notificationholder = document.getElementById("notificationholder");
    notificationholder.innerHTML = "";

    if (notifications.length === 0) {
        notificationholder.innerHTML = `
            <div class="nothing">
                <svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="24" height="24" viewBox="0 0 172 172"><g fill="none" fill-rule="nonzero" stroke="none" stroke-width="1" stroke-linecap="butt" stroke-linejoin="miter" stroke-miterlimit="10" stroke-dasharray="" stroke-dashoffset="0" font-family="none" font-size="none" style="mix-blend-mode: normal"><path d="M0,172v-172h172v172z" fill="none"></path><g fill="#FFFFFF73"><path d="M135.45,0c-6.97406,0.16125 -13.82719,2.33813 -19.6725,6.45l-28.4875,20.1025l7.8475,11.395l28.595,-20.21c4.79719,-3.35937 10.57531,-4.66281 16.34,-3.655c5.76469,1.00781 10.83063,4.23281 14.19,9.03c3.37281,4.79719 4.66281,10.57531 3.655,16.34c-1.00781,5.76469 -4.23281,10.81719 -9.03,14.19c-2.09625,1.47813 -23.08562,16.20563 -24.51,17.2c-5.46906,3.84313 -13.07469,5.76469 -22.4675,1.6125l-12.255,8.4925c13.06125,10.13188 30.46281,9.70188 42.57,1.1825c1.42438,-0.99437 22.54813,-15.84281 24.6175,-17.3075c7.78031,-5.4825 12.98063,-13.61219 14.62,-23.005c1.63938,-9.37937 -0.44344,-18.86625 -5.9125,-26.66c-5.4825,-7.79375 -13.71969,-12.98062 -23.1125,-14.62c-2.33812,-0.40312 -4.66281,-0.59125 -6.9875,-0.5375zM9.5675,3.3325c-2.67406,0.25531 -4.95844,2.05594 -5.83187,4.59563c-0.88688,2.55312 -0.20156,5.375 1.74687,7.22937l37.84,37.84c1.65281,2.05594 4.32688,2.98313 6.89344,2.39188c2.58,-0.59125 4.58219,-2.59344 5.17344,-5.17344c0.59125,-2.56656 -0.33594,-5.24062 -2.39188,-6.89344l-37.84,-37.84c-1.29,-1.35719 -3.07719,-2.13656 -4.945,-2.15c-0.215,-0.01344 -0.43,-0.01344 -0.645,0zM26.5525,87.29l-20.1025,28.595c-5.4825,7.79375 -7.55188,17.17313 -5.9125,26.5525c1.63938,9.39281 6.82625,17.63 14.62,23.1125c7.79375,5.4825 17.28063,7.55188 26.66,5.9125c9.37938,-1.63937 17.53594,-6.83969 23.005,-14.62c1.46469,-2.06937 16.31313,-23.19312 17.3075,-24.6175c8.51938,-12.10719 8.94938,-29.40125 -1.1825,-42.4625l-8.4925,12.1475c4.15219,9.39281 2.23063,16.99844 -1.6125,22.4675c-0.99437,1.42438 -15.72187,22.52125 -17.2,24.6175c-3.37281,4.79719 -8.42531,7.91469 -14.19,8.9225c-5.76469,1.00781 -11.54281,-0.28219 -16.34,-3.655c-4.79719,-3.35937 -8.02219,-8.42531 -9.03,-14.19c-1.00781,-5.76469 0.29563,-11.54281 3.655,-16.34l20.21,-28.595z"></path></g></g></svg>
                Nothing here
            </div>
        `;
    }

    let unread = false;
    let realmsg = 0

    notifications.forEach(message => {
        if (message.title === undefined && message.text === undefined) return;
        if (message.read === false || message.read === undefined) unread = true;

        let notidiv = document.createElement("div");
        notidiv.classList.add("noti");
        notidiv.innerHTML += `
                <div class="notisplit">
                    <div class="content">
                        <div class="notificaioncontent">
                            <h4>${message.title}</h4>
                        </div>
                        <div class="notificationcontent">
                            <h6>${message.text.replaceAll("\n", "<br>")}</h6>
                        </div>
                    </div>
                    <svg onclick="deletenoti('${message.id}')" xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="24" height="24" viewBox="0 0 172 172"><g fill="none" fill-rule="nonzero" stroke="none" stroke-width="1" stroke-linecap="butt" stroke-linejoin="miter" stroke-miterlimit="10" stroke-dasharray="" stroke-dashoffset="0" font-family="none" font-size="none" style="mix-blend-mode: normal"><path d="M0,172v-172h172v172z" fill="none"></path><g fill="#FFFFFF73"><path d="M71.66667,14.33333l-7.16667,7.16667h-28.66667c-4.3,0 -7.16667,2.86667 -7.16667,7.16667c0,4.3 2.86667,7.16667 7.16667,7.16667h14.33333h71.66667h14.33333c4.3,0 7.16667,-2.86667 7.16667,-7.16667c0,-4.3 -2.86667,-7.16667 -7.16667,-7.16667h-28.66667l-7.16667,-7.16667zM35.83333,50.16667v93.16667c0,7.88333 6.45,14.33333 14.33333,14.33333h71.66667c7.88333,0 14.33333,-6.45 14.33333,-14.33333v-93.16667zM64.5,64.5c4.3,0 7.16667,2.86667 7.16667,7.16667v64.5c0,4.3 -2.86667,7.16667 -7.16667,7.16667c-4.3,0 -7.16667,-2.86667 -7.16667,-7.16667v-64.5c0,-4.3 2.86667,-7.16667 7.16667,-7.16667zM107.5,64.5c4.3,0 7.16667,2.86667 7.16667,7.16667v64.5c0,4.3 -2.86667,7.16667 -7.16667,7.16667c-4.3,0 -7.16667,-2.86667 -7.16667,-7.16667v-64.5c0,-4.3 2.86667,-7.16667 7.16667,-7.16667z"></path></g></g></svg>
                </div>
        `;

        notificationholder.appendChild(notidiv);
        realmsg++;
    });

    if (unread) {
        document.getElementById("newnoti").style.display = "block";
        document.getElementById("nonoti").style.display = "none";
    } else {
        document.getElementById("newnoti").style.display = "none";
        document.getElementById("nonoti").style.display = "flex";
    }

    document.getElementById("notification").innerHTML = `
        ${realmsg} Notification${realmsg === 1 ? "" : "s"}
        <svg onclick="deleteall()" xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="24" height="24" viewBox="0 0 172 172"><g fill="none" fill-rule="nonzero" stroke="none" stroke-width="1" stroke-linecap="butt" stroke-linejoin="miter" stroke-miterlimit="10" stroke-dasharray="" stroke-dashoffset="0" font-family="none" font-size="none" style="mix-blend-mode: normal"><path d="M0,172v-172h172v172z" fill="none"></path><g fill="#FFFFFF73"><path d="M71.66667,14.33333l-7.16667,7.16667h-28.66667c-4.3,0 -7.16667,2.86667 -7.16667,7.16667c0,4.3 2.86667,7.16667 7.16667,7.16667h14.33333h71.66667h14.33333c4.3,0 7.16667,-2.86667 7.16667,-7.16667c0,-4.3 -2.86667,-7.16667 -7.16667,-7.16667h-28.66667l-7.16667,-7.16667zM35.83333,50.16667v93.16667c0,7.88333 6.45,14.33333 14.33333,14.33333h71.66667c7.88333,0 14.33333,-6.45 14.33333,-14.33333v-93.16667zM64.5,64.5c4.3,0 7.16667,2.86667 7.16667,7.16667v64.5c0,4.3 -2.86667,7.16667 -7.16667,7.16667c-4.3,0 -7.16667,-2.86667 -7.16667,-7.16667v-64.5c0,-4.3 2.86667,-7.16667 7.16667,-7.16667zM107.5,64.5c4.3,0 7.16667,2.86667 7.16667,7.16667v64.5c0,4.3 -2.86667,7.16667 -7.16667,7.16667c-4.3,0 -7.16667,-2.86667 -7.16667,-7.16667v-64.5c0,-4.3 2.86667,-7.16667 7.16667,-7.16667z"></path></g></g></svg>
    `

}

function toggleNotification() {
    if (activeBurger && !activeNotification) {
        toggleBurger()
    }

    document.getElementById("notification").style.display = activeNotification ? "" : "flex";
    document.getElementById("notificationholder").style.display = activeNotification ? "" : "flex";
    document.getElementById("tabs").style.display = activeNotification ? "" : "flex";
    activeNotification = !activeNotification;

    document.getElementById("nonoti").style.display = "flex";
    document.getElementById("newnoti").style.display = "none";

    markRead();
}

function markRead() {
    let unreads = false
    notifications.forEach(notification => {
        if (!notification.read) {
            unreads = true;
            notification.read = true;
        };
    });

    if (unreads) {
        fetch(activedata + "notifications/users/@me", {
            method: "PUT",
            credentials: "include",
            headers: {
                "password": password,
            },
            body: JSON.stringify({ read: true })
        }).catch(err => {
            error(err);
        });
    };
}

function deleteall() {
    fetch(activedata + "notifications/users/@me", {
        method: "DELETE",
        credentials: "include",
        headers: {
            "password": password,
        }
    }).catch(err => {
        error(err);
    });
    notifications = [];
    setNotifications();
}

function removeNotification(id) {
    notifications = notifications.filter(notification => notification.id !== id);
    setNotifications();
}

function deletenoti(id) {
    fetch(activedata + "notifications/" + id + "/users/@me", {
        method: "DELETE",
        credentials: "include",
        headers: {
            "password": password,
        }
    }).then(() => {
        removeNotification(id);
        setNotifications()
    }).catch(err => {
        error(err);
    });
}