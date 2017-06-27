/* global require: false */
'use strict';
var config = require('./config');
var gulp = require('gulp');
var less = require('gulp-less');
var concat = require('gulp-concat');
var uglify = require('gulp-uglify');
var autoprefixer = require('gulp-autoprefixer');
var ngAnnotate = require('gulp-ng-annotate');
var del = require('del');
var LessPluginCleanCSS = require('less-plugin-clean-css');
var cleancss = new LessPluginCleanCSS({advanced: true});
var ngTemplates = require('gulp-ng-templates');
var sourcemaps = require('gulp-sourcemaps');

gulp.task('templates', function () {
    return gulp.src(config.templates.src)
        .pipe(ngTemplates({
            filename: config.templates.out,
            module: config.templates.moduleName,
            path: function (path, base) {
                return path.replace(base, '');
            }
        }))
        .pipe(gulp.dest(config.templates.dest));
});

gulp.task('styles', function () {
    return gulp.src(config.styles.src)
        .pipe(less({
            plugins: [cleancss],
            paths: config.styles.paths
        }))
        .pipe(autoprefixer(config.styles.browsers))
        .pipe(concat(config.styles.out))
        .pipe(gulp.dest(config.styles.dest));
});

gulp.task('scripts', function () {
    return gulp.src(config.scripts.src)
        .pipe(sourcemaps.init())
        .pipe(ngAnnotate())
        .pipe(concat(config.scripts.out))
        .pipe(uglify())
        .pipe(gulp.dest(config.scripts.dest));
});

gulp.task('copy', function () {
    return config.copy.forEach(function (row) {
        gulp.src(row[0]).pipe(gulp.dest(row[1]));
    });
});

gulp.task('clean', function (cb) {
    del(config.clean, {force: true}).then(paths => cb());
});

gulp.task('build', ['styles', 'scripts', 'templates', 'copy']);

gulp.task('watch', ['build'], function () {
    gulp.watch(config.styles.src, ['styles']);
    gulp.watch(config.scripts.src, ['scripts']);
    gulp.watch(config.templates.src, ['templates']);
    gulp.watch(config.copy.map(el => el[0]), ['copy']);
});
