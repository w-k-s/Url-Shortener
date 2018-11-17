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
            options: {  
                    compress: true  
            },  
            applib: {  
                src:['js/index.js'],
                dest:'build/js/index.min.js'
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
    grunt.registerTask('default', ['uglify', 'cssmin','processhtml','copy']);  
};