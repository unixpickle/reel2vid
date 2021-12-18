# reel2vid

Convert a video that is encoded into one large image as many images side-by-side into a video file.

For example, suppose we have the following file:

![Several images of a yellow vase](example/reel.png)

We could run the following command:

```
go run . -width 64 -height 64 -fps 5 reel.png output.mp4
```

This would turn all of the 64x64 frames of the above image into a video like so:

![A video of a vase spinning](example/vid.gif)
