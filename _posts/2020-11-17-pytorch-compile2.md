# Compiling PyTorch 1.7.0 for Raspberry Pi 3B (Part 2)

Okay new plan required.. didn't bother upgrading the compiler on the raspberry pi 3, that looks pretty hard.
It turns out that the GCC running on my raspberry pi is 6.3.0. I tried running `sudo apt-get install -t stretch gcc` but I don't think they are currently shipping a later version of the GCC to rpi.



Instead I tried to hack XNNPACK to make it not use an architecture the compiler does not support.
Whoa whoa hold up.. you did what? Well reading the build output I discovered that a particular library being built called XNNPACK was failing.
XNNPACK is this: [XNNPACK is a highly optimized library of floating-point neural network inference operators for ARM, WebAssembly, and x86 platforms.](https://github.com/google/XNNPACK)... yeah I didn't write that.

So anyways, inside the library I had a look in the CMakeLists.txt file. I found the `-march` flags set to `armv8.2-a+dotprod` and tried changing them to `armv8-a`.
Guys.. I have no idea what I'm doing! I think an architecture is an instruction set that the compiler can omit and I thought "heck, using that fancy smanshy architecture is probably just an optimization... why not try downgrading it!"

Well it had some effect... the code didn't fail in the some way this time. However a new fun error to play with:

```
/home/pi/Downloads/pytorch/third_party/XNNPACK/src/qs8-gemm/gen/1x8c4-minmax-neondot.c: In function ‘xnn_qs8_gemm_minmax_ukernel_1x8c4__neondot’:
/home/pi/Downloads/pytorch/third_party/XNNPACK/src/qs8-gemm/gen/1x8c4-minmax-neondot.c:62:20: warning: implicit declaration of function ‘vdotq_lane_s32’ [-Wimplicit-function-declaration]
       vacc0x0123 = vdotq_lane_s32(vacc0x0123, vb0123x0123, va0x01234567, 0);
                    ^~~~~~~~~~~~~~
/home/pi/Downloads/pytorch/third_party/XNNPACK/src/qs8-gemm/gen/1x8c4-minmax-neondot.c:62:18: error: incompatible types when assigning to type ‘int32x4_t’ from type ‘int’
       vacc0x0123 = vdotq_lane_s32(vacc0x0123, vb0123x0123, va0x01234567, 0);
                  ^
/home/pi/Downloads/pytorch/third_party/XNNPACK/src/qs8-gemm/gen/1x8c4-minmax-neondot.c:63:18: error: incompatible types when assigning to type ‘int32x4_t’ from type ‘int’
       vacc0x4567 = vdotq_lane_s32(vacc0x4567, vb0123x4567, va0x01234567, 0);
                  ^
/home/pi/Downloads/pytorch/third_party/XNNPACK/src/qs8-gemm/gen/1x8c4-minmax-neondot.c:64:18: error: incompatible types when assigning to type ‘int32x4_t’ from type ‘int’
       vacc0x0123 = vdotq_lane_s32(vacc0x0123, vb4567x0123, va0x01234567, 1);
                  ^
/home/pi/Downloads/pytorch/third_party/XNNPACK/src/qs8-gemm/gen/1x8c4-minmax-neondot.c:65:18: error: incompatible types when assigning to type ‘int32x4_t’ from type ‘int’
       vacc0x4567 = vdotq_lane_s32(vacc0x4567, vb4567x4567, va0x01234567, 1);
                  ^
/home/pi/Downloads/pytorch/third_party/XNNPACK/src/qs8-gemm/gen/1x8c4-minmax-neondot.c:79:18: error: incompatible types when assigning to type ‘int32x4_t’ from type ‘int’
       vacc0x0123 = vdotq_lane_s32(vacc0x0123, vb0123x0123, va0x01234567, 0);
                  ^
/home/pi/Downloads/pytorch/third_party/XNNPACK/src/qs8-gemm/gen/1x8c4-minmax-neondot.c:80:18: error: incompatible types when assigning to type ‘int32x4_t’ from type ‘int’
       vacc0x4567 = vdotq_lane_s32(vacc0x4567, vb0123x4567, va0x01234567, 0);
                  ^
/home/pi/Downloads/pytorch/third_party/XNNPACK/src/qs8-gemm/gen/1x8c4-minmax-neondot.c:88:20: error: incompatible types when assigning to type ‘int32x4_t’ from type ‘int’
         vacc0x0123 = vdotq_lane_s32(vacc0x0123, vb4567x0123, va0x01234567, 1);
                    ^
/home/pi/Downloads/pytorch/third_party/XNNPACK/src/qs8-gemm/gen/1x8c4-minmax-neondot.c:89:20: error: incompatible types when assigning to type ‘int32x4_t’ from type ‘int’
```

Something tells me I'm going to be here for a while.

Next stop: stop trying to compile pytorch 1.7.0 on rpi by myself. That was fraught with error from the beginning. Good job me for discovering XNNPACK and understanding some cool things about it. Good job me for realizing right now not go done this path because its probably harder to steer an active project that doesn't care about raspberry pi, without holding 90% of the domain knowledge -> that my friend's is called a rabbit hole. 
Instead going to go for something much more solid and probably follow [this tutorial](https://medium.com/secure-and-private-ai-writing-challenge/a-step-by-step-guide-to-installing-pytorch-in-raspberry-pi-a1491bb80531)

Looks like 1.7.0 is [possibly supported](https://forums.developer.nvidia.com/t/pytorch-for-jetson-version-1-7-0-now-available/72048). Will find out next time!
