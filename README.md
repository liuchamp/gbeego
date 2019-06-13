# MUL
一个基于beego的api接口服务。

## 使用方式
在使用前，获取源码。做如下修改：
1. 更改正确的环境配置，在runscript.sh中，更改连接信息，数据库db信息等。
2. 正确打包
### 更改连接配置
在runscript.sh中有如下内容，配置时，需要请注意连接信息，尤其是分库处理的服务。
```Bash
export MGO_HOSTS= # MongoDB连接信息
export MGO_DATABASE= # MongoDB连接到那个数据库
export MGO_USERNAME=
export MGO_PASSWORD=
export USER_DATABASE= # user接口对应的表名
```
### 正确打包

#### 普通打包
只需要执行:
```Bash
bee pack
```
就回得到一个名为mul.tar.gz的压缩文件.
#### Docker打包
通用的docker生成image的方式。源码文件中提供了对应的Dockerfile。
在使用时，可以强制控制环境变量，也就是上面的runscript.sh脚本中设置的环境变量，保证数据库连接正常。
## 开发方式
执行如下命令：
```Bash
source runscript.sh 
bee run -downdoc=true -gendoc=true
```

