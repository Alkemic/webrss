/* global module: false */
const destDir = "../static"
const nodeDir = "./node_modules"

module.exports = {
    styles: {
        src: [
            nodeDir + "/bootstrap/less/bootstrap.less",
            "./less/default.less",
        ],
        dest: destDir + "/css",
        paths: [
            nodeDir + "/bootstrap/less/",
        ],
        out: "webrss.css",
        browsers: ["last 2 version", ">5%"],
    },
    vendorScripts: {
        src: [
            nodeDir + "/jquery/dist/jquery.js",
            nodeDir + "/bootstrap/dist/js/bootstrap.js",
            nodeDir + "/angular/angular.js",
            nodeDir + "/angular-ui-bootstrap/dist/ui-bootstrap.js",
            nodeDir + "/angular-ui-bootstrap/dist/ui-bootstrap-tpls.js",
            nodeDir + "/angular-resource/angular-resource.js",
            nodeDir + "/angular-sanitize/angular-sanitize.js",
        ],
        out: "vendor.js",
        dest: destDir + "/js",
    },
    scripts: {
        src: [
            "./js/main.js",
            "./js/*.js",
        ],
        out: "webrss.js",
        dest: destDir + "/js",
    },
    templates: {
        src: "./templates/*.html",
        out: "webrss.templates.js",
        dest: destDir + "/js",
        moduleName: "webrssApp.templates",
    },
    copy: [
        {src: "./images/favicon.ico", dest: destDir + "/images"},
        {src: nodeDir + "/bootstrap/dist/fonts/glyphicons-halflings-regular.eot", dest: destDir + "/fonts"},
        {src: nodeDir + "/bootstrap/dist/fonts/glyphicons-halflings-regular.eot", dest: destDir + "/fonts"},
        {src: nodeDir + "/bootstrap/dist/fonts/glyphicons-halflings-regular.woff2", dest: destDir + "/fonts"},
        {src: nodeDir + "/bootstrap/dist/fonts/glyphicons-halflings-regular.woff", dest: destDir + "/fonts"},
        {src: nodeDir + "/bootstrap/dist/fonts/glyphicons-halflings-regular.ttf", dest: destDir + "/fonts"},
        {src: nodeDir + "/bootstrap/dist/fonts/glyphicons-halflings-regular.svg", dest: destDir + "/fonts"},
    ],
    clean: [
        destDir + "/images",
        destDir + "/css",
        destDir + "/js",
        destDir + "/fonts",
    ],
}
