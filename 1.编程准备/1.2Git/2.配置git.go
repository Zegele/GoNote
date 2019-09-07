//配置user信息
//打开git Bush Here

git config --global user.name 'your_name'
git config --global user.email 'your_email@domain.com'
//例子：
git config --global user.name 'Ang'//设置一个你的名字
git config --global user.email 'studio@gmial.com'//设置一个你可以使用的邮箱。（这个邮箱是我随意写的）


//查看设置
git config --list global

//衍生
git config --local //只对某个仓库有效。注意：没有设置config的话等同于local
git config --global //对当前用户所有仓库有效，意思是对这台计算机的所有仓库有效。
git config --system //对系统所有登录的用户有效（一般也不怎么用）

//显示config的配置， 加--list
git config --list --local //显示local的设置
git config --list --global //显示global的设置
git config --list --system //显示system的设置

//注意：
git config --local --list//也可以查看，效果同上