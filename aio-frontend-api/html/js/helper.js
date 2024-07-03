function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

function getCookie(name) {
    function escape(s) { return s.replace(/([.*+?\^$(){}|\[\]\/\\])/g, '\\$1'); }
    let match = document.cookie.match(RegExp('(?:^|;\\s*)' + escape(name) + '=([^;]*)'));
    return match ? match[1] : null;
}

function removeCookie(name) {
    document.cookie = name + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT; domain=.copped-inc.com';
    document.cookie = name + '=; expires=Thu, 01 Jan 1970 00:00:01 GMT';
}
