App.controller('RSSCtrl', function($scope, $http, $sce, $modal) {
    'use strict';
    $scope.feeds = {
        categories: [], // all feeds
        selected: null, // currently selected feed
        entries: {
            list: [], // entries in feed
            current: null, // currently selected entry
        },
    };

    $scope.createUpdateCategory = function(category) {
        $modal.open({
            templateUrl: 'category_form.html',
            controller: 'RSSAddEditCategoryCtrl',
            size: 'small',
            resolve: {
                category: function() { return category || []; },
                parentScope: function() { return $scope; }
            }
        });
    };

    $scope.deleteCategory = function(category) {
        $modal.open({
            templateUrl: 'category_delete.html',
            controller: 'RSSDeleteCategoryCtrl',
            size: 'small',
            resolve: {
                category: function() { return category; },
                parentScope: function() { return $scope; }
            }
        });
    };

    $scope.createFeed = function(category) {
        $modal.open({
            templateUrl: 'feed_create.html',
            controller: 'RSSCreateFeedCtrl',
            size: 'small',
            resolve: {
                category: function() { return category; },
                parentScope: function() { return $scope; }
            }
        });
    };

    $scope.updateFeed = function(feed) {
        $modal.open({
            templateUrl: 'feed_update.html',
            controller: 'RSSUpdateFeedCtrl',
            size: 'small',
            resolve: {
                feed: function() { return feed; },
                parentScope: function() { return $scope; }
            }
        });
    };

    $scope.deleteFeed = function(feed) {
        $modal.open({
            templateUrl: 'feed_delete.html',
            controller: 'RSSDeleteFeedCtrl',
            size: 'small',
            resolve: {
                feed: function() { return feed; },
                parentScope: function() { return $scope; }
            }
        });
    };

    $scope.moveUpCategory = function(category) {
        $http.post('/api/category/' + category.id + '/move_up')
            .then(function() {
                $scope.loadCategories();
            });
    };

    $scope.moveDownCategory = function(category) {
        $http.post('/api/category/' + category.id + '/move_down')
            .then(function() {
                $scope.loadCategories();
            });
    };

    $scope.loadCategories = function() {
        $http.get('/api/category')
            .then(function(res) {
                $scope.feeds.categories = res.data;
            });

        setTimeout($scope.loadCategories, 60000);
    };
    $scope.loadCategories();

    $scope.$watch('feeds.selected', function(newValue, oldValue) {
        if (!$scope.feeds.selected) return;
        $scope.feeds.selected.new_entries = false;

        $http.get('/api/entry/?feed=' + $scope.feeds.selected.id)
            .then(function(res) {
                $scope.feeds.entries.list = res.data;
                $scope.feeds.entries.current = null;
            });
    });

    $scope.$watch('feeds.entries.current', function() {
        if (!$scope.feeds.entries.current) return;

        $scope.feeds.entries.current.new_entry = false;

        $http.get('/api/entry/' + $scope.feeds.entries.current.id)
            .then(function(res) {
                var feed, feeds = [];

                _.forEach($scope.feeds.categories.objects, function(category) {
                    feeds.push.apply(feeds, category.feeds);
                });

                feed = _.find(feeds, function(obj) {
                    return obj.id === $scope.feeds.entries.current.feed.id;
                });

                if (feed)
                    feed.un_read -= 1;

                if (!$scope.feeds.entries.current.read_at)
                    $scope.feeds.entries.current.read_at = new Date();
            });
    });

    $scope.loadMore = function(feedUrl) {
        $http.get(feedUrl)
            .then(function(res) {
                var feedEntries = $scope.feeds.entries.list.objects;
                feedEntries.push.apply(feedEntries, res.data.objects);
                $scope.feeds.entries.list.objects = feedEntries;
                $scope.feeds.entries.list.meta = res.data.meta;
            });
    };

    $scope.safe = $sce.trustAsHtml;

    $scope.doSearch = function() {
        $http.get('/api/entry/?title__ilike=%' + $scope.search + '%')
            .then(function(res) {
                $scope.feeds.entries.list = res.data;
                $scope.feeds.selected = null;
                $scope.feeds.entries.current = null;
            });

    };
}).controller('RSSAddEditCategoryCtrl',
function($scope, $modalInstance, $http, category, parentScope) {
    'use strict';
    $scope.category = category;
    $scope.form = angular.copy(category);

    $scope.save = function() {
        var method;
        if (category.id === undefined) {
            method = $http.post(
                '/api/category/',
                {title: $scope.form.title}
            );
        } else {
            method = $http.post(
                '/api/category/' + category.id + '/',
                {title: $scope.form.title}
            );
        }

        method.then(function(res) {
            parentScope.loadCategories();
            $modalInstance.close();
        }, function() {
            $scope.error = 'Something went wrong';
        });
    };

    $scope.cancel = function() {
        $modalInstance.dismiss();
    };
}).controller('RSSCreateFeedCtrl',
function($scope, $modalInstance, $http, category, parentScope) {
    'use strict';
    $scope.form = {feed_url: '', category: ''};

    if (category !== undefined && category)
        $scope.form.category = category.id;

    $scope.categories = parentScope.feeds.categories.objects;

    $scope.save = function() {
        $http.post('/api/feed/', $scope.form)
            .then(function(res) {
                parentScope.loadCategories();
                $modalInstance.close();
            }, function() {
                $scope.error = 'Something went wrong';
            });
    };

    $scope.cancel = function() {
        $modalInstance.dismiss();
    };
}).controller('RSSUpdateFeedCtrl',
function($scope, $modalInstance, $http, feed, parentScope) {
    'use strict';
    $scope.feed = feed;
    $scope.form = angular.copy(feed);
    $scope.form.category = $scope.form.category.toString();
    delete $scope.form.un_read;
    delete $scope.form.new_entries;

     $scope.categories = parentScope.feeds.categories.objects;

    $scope.save = function() {
        $http.put('/api/feed/' + feed.id + '/', $scope.form)
            .then(function(res) {
                parentScope.loadCategories();
                $modalInstance.close();
            }, function() {
                $scope.error = 'Something went wrong';
            });
    };

    $scope.cancel = function() {
        $modalInstance.dismiss();
    };
}).controller('RSSDeleteCategoryCtrl',
function($scope, $modalInstance, $http, category, parentScope) {
    'use strict';
    $scope.category = category;
    $scope.ok = function() {
        $http.delete('/api/category/' + category.id + '/')
            .then(function(res) {
                parentScope.loadCategories();
                $modalInstance.close();
            }, function() {
                $scope.error = 'Something went wrong';
            });
    };

    $scope.cancel = function() {
        $modalInstance.dismiss();
    };
}).controller('RSSDeleteFeedCtrl',
function($scope, $modalInstance, $http, feed, parentScope) {
    'use strict';
    $scope.feed = feed;
    $scope.ok = function() {
        $http.delete('/api/feed/' + feed.id + '/')
            .then(function(res) {
                parentScope.loadCategories();
                $modalInstance.close();
            }, function() {
                $scope.error = 'Something went wrong';
            });
    };

    $scope.cancel = function() {
        $modalInstance.dismiss();
    };
});
