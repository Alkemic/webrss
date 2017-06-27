/* global module: false */
var destDir = '../webrss/static';
var nodeDir = './node_modules';

module.exports = {
    styles: {
        src: [
            './node_modules/bootstrap/less/bootstrap.less',
            './less/default.less',
        ],
        dest: destDir + '/css',
        paths: [
            './node_modules/bootstrap/less/',
        ],
        out: 'webrss.css',
        browsers: ['last 2 version', '>5%']
    },
    vendorScripts: {
        src: [
            './node_modules/jquery/dist/jquery.js',
            './node_modules/underscore/underscore.js',
            './node_modules/bootstrap/dist/js/bootstrap.js',
            './node_modules/angular/angular.js',
            './node_modules/angular-ui-bootstrap/dist/ui-bootstrap.js',
            './node_modules/angular-ui-bootstrap/dist/ui-bootstrap-tpls.js',
            './node_modules/angular-resource/angular-resource.js',
            './node_modules/angular-sanitize/angular-sanitize.js',
        ],
        out: 'vendor.js',
        dest: destDir + '/js'
    },
    scripts: {
        src: [
            './js/main.js',
            './js/*.js',
        ],
        out: 'webrss.js',
        dest: destDir + '/js'
    },
    templates: {
        src: './templates/*.html',
        out: 'webrss.templates.js',
        dest: destDir + '/js',
        moduleName: 'webrssApp.templates'
    },
    copy: [
        ['./images/favicon.ico', destDir + '/images'],
        ['./node_modules/bootstrap/dist/fonts/glyphicons-halflings-regular.eot', destDir + '/fonts'],
        ['./node_modules/bootstrap/dist/fonts/glyphicons-halflings-regular.eot', destDir + '/fonts'],
        ['./node_modules/bootstrap/dist/fonts/glyphicons-halflings-regular.woff2', destDir + '/fonts'],
        ['./node_modules/bootstrap/dist/fonts/glyphicons-halflings-regular.woff', destDir + '/fonts'],
        ['./node_modules/bootstrap/dist/fonts/glyphicons-halflings-regular.ttf', destDir + '/fonts'],
        ['./node_modules/bootstrap/dist/fonts/glyphicons-halflings-regular.svg', destDir + '/fonts'],
    ],
    clean: [
        destDir + '/images',
        destDir + '/css',
        destDir + '/js',
        destDir + '/fonts',
    ]
};
