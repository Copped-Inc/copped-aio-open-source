let activeBurger = false;

function toggleBurger() {
    if (activeNotification && !activeBurger) {
        toggleNotification();
    }

    document.getElementById("tabs").style.display = activeBurger ? "" : "flex";
    activeBurger = !activeBurger;
}