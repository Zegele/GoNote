1. 在系统之家下载合适的系统，然后阅读安装说明，进行安装。
网页上的说明已经很详细了。不过建议用手机（或另一台电脑）打开看着步骤指导电脑的重装。
如果遇到问题。查阅文件解决。

**（我只用了硬盘装机法，U盘装机还没鼓捣过。）**

重装过程中，我遇到0xc000000f问题

方法一：
开始-运行-输入msconfig，弹出系统配置选项表，点击“引导”，选择你的系统，并设置成为默认值，点确定即可。


方法二：需要系统U盘，可参考：https://blog.csdn.net/jessical1020/article/details/76461317


//用一种方法安装成功即可，以下也可以不用看。

/*
U盘怎样安装系统
# 64位系统u盘安装怎么装|u盘启动安装64位系统教程

怎么用U盘安装64位系统？现在电脑硬件性能越来越强，需要装64位系统才能发挥硬件性能。绝大多数装机员都是使用U盘来安装系统，只要电脑有usb接口，都可以通过U盘安装系统。不过64位系统镜像普遍大于4G，用传统的U盘启动盘会比较麻烦，所以有不少人不知道64位系统怎么U盘安装。下面系统城小编以安装[win7 64位](http://www.xitongcheng.com/win7/64/)为例，跟大家介绍U盘启动安装64位系统教程。

**安装须知：**

1、本文介绍legacy模式U盘安装64位系统方法，硬盘分区表是MBR

2、如果是uefi机型，打算uefi模式下安装，参考这个教程：[怎么用wepe装win10系统](http://www.xitongcheng.com/jiaocheng/xtazjc_article_42355.html)

**一、安装准备**

1、8G或更大容量U盘

2、制作微pe启动盘：[微pe工具箱怎么制作u盘启动盘](http://www.xitongcheng.com/jiaocheng/xtazjc_article_42199.html)

3、系统镜像下载：[技术员联盟ghost win7 64位旗舰最新版V2018.05](http://www.xitongcheng.com/win7/jsylm_xiazai_4746.html)

**二、U盘启动设置**：[bios设置u盘启动方法](http://www.xitongcheng.com/jiaocheng/xtazjc_article_41832.html)

![64位系统u盘安装怎么装|u盘启动安装64位系统教程](https://upload-images.jianshu.io/upload_images/4370290-295e7d41fc5dde55.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

**三、64位系统U盘安装步骤如下**

1、根据安装准备的教程制作好微pe启动U盘，然后把下载的64位系统镜像iso直接复制到U盘中，镜像大于4G也可以放进去；

![64位系统u盘安装怎么装|u盘启动安装64位系统教程](https://upload-images.jianshu.io/upload_images/4370290-6bae09ce374529a7.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

2、在需要装64位系统的电脑上插入微pe工具箱，重启过程中不停按F12或F11或Esc等启动快捷键调出启动菜单，选择识别到的U盘选项，一般是带有USB的选项，或者是U盘的品牌名称，比如Toshiba、Sandisk或者Generic Flash Disk。如果同时出现两个U盘项，选择不带uefi的项，表示在legacy模式下安装，选择之后按回车键；

![64位系统u盘安装怎么装|u盘启动安装64位系统教程](https://upload-images.jianshu.io/upload_images/4370290-75568c370bde6012.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

3、进入到pe系统，如果不需要全盘重新分区，直接执行第6步，如果打算重新分盘，双击桌面上的【分区工具DiskGenius】，右键HD0整个硬盘，选择【快速分区】；

![64位系统u盘安装怎么装|u盘启动安装64位系统教程](https://upload-images.jianshu.io/upload_images/4370290-42892b99336fb577.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

4、默认是MBR分区表类型，设置分区数目、分区大小，卷标为【系统】的表示系统盘(C盘)，建议50G以上，如果是固态硬盘，勾选【对齐分区到此扇区数的整数倍】，默认2048即可4k对齐，选择4096也可以，最后点击确定；

![64位系统u盘安装怎么装|u盘启动安装64位系统教程](https://upload-images.jianshu.io/upload_images/4370290-a65f6eca8fd3082b.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

5、执行硬盘重新分区过程，等待一会儿即可，分区之后，如图所示，如果盘符错乱，右键选择更改驱动器路径，自行修改；

![64位系统u盘安装怎么装|u盘启动安装64位系统教程](https://upload-images.jianshu.io/upload_images/4370290-0ac0be7ffdc7b154.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

6、完成分区之后，打开此电脑—微pe工具箱，右键系统iso镜像，选择【装载】，如果没有装载，则右键—打开方式—资源管理器打开；

![64位系统u盘安装怎么装|u盘启动安装64位系统教程](https://upload-images.jianshu.io/upload_images/4370290-3a79d4134ca975b2.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

7、win10pe支持直接打开iso镜像，如图所示，运行绿色图标【双击安装系统】；

![64位系统u盘安装怎么装|u盘启动安装64位系统教程](https://upload-images.jianshu.io/upload_images/4370290-228cc8677a443a40.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

8、选择【还原分区】，GHO WIM ISO映像路径自动提取到gho文件，安装位置是通常是C盘，建议看仔细点，可能不是显示C盘，可以根据卷标或总大小来判断，最后点击确定；

![64位系统u盘安装怎么装|u盘启动安装64位系统教程](https://upload-images.jianshu.io/upload_images/4370290-d4a098bffc6ac304.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

9、弹出提示框，勾选【完成后重启】和【引导修复】，点击是；

![64位系统u盘安装怎么装|u盘启动安装64位系统教程](https://upload-images.jianshu.io/upload_images/4370290-330751f62d4ccb81.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

10、在这个界面中，执行系统安装部署到C盘的过程，等待进度条；

[图片上传中...(image-f7c115-1565055085329-2)]

11、操作完成后自动重启，**重启时拔出U盘**，进入这个界面，执行系统组件、驱动安装、系统配置和激活过程；

![64位系统u盘安装怎么装|u盘启动安装64位系统教程](https://upload-images.jianshu.io/upload_images/4370290-619ab481d0f483c7.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

12、之后还会再重启一次，为首次使用计算机做准备，最后重启进入系统桌面，64位系统安装完成。

![64位系统u盘安装怎么装|u盘启动安装64位系统教程](https://upload-images.jianshu.io/upload_images/4370290-e5f037af33d1b59a.png?imageMogr2/auto-orient/strip%7CimageView2/2/w/1240)

U盘装64位系统的方法就为大家介绍到这边，U盘安装不会难，比较关键的一步是设置U盘启动，后面都好办。
*/