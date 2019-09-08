//Windows系统
//在git中输入以下命令
 ssh-keygen.exe
/*然后连续3次回车，不用设置，生成一对秘钥
id_rsa是私钥
id_rsa.pub是公钥
在系统 “用户 -- .ssh” 的文件夹内可以看到。*/


//查看密钥，在git中输入
方法一
cd id_rsa//可以查看私钥//cd可以打开文件夹，也可以打开文件。
cd id_rsa.pub//可以查公钥

方法二
cat ~/.ssh/id_rsa.pub  
（cat命令是打印出文件内容。id_rsa.pub是上面刚生成的公钥。~代表“用户”地址，一般叫做home？？？）

方法三
进入“用户”地址，打开.ssh文件，打开key文件，复制key

//把公钥添加到github中，就可以连接到github中了。（注意：复制公钥可以不复制公钥中的计算机名。）
//这样就可以从github中clone等操作。


//linux系统
//密钥是放在 ~/.ssh目录中的




与github关联
//复制key后，进入github的settings，选中SSH and GPG keys

在title中输入标题（可以是任意内容）
在key中粘贴之前复制的key

点击 add SSH key

这样你的git就和你的github关联了。

参考：
[https://www.cnblogs.com/MrReed/p/6373988.html](https://www.cnblogs.com/MrReed/p/6373988.html)

[https://blog.csdn.net/weixin_42063071/article/details/80999690](https://blog.csdn.net/weixin_42063071/article/details/80999690)
