#!/bin/bash

umount ./merged
rm upper lower merged work -r

mkdir upper lower merged work

echo "I'm from lower!" > lower/in_lower.txt
echo "I'm from upper!" > upper/in_upper.txt
# `in_both` is in both directories
echo "I'm from lower!" > lower/in_both.txt
echo "I'm from upper!" > upper/in_both/txt

sudo mount -t overlay overlay -o lowerdir=./lower,upperdir=./upper,workdir=./work ./merged

# 最底下一层是不会被修改的，只读，overlay是支持多个lowerdir, 多个lower用冒号隔开
# upper，它是被mount两层目录中上面的这一层，在overlayfs中，如果有文件创建，修改，删除，那么都会在这一层反映出来。它是可读写的
# merged，挂载点目录，也就是用户看到的目录，实际操作也是在这里
# work，存放临时文件的目录，OverlayFS中如果有文件修改，就会在中间过程中临时存放文件在这里。