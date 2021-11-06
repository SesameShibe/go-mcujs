package main

/*
#include "mcujs-lcd.h"
*/
import "C"

type GfxBase interface {
	getPtr() *C.MyLcdGfx_t
}

type Gfx struct {
	ptr *C.MyLcdGfx_t
}

func (gfx *Gfx) getPtr() *C.MyLcdGfx_t {
	return gfx.ptr
}

func (gfx *Gfx) DrawPixel(x, y, color int) {
	C.gfxDrawPixel(C.cgfx, C.short(x), C.short(y), C.ushort(color))
}

func (gfx *Gfx) DrawLine(x1, y1, x2, y2, color int) {
	C.gfxDrawLine(C.cgfx, C.short(x1), C.short(y1), C.short(x2), C.short(y2), C.ushort(color))
}

func (gfx *Gfx) DrawRect(x, y, w, h, color int) {
	C.gfxDrawRect(C.cgfx, C.short(x), C.short(y), C.short(w), C.short(h), C.ushort(color))
}

func (gfx *Gfx) FillRect(x, y, w, h, color int) {
	C.gfxFillRect(C.cgfx, C.short(x), C.short(y), C.short(w), C.short(h), C.ushort(color))
}

func (gfx *Gfx) DrawCircle(x, y, radius, color int) {
	C.gfxDrawCircle(C.cgfx, C.short(x), C.short(y), C.short(radius), C.ushort(color))
}

func (gfx *Gfx) DrawCircleHelper(x, y, radius, cornername, color int) {
	C.gfxDrawCircleHelper(C.cgfx, C.short(x), C.short(y), C.short(radius), C.uchar(cornername), C.ushort(color))
}

func (gfx *Gfx) FillCircle(x, y, radius, color int) {
	C.gfxFillCircle(C.cgfx, C.short(x), C.short(y), C.short(radius), C.ushort(color))
}

func (gfx *Gfx) FillCircleHelper(x, y, radius, cornername, delta, color int) {
	C.gfxFillCircleHelper(C.cgfx, C.short(x), C.short(y), C.short(radius), C.uchar(cornername), C.short(delta), C.ushort(color))
}

func (gfx *Gfx) DrawTriangle(x1, y1, x2, y2, x3, y3, color int) {
	C.gfxDrawTriangle(C.cgfx, C.short(x1), C.short(y1), C.short(x2), C.short(y2), C.short(x3), C.short(y3), C.ushort(color))
}

func (gfx *Gfx) FillTriangle(x1, y1, x2, y2, x3, y3, color int) {
	C.gfxFillTriangle(C.cgfx, C.short(x1), C.short(y1), C.short(x2), C.short(y2), C.short(x3), C.short(y3), C.ushort(color))
}

func (gfx *Gfx) DrawRoundRect(x, y, w, h, radius, color int) {
	C.gfxDrawRoundRect(C.cgfx, C.short(x), C.short(y), C.short(w), C.short(h), C.short(radius), C.ushort(color))
}

func (gfx *Gfx) FillRoundRect(x, y, w, h, radius, color int) {
	C.gfxFillRoundRect(C.cgfx, C.short(x), C.short(y), C.short(w), C.short(h), C.short(radius), C.ushort(color))
}

var gfx *Gfx = &Gfx{
	ptr: C.cgfx,
}

func GfxInit() {
	vm.Set("gfx", gfx)
}
