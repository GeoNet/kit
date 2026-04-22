/* Projects that import version 2 of the GeoNet header basic must contain
the following JS to initialise the header functions */

import { initGeoNetHeader } from "/dependencies/geonet-design-system/js/geonet-design-system.js";

document.addEventListener("DOMContentLoaded", function () {
  initGeoNetHeader();
});
