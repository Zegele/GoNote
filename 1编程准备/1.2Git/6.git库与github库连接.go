//在github中创建好库后，会显示入如下代码：
git remote add origin git@github.com:Zegele/GoNote.git
//复制改代码到git中执行，表示将github仓库与本机仓库连接。
git push -u origin master
//将文档从本地，推到github仓库。
//第一次需要这样写完整，后面的可以直接写git push，即可将代码推到github


git remote add origin git@github.com:Zegele/GoNote.git
//将这段代码粘贴进git中（此时的git是项目的位置），这表示本地代码库和github中名叫GoNote的库联系在一起了。

//然后输入以下代码：
git push -u origin master
//等待上传，git中显示上传进度，上传完毕后，刷新github，则可以在github中刚看到相应的代码文件。

//以后更新代码后，git add .  --> git commit -m'xxx' --> git push即可上传到github
//以后的git push后可以不用跟分支，默认push当前支。