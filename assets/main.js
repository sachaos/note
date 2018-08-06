window.onload = function () {
    if (!window["WebSocket"]) {
        alert("Error! Your browser is out of support.");
    } else {
        socket = new WebSocket("ws://localhost:1129/ws");
        socket.onclose = function() {
            console.log("Connection closed.");
        };
        socket.onmessage = function(e) {
            var content = document.getElementById('content');
            content.innerHTML = e.data;
            var codeBlocks = document.querySelectorAll('pre code');
            Array.prototype.forEach.call(codeBlocks, function(item) {
                hljs.highlightBlock(item);
            });
        };
    }
};
