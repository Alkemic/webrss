'use strict';

App.controller("RSSCtrl", function ($scope, $http, $sce, $modal) {
    $scope.feedEntries = [];
    $scope.categories = [];
    $scope.test = 'asd';

    $scope.createUpdateCategory = function (category) {
        $modal.open({
            templateUrl: 'category_create.html',
            controller: 'RSSAddEditCategoryCtrl',
            size: 'small',
            resolve: {
                category: function () { return category || []; },
                parentScope: function () { return $scope; }
            }
        });
    };

    $scope.deleteCategory = function (category) {
        $modal.open({
            templateUrl: 'category_delete.html',
            controller: 'RSSDeleteCategoryCtrl',
            size: 'small',
            resolve: {
                category: function () { return category; },
                parentScope: function () { return $scope; }
            }
        });
    };

    $scope.createFeed = function (category) {
        $modal.open({
            templateUrl: 'feed_create.html',
            controller: 'RSSCreateFeedCtrl',
            size: 'small',
            resolve: {
                category: function () { return category; },
                parentScope: function () { return $scope; }
            }
        });
    };

    $scope.updateFeed = function ($event, feed) {
        $event.preventDefault();

        $modal.open({
            templateUrl: 'feed_update.html',
            controller: 'RSSUpdateFeedCtrl',
            size: 'small',
            resolve: {
                feed: function () { return feed; },
                parentScope: function () { return $scope; }
            }
        });
    };

    $scope.deleteFeed = function ($event, feed) {
        $event.preventDefault();

        $modal.open({
            templateUrl: 'feed_delete.html',
            controller: 'RSSDeleteFeedCtrl',
            size: 'small',
            resolve: {
                feed: function () { return feed; },
                parentScope: function () { return $scope; }
            }
        });
    };

    $scope.moveUpCategory = function (category) {
        $http.post('/api/category/' + category.id + '/move_up')
            .then(function () {
                $scope.loadCategories();
            });
    };

    $scope.moveDownCategory = function (category) {
        $http.post('/api/category/' + category.id + '/move_down')
            .then(function () {
                $scope.loadCategories();
            });
    };

    $scope.loadCategories = function () {
        $http.get('/api/category')
            .then(function (res) {
                $scope.categories = res.data;
            });

        setTimeout($scope.loadCategories, 60000);
    };
    $scope.loadCategories();

    $scope.loadFeed = function (feed) {
        _.each($scope.categories.objects, function (category) {
            _.each(category.feeds, function (feed) {
                feed.active = false;
            });
        });
        feed.active = true;
        feed.new_entries = false;

        $http.get('/api/entry/?feed=' + feed.id)
            .then(function (res) {
                $scope.feedEntries = res.data;
            });
    };

    $scope.loadMore = function (feedUrl) {
        $http.get(feedUrl)
            .then(function (res) {
                var feedEntries = $scope.feedEntries.objects;
                feedEntries.push.apply(feedEntries, res.data.objects);
                $scope.feedEntries.objects = feedEntries;
                $scope.feedEntries.meta = res.data.meta;
            });
    };

    $scope.loadEntry = function (entry) {
        _.each($scope.feedEntries.objects, function (entry) {
            entry.active = false;
        });
        entry.active = true;

        $http.get('/api/entry/' + entry.id)
            .then(function (res) {
                $scope.feedEntry = res.data;

                var feeds = [];
                _.forEach($scope.categories.objects, function (category) {
                    feeds.push.apply(feeds, category.feeds);
                });

                var feed = _.find(feeds, function (obj) {
                    return obj.id === entry.feed.id;
                });

                if(feed && !entry.read_at) feed.un_read -= 1;
                entry.read_at = new Date();
            });
    };

    $scope.safe = $sce.trustAsHtml;

    $scope.doSearch = function () {
        _.each($scope.categories.objects, function (category) {
            _.each(category.feeds, function (feed) {
                feed.active = false;
            });
        });

        $http.get('/api/entry/?title__ilike=%' + $scope.search + '%')
            .then(function (res) {
                $scope.feedEntries = res.data;
            });

    };
}).controller("RSSAddEditCategoryCtrl", function ($scope, $modalInstance, $http, category, parentScope) {
    $scope.category = category;
    $scope.save = function () {
        var method;
        if(category.id === undefined)
            method = $http.post('/api/category/', {'title': $scope.category.title});
        else
            method = $http.post('/api/category/' + category.id + '/' , {'title': $scope.category.title});

        method.then(function (res) {
            parentScope.loadCategories();
            $modalInstance.close();
        }, function () {
            $scope.error = 'Something went wrong';
        });
    };

    $scope.cancel = function () {
        $modalInstance.dismiss();
    };
}).controller("RSSCreateFeedCtrl", function ($scope, $modalInstance, $http, category, parentScope) {
    $scope.form = {feed_url: '', category: ''};

    if(category !== undefined && category)
        $scope.form.category = category.id;

    $scope.categories = parentScope.categories.objects;

    $scope.save = function () {
        $http.post('/api/feed/', $scope.form)
            .then(function (res) {
                parentScope.loadCategories();
                $modalInstance.close();
            }, function () {
                $scope.error = 'Something went wrong';
            });
    };

    $scope.cancel = function () {
        $modalInstance.dismiss();
    };
}).controller('RSSUpdateFeedCtrl', function ($scope, $modalInstance, $http, feed, parentScope) {
    $scope.form = angular.copy(feed);
    $scope.form.category = $scope.form.category.toString();
    delete $scope.form.un_read;
    delete $scope.form.active;
    delete $scope.form.new_entries;

    $scope.categories = parentScope.categories.objects;

    $scope.save = function () {
        $http.put('/api/feed/' + feed.id + '/', $scope.form)
            .then(function (res) {
                parentScope.loadCategories();
                $modalInstance.close();
            }, function () {
                $scope.error = 'Something went wrong';
            });
    };

    $scope.cancel = function () {
        $modalInstance.dismiss();
    };
}).controller("RSSDeleteCategoryCtrl", function ($scope, $modalInstance, $http, category, parentScope) {
    $scope.category = category;
    $scope.ok = function () {
        $http.delete('/api/category/' + category.id + '/')
            .then(function (res) {
                parentScope.loadCategories();
                $modalInstance.close();
            }, function () {
                $scope.error = 'Something went wrong';
            });
    };

    $scope.cancel = function () { $modalInstance.dismiss(); };
}).controller("RSSDeleteFeedCtrl", function ($scope, $modalInstance, $http, feed, parentScope) {
    $scope.feed = feed;
    $scope.ok = function () {
        $http.delete('/api/feed/' + feed.id + '/')
            .then(function (res) {
                parentScope.loadCategories();
                $modalInstance.close();
            }, function () {
                $scope.error = 'Something went wrong';
            });
    };

    $scope.cancel = function () { $modalInstance.dismiss(); };
});
