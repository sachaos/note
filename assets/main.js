window.onload = function () {
    if (!window["WebSocket"]) {
        alert("Error! Your browser is out of support.");
    } else {
        socket = new WebSocket("ws://localhost:1129/ws");
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
                hljs.highlightBlock(item);
            });
        };
    }
};
