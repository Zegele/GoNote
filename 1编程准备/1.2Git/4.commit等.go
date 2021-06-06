//代码新编辑后,在该项目中打开git，输入：
1.git add 
git add .//添加所有的文件。 . 表示所有文件
git add name//表示添加name这个文件。
//git status查看被提交文件的状态。
2.
git commit -m'描述'//注意：描述是必须写的，描述这次提交的标志。

连接到github可参考：
echo "# xuexiqianfeng" >> README.md
git init
git add README.md
git commit -m "first commit"
git branch -M main
git remote add origin git@github.com:Zegele/xuexiqianfeng.git
git push -u origin main