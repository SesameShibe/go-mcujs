#pragma once
#include <SDL.h>
#include <stdint.h>
#include "lcd-driver.h"

#define LCD_WIDTH 320
#define LCD_HEIGHT 240

extern MyLcdGfx_t * cgfx;

void mjsInit();
int mjsEventLoop();

#define MCUJS_EMU 1