/* Projects that import the GeoNet footer must contain the following JS so that
the app link when on mobile changes based on the device's operating system. */

/**
 * If the user is using a mobile device, this sets the link to the app to either
 * the Android Play Store or the Apple App Store depending on the operating system of
 * the device. If it's an unknown OS, the link is hidden.
 */
function setLinkToApp() {
    var linkContainer = document.getElementById('appLinkContainer');
    var linkEl = document.getElementById('appLink');
    var osIconEl = document.getElementById('osIcon');
    if (!linkEl || !osIconEl ) {
        return;
    }
    var os = getMobileOperatingSystem();
    switch(os) {
        case "Android":
            linkEl.setAttribute('href', 'https://play.google.com/store/apps/details?id=nz.org.geonet.quake&hl=en');
            osIconEl.classList = 'fa-brands fa-android fa-1';
            break;
        case "iOS":
            linkEl.setAttribute('href', 'https://itunes.apple.com/nz/app/geonet-quake/id533054360?mt=8');
            osIconEl.classList = 'fa-brands fa-apple fa-1';
            break;
        default:
            linkContainer.classList.replace("d-md-none", "d-none");
            break;
    }
};

/**
 * Determines the mobile OS, returning 'iOS', 'Android', or 'unknown'.
 * @returns {string} "Android", "iOS", or "unknown"
 */
function getMobileOperatingSystem() {
    var userAgent = navigator.userAgent || navigator.vendor || window.opera;

    if (/android/i.test(userAgent)) {
        return "Android";
    }
    if (/iPad|iPhone|iPod/.test(userAgent) && !window.MSStream) {
        return "iOS";
    }
    return "unknown";
}