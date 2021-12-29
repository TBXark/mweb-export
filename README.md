# mweb-export

[![Build Release](https://github.com/TBXark/mweb-export/actions/workflows/Release.yml/badge.svg)](https://github.com/TBXark/mweb-export/actions/workflows/Release.yml)

用于生成 Mweb `git`备份 `README` 目录文件的工具, MWeb可以使用iCloud备份[https://zh.mweb.im/mweb4qa.html](https://zh.mweb.im/mweb4qa.html)，同时可以将`~/Library/Containers/com.coderforart.MWeb3/Data/Library/Application Support/MWebLibrary`使用git备份到github，但是实际使用中，备份到github之后文件夹结构没有了，文件名也变成了`UUID`。完全没有可读性。
所以`mweb-export`就是一个解析Mweb数据库并生成一个包含整个目录结构和索引跳转的README文件。




## Install

### Go
```shell
go install github.com/TBXark/mweb-export@latest
```

### Brew
```shell
brew install --build-from-source tbxark/repo/xcode-tool
```



## Usage

```
Usage of mweb-export:
  -help
    	show usage
  -mode string
    	'file': save file, 'debug': print only (default "debug")
  -path string
    	path to MWebLibrary (default "~/Library/Containers/com.coderforart.MWeb3/Data/Library/Application Support/MWebLibrary")
  -target string
    	export README.md directory (default "$(pwd)")
```


实际使用中不推荐直接在`MWebLibrary`原始文件夹中使用`git`，避免多台设备同步`git`与iCloud发生冲突，可以另外新建新建一个`git`文件夹，在`git`目录下添加下面脚本,需要同步的时候执行下面脚本，会自动将`MWebLibrary`拷贝到`git`文件夹，并生成`README.md`文件并且自动`push`。


```shell
#!/bin/bash

git stash --include-untracked
git pull
/bin/rm -rf docs # 删除所有文件，重新拷贝，避免有已删除的文件继续留在git
/bin/rm -f mainlib.db
/bin/rm -rf metadata
cp ~/Library/Containers/com.coderforart.MWeb3/Data/Library/Application\ Support/MWebLibrary/mainlib.db mainlib.db
cp -R ~/Library/Containers/com.coderforart.MWeb3/Data/Library/Application\ Support/MWebLibrary/docs docs
cp -R ~/Library/Containers/com.coderforart.MWeb3/Data/Library/Application\ Support/MWebLibrary/metadata metadata
mweb-export -mode=save
git add .
git commit -a -m $(date -u +%Y-%m-%dT%H:%M:%SZ)
git push origin master
git stash clear
```
