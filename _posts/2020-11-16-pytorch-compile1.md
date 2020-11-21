# Compiling PyTorch 1.7.0 for Raspberry Pi 3B

Been following these articles here and here. Decided to mostly follow that second one because its a bit more specific to raspberry pi. I followed all the environment variable overrides and left the build running overnight. It got 60% through the compilation! I was quite surprised by this because I only have 600mb of space left on my flash drive so I was wondering if I would run out of space before the build finished. Other thing is I have not created a 2Gb swap RAM and just stuck with the default 100Mb -> again for space reasons.

Found that my build got stuck here:
```
[ 59%] Building C object confu-deps/XNNPACK/CMakeFiles/XNNPACK.dir/src/qs8-gemm/gen/1x8c4-minmax-neondot.c.o
cc: error: unrecognized argument in option ‘-march=armv8.2-a+dotprod’
cc: note: valid arguments to ‘-march=’ are: armv2 armv2a armv3 armv3m armv4 armv4t armv5 armv5e armv5t armv5te armv6 armv6-m armv6j armv6k armv6kz armv6s-m armv6t2 armv6z armv6zk armv7 armv7-a armv7-m armv7-r armv7e-m armv7ve armv8-a armv8-a+crc armv8.1-a armv8.1-a+crc iwmmxt iwmmxt2 native
confu-deps/XNNPACK/CMakeFiles/XNNPACK.dir/build.make:21206: recipe for target 'confu-deps/XNNPACK/CMakeFiles/XNNPACK.dir/src/qs8-gemm/gen/1x8c4-minmax-neondot.c.o' failed
make[2]: *** [confu-deps/XNNPACK/CMakeFiles/XNNPACK.dir/src/qs8-gemm/gen/1x8c4-minmax-neondot.c.o] Error 1
```

Doing some research I have tracked down that the issue is unsupported architecture on the GCC compiler. I think I need to upgrade my GCC to 9.3 at least. I found this by simply scanning the gcc docs for the problem architecture armv8.2-a+dotprod.

Not present in the 7.2.0 docs: https://gcc.gnu.org/onlinedocs/gcc-7.2.0/gcc/ARM-Options.html#ARM-Options
Is present in the 9.3.0 docs: https://gcc.gnu.org/onlinedocs/gcc-9.3.0/gcc/ARM-Options.html#ARM-Options
Will have a crack at updating my GCC and see where that leads!
