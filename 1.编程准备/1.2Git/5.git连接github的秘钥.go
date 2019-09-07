//Windows系统
//在git中输入以下命令
 ssh-keygen.exe
/*然后连续3次回车，不用设置，生成一对秘钥
id_rsa是私钥
id_rsa.pub是公钥
在系统 “用户 -- .ssh” 的文件夹内可以看到。*/

//在git中输入
cd id_rsa//可以查看私钥
cd id_rsa.pub//可以查公钥

//把公钥添加到github中，就可以连接到github中了。（注意：复制公钥可以不复制公钥中的计算机名。）
//这样就可以从github中clone等操作。


//linux系统
//密钥是放在 ~/.ssh目录中的