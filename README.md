# gomage

#### Image optimisation server written in golang

#### API usage

GET http://localhost:3300/v1/optimize/s.jpg?filters

##### Check api usage with sample.jpg file

http://localhost:3300/v1/optimize/s.jpg?zoom=2

### Following filters are supported

1. **format** = jpg | png | webp
   - output file format
   - e.g. format=jpg
2. **scale** = number
   - scale the image and maintain the ratio
   - number can be between 1 to 10
3. **smartcrop** = width,height,interest
   - crop the image
   - width and height divided by comma in px
   - interest can be any of these
     -- entropy
     -- center
     -- attention
     -- high
     -- low
     -- last
     -- none
   - e.g. crop=50,50,center
4. **rotate** = angle
   - rotate the image
   - allowed angles are
     -- 90
     -- 180
     -- 270
     -- 0
   - e.g. rotate=90
5. **flip** = h or v
   - flip the image vertically or horizontally
   - h for horizontal and v for vertical
   - e.g. flip=h or flip=v
6. **sharpen** = sigma<number>,threshold<number>?,slope<number>?
   - sharpen the image
   - sigma = the sigma of the Gaussian mask
   - threshols = the level of sharpening to apply to "flat" areas.
   - slope = the level of sharpening to apply to "jagged" areas.
   - e.g. sharpen=5,1.0,2.0
7. **modulate** = brightness<number>,saturation<number>?,hue<degree>?
   - transforms the image using brightness, saturation, hue rotation
   - e.g = modulate=2,2,180
8. **label** = text<string>,font<stringr>,width<number>,height<number>,x-location<number>,y-location<number>,opacity<float>,color<r:g:b>
   - add text on top of image
   - e.g. label=hello,arial,50,40,10,10,1,185:53:53
9. **pixlate** = number
   - pixlate the image
   - number can be between 1 to 100
   - e.g. pixlate=10
10. **scale** = float-number
    - scale the image and maintain the ratio
    - float number can be between 0 to 10
    - e.g. scale=2.5
11. **repeat** = x-number,y-number
    - x-number is number of times to repeat horizontally
    - y-number is number of times to repeat vertically
    - number can be between 1 to 10
    - e.g. repeat=5,2
