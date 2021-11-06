#include <stdint.h>

typedef struct MyLcdGfx_t MyLcdGfx_t;

MyLcdGfx_t *gfxNewObject(int16_t w, int16_t h);
uint8_t * gfxGetFrameBuffer(MyLcdGfx_t *obj);
void gfxSetFrameBuffer(MyLcdGfx_t *obj, uint8_t *fb);
void gfxSetRotation(MyLcdGfx_t *gfx, uint8_t r);
void gfxInvert(MyLcdGfx_t *gfx, uint8_t i);
void gfxDrawPixel(MyLcdGfx_t *gfx, int16_t x, int16_t y, uint16_t color);
void gfxDrawLine(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t x1, int16_t y1, uint16_t color);
void gfxDrawRect(MyLcdGfx_t *gfx, int16_t x, int16_t y, int16_t w, int16_t h, uint16_t color);
void gfxFillRect(MyLcdGfx_t *gfx, int16_t x, int16_t y, int16_t w, int16_t h, uint16_t color);
void gfxFillScreen(MyLcdGfx_t *gfx, uint16_t color);
void gfxDrawCircle(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t r, uint16_t color);
void gfxDrawCircleHelper(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t r, uint8_t cornername,
                         uint16_t color);
void gfxFillCircle(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t r, uint16_t color);
void gfxFillCircleHelper(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t r, uint8_t cornername,
                         int16_t delta, uint16_t color);
void gfxDrawTriangle(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t x1, int16_t y1, int16_t x2,
                     int16_t y2, uint16_t color);
void gfxFillTriangle(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t x1, int16_t y1, int16_t x2,
                     int16_t y2, uint16_t color);
void gfxDrawRoundRect(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t w, int16_t h,
                      int16_t radius, uint16_t color);
void gfxFillRoundRect(MyLcdGfx_t *gfx, int16_t x0, int16_t y0, int16_t w, int16_t h,
                      int16_t radius, uint16_t color);
void gfxDrawBitmap(MyLcdGfx_t *gfx, int16_t x, int16_t y, const uint8_t *bitmap, int16_t w,
                     int16_t h, uint16_t color);
void gfxDrawXBitmap(MyLcdGfx_t *gfx, int16_t x, int16_t y, const uint8_t *bitmap, int16_t w,
                        int16_t h, uint16_t color);
void gfxDrawGrayscaleBitmap(MyLcdGfx_t *gfx, int16_t x, int16_t y, const uint8_t *bitmap,
                            int16_t w, int16_t h);
void gfxDrawRGBBitmap(MyLcdGfx_t *gfx, int16_t x, int16_t y, const uint16_t *bitmap, int16_t w,
                        int16_t h);