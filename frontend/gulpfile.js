/* global require: false */
const config = require("./config")
const gulp = require("gulp")
const less = require("gulp-less")
const concat = require("gulp-concat")
const autoprefixer = require("gulp-autoprefixer")
const del = require("del")
const LessPluginCleanCSS = require("less-plugin-clean-css")
const cleancss = new LessPluginCleanCSS({advanced: true})
const ngTemplates = require("gulp-ng-templates")
const sourcemaps = require("gulp-sourcemaps")
const gulpIf = require("gulp-if")
const babel = require("gulp-babel")

const production = typeof process.env.PRODUCTION !== "undefined" && process.env.PRODUCTION === "true"

const templates = () => gulp
    .src(config.templates.src)
    .pipe(ngTemplates({
        filename: config.templates.out,
        module: config.templates.moduleName,
        path: (path, base) => path.replace(base+"/", ""),
    }))
    .pipe(gulp.dest(config.templates.dest))


const styles = () => gulp
    .src(config.styles.src)
    .pipe(gulpIf(!production, sourcemaps.init()))
    .pipe(less({
        plugins: production ? [cleancss] : [],
        paths: config.styles.paths
    }))
    .pipe(autoprefixer(config.styles.browsers))
    .pipe(concat(config.styles.out))
    .pipe(gulpIf(!production, sourcemaps.write()))
    .pipe(gulp.dest(config.styles.dest))


const scripts = () => gulp
    .src(config.scripts.src)
    .pipe(gulpIf(!production, sourcemaps.init()))
    .pipe(concat(config.scripts.out))
    .pipe(gulpIf(production, babel({
        "presets": ["minify", {comments: false}],
        "plugins": ["angularjs-annotate"]
    })))
    .pipe(gulpIf(!production, sourcemaps.write()))
    .pipe(gulp.dest(config.scripts.dest))


const vendorScripts = () => gulp
    .src(config.vendorScripts.src)
    .pipe(gulpIf(!production, sourcemaps.init()))
    .pipe(concat(config.vendorScripts.out))
    .pipe(gulpIf(production, babel({
        "presets": ["minify", {comments: false}],
        "plugins": ["angularjs-annotate"]
    })))
    .pipe(gulpIf(!production, sourcemaps.write()))
    .pipe(gulp.dest(config.vendorScripts.dest))


const copy = (cb) => {
    config.copy.forEach(file => gulp.src(file.src).pipe(gulp.dest(file.dest)))
    cb()
}

const clean = (cb) => {
    del(config.clean, {force: true}).then(paths => cb())
}

const build = gulp.series(clean, gulp.parallel(styles, vendorScripts, scripts, templates, copy))

const watch = () => {
    gulp.watch(config.styles.src, styles)
    gulp.watch(config.vendorScripts.src, vendorScripts)
    gulp.watch(config.scripts.src, scripts)
    gulp.watch(config.templates.src, templates)
    gulp.watch(config.copy.map(el => el.src), copy)
}

exports.clean = clean
exports.styles = styles
// exports.vendorStyles = vendorStyles
exports.templates = templates
exports.scripts = scripts
exports.vendorScripts = vendorScripts
exports.copy = copy
exports.watch = watch
exports.build = build
exports.default = build
