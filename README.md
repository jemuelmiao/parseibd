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

btree_xxx：索引btree page关系

extents：全局extent列表，包括：空闲列表、部分使用列表、全部使用列表，由page fsp管理

inodes：全局page inode列表，包括：部分使用列表、全部使用列表，由page fsp管理

pages：所有page的编号、类型

rec_xxx：索引记录，包括：聚簇索引非叶子记录、聚簇索引叶子记录、二级索引非叶子记录、二级索引叶子记录

segments：所有page inode管理的segment列表及segment管理的extent列表

### 示例
表结构

![image](https://user-images.githubusercontent.com/28854032/199264218-0361c68e-3f1e-44e2-8393-71a99d367c22.png)

输出结果文件

- btree_card_id

![image](https://user-images.githubusercontent.com/28854032/199232479-2a41650d-8d6e-4878-a7a7-e88ffde6c671.png)

- btree_id

![image](https://user-images.githubusercontent.com/28854032/199257672-e30e6f04-752c-4387-9014-8a7ee263fef3.png)

- extents

![image](https://user-images.githubusercontent.com/28854032/199257930-eb94cf7e-e941-47f6-ab5b-6cd6faed1ef5.png)

- inodes

![image](https://user-images.githubusercontent.com/28854032/199258039-78304f04-39d7-4517-99a3-5c71983732e3.png)

- pages

![image](https://user-images.githubusercontent.com/28854032/199258116-53c982e8-f7d3-43e9-9ee8-841beed85697.png)

- rec_card_id

![image](https://user-images.githubusercontent.com/28854032/199258241-bbfae17c-5fe1-470d-8c38-c7484da316b0.png)

- rec_id

![image](https://user-images.githubusercontent.com/28854032/199258323-a204fe3b-73ed-4c14-87a8-79b22f8f4f5b.png)

- segments

![image](https://user-images.githubusercontent.com/28854032/199258503-6bc630d6-0dcf-432a-86f5-0f9c2f356725.png)

### TODO

- 解析数据字典、frm文件，去掉连接mysql读取元数据的依赖
- 部分不常用的数据类型解析，如point、geometry等
- 前端可视化展示结果

### 相关阅读
http://www.miaozhouguang.com/
