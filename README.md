# 解析ibd文件，方便查看文件结构

### ibd管理结构
space：表空间，管理segment和page

segment：段，管理extent

extent：区，管理page

page：页，最小存储单元，默认16KB

### 前提
parseibd默认page采用16KB大小，可通过变量innodb_page_size查看，并且表空间都是独立文件，可通过变量innodb_file_per_table查看

### 执行
`git clone https://github.com/jemuelmiao/parseibd.git`

`cd parseibd`

`./build.sh`

`cd cmd`

`./parseibd -h {mysql地址} -u {用户名} -p {密码} -d {数据库名} -t {表名} -f {ibd文件路径}`

### 输出文件说明
结果文件存放在output目录中

btree_xxx：索引btree page关系，可以方便查看btree的层级及page之间的连接关系

extents：全局extent列表，包括：空闲列表、部分使用列表、全部使用列表，由首页page fsp管理

inodes：全局page inode列表，包括：部分使用列表、全部使用列表，由首页page fsp管理

pages：所有page的编号、类型

rec_xxx：索引记录，包括：聚簇索引非叶子记录、聚簇索引叶子记录、二级索引非叶子记录、二级索引叶子记录

segments：所有page inode管理的segment列表及segment管理的extent列表

### 示例
表结构

![image](https://user-images.githubusercontent.com/28854032/199284381-90e75dab-838b-4786-98fe-871e79b1a9d5.png)

输出结果文件

- btree_card_id

![image](https://user-images.githubusercontent.com/28854032/199284776-cf58c3a2-8937-479c-b4af-389b11ac20c5.png)

- btree_id

![image](https://user-images.githubusercontent.com/28854032/199284848-9b72b573-dd73-4537-9172-c65c784d521b.png)

- extents

![image](https://user-images.githubusercontent.com/28854032/199284933-6c6b858e-c302-49b1-9400-e0f6a1086013.png)

- inodes

![image](https://user-images.githubusercontent.com/28854032/199285001-44320e88-e887-4c48-b4cd-08d5c786f337.png)

- pages

![image](https://user-images.githubusercontent.com/28854032/199285070-cf147f5f-506f-44dd-b59d-1fe0c109c4b1.png)

- rec_card_id

![image](https://user-images.githubusercontent.com/28854032/199285199-e2078851-2ced-46cf-8ed0-606bb643ef77.png)

- rec_id

![image](https://user-images.githubusercontent.com/28854032/199285268-f6f5b60a-681c-4169-be72-f198b45c2dc2.png)

- segments

![image](https://user-images.githubusercontent.com/28854032/199285339-6d852ab3-0e76-4ce5-9565-bb29348e52db.png)

### TODO

- 解析数据字典、frm文件，去掉连接mysql读取元数据的依赖
- 部分不常用的数据类型解析，如point、geometry等
- 前端可视化展示结果

### 相关阅读
http://www.miaozhouguang.com/
