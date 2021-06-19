var s = http.createServer(8080,"","/",0)
s.onrequest = function(req, resp) {
    resp.write("hello world")
    resp.close()
}