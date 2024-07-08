Fork info:
This fork implements generating video thumbnails.

Inspired by this fork https://github.com/cody82/filebrowser/tree/videopreview
and further customized.

customaizations include:
* Generation into a .thumbnails subfolder with files in addition to cache mechanism.
* Intel hardware QSV acceleration.
* Thumbnails will stop generating when leaving a folder to allow generating for the next folder before finishing the last one.

** Obvoisly requires ffmpeg, an intel card + drivers.

** Uses the ctx utility to impelent the feature of stoppeing generation, so probably not compatible with windows.

Modified files are:

/http/preview.go

/src/components/files/ListingItem.vue

ChadGPT havevily used in generating the code.

Although I have some programming skills and even a software engeneering degree, 
my career path led me to IT and System Administration.

This is basically my first experience with golang, so don't be harsh.

<p align="center">
  <img src="https://raw.githubusercontent.com/filebrowser/logo/master/banner.png" width="550"/>
</p>

![Preview](https://user-images.githubusercontent.com/5447088/50716739-ebd26700-107a-11e9-9817-14230c53efd2.gif)

[![Build](https://github.com/filebrowser/filebrowser/actions/workflows/main.yaml/badge.svg)](https://github.com/filebrowser/filebrowser/actions/workflows/main.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/filebrowser/filebrowser?style=flat-square)](https://goreportcard.com/report/github.com/filebrowser/filebrowser)
[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/filebrowser/filebrowser)
[![Version](https://img.shields.io/github/release/filebrowser/filebrowser.svg?style=flat-square)](https://github.com/filebrowser/filebrowser/releases/latest)
[![Chat IRC](https://img.shields.io/badge/freenode-%23filebrowser-blue.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23filebrowser)

filebrowser provides a file managing interface within a specified directory and it can be used to upload, delete, preview, rename and edit your files. It allows the creation of multiple users and each user can have its own directory. It can be used as a standalone app.

## Demo

url: https://demo.filebrowser.org/

credentials: `demo`/`demo`

## Features

Please refer to our docs at [https://filebrowser.org/features](https://filebrowser.org/features)

## Install

For installation instructions please refer to our docs at [https://filebrowser.org/installation](https://filebrowser.org/installation).

## Configuration

[Authentication Method](https://filebrowser.org/configuration/authentication-method) - You can change the way the user authenticates with the filebrowser server

[Command Runner](https://filebrowser.org/configuration/command-runner) - The command runner is a feature that enables you to execute any shell command you want before or after a certain event.

[Custom Branding](https://filebrowser.org/configuration/custom-branding) - You can customize your File Browser installation by change its name to any other you want, by adding a global custom style sheet and by using your own logotype if you want.

## Contributing

If you're interested in contributing to this project, our docs are best places to start [https://filebrowser.org/contributing](https://filebrowser.org/contributing).
