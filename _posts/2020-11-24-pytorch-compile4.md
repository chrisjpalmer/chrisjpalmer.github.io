# Compiling PyTorch 1.7.0 for Raspberry Pi 3B (Part 4)

Still following these two articles BUT this time a different setup.
1. https://nmilosev.svbtle.com/compling-arm-stuff-without-an-arm-board-build-pytorch-for-the-raspberry-pi
2. https://www.spinellis.gr/blog/20200317/index.html

Setup:
1. Overpowered windows computer with skylake 8700k 6 core intel processor (which means 12 cores thanks to hyperthreading), 16 Gb of RAM
2. virtual box running fedora
3. running the following script in virtual box

```sh
sudo dnf install qemu-system-arm qemu-user-static virt-manager
sudo dnf install --releasever=30 --installroot=/virt/F30ARM --forcearch=armv7hl --repo=fedora --repo=updates systemd passwd dnf fedora-release vim-minimal openblas-devel blas-devel m4 cmake python3-Cython python3-devel python3-yaml python3-pillow python3-setuptools python3-numpy python3-cffi python3-wheel gcc-c++ tar gcc git make tmux -y

sudo chroot /virt/F30ARM

sed -i "s/'armv7hnl', 'armv8hl'/'armv7hnl', 'armv7hcnl', 'armv8hl'/" /usr/lib/python3.7/site-packages/dnf/rpm/__init__.py
alias dnf='dnf --releasever=30 --forcearch=armv7hl --repo=fedora --repo=updates'
alias python=python3
echo 'nameserver 8.8.8.8' > /etc/resolv.conf

git clone --depth=1 --recursive --branch=v1.7.0 https://github.com/pytorch/pytorch
cd pytorch

export MAX_JOBS=6 # because I gave my virtual box machine 6 cores of my 12 core system

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

And for torch vision:
```sh
sudo mount -o bind /dev /virt/F30ARM/dev # run in a different terminal not in the chroot

# install pytorch from the wheel we just made
cd dist
pip3 install torch-1.7.0a0-cp37-cp37m-linux_armv7l.whl
cd ../../

# make the vision wheel
git clone  --depth=1 --recursive --branch=v0.8.1 https://github.com/pytorch/vision
cd vision
python setup.py bdist_wheel
```
