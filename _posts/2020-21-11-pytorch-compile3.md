# Building PyTorch 1.7.0 for Raspberry Pi 3B


Time to do witchcraft (aka cross compilation for arm). And what's more.. using docker!

Here's my recipe based on the following two articles:
1. https://nmilosev.svbtle.com/compling-arm-stuff-without-an-arm-board-build-pytorch-for-the-raspberry-pi
2. https://www.spinellis.gr/blog/20200317/index.html

```
docker run -it fedora sh
sudo dnf install qemu-system-arm qemu-user-static virt-manager
sudo dnf install --releasever=30 --installroot=/tmp/F30ARM --forcearch=armv7hl --repo=fedora --repo=updates systemd passwd dnf fedora-release vim-minimal openblas-devel blas-devel m4 cmake python3-Cython python3-devel python3-yaml python3-pillow python3-setuptools python3-numpy python3-cffi python3-wheel gcc-c++ tar gcc git make tmux -y

sudo chroot /tmp/F30ARM

sed -i "s/'armv7hnl', 'armv8hl'/'armv7hnl', 'armv7hcnl', 'armv8hl'/" /usr/lib/python3.7/site-packages/dnf/rpm/__init__.py
alias dnf='dnf --releasever=30 --forcearch=armv7hl --repo=fedora --repo=updates'

alias python=python3
echo 'nameserver 8.8.8.8' > /etc/resolv.conf

# dnf install git-all

git clone --depth=1 --recursive --branch=1.7.0 https://github.com/pytorch/pytorch
cd pytorch
# git submodule update --remote third_party/protobuf

# curl https://codeload.github.com/pytorch/pytorch/zip/v1.7.0 --output pytorch-1.7.0.zip
# curl -L https://github.com/pytorch/pytorch/archive/v1.7.0.zip --output pytorch-1.7.0.zip
# unzip pytorch-1.7.0.zip -d .
# cd pytorch-1.7.0

# Limit the number of parallel jobs in a 1MB Pi to prevent thrashing
export MAX_JOBS=4

# Disable features that don't make sense on a Pi
export USE_CUDA=0
export USE_CUDNN=0
export USE_MKLDNN=0
export USE_NNPACK=0
export USE_QNNPACK=0
export USE_DISTRIBUTED=0

# Disable testing, which takes ages
export BUILD_TEST=0

python setup.py bdist_wheel
```

Getting this error:
```
qemu: uncaught target signal 11 (Segmentation fault) - core dumped
```


Hmm maybe the version of git is a bit out of date... Let's try the latest fedora release 33...

```
docker run -it fedora:33 sh
sudo dnf install qemu-system-arm qemu-user-static virt-manager
sudo dnf install --releasever=33 --installroot=/tmp/F33ARM --forcearch=armv7hl --repo=fedora --repo=updates systemd passwd dnf fedora-release vim-minimal openblas-devel blas-devel m4 cmake python3-Cython python3-devel python3-yaml python3-pillow python3-setuptools python3-numpy python3-cffi python3-wheel gcc-c++ tar gcc git make tmux -y

sudo chroot /tmp/F33ARM

sed -i "s/'armv7hnl', 'armv8hl'/'armv7hnl', 'armv7hcnl', 'armv8hl'/" /usr/lib/python3.9/site-packages/dnf/rpm/__init__.py
alias dnf='dnf --releasever=33 --forcearch=armv7hl --repo=fedora --repo=updates'

alias python=python3
rm -rf /etc/resolv.conf
echo 'nameserver 8.8.8.8' > /etc/resolv.conf

# dnf install git-all

git clone --depth=1 --recursive --branch=v1.7.0 https://github.com/pytorch/pytorch
cd pytorch
# git submodule update --remote third_party/protobuf

# Limit the number of parallel jobs in a 1MB Pi to prevent thrashing
export MAX_JOBS=4

# Disable features that don't make sense on a Pi
export USE_CUDA=0
export USE_CUDNN=0
export USE_MKLDNN=0
export USE_NNPACK=0
export USE_QNNPACK=0
export USE_DISTRIBUTED=0

# Disable testing, which takes ages
export BUILD_TEST=0

python setup.py bdist_wheel
```

Still getting the same error:
```
qemu: uncaught target signal 11 (Segmentation fault) - core dumped
```

Conclusion, this probably doesn't work inside docker.. wasn't expecting it to anyway.
Will try in a Fedora VM using Virtual Box later... on my windows Pc ;)
