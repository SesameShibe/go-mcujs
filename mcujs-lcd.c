#include "mcujs-lcd.h"
#include <stdio.h>

// include for boolean
#include <stdbool.h>

SDL_Window *window;
SDL_Surface *screenSurf;
SDL_Surface *fbSurf;
uint16_t *lcdFB;
uint8_t *oneBppFB;
MyLcdGfx_t *cgfx;

static bool left_button_down = false;
static int16_t last_x = 0;
static int16_t last_y = 0;

void mouse_handler(SDL_Event *event) {
  const int MONITOR_ZOOM = 1;
  switch (event->type) {
  case SDL_MOUSEBUTTONUP:
    if (event->button.button == SDL_BUTTON_LEFT)
      left_button_down = false;
    break;
  case SDL_MOUSEBUTTONDOWN:
    if (event->button.button == SDL_BUTTON_LEFT) {
      left_button_down = true;
      last_x = event->motion.x / MONITOR_ZOOM;
      last_y = event->motion.y / MONITOR_ZOOM;
    }
    break;
  case SDL_MOUSEMOTION:
    last_x = event->motion.x / MONITOR_ZOOM;
    last_y = event->motion.y / MONITOR_ZOOM;
    break;

  // case SDL_FINGERUP:
  //   left_button_down = false;
  //   last_x = LV_HOR_RES * event->tfinger.x / MONITOR_ZOOM;
  //   last_y = LV_VER_RES * event->tfinger.y / MONITOR_ZOOM;
  //   break;
  // case SDL_FINGERDOWN:
  //   left_button_down = true;
  //   last_x = LV_HOR_RES * event->tfinger.x / MONITOR_ZOOM;
  //   last_y = LV_VER_RES * event->tfinger.y / MONITOR_ZOOM;
  //   break;
  // case SDL_FINGERMOTION:
  //   last_x = LV_HOR_RES * event->tfinger.x / MONITOR_ZOOM;
  //   last_y = LV_VER_RES * event->tfinger.y / MONITOR_ZOOM;
  //   break;
  }
}

// include file io
#include <stdio.h>

void map1BppFB() {
  for (int x = 0; x < LCD_WIDTH; x++) {
    for (int y = 0; y < LCD_HEIGHT; y++) {
      int index = x + (y / 8) * LCD_WIDTH;
      int bit = (1 << (y & 7));
      uint16_t pix = (oneBppFB[index] & bit) ? 0xFFFF : 0;
      lcdFB[x + y * LCD_WIDTH] = pix;
    }
  }
}

int mjsEventLoop() {
  SDL_Event event = {0};
  while (SDL_PollEvent(&event) != 0) {
    if (event.type == SDL_QUIT) {
      return 1;
    }
    mouse_handler(&event);
  }

  map1BppFB();
  SDL_BlitSurface(fbSurf, NULL, screenSurf, NULL);
  SDL_UpdateWindowSurface(window);

  return 0;
}

void mjsExit() {}

void mjsInit() {
  cgfx = gfxNewObject(LCD_WIDTH, LCD_HEIGHT);

  SDL_Init(SDL_INIT_EVERYTHING);
  window = SDL_CreateWindow("go-mcujs", SDL_WINDOWPOS_UNDEFINED,
                            SDL_WINDOWPOS_UNDEFINED, LCD_WIDTH, LCD_HEIGHT, 0);
  screenSurf = SDL_GetWindowSurface(window);
  fbSurf = SDL_CreateRGBSurfaceWithFormat(0, LCD_WIDTH, LCD_HEIGHT, 16,
                                          SDL_PIXELFORMAT_RGB565);
  lcdFB = (uint16_t *)(fbSurf->pixels);
  memset(lcdFB, 0x0F, LCD_WIDTH * LCD_HEIGHT * 2);

  oneBppFB = (uint8_t*)malloc(LCD_WIDTH * LCD_HEIGHT / 8);
  memset(oneBppFB, 0xFF, LCD_WIDTH * LCD_HEIGHT / 8);
  gfxSetFrameBuffer(cgfx, oneBppFB);
}