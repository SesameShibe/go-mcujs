#pragma once
#include <stdint.h>
#include "Adafruit_GFX.h"
#include <stdio.h>

#define LCD_COLOR_WHITE (1)
#define LCD_COLOR_BLACK (0)
#define LCD_COLOR_INVERSE (2)

class MyLcdGfx : public Adafruit_GFX
{
private:
    uint8_t *lcdFB;

public:
    void drawPixel(int16_t x, int16_t y, uint16_t color);
    MyLcdGfx(int16_t w, int16_t h) : Adafruit_GFX(w, h) {}

    uint8_t *getFrameBuffer() { return this->lcdFB; }
    void setFrameBuffer(uint8_t *fb) { this->lcdFB = fb; }
};

void MyLcdGfx::drawPixel(int16_t x, int16_t y, uint16_t color) {
  int16_t t;
  switch (rotation) {
  case 1:
    t = x;
    x = WIDTH - 1 - y;
    y = t;
    break;
  case 2:
    x = WIDTH - 1 - x;
    y = HEIGHT - 1 - y;
    break;
  case 3:
    t = x;
    x = y;
    y = HEIGHT - 1 - t;
    break;
  }

  if ((x < 0) || (y < 0) || (x >= WIDTH) || (y >= HEIGHT))
    return;

  switch (color) {
  case LCD_COLOR_WHITE:
    lcdFB[x + (y / 8) * WIDTH] |= (1 << (y & 7));
    break;
  case LCD_COLOR_BLACK:
    lcdFB[x + (y / 8) * WIDTH] &= ~(1 << (y & 7));
    break;
  case LCD_COLOR_INVERSE:
    lcdFB[x + (y / 8) * WIDTH] ^= (1 << (y & 7));
    break;
  }
}