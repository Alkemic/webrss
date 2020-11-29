App.controller("RSSCtrl", ($scope, $http, $sce, $uibModal, $location) => {
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

    $scope.toggleEntrySelect = entry => {
        if ($scope.feeds.entries.current === entry) {
            $scope.feeds.entries.current = null
        } else {
            $scope.feeds.entries.current = entry
        }
    }

    $scope.createUpdateCategory = category => {
        $uibModal.open({
            templateUrl: "category_form.html",
            controller: "RSSAddEditCategoryCtrl",
            size: "small",
            resolve: {
                category: () => category || [],
                parentScope: () => $scope,
            },
        })
    }

    $scope.deleteCategory = category => {
        $uibModal.open({
            templateUrl: "category_delete.html",
            controller: "RSSDeleteCategoryCtrl",
            size: "small",
            resolve: {
                category: () => category,
                parentScope: () => $scope,
            },
        })
    }

    $scope.createFeed = category => {
        $uibModal.open({
            templateUrl: "feed_create.html",
            controller: "RSSCreateFeedCtrl",
            size: "small",
            resolve: {
                category: () => category,
                parentScope: () => $scope,
            }
        })
    }

    $scope.editUser = user => {
        $uibModal.open({
            templateUrl: "edit_user.html",
            controller: "EditUserCtrl",
            size: "small",
            resolve: {
                user: () => user,
                parentScope: () => $scope,
            }
        })
    }

    $scope.updateFeed = feed => {
        $uibModal.open({
            templateUrl: "feed_update.html",
            controller: "RSSUpdateFeedCtrl",
            size: "small",
            resolve: {
                feed: () => feed,
                parentScope: () => $scope,
            },
        })
    }

    $scope.deleteFeed = feed => {
        $uibModal.open({
            templateUrl: "feed_delete.html",
            controller: "RSSDeleteFeedCtrl",
            size: "small",
            resolve: {
                feed: () => feed,
                parentScope: () => $scope,
            }
        })
    }

    $scope.moveUpCategory = category => {
        $http.post(`/api/category/${category.id}/move_up`)
            .then(() => {
                $scope.loadCategories(false)
            })
    }

    $scope.moveDownCategory = category => {
        $http.post(`/api/category/${category.id}/move_down`)
            .then(() => {
                $scope.loadCategories(false)
            })
    }

    $scope.$watch("feeds.selected", feed => {
        if (!feed) return
        feed.new_entries = false

        let slug = `${feed.id}-${feed.feed_title.toLowerCase().replace(/ /g, "-")}`
        if ($location.url() !== slug) $location.url(slug)
    })

    let onChangeUrl = () => {
        if (!!(match = /^\/(\d+)-.*/.exec($location.url()))) {
            if (!$scope.feeds.categories.objects) return

            let feed = null, feedId = parseInt(match[1])
            $scope.feeds.categories.objects.forEach(category => {
                category.feeds.forEach(_feed => {
                    if (_feed.id === feedId) {
                        feed = _feed
                    }
                })
            })
            $scope.feeds.selected = feed
            $http.get(`/api/entry/?feed=${feed.id}`)
                .then(res => {
                    $scope.feeds.entries.list = res.data
                    $scope.feeds.entries.current = null
                })
        } else {
            $scope.feeds.selected = null
        }

        if (!!(match = /^\/search=(.*)/.exec($location.url()))) {
            let phrase = decodeURI(match[1])
            $scope.search = phrase
            $http.get(`/api/entry/search/?phrase=${phrase}`)
                .then(res => {
                    $scope.feeds.entries.list = res.data
                    $scope.feeds.selected = null
                    $scope.feeds.search = true
                    $scope.feeds.entries.current = null
                    $scope.loading = false
                }, err => {
                    alert("Error fetching search results.")
                    console.error(err)
                    $scope.loading = false
                })
        }
    }

    let initialLoading = true
    $scope.loadCategories = quiet => {
        quiet = typeof quiet !== "undefined" ? quiet : true

        $scope.loading = !quiet
        $http.get("/api/category/")
            .then(res => {
                $scope.failedLoadCategories = false
                $scope.feeds.categories = res.data
                $scope.loading = false
                if (initialLoading) {
                    onChangeUrl()
                    initialLoading = false
                    $scope.$on("$locationChangeSuccess", onChangeUrl)
                }
                setTimeout($scope.loadCategories, 60000)
            }, data => {
                $scope.failedLoadCategories = true
                console.error("Error loading data", data)
                $scope.loading = false
            })
    }
    $scope.loadCategories(false)

    $scope.loadCategoriesFn = e => {
        e.preventDefault()
        $scope.loadCategories(false)
        $scope.failedLoadCategories = false
    }

    $scope.$watch("feeds.entries.current", entry => {
        if (!entry) return

        if (!entry.read_at) {
            entry.read_at = new Date()
            $scope.feeds.selected.un_read -= 1
        }
        if (entry.new_entry) {
            entry.new_entry = false
        }
        $http.get(`/api/entry/${entry.id}`)
    })

    $scope.loadMore = feedUrl => {
        $http.get(feedUrl)
            .then(res => {
                let feedEntries = $scope.feeds.entries.list.objects
                feedEntries.push.apply(feedEntries, res.data.objects)
                $scope.feeds.entries.list.objects = feedEntries
                $scope.feeds.entries.list.meta = res.data.meta
            })
    }

    $scope.safe = $sce.trustAsHtml

    $scope.doSearch = () => {
        if (!$scope.search) return
        $location.url(`search=${$scope.search}`)
    }
}).controller("RSSAddEditCategoryCtrl", ($scope, $uibModalInstance, $http, category, parentScope) => {
    $scope.category = category
    $scope.form = angular.copy(category)

    $scope.save = () => {
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

        method.then(() => {
            parentScope.loadCategories(false)
            $uibModalInstance.close()
        }, () => {
            $scope.error = "Something went wrong"
        })
    }

    $scope.cancel = $uibModalInstance.dismiss
}).controller("RSSCreateFeedCtrl", ($scope, $uibModalInstance, $http, category, parentScope) => {
    $scope.form = {feed_url: "", category: ""}

    if (category !== undefined && category)
        $scope.form.category = category.id

    $scope.categories = parentScope.feeds.categories.objects

    $scope.save = () => {
        $http.post("/api/feed/", $scope.form)
            .then(() => {
                parentScope.loadCategories(false)
                $uibModalInstance.close()
            }, () => {
                $scope.error = "Something went wrong"
            })
    }

    $scope.cancel = $uibModalInstance.dismiss
}).controller("RSSUpdateFeedCtrl", ($scope, $uibModalInstance, $http, feed, parentScope) => {
    $scope.categories = parentScope.feeds.categories.objects
    $scope.feed = feed
    $scope.form = angular.copy(feed)
    console.log(feed.category_id)
    $scope.form.category = $scope.categories.find(c => c.id === feed.category_id)
    delete $scope.form.un_read
    delete $scope.form.new_entries

    $scope.save = () => {
        let formData = angular.copy($scope.form)
        formData.category = $scope.form.category.id
        $http.put(`/api/feed/${feed.id}/`, formData)
            .then(() => {
                parentScope.loadCategories(false)
                $uibModalInstance.close()
            }, () => {
                $scope.error = "Something went wrong"
            })
    }

    $scope.cancel = $uibModalInstance.dismiss
}).controller("RSSDeleteCategoryCtrl", ($scope, $uibModalInstance, $http, category, parentScope) => {
    $scope.category = category
    $scope.ok = () => {
        $http.delete(`/api/category/${category.id}/`)
            .then(() => {
                parentScope.loadCategories(false)
                $uibModalInstance.close()
            }, () => {
                $scope.error = "Something went wrong"
            })
    }

    $scope.cancel = () => $uibModalInstance.dismiss
}).controller("RSSDeleteFeedCtrl", ($scope, $uibModalInstance, $http, feed, parentScope) => {
    $scope.feed = feed
    $scope.ok = () => {
        $http.delete(`/api/feed/${feed.id}/`)
            .then(() => {
                parentScope.loadCategories(false)
                $uibModalInstance.close()
            }, () => {
                $scope.error = "Something went wrong"
            })
    }

    $scope.cancel = $uibModalInstance.dismiss
}).controller("EditUserCtrl", ($scope, $uibModalInstance, $http, $window, user) => {
    $scope.form = angular.copy(user)
    $scope.save = () => {
        let postData = `name=${encodeURIComponent($scope.form.name)}&email=${encodeURIComponent($scope.form.email)}&password=${$scope.form.password?encodeURIComponent($scope.form.password):''}`
        $http({
            method: 'POST',
            url: `/api/user/`,
            data: postData,
            headers: {'Content-Type': 'application/x-www-form-urlencoded'}
        }).then(() => {
            $window.location.reload()
            $uibModalInstance.close()
        }, () => {
            $scope.error = "Something went wrong"
        })
    }

    $scope.cancel = $uibModalInstance.dismiss
})
