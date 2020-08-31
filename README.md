# 使用前需要先修改项目下的bat批处理文件中的对应配置

bat批处理文件:
@echo off

::将excel中的前四列转化为struct
::第一列字段类型		如 int
::第二列字段备注		如 符文等级上限
::第三列字段名		 如 id
::第四列s,c,all	 s表示服务端使用, c表示客户端使用, all表示都使用

::生成objs.go文件的路径
set savePath=F:\code\src\generateStruct\tool
::目标excel文件路径
set readPath=F:\code\src\generateStruct
::所有的字段类型
set allType=int,IntSlice,IntSlice2,IntSlice3,IntMap,string,StringSlice,float64,Condition,Conditions,ItemInfo,ItemInfos,ItemInfosSlice,PropInfo,PropInfos,ProbItem,ProbItems,HmsTime,HmsTimes,Defenderweights
echo 开始生成objs.go文件 注意:生成文件时需要关闭对应的excel文件
::%readPath%\generateStruct -savePath=%savePath% -readPath=%readPath% -allType=%allType%
generateStruct -savePath=%savePath% -readPath=%readPath% -allType=%allType%

echo 生成完毕，按任意键继续
TIMEOUT /T 999

::
::	项目中所有的字段类型(注意区分类型的大小写)
::	int				整型 				例如:1
::	IntSlice		整型的一维数组 	 	例如:1,2,3
::	IntSlice2		整型的二维数组  	例如:1,2,3;4,5,6
::	IntSlice3		整型的三维数组  	例如:1,2;3,4|5,6;7,8
::	IntMap			k和v都是int的集合  	例如:1,2;3,4;5,6
::	string			字符型  			例如:"无生无灭炉"
::	StringSlice		字符型的一维数组  	例如:"无生无灭炉","雕花青铜炉"
::	float64			浮点型  			例如:1.5
::	Condition		k,v,map[int]int  	例如:2,0
::	Conditions		Condition的一维数组 例如:2,0;2,0
::	ItemInfo		物品信息ItemId int, Count int 	例如:3200071,1
::	ItemInfos		ItemInfo的一维数组 	例如:3200071,1;3200072,1
::	ItemInfosSlice	ItemInfos的一维数组 例如:[[3200071,1;3200072,1],[3200071,1;3200072,1]]
::	PropInfo		k int,v int  		例如:3200071,1
::	PropInfos		PropInfo的一维数组  例如:3200071,1;3200072,1
::	ProbItem		物品掉落几率 id,count,prob 例如:3200071,1,10
::	ProbItems		ProbItem的一维数组  例如:3200071,1,10;3200072,1,5
::	HmsTime			时间类型 			例如:06:00:00
::	HmsTimes		HmsTime的一维数组 	例如:06:00:00;07:00:00
::	Defenderweights	物品掉落权重(序号，掉落最小值，掉落最大值，权重，CD秒) 	例如:1,41,50,9,28800;1,41,50,9,28800
