window.onload = function () {
    if (!window["WebSocket"]) {
        alert("エラー : WebSocketに対応していないブラウザです。");
    } else {
        socket = new WebSocket("ws://localhost:1129/ws");
        socket.onclose = function() {
            alert(" 接続が終了しました。");
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
