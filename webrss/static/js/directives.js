var App = angular.module('webrssApp');

App.directive('feedSelect', function() {
    return {
        transclude: false,
        restrict: 'E',
        scope: {
            feeds: '=',
            selected: '=ngModel',
            updateAction: '=',
            deleteAction: '=',
        },
        templateUrl: 'feed-select.html',
        controller: function($scope) {
            $scope.doSelect = function(feed) {
                $scope.selected = feed;
            };
        },
    };
});
