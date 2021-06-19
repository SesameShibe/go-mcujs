#include "mcujs.h"
#include "lvgl.h"
#include <stdio.h>

SDL_Window *window;
SDL_Surface *screenSurf;
SDL_Surface *fbSurf;
uint16_t *lcdFB;
lv_color_t drawBuf[LCD_WIDTH * LCD_HEIGHT];

lv_disp_draw_buf_t draw_buf;
lv_disp_drv_t disp_drv;

void my_disp_flush(lv_disp_drv_t *disp_drv, const lv_area_t *area,
                   lv_color_t *color_p) {

  for (int y = area->y1; y <= area->y2; ++y) {
    for (int x = area->x1; x <= area->x2; ++x) {
      lcdFB[y * disp_drv->hor_res + x] = lv_color_to16(*color_p);
      color_p++;
    }
  }

  lv_disp_flush_ready(disp_drv); /* Indicate you are ready with the flushing*/
}

void mjsInit() {
  SDL_Init(SDL_INIT_EVERYTHING);
  window = SDL_CreateWindow("go-mcujs", SDL_WINDOWPOS_UNDEFINED,
                            SDL_WINDOWPOS_UNDEFINED, LCD_WIDTH, LCD_HEIGHT, 0);
  screenSurf = SDL_GetWindowSurface(window);
  fbSurf = SDL_CreateRGBSurfaceWithFormat(0, LCD_WIDTH, LCD_HEIGHT, 16,
                                          SDL_PIXELFORMAT_RGB565);
  lcdFB = (uint16_t *)(fbSurf->pixels);
  memset(lcdFB, 0x0F, LCD_WIDTH * LCD_HEIGHT * 2);
  lv_init();

  lv_disp_draw_buf_init(&draw_buf, drawBuf, NULL, LCD_WIDTH * LCD_HEIGHT);
  lv_disp_drv_init(&disp_drv);       /*Basic initialization*/
  disp_drv.flush_cb = my_disp_flush; /*Set your driver function*/
  disp_drv.draw_buf = &draw_buf;     /*Assign the buffer to the display*/
  disp_drv.hor_res = LCD_WIDTH; /*Set the horizontal resolution of the display*/
  disp_drv.ver_res =
      LCD_HEIGHT; /*Set the verizontal resolution of the display*/

  lv_disp_drv_register(&disp_drv);
  /*
    lv_obj_t* btn = lv_btn_create(lv_scr_act());
    lv_obj_t* label = lv_label_create(btn);
    lv_label_set_text(label, "Hello");*/
}

int mjsEventLoop() {
  SDL_Event event = {0};
  while (SDL_PollEvent(&event) != 0) {
    if (event.type == SDL_QUIT) {
      return 1;
    }
  }

  lv_task_handler();
  SDL_BlitSurface(fbSurf, NULL, screenSurf, NULL);
  SDL_UpdateWindowSurface(window);
  
 
  return 0;
}

void mjsExit() {}
