@echo off

cd .\gfx\
g++ -c .\lcd-driver.cpp ..\lib\Adafruit\Adafruit_GFX.cpp -I..\lib\Adafruit\
ar -crs liblcd.a .\lcd-driver.o .\Adafruit_GFX.o
cd ..