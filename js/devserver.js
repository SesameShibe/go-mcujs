this.devserver = (function() {
    var global = this


    function suggest(line) {
    }

    var activeWSConn
    var devServer = http.createServer(9971, "devtools", "/api/", 0)
    devServer.onwebsocket = function(ws) {
        console.log("DevTools connected...")
        if (activeWSConn) {
            activeWSConn.close()
        }
        activeWSConn = ws
        ws.ontext = function(ws, text) {
            try {
                var ret = global.eval(text)
                ws.sendText('R' + ret)
            } catch (error) {
                ws.sendText('E' + error)
            }
        }
    }

    return {
        suggest: suggest
    }
})()