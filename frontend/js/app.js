let App = angular.module(
    "webrssApp",
    ["ngSanitize", "ngResource", "ui.bootstrap", "webrssApp.templates"]
)

App.config($locationProvider => {
    $locationProvider.hashPrefix("")
})
