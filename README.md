# mweb-export

用于生成 Mweb `git`备份 `README` 目录文件的工具, MWeb可以使用iCloud备份[https://zh.mweb.im/mweb4qa.html](https://zh.mweb.im/mweb4qa.html)，同时可以将`~/Library/Containers/com.coderforart.MWeb3/Data/Library/Application Support/MWebLibrary`使用git备份到github，但是实际使用中，备份到github之后文件夹结构没有了，文件名也变成了`UUID`。完全没有可读性。
所以`mweb-export`就是一个解析Mweb数据库并声称一个包含整个目录结构和索引跳转的README文件。


## Install


```shell
go install github.com/TBXark/mweb-export@latest
```


## Usage

实际使用中不推荐直接在`MWebLibrary`原始文件夹中使用`git`，避免多台设备同步`git`与iCloud发生冲突，可以另外新建新建一个`git`文件夹，需要同步的时候执行下面脚本，将`MWebLibrary`拷贝到`git`文件夹，并生成`README.md`文件并且自动push。


```shell
#!/bin/bash

/bin/rm -rf docs
/bin/rm -f mainlib.lib
/bin/rm -rf metadata
cp ~/Library/Containers/com.coderforart.MWeb3/Data/Library/Application\ Support/MWebLibrary/docs docs
cp ~/Library/Containers/com.coderforart.MWeb3/Data/Library/Application\ Support/MWebLibrary/mainlib.lib mainlib.lib
cp ~/Library/Containers/com.coderforart.MWeb3/Data/Library/Application\ Support/MWebLibrary/metadata metadata
mweb-export -path=$(pwd)
git add .
git commit -a -m $(date -u +%Y-%m-%dT%H:%M:%SZ)
git push origin master
```