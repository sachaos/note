note
===

Simple Real-time Markdown previewer & editor.

## Description

`note` is a command written in Golang.
This command improve your Markdown editing experience on your favorite editor(emacs, vim etc) by rendering on Web browser in real time.

This software watch markdown file which you are editing, and serve that markuped HTML to browser through WebSocket when that file changed.

## Features

* Style like GitHub
* No dependencies on specific editor
    * `note` use editor which is set in `$EDITOR`
* Support PlantUML

## Demo

![note2 mp4](https://user-images.githubusercontent.com/6121271/43771050-f421ce64-9a78-11e8-9457-256234365032.gif)

### PlantUML

![plantuml mp4](https://user-images.githubusercontent.com/6121271/43970738-84e33dba-9d09-11e8-94d4-909681344824.gif)

## Install

### Binary

Go to [release page](https://github.com/sachaos/note/releases) and download.

```shell
$ wget https://github.com/sachaos/note/releases/download/v0.2.0/note_darwin_amd64 -O /usr/local/bin/note
$ chmod +x /usr/local/bin/note
```

### Manually Build

You need Golang compiler, and [golang/dep: Go dependency management tool](https://github.com/golang/dep), and npm.

```shell
$ git clone https://github.com/sachaos/note.git
$ make install
```

## How to use

```shell
$ note {filename}
```
