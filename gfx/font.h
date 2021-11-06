#pragma once

#define MCUJS_EMU 1

#ifdef MCUJS_EMU
#include <stdio.h>
#else
#include "esp_partition.h"
spi_flash_mmap_handle_t fontMmapHandle;
#endif
// #include "mcujs-lcd.h"
#include "string.h"
#include "stdlib.h"
#include "lcd-driver.h"

uint16_t fbPosX;
uint16_t fbPosY;
uint16_t fbFgColor;

typedef struct _FB_FONT {
  int valid;
  int charWidth;
  int charHeight;
  int charDataSize;
  int pageSize;
  uint8_t *pData;
  uint8_t *pIndex;
  uint8_t *pCharData;
} FB_FONT;

FB_FONT fbFontAscii16;
FB_FONT fbFontCJK16;

FB_FONT *fbCurrentFont;

static inline void fbSetP(MyLcdGfx *gfx, uint16_t x, uint16_t y, uint16_t v) { gfxDrawPixel(gfx, x, y, v); }

static void fbInitFont() {
  FB_FONT *pFont = &fbFontCJK16;
  const void *fontData = 0;

  memset(pFont, 0, sizeof(FB_FONT));

#ifdef MCUJS_EMU
  FILE *fp = fopen("font.bin", "rb");
  fseek(fp, 0, SEEK_END);
  long fsize = ftell(fp);
  fseek(fp, 0, SEEK_SET);
  fontData = (const void *)malloc(fsize);
  fread((void *)fontData, 1, fsize, fp);
  fclose(fp);
#else
  auto part = esp_partition_find_first((esp_partition_type_t)0x40,
                                       (esp_partition_subtype_t)0, "font");
  auto ret = esp_partition_mmap(part, 0, part->size, SPI_FLASH_MMAP_DATA,
                                &fontData, &fontMmapHandle);
  if (ret != 0) {
    printf("mmap font failed...\n");
    return;
  }
#endif

  pFont->valid = 1;
  pFont->charWidth = *(uint8_t *)(fontData + 13);
  pFont->charHeight = *(uint8_t *)(fontData + 14);
  pFont->charDataSize = 1 + ((pFont->charWidth + 7) / 8) * pFont->charHeight;
  pFont->pageSize = pFont->charDataSize * 256;
  pFont->pData = (uint8_t *)fontData;
  pFont->pIndex = pFont->pData + 16;
  pFont->pCharData = pFont->pIndex + 256;

  fbCurrentFont = pFont;
}

static int fbDrawUnicodeRune(MyLcdGfx *gfx, uint32_t rune) {
  if (!fbCurrentFont) {
    return 0;
  }
  int fontW = fbCurrentFont->charWidth;
  int fontH = fbCurrentFont->charHeight;
  int fontCharSize = fbCurrentFont->charDataSize;
  int fontPageSize = fbCurrentFont->pageSize;

  int screenW = gfx->width();
  int screenH = gfx->height();
  rune = (uint16_t)(rune);
  if (rune == '\n') {
    fbPosY += fontH + 1;
    fbPosX = 0;
    return 0;
  }
  uint8_t pgOffset = fbCurrentFont->pIndex[rune >> 8];
  if (pgOffset == 0xFF) {
    return 0;
  }
  uint8_t *ptr = fbCurrentFont->pCharData + fontPageSize * pgOffset +
            fontCharSize * (rune & 0xff);
  uint8_t width = *ptr;
  ptr++;

  if (fbPosX + width >= screenW) {
    fbPosY += fontH + 1;
    fbPosX = 0;
  }
  if (fbPosY + fontH >= screenH) {
    return 0;
  }
  for (uint8_t y = 0; y < fontH; y++) {
    for (uint8_t x = 0; x < width; x++) {
      uint8_t pix = ptr[y * 2 + x / 8] & (1 << (x % 8));
      if (pix) {
        fbSetP(gfx, fbPosX + x, fbPosY + y, 1);
      }
    }
  }
  fbPosX += width + 1;
  if (fbPosX >= screenW) {
    fbPosY += fontH + 1;
    fbPosX = 0;
  }
  return width;
}

static void fbDrawUtf8String(MyLcdGfx *gfx, const char *utf8Str) {
  uint8_t *p = (uint8_t *)utf8Str;
  uint16_t rune = 0;
  while (*p) {
    rune = 0;
    uint8_t byte1 = *p;
    p++;
    if ((byte1 & 0x80) == 0) {
      rune = byte1;
    } else {
      uint8_t byte2 = *p;
      p++;
      if (byte2 == 0) {
        break;
      }
      if ((byte1 & 0xE0) == 0xC0) {
        rune = ((byte1 & 0x1F) << 6) | (byte2 & 0x3F);
      } else {
        uint8_t byte3 = *p;
        p++;
        if (byte3 == 0) {
          break;
        }
        if ((byte1 & 0xf0) == 0xE0) {
          rune =
              ((byte1 & 0x0F) << 12) | ((byte2 & 0x3F) << 6) | (byte3 & 0x3F);
        } else {
          break;
        }
      }
    }
    fbDrawUnicodeRune(gfx, rune);
  }
}