note
===

Simple realtime Markdown previewer & editor.

## Description

`note` is a command written in Golang.
This command improve your Markdown editing experience on your favorite editor(emacs, vim etc) by rendering on Web browser in realtime.

This software watch markdown file which you are editing, and serve that markuped HTML to browser through WebSocket when that file changed.

## Demo

![note2 mp4](https://user-images.githubusercontent.com/6121271/43771050-f421ce64-9a78-11e8-9457-256234365032.gif)

## Install

```shell
$ make install
```

## How to use

```shell
$ note {filename}
```
