var App = angular.module('webrssApp');

App.filter('stripTags', function() {
    return function(text) {
        if (!text) {
            return '';
        }
        return String(text).replace(/<[^>]+>/gm, '').replace(/&[^;]+;/gm, '');
    };
});
