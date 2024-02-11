# lab1
## MapReduce基本逻辑：
有一个主节点和其他从属工作节点，主节点负责获取文档内容以及子任务划分与分派。工作节点负责对子任务的处理（MapReduce操作）。
map：此项目为把文档每个词分割为一个字并存储在keyvalue的map中，并通过哈希函数进行分发到每个reduce。
reduce：reduce负责对map的数据进行处理，此项目为统计字符出现频率