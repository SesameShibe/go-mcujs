var scr = lv.scr_act()
var btn = lv.btn(scr)
btn.on("click", function (e) {
    print("clicked!")
})
lv.label(btn).set_text("OK")
btn.set_pos(100, 100)
var label = lv.label(scr)
label.set_text("Hello world!")

