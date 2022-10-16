/* Projects that import the GeoNet header must contain the following JS so that
the search inputs function. This is on top of the standard
GeoNet Bootstrap.js file (which also needs to be present).*/

document.addEventListener("DOMContentLoaded", function() {

    // To have more control of when the search text field collapses,
    // we remove the collapse target after it's shown, and re-add it
    // if the text field is blank. This means it won't collapse when a
    // valid search is made. This might be hacky but it works.

    const searchContainer = document.getElementById('searchContainer');
    searchContainer.addEventListener('shown.bs.collapse', function(e) {
        const searchBtn = document.getElementById("searchBtn");
        searchBtn.dataset.bsTarget = "";
    })

    // Set value to blank so when field is shown because when the user
    // navigates back from the search page, the value is still there, which
    // triggers a search when it's first shown - we don't want this.
    searchContainer.addEventListener('show.bs.collapse', function(e) {
        document.getElementById("search_query").value = "";
    })

    const form = document.forms["search_form"];
    form.onsubmit = function(e) {
        const query = document.getElementById("search_query").value;
        if (query !== "") {
            // Make search
            return true;
        }
        // Re-add collapse behaviour, trigger collapse manually to hide text
        // field, don't search.
        const searchBtn = document.getElementById("searchBtn");
        searchBtn.dataset.bsTarget = "#searchContainer";
        const collapsibleInstance = bootstrap.Collapse.getInstance(searchContainer);
        collapsibleInstance.hide();
        return false;
    }

    // Just before the page leaves to go to search page, but after it's hidden to the
    // user, collapse the search form so that it's collapsed when navigating back.
    // Note: seems to just be needed for Safari.
    document.addEventListener('visibilitychange', event => {
        if (event.target.visibilityState == "hidden") {
            const searchBtn = document.getElementById("searchBtn");
            searchBtn.dataset.bsTarget = "#searchContainer";
            const collapsibleInstance = bootstrap.Collapse.getInstance(searchContainer);
            collapsibleInstance.hide();
        }
    });

    const formMobile = document.forms["search_form_2"];
    formMobile.onsubmit = function(e) {
        return document.getElementById("search_query_2").value !== "";
    }
})


