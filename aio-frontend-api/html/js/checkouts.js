function setCheckouts() {
    let date = getDate();
    data.checkouts = data.checkouts == null ? [] : data.checkouts;

    let checkoutlist = document.getElementById("checkoutlist");
    checkoutlist.innerHTML = "";

    for (let i = 0; i < data.checkouts.length; i++) {
        if (data.checkouts[i].date.getTime() <= date.getTime()) continue;
        let checkout = data.checkouts[i];
        let checkoutdiv = document.createElement("div");
        checkoutdiv.classList.add("checkout");
        checkoutdiv.innerHTML = `
            <img src="${checkout.image}" alt="">
            <div class="description">
                <div class="split">
                    <h4>${checkout.name}</h4>
                    <h5 class="date">${checkout.date.toLocaleDateString()}</h5>
                </div>
                <div class="attributes">
                    <h6>Size: ${checkout.size}</h6>
                    <h6>Paid: ${checkout.price}€</h6>
                </div>
            </div>
        `;
        checkoutlist.appendChild(checkoutdiv);
    }

    if (checkoutlist.children.length === 0) {
        checkoutlist.innerHTML = `
            <div class="nothing-list">
                <svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="24" height="24" viewBox="0 0 172 172"><g fill="none" fill-rule="nonzero" stroke="none" stroke-width="1" stroke-linecap="butt" stroke-linejoin="miter" stroke-miterlimit="10" stroke-dasharray="" stroke-dashoffset="0" font-family="none" font-size="none" style="mix-blend-mode: normal"><path d="M0,172v-172h172v172z" fill="none"></path><g fill="#FFFFFF73"><path d="M135.45,0c-6.97406,0.16125 -13.82719,2.33813 -19.6725,6.45l-28.4875,20.1025l7.8475,11.395l28.595,-20.21c4.79719,-3.35937 10.57531,-4.66281 16.34,-3.655c5.76469,1.00781 10.83063,4.23281 14.19,9.03c3.37281,4.79719 4.66281,10.57531 3.655,16.34c-1.00781,5.76469 -4.23281,10.81719 -9.03,14.19c-2.09625,1.47813 -23.08562,16.20563 -24.51,17.2c-5.46906,3.84313 -13.07469,5.76469 -22.4675,1.6125l-12.255,8.4925c13.06125,10.13188 30.46281,9.70188 42.57,1.1825c1.42438,-0.99437 22.54813,-15.84281 24.6175,-17.3075c7.78031,-5.4825 12.98063,-13.61219 14.62,-23.005c1.63938,-9.37937 -0.44344,-18.86625 -5.9125,-26.66c-5.4825,-7.79375 -13.71969,-12.98062 -23.1125,-14.62c-2.33812,-0.40312 -4.66281,-0.59125 -6.9875,-0.5375zM9.5675,3.3325c-2.67406,0.25531 -4.95844,2.05594 -5.83187,4.59563c-0.88688,2.55312 -0.20156,5.375 1.74687,7.22937l37.84,37.84c1.65281,2.05594 4.32688,2.98313 6.89344,2.39188c2.58,-0.59125 4.58219,-2.59344 5.17344,-5.17344c0.59125,-2.56656 -0.33594,-5.24062 -2.39188,-6.89344l-37.84,-37.84c-1.29,-1.35719 -3.07719,-2.13656 -4.945,-2.15c-0.215,-0.01344 -0.43,-0.01344 -0.645,0zM26.5525,87.29l-20.1025,28.595c-5.4825,7.79375 -7.55188,17.17313 -5.9125,26.5525c1.63938,9.39281 6.82625,17.63 14.62,23.1125c7.79375,5.4825 17.28063,7.55188 26.66,5.9125c9.37938,-1.63937 17.53594,-6.83969 23.005,-14.62c1.46469,-2.06937 16.31313,-23.19312 17.3075,-24.6175c8.51938,-12.10719 8.94938,-29.40125 -1.1825,-42.4625l-8.4925,12.1475c4.15219,9.39281 2.23063,16.99844 -1.6125,22.4675c-0.99437,1.42438 -15.72187,22.52125 -17.2,24.6175c-3.37281,4.79719 -8.42531,7.91469 -14.19,8.9225c-5.76469,1.00781 -11.54281,-0.28219 -16.34,-3.655c-4.79719,-3.35937 -8.02219,-8.42531 -9.03,-14.19c-1.00781,-5.76469 0.29563,-11.54281 3.655,-16.34l20.21,-28.595z"></path></g></g></svg>
                Nothing here
            </div>
        `;
    }
}

function searchCheckouts() {
    let date = getDate();
    let checkoutsinput = document.getElementById("checkoutsinput").value.toLowerCase();

    data.checkouts = data.checkouts == null ? [] : data.checkouts;

    let checkoutlist = document.getElementById("checkoutlist");
    checkoutlist.innerHTML = "";

    if (data.checkouts.length === 0) {
        checkoutlist.innerHTML = `
            <div class="nothing-list">
                <svg xmlns="http://www.w3.org/2000/svg" x="0px" y="0px" width="24" height="24" viewBox="0 0 172 172"><g fill="none" fill-rule="nonzero" stroke="none" stroke-width="1" stroke-linecap="butt" stroke-linejoin="miter" stroke-miterlimit="10" stroke-dasharray="" stroke-dashoffset="0" font-family="none" font-size="none" style="mix-blend-mode: normal"><path d="M0,172v-172h172v172z" fill="none"></path><g fill="#FFFFFF73"><path d="M135.45,0c-6.97406,0.16125 -13.82719,2.33813 -19.6725,6.45l-28.4875,20.1025l7.8475,11.395l28.595,-20.21c4.79719,-3.35937 10.57531,-4.66281 16.34,-3.655c5.76469,1.00781 10.83063,4.23281 14.19,9.03c3.37281,4.79719 4.66281,10.57531 3.655,16.34c-1.00781,5.76469 -4.23281,10.81719 -9.03,14.19c-2.09625,1.47813 -23.08562,16.20563 -24.51,17.2c-5.46906,3.84313 -13.07469,5.76469 -22.4675,1.6125l-12.255,8.4925c13.06125,10.13188 30.46281,9.70188 42.57,1.1825c1.42438,-0.99437 22.54813,-15.84281 24.6175,-17.3075c7.78031,-5.4825 12.98063,-13.61219 14.62,-23.005c1.63938,-9.37937 -0.44344,-18.86625 -5.9125,-26.66c-5.4825,-7.79375 -13.71969,-12.98062 -23.1125,-14.62c-2.33812,-0.40312 -4.66281,-0.59125 -6.9875,-0.5375zM9.5675,3.3325c-2.67406,0.25531 -4.95844,2.05594 -5.83187,4.59563c-0.88688,2.55312 -0.20156,5.375 1.74687,7.22937l37.84,37.84c1.65281,2.05594 4.32688,2.98313 6.89344,2.39188c2.58,-0.59125 4.58219,-2.59344 5.17344,-5.17344c0.59125,-2.56656 -0.33594,-5.24062 -2.39188,-6.89344l-37.84,-37.84c-1.29,-1.35719 -3.07719,-2.13656 -4.945,-2.15c-0.215,-0.01344 -0.43,-0.01344 -0.645,0zM26.5525,87.29l-20.1025,28.595c-5.4825,7.79375 -7.55188,17.17313 -5.9125,26.5525c1.63938,9.39281 6.82625,17.63 14.62,23.1125c7.79375,5.4825 17.28063,7.55188 26.66,5.9125c9.37938,-1.63937 17.53594,-6.83969 23.005,-14.62c1.46469,-2.06937 16.31313,-23.19312 17.3075,-24.6175c8.51938,-12.10719 8.94938,-29.40125 -1.1825,-42.4625l-8.4925,12.1475c4.15219,9.39281 2.23063,16.99844 -1.6125,22.4675c-0.99437,1.42438 -15.72187,22.52125 -17.2,24.6175c-3.37281,4.79719 -8.42531,7.91469 -14.19,8.9225c-5.76469,1.00781 -11.54281,-0.28219 -16.34,-3.655c-4.79719,-3.35937 -8.02219,-8.42531 -9.03,-14.19c-1.00781,-5.76469 0.29563,-11.54281 3.655,-16.34l20.21,-28.595z"></path></g></g></svg>
                Nothing here
            </div>
        `;
    }

    for (let i = 0; i < data.checkouts.length; i++) {
        if (data.checkouts[i].date.getTime() <= date.getTime()) continue;
        let name = data.checkouts[i].name.toLowerCase().replaceAll(" ", "");
        let store = data.checkouts[i].store.toLowerCase().replaceAll(" ", "");
        let size = data.checkouts[i].size.toLowerCase().replaceAll(" ", "");
        let price = data.checkouts[i].price.toString().toLowerCase().replaceAll(" ", "");
        let estsell = data.checkouts[i].est_sell.toString().toLowerCase().replaceAll(" ", "");

        if (name.includes(checkoutsinput) || store.includes(checkoutsinput) || size.includes(checkoutsinput) || price.includes(checkoutsinput) || estsell.includes(checkoutsinput)) {
            let checkout = data.checkouts[i];
            let checkoutdiv = document.createElement("div");
            checkoutdiv.classList.add("checkout");
            checkoutdiv.innerHTML = `
                <img src="${checkout.image}" alt="">
                <div class="description">
                    <h4>${checkout.name}</h4>
                    <h6>Size: ${checkout.size}</h6>
                    <h6>Paid: ${checkout.price}€</h6>
                </div>
            `;

            checkoutlist.appendChild(checkoutdiv);
        }
    }
}

function getDate() {
    let date = new Date();

    if (document.getElementById("daycheckouts").classList.contains("switchactive")) {
        date.setHours(date.getHours() - 12);
    } else if (document.getElementById("weekcheckouts").classList.contains("switchactive")) {
        date.setDate(date.getDate() - 7);
    } else if (document.getElementById("monthcheckouts").classList.contains("switchactive")) {
        date.setDate(date.getDate() - 28);
    }

    return date;
}

function parseDate() {
    data.checkouts = data.checkouts == null ? [] : data.checkouts;
    for (let i = 0; i < data.checkouts.length; i++) {
        data.checkouts[i].date = new Date(data.checkouts[i].date);
    }
}

function changeCheckouts(name) {
    document.getElementById("daycheckouts").classList.remove("switchactive");
    document.getElementById("weekcheckouts").classList.remove("switchactive");
    document.getElementById("monthcheckouts").classList.remove("switchactive");
    document.getElementById(name + "checkouts").classList.add("switchactive");
    setCheckouts();
}
