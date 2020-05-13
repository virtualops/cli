package config

const DefaultConfigFile = `project: breeze

# This is your global build settings, telling Breeze how to build the Docker image.
# You may customise this per each environment if you would like.
build:
  preset: laravel

deploy:
  path: /

environments:
  production:
    # specify an image if you have a custom build system. If present, the ` + "`build`" + `
    # settings will be ignored, and this image will be used instead.
    # image: yourorg/homepage:latest
`
