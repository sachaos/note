import PNotify from 'pnotify/dist/es/PNotify.js';
import hljs from 'highlight.js';
import plantumlEncoder from 'plantuml-encoder';

import 'pnotify/dist/PNotifyBrightTheme.css';
import 'github-markdown-css/github-markdown.css';
import 'highlight.js/styles/github.css';
import './main.css';

if (!window["WebSocket"]) {
    alert("Error! Your browser is out of support.");
} else {
    let socket = new WebSocket("ws://localhost:1129/ws");
    socket.onclose = function() {
        PNotify.error("Connection closed.");
    };
    socket.onopen = function() {
        PNotify.success("Connection open! Enjoy!");
    };
    socket.onmessage = function(e) {
        var data = JSON.parse(e.data);

        var content = document.getElementById('content');
        content.innerHTML = data.html;

        var title = document.getElementsByTagName('title')[0];
        title.textContent = data.title;

        var codeBlocks = document.querySelectorAll('pre code');
        Array.prototype.forEach.call(codeBlocks, function(item) {
            if (item.className == "language-plantuml" || item.className == "language-uml") {
                let encoded = plantumlEncoder.encode(item.textContent);
                let url = 'http://www.plantuml.com/plantuml/img/' + encoded;
                let imgTag = document.createElement('img');
                imgTag.setAttribute("src", url);

                item.parentElement.replaceWith(imgTag);
            } else {
                hljs.highlightBlock(item);
            }
        });
    };
}
