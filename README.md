# 解析ibd文件，方便查看文件结构

### ibd管理结构
space：表空间

segment：段

extent：区

page：页

表空间管理段，段管理区，区管理页

### 执行
`git clone https://github.com/jemuelmiao/parseibd.git`

`cd parseibd`

`./build.sh`

`cd cmd`

`./parseibd`

### 输出文件说明
结果文件存放在output目录中

btree_xxx：索引btree page关系

extents：全局extent列表，包括：空闲列表、部分使用列表、全部使用列表，由page fsp管理

inodes：全局page inode列表，包括：部分使用列表、全部使用列表，由page fsp管理

pages：所有page的编号、类型

rec_xxx：索引记录，包括：聚簇索引非叶子记录、聚簇索引叶子记录、二级索引非叶子记录、二级索引叶子记录

segments：所有page inode管理的segment列表及segment管理的extent列表

### 相关阅读
http://www.miaozhouguang.com/
