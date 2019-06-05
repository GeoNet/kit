/*********************************************************
 * GeoNet GNSS data plot client application	******
 * -- leaflet map showing sites
 * -- static plot for selected site and data type
 *
 *
 * baishan 30/7/2017
 **********************************************************/

var $ = jQuery.noConflict();

/****** all the chart functionalities defined here ******/
var gnssPlotClient = {
    //### 1. constants and vars
    NZ_CENTRE: new L.LatLng(-40.5, 174.5),
    //sites data by param
    SITES_DATA_URL: "/p/fits/site?typeID=e&methodID=gamit",
    allSitesData: null,
    lftMap: null,

    //### 2. functions
    /***
     * init parameters, called from page
     * ***/
    initParams: function () {
        this.initBaseMap();
        this.initControlFunctions();
        this.showSites();
        this.initEvents();
    },

    /***
     * init leaflet basemap
     * ***/
    initBaseMap: function () {
        var osmUrl = '//{s}.geonet.org.nz/osm/v2/{z}/{x}/{y}.png',
                osmLayer = new L.TileLayer(osmUrl, {
                    minZoom: 1,
                    maxZoom: 18,
                    errorTileUrl: '//static.geonet.org.nz/osm/images/logo_geonet.png',
                    subdomains: ['static1', 'static2', 'static3', 'static4', 'static5']
                });

        this.lftMap = L.map('gnss-sites-map', {
            attributionControl: false,
            zoom: 18,
            layers: [osmLayer]
        });

        //L.control.layers(baseLayers).addTo(this.lftMap);
        this.lftMap.setView(this.NZ_CENTRE, 5);
    },

    /**
     * get the copyrights popover working
     */
    initControlFunctions: function () {
        jQuery("a[rel=popover]")
                .popover({
                    html: true,
                    trigger: 'hover',
                    delay: {
                        hide: 500,
                        show: 80
                    },
                    placement: function (tip, element) {
                        var offset = jQuery(element).offset();
                        var height = jQuery(document).outerHeight();
                        var width = jQuery(document).outerWidth();
                        var vert = 0.5 * height - offset.top;
                        var vertPlacement = vert > 0 ? 'bottom' : 'top';
                        var horiz = 0.5 * width - offset.left;
                        var horizPlacement = horiz > 0 ? 'right' : 'left';
                        var placement = Math.abs(horiz) > Math.abs(vert) ? horizPlacement : vertPlacement;
                        return placement;
                    }
                })
                .click(function (e) {
                    e.preventDefault();
                });

    },

    /***
     * show sites on map
     * ***/
    showSitesDataOnMap: function (sitesJson) {
        var sitesLayer = new L.GeoJSON1(sitesJson, {

            pointToLayer: function (feature, latlng) {
                var iconSize = 12;
                var iconSizeLarge = 18;
                var svgURL = gnssPlotClient.getSiteIconSVG('#3366AA');
                var svgURLSelected = gnssPlotClient.getSiteIconSVG('#f4f450');
                gnssPlotClient.sVGIcon = L.icon({
                    iconUrl: svgURL,
                    iconSize: [iconSize, iconSize],
                    shadowSize: [0, 0],
                    iconAnchor: [0.5 * iconSize, 0.5 * iconSize],
                    popupAnchor: [0, 0]
                });
                gnssPlotClient.sVGIconSelected = L.icon({
                    iconUrl: svgURLSelected,
                    iconSize: [iconSizeLarge, iconSizeLarge],
                    shadowSize: [0, 0],
                    iconAnchor: [0.5 * iconSizeLarge, 0.5 * iconSizeLarge],
                    popupAnchor: [0, 0]
                });

                var marker = L.marker(latlng, {
                    icon: gnssPlotClient.sVGIcon,
                    riseOnHover: true,
                    zIndexOffset: 150
                });

                marker.on('click', function () {
                    if (gnssPlotClient.selectedMarker) {//if any marker is already selected, set to normal
                        gnssPlotClient.selectedMarker.setIcon(gnssPlotClient.sVGIcon);
                        gnssPlotClient.selectedMarker.setZIndexOffset(150);
                    }
                    gnssPlotClient.selectedMarker = this;
                    gnssPlotClient.showSitePlots(feature.properties.siteID, feature.properties.networkID, feature.properties.name);
                    this.setIcon(gnssPlotClient.sVGIconSelected);
                    this.setZIndexOffset(160);
                });

                return marker;
            }
        });
        this.lftMap.addLayer(sitesLayer);
        sitesLayer.checkFeatureLocation();
    },

    getSiteIconSVG: function (iconColor, borderColor) {
        var iconSize = 6,
                borderWidth = 1;
        if (!iconColor) {
            iconColor = '#3366AA';
        }
        if (!borderColor) {
            borderColor = '#111c1c';
        }
        // here's the trick, base64 encode the URL
        var svgIcon = "<svg xmlns='http://www.w3.org/2000/svg' version='1.1'"
                + " width='" + (2 * (iconSize + borderWidth)) + "' height='" + (2 * (iconSize + borderWidth)) + "'>"
                + "<path d='M" + (iconSize + borderWidth) + " 0 L0 " + 1.8 * (iconSize + borderWidth) + "L" + 2 * (iconSize + borderWidth) + " " + 1.8 * (iconSize + borderWidth) + " Z' stroke='"
                + borderColor + "' stroke-width='" + borderWidth + "' fill='" + iconColor + "' /></svg>";

        return  "data:image/svg+xml;base64," + btoa(svgIcon);

    },

    showSitePlots: function (siteId, networkId, siteName) {
        if (siteId && siteName) {
            this.siteId = siteId;
            this.networkId = networkId;
            this.siteName = siteName;
        }
        var plottype = jQuery("input[name='plotRadios']:checked").val();
        var plotTypes;
        var title;
        if (plottype === 'displacement') {
            plotTypes = ['e', 'n', 'u'];
        } else {
            plotTypes = ['mp1', 'mp2'];
        }
        if (this.siteId) {
            $('#gnss-chart').empty();
            this.allImages = plotTypes.length;
            this.imagesLoaded = 0;
            for (var i = 0; i < plotTypes.length; i++) {
                var imgurl = '/data/gnss/plot/' + this.siteId + '/' + plotTypes[i];
                if (this.networkId) {
                    imgurl += '/' + this.networkId;
                }
                //console.log("imgurl " + imgurl);
                $('<img/>', {
                    'src': imgurl,
                    'width': "100%"
                }).on('load', function () {
                    gnssPlotClient.onImageLoad();
                }).appendTo('#gnss-chart');

                $('<br/>').appendTo('#gnss-chart');
            }
            //change layout
            jQuery("#gnss-map-container").removeClass("col-lg-10").addClass("col-lg-5");
            this.lftMap.invalidateSize();
            //show plots
            jQuery("#gnss-plot-container").css("display", "block");
            title = "GNSS Time Series Plot - " + this.siteName;
            jQuery("#gnss-plot-header").html(title);
        }

    },

    onImageLoad: function () {
        this.imagesLoaded++;
        if (this.imagesLoaded >= this.allImages) {//when all images have been loaded
            //scroll to plot on small screen size
            if (jQuery(window).width() < 992) {//
                var element = jQuery("#gnss-plot-container");
                jQuery('html, body').scrollTop(element.offset().top);
            }
        }
    },

    initEvents: function () {
        jQuery("input[name='plotRadios']").on('change', function () {
            gnssPlotClient.showSitePlots();
        });

    },

    /***
     * query sites from http
     * ***/
    showSites: function () {
        var sitesData = this.allSitesData;
        //console.log("sitesData " + sitesData);
        if (sitesData) {
            this.showSitesDataOnMap(sitesData);
        } else {
            var url = this.SITES_DATA_URL;
            //console.log("show sites " + " url " + url);
            jQuery.getJSON(url, function (data) {
                //console.log(JSON.stringify(data));
                gnssPlotClient.allSitesData = data;
                gnssPlotClient.showSitesDataOnMap(data);
            });
        }
    }
};

