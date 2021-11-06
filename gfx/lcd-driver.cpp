#include "lcd-driver.hpp"

extern "C"
{
  #include "lcd-driver.h"
}

struct MyLcdGfx_t : MyLcdGfx
{
  MyLcdGfx_t(uint16_t width, uint16_t height) : MyLcdGfx(width, height) {}
  ~MyLcdGfx_t() {}
};

MyLcdGfx_t *gfxNewObject(int16_t w, int16_t h)
{
  auto gfx = new MyLcdGfx_t(w, h);
  return gfx;
}

uint8_t * gfxGetFrameBuffer(MyLcdGfx_t *gfx)
{
  return gfx->getFrameBuffer();
}

void gfxSetFrameBuffer(MyLcdGfx_t *gfx, uint8_t *fb)
{
  gfx->setFrameBuffer(fb);
}

void gfxSetRotation(MyLcdGfx_t *gfx, uint8_t r)
{
  gfx->setRotation(r);
}

void gfxInvert(MyLcdGfx_t *gfx, uint8_t i)
{
  gfx->invertDisplay(i);
}

void gfxDrawPixel(MyLcdGfx_t *gfx, int16_t x, int16_t y, uint16_t color)
{
  gfx->drawPixel(x, y, color);
}

void gfxDrawLine(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t x1, int16_t y1, uint16_t color)
{
  gfx->drawLine(x0, y0, x1, y1, color);
}

void gfxDrawRect(MyLcdGfx_t *gfx, int16_t x, int16_t y, int16_t w, int16_t h, uint16_t color)
{
  gfx->drawRect(x, y, w, h, color);
}

void gfxFillRect(MyLcdGfx_t *gfx, int16_t x, int16_t y, int16_t w, int16_t h, uint16_t color)
{
  gfx->fillRect(x, y, w, h, color);
}

void gfxFillScreen(MyLcdGfx_t *gfx, uint16_t color)
{
  gfx->fillScreen(color);
}

void gfxDrawCircle(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t r, uint16_t color)
{
  gfx->drawCircle(x0, y0, r, color);
}

void gfxDrawCircleHelper(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t r, uint8_t cornername,
                         uint16_t color)
{
  gfx->drawCircleHelper(x0, y0, r, cornername, color);
}

void gfxFillCircle(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t r, uint16_t color)
{
  gfx->fillCircle(x0, y0, r, color);
}

void gfxFillCircleHelper(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t r, uint8_t cornername,
                         int16_t delta, uint16_t color)
{
  gfx->fillCircleHelper(x0, y0, r, cornername, delta, color);
}

void gfxDrawTriangle(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t x1, int16_t y1, int16_t x2,
                     int16_t y2, uint16_t color)
{
  gfx->drawTriangle(x0, y0, x1, y1, x2, y2, color);
}

void gfxFillTriangle(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t x1, int16_t y1, int16_t x2,
                     int16_t y2, uint16_t color)
{
  gfx->fillTriangle(x0, y0, x1, y1, x2, y2, color);
}

void gfxDrawRoundRect(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t w, int16_t h,
                      int16_t radius, uint16_t color)
{
  gfx->drawRoundRect(x0, y0, w, h, radius, color);
}

void gfxFillRoundRect(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t w, int16_t h,
                      int16_t radius, uint16_t color)
{
  gfx->fillRoundRect(x0, y0, w, h, radius, color);
}

void gfxDrawBitmap(MyLcdGfx_t *gfx, int16_t x, int16_t y, const uint8_t *bitmap, int16_t w,
                   int16_t h, uint16_t color)
{
  gfx->drawBitmap(x, y, bitmap, w, h, color);
}

void gfxDrawXBitmap(MyLcdGfx_t *gfx, int16_t x, int16_t y, const uint8_t *bitmap, int16_t w,
                    int16_t h, uint16_t color)
{
  gfx->drawXBitmap(x, y, bitmap, w, h, color);
}

void gfxDrawGrayscaleBitmap(MyLcdGfx_t *gfx, int16_t x, int16_t y, const uint8_t *bitmap,
                            int16_t w, int16_t h)
{
  gfx->drawGrayscaleBitmap(x, y, bitmap, w, h);
}

void gfxDrawRGBBitmap(MyLcdGfx_t *gfx, int16_t x, int16_t y, const uint16_t *bitmap, int16_t w,
                      int16_t h)
{
  gfx->drawRGBBitmap(x, y, bitmap, w, h);
}