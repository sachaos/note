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

            if (data.lines.length != 0) {
                var maxLine = Math.max.apply(null, data.lines);;
                var element = document.getElementById(maxLine);
                if (element) {
                    window.scroll(0, element.offsetTop);
                }
            }

            var codeBlocks = document.querySelectorAll('pre code');
            Array.prototype.forEach.call(codeBlocks, function(item) {
                hljs.highlightBlock(item);
            });
        };
    }
};
