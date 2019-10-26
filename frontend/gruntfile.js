module.exports = function (grunt) {  
    require("matchdep").filterDev("grunt-*").forEach(grunt.loadNpmTasks);  
    // Project configuration.  
    grunt.initConfig({  
        pkg: grunt.file.readJSON('package.json'),  
        cssmin: {  
            minify: {  
                files: {  
                    'build/css/index.min.css': [  
                        'css/*.css'
           			]
                }  
            }
        },  
        uglify: {
            local:{
                options: {  
                    compress: true  
                },  
                files: {  
                    'build/js/index.min.js':['js/config.default.js','js/index.js']
                }
            },
            production:{
                options: {  
                    compress: true  
                },  
                files: {  
                    'build/js/index.min.js':['js/config.production.js','js/index.js']
                }
            }
        },
        processhtml: {
            dist: {
                options: {
                    process: true,
                    data:{
                        title: 'My app',
                        message: 'This is production distribution'
                    }
                },
                files: {
                    'build/index.html': ['index.html']
                }
            }
        },
        copy: {
          main: {
            expand: true,
            src: 'assets/*/**',
            dest: 'build/',
          },
        }
    });  
    // Default task.  
    grunt.registerTask('default', ['uglify:local', 'cssmin','processhtml','copy']);  
    grunt.registerTask('production', ['uglify:production', 'cssmin','processhtml','copy']);  
};