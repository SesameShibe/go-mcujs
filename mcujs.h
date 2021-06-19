#pragma once
#include <SDL.h>
#include <stdint.h>
#include "lvgl.h"

#define LCD_WIDTH 320
#define LCD_HEIGHT 240

void mjsInit();
int mjsSetLvCallback(lv_obj_t * obj, int eCode);
int mjsEventLoop();