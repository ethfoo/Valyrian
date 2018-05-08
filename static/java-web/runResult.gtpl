<html>
    <head></head>
    <body>
        <div><h2>运行结果: </h2></div>

        <h1>WebSocket Echo Test</h1>
        <form>
        <p>Message: <input id="message" type="text" value="Hello, world!"></p>
        </form>
        <button onclick="send();">Send Message</button>
        <script type="text/javascript">
            var sock = null;
            var wsuri = "ws://127.0.0.1:8848/ws";
            window.onload = function() {
                console.log("onload");
                sock = new WebSocket(wsuri);
                sock.onopen = function() {
                    console.log("connected to " + wsuri);
                }
                sock.onclose = function(e) {
                    console.log("connection closed (" + e.code + ")");
                }
                sock.onmessage = function(e) {
                    console.log("message received: " + e.data);
                }
            };

            function send() {
                var msg = document.getElementById('message').value;
                sock.send(msg);
            };
        </script>
    </body>
</html>