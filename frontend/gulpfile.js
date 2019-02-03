/* global require: false */
"use strict"
const config = require("./config")
const gulp = require("gulp")
const less = require("gulp-less")
const concat = require("gulp-concat")
const autoprefixer = require("gulp-autoprefixer")
const ngAnnotate = require("gulp-ng-annotate")
const del = require("del")
const LessPluginCleanCSS = require("less-plugin-clean-css")
const cleancss = new LessPluginCleanCSS({advanced: true})
const ngTemplates = require("gulp-ng-templates")
const sourcemaps = require("gulp-sourcemaps")
const gulpIf = require("gulp-if")
const babel = require("gulp-babel")

const production = typeof process.env.PRODUCTION !== "undefined" && process.env.PRODUCTION === "true"

gulp.task("templates", function () {
    return gulp.src(config.templates.src)
        .pipe(ngTemplates({
            filename: config.templates.out,
            module: config.templates.moduleName,
            path: function (path, base) {
                return path.replace(base, "")
            }
        }))
        .pipe(gulp.dest(config.templates.dest))
})

gulp.task("styles", function () {
    return gulp.src(config.styles.src)
        .pipe(gulpIf(!production, sourcemaps.init()))
        .pipe(less({
            plugins: production ? [cleancss] : [],
            paths: config.styles.paths
        }))
        .pipe(autoprefixer(config.styles.browsers))
        .pipe(concat(config.styles.out))
        .pipe(gulpIf(!production, sourcemaps.write()))
        .pipe(gulp.dest(config.styles.dest))
})

gulp.task("scripts", function () {
    return gulp.src(config.scripts.src)
        .pipe(gulpIf(!production, sourcemaps.init()))
        .pipe(ngAnnotate())
        .pipe(concat(config.scripts.out))
        .pipe(gulpIf(production, babel({presets: ["babili"]})))
        .pipe(gulpIf(!production, sourcemaps.write()))
        .pipe(gulp.dest(config.scripts.dest))
})

gulp.task("vendorScripts", function () {
    return gulp.src(config.vendorScripts.src)
        .pipe(gulpIf(!production, sourcemaps.init()))
        .pipe(ngAnnotate())
        .pipe(concat(config.vendorScripts.out))
        .pipe(gulpIf(production, babel({presets: ["babili"]})))
        .pipe(gulpIf(!production, sourcemaps.write()))
        .pipe(gulp.dest(config.vendorScripts.dest))
})

gulp.task("copy", function () {
    return config.copy.forEach(function (row) {
        gulp.src(row[0]).pipe(gulp.dest(row[1]))
    })
})

gulp.task("clean", function (cb) {
    del(config.clean, {force: true}).then(paths => cb())
})

gulp.task("build", ["styles", "vendorScripts", "scripts", "templates", "copy"])

gulp.task("watch", ["build"], function () {
    gulp.watch(config.styles.src, ["styles"])
    gulp.watch(config.vendorScripts.src, ["vendorScripts"])
    gulp.watch(config.scripts.src, ["scripts"])
    gulp.watch(config.templates.src, ["templates"])
    gulp.watch(config.copy.map(el => el[0]), ["copy"])
})
