App.directive("feedSelect", () => ({
    transclude: false,
    restrict: "E",
    scope: {
        feeds: "=",
        selected: "=ngModel",
        updateAction: "=",
        deleteAction: "="
    },
    templateUrl: "feed-select.html",
    controller: $scope => {
        $scope.doSelect = feed => {
            $scope.selected = feed
        }
    },
}))
