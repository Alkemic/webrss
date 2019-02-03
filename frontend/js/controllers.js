App.controller("RSSCtrl", function($scope, $http, $sce, $uibModal, $location) {
    "use strict"
    $scope.loading = false
    $scope.failedLoadCategories = false
    $scope.feeds = {
        categories: [], // all feeds
        selected: null, // currently selected feed
        entries: {
            list: [], // entries in feed
            current: null, // currently selected entry
        },
    }

    $scope.toggleEntrySelect = function(entry) {
        if ($scope.feeds.entries.current === entry) {
            $scope.feeds.entries.current = null
        } else {
            $scope.feeds.entries.current = entry
        }
    }

    $scope.createUpdateCategory = function(category) {
        $uibModal.open({
            templateUrl: "category_form.html",
            controller: "RSSAddEditCategoryCtrl",
            size: "small",
            resolve: {
                category: function() { return category || [] },
                parentScope: function() { return $scope }
            }
        })
    }

    $scope.deleteCategory = function(category) {
        $uibModal.open({
            templateUrl: "category_delete.html",
            controller: "RSSDeleteCategoryCtrl",
            size: "small",
            resolve: {
                category: function() { return category },
                parentScope: function() { return $scope }
            }
        })
    }

    $scope.createFeed = function(category) {
        $uibModal.open({
            templateUrl: "feed_create.html",
            controller: "RSSCreateFeedCtrl",
            size: "small",
            resolve: {
                category: function() { return category },
                parentScope: function() { return $scope }
            }
        })
    }

    $scope.updateFeed = function(feed) {
        $uibModal.open({
            templateUrl: "feed_update.html",
            controller: "RSSUpdateFeedCtrl",
            size: "small",
            resolve: {
                feed: function() { return feed },
                parentScope: function() { return $scope }
            }
        })
    }

    $scope.deleteFeed = function(feed) {
        $uibModal.open({
            templateUrl: "feed_delete.html",
            controller: "RSSDeleteFeedCtrl",
            size: "small",
            resolve: {
                feed: function() { return feed },
                parentScope: function() { return $scope }
            }
        })
    }

    $scope.moveUpCategory = function(category) {
        $http.post(`/api/category/${category.id}/move_up`)
            .then(function() {
                $scope.loadCategories(false)
            })
    }

    $scope.moveDownCategory = function(category) {
        $http.post(`/api/category/${category.id}/move_down`)
            .then(function() {
                $scope.loadCategories(false)
            })
    }

    $scope.$watch("feeds.selected", (feed) => {
        if (!feed) return
        feed.new_entries = false

        let slug = `${feed.id}-${feed.feed_title.toLowerCase().replace(/ /g, "-")}`
        if($location.url() !== slug) $location.url(slug)
    })

    let onChangeFeed = () => {
        let match = /^\/(\d+)-.*/.exec($location.url())
        if(!!match) {
            if(!$scope.feeds.categories.objects) return
            let feed = null, feedId = parseInt(match[1])
            $scope.feeds.categories.objects.forEach((category) => {
                category.feeds.forEach((_feed) => {
                    if(_feed.id === feedId) {
                        feed = _feed
                    }
                })
            })
            $scope.feeds.selected = feed
            $http.get(`/api/entry/?feed=${feed.id}`)
                .then(function(res) {
                    $scope.feeds.entries.list = res.data
                    $scope.feeds.entries.current = null
                })
        } else {
            $scope.feeds.selected = null
        }
    }

    let initialLoading = true
    $scope.loadCategories = function (quiet) {
        quiet = typeof quiet !== "undefined" ? quiet : true

        $scope.loading = !quiet
        $http.get("/api/category/")
            .then(function(res) {
                $scope.failedLoadCategories = false
                $scope.feeds.categories = res.data
                $scope.loading = false
                if(initialLoading) {
                    onChangeFeed()
                    initialLoading = false
                    $scope.$on("$locationChangeSuccess", onChangeFeed)
                }
                setTimeout($scope.loadCategories, 60000)
            }, function(data) {
                $scope.failedLoadCategories = true
                console.error("Error loading data", data)
                $scope.loading = false
            })
    }
    $scope.loadCategories(false)

    $scope.loadCategoriesFn = function (e) {
        e.preventDefault()
        $scope.loadCategories(false)
        $scope.failedLoadCategories = false
    }

    $scope.$watch("feeds.entries.current", function() {
        if (!$scope.feeds.entries.current) return

        $scope.feeds.entries.current.new_entry = false

        $http.get(`/api/entry/${$scope.feeds.entries.current.id}`)
            .then(function(res) {
                let feed, feeds = []

                _.forEach($scope.feeds.categories.objects, function(category) {
                    feeds.push.apply(feeds, category.feeds)
                })

                feed = _.find(feeds, function(obj) {
                    return obj.id === $scope.feeds.entries.current.feed.id
                })

                if (feed)
                    feed.un_read -= 1

                if (!$scope.feeds.entries.current.read_at)
                    $scope.feeds.entries.current.read_at = new Date()
            })
    })

    $scope.loadMore = function(feedUrl) {
        $http.get(feedUrl)
            .then(function(res) {
                let feedEntries = $scope.feeds.entries.list.objects
                feedEntries.push.apply(feedEntries, res.data.objects)
                $scope.feeds.entries.list.objects = feedEntries
                $scope.feeds.entries.list.meta = res.data.meta
            })
    }

    $scope.safe = $sce.trustAsHtml

    $scope.doSearch = function() {
        $scope.loading = true
        $http.get(`/api/entry/?title__ilike=%${$scope.search}%`)
            .then(function(res) {
                $scope.feeds.entries.list = res.data
                $scope.feeds.selected = null
                $scope.feeds.search = true
                $scope.feeds.entries.current = null
                $scope.loading = false
            }, function (err) {
                alert("Error fetching search results.")
                console.error(err)
                $scope.loading = false
            })

    }
}).controller("RSSAddEditCategoryCtrl",
function($scope, $uibModalInstance, $http, category, parentScope) {
    "use strict"
    $scope.category = category
    $scope.form = angular.copy(category)

    $scope.save = function() {
        let method
        if (category.id === undefined) {
            method = $http.post(
                "/api/category/",
                {title: $scope.form.title}
            )
        } else {
            method = $http.post(
                `/api/category/${category.id}/`,
                {title: $scope.form.title}
            )
        }

        method.then(function(res) {
            parentScope.loadCategories(false)
            $uibModalInstance.close()
        }, function() {
            $scope.error = "Something went wrong"
        })
    }

    $scope.cancel = function() {
        $uibModalInstance.dismiss()
    }
}).controller("RSSCreateFeedCtrl",
function($scope, $uibModalInstance, $http, category, parentScope) {
    "use strict"
    $scope.form = {feed_url: "", category: ""}

    if (category !== undefined && category)
        $scope.form.category = category.id

    $scope.categories = parentScope.feeds.categories.objects

    $scope.save = function() {
        $http.post("/api/feed/", $scope.form)
            .then(function(res) {
                parentScope.loadCategories(false)
                $uibModalInstance.close()
            }, function() {
                $scope.error = "Something went wrong"
            })
    }

    $scope.cancel = function() {
        $uibModalInstance.dismiss()
    }
}).controller("RSSUpdateFeedCtrl",
function($scope, $uibModalInstance, $http, feed, parentScope) {
    "use strict"
    $scope.categories = parentScope.feeds.categories.objects
    $scope.form = angular.copy(feed)
    $scope.form.category = $scope.categories.filter(c => c.id === feed.category)[0]
    delete $scope.form.un_read
    delete $scope.form.new_entries

    $scope.save = function() {
        let formData = angular.copy($scope.form)
        formData.category = $scope.form.category.id
        $http.put(`/api/feed/${feed.id}/`, formData)
            .then(function(res) {
                parentScope.loadCategories(false)
                $uibModalInstance.close()
            }, function() {
                $scope.error = "Something went wrong"
            })
    }

    $scope.cancel = $uibModalInstance.dismiss
}).controller("RSSDeleteCategoryCtrl",
function($scope, $uibModalInstance, $http, category, parentScope) {
    "use strict"
    $scope.category = category
    $scope.ok = function() {
        $http.delete(`/api/category/${category.id}/`)
            .then(function(res) {
                parentScope.loadCategories(false)
                $uibModalInstance.close()
            }, function() {
                $scope.error = "Something went wrong"
            })
    }

    $scope.cancel = function() {
        $uibModalInstance.dismiss()
    }
}).controller("RSSDeleteFeedCtrl",
function($scope, $uibModalInstance, $http, feed, parentScope) {
    "use strict"
    $scope.feed = feed
    $scope.ok = function() {
        $http.delete(`/api/feed/${feed.id}/`)
            .then(function(res) {
                parentScope.loadCategories(false)
                $uibModalInstance.close()
            }, function() {
                $scope.error = "Something went wrong"
            })
    }

    $scope.cancel = function() {
        $uibModalInstance.dismiss()
    }
})
