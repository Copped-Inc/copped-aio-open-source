function error(error) {
    console.log(error)
    let signalholder = document.getElementById("signal");
    document.getElementById("signaltext").innerHTML = error;
    signalholder.classList.add("bad");
    signalholder.classList.remove("good");
    signalholder.style.display = "flex";
}

function signal(signal) {
    let signalholder = document.getElementById("signal");
    document.getElementById("signaltext").innerHTML = signal;
    signalholder.classList.add("good");
    signalholder.classList.remove("bad");
    signalholder.style.display = "flex";

    sleep(3000).then(() => {
        closesignal();
    });
}

function closesignal() {
    document.getElementById("signal").style.display = "none";
}
