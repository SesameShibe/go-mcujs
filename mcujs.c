#include "mcujs.h"
#include <stdio.h>

SDL_Window *window;
SDL_Surface *screenSurf;
SDL_Surface *fbSurf;
uint16_t *lcdFB;
lv_color_t drawBuf[LCD_WIDTH * LCD_HEIGHT];

lv_disp_draw_buf_t draw_buf;
lv_disp_drv_t disp_drv;
lv_indev_drv_t indev_drv;


static bool left_button_down = false;
static int16_t last_x = 0;
static int16_t last_y = 0;

int mjsSetLvCallback(lv_obj_t *obj, int eCode) {
  extern void lvEventCallback(lv_event_t * evt);
  lv_obj_add_event_cb(obj, lvEventCallback, (lv_event_code_t)eCode, NULL);
}

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

void mouse_read(lv_indev_drv_t *indev_drv, lv_indev_data_t *data) {
  (void)indev_drv; /*Unused*/

  /*Store the collected data*/
  data->point.x = last_x;
  data->point.y = last_y;
  data->state =
      left_button_down ? LV_INDEV_STATE_PRESSED : LV_INDEV_STATE_RELEASED;
}

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

  case SDL_FINGERUP:
    left_button_down = false;
    last_x = LV_HOR_RES * event->tfinger.x / MONITOR_ZOOM;
    last_y = LV_VER_RES * event->tfinger.y / MONITOR_ZOOM;
    break;
  case SDL_FINGERDOWN:
    left_button_down = true;
    last_x = LV_HOR_RES * event->tfinger.x / MONITOR_ZOOM;
    last_y = LV_VER_RES * event->tfinger.y / MONITOR_ZOOM;
    break;
  case SDL_FINGERMOTION:
    last_x = LV_HOR_RES * event->tfinger.x / MONITOR_ZOOM;
    last_y = LV_VER_RES * event->tfinger.y / MONITOR_ZOOM;
    break;
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

  lv_task_handler();
  SDL_BlitSurface(fbSurf, NULL, screenSurf, NULL);
  SDL_UpdateWindowSurface(window);

  return 0;
}

void mjsExit() {}

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
  lv_indev_drv_init(&indev_drv);          /*Basic initialization*/
  indev_drv.type = LV_INDEV_TYPE_POINTER; /*See below.*/
  indev_drv.read_cb = mouse_read;         /*See below.*/
  /*Register the driver in LVGL and save the created input device object*/
  lv_indev_t *my_indev = lv_indev_drv_register(&indev_drv);
}