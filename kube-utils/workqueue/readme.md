### 结构
1：队列

2：processing set

3：dirty set

### 插入
1：当一个请求到来，首先加入dirty set或者先把dirty set中相同对象删除。

2：加入队列中，如果对象不存在processing set中。

### 处理
1：reconcile loop会从队列前端获取对象，将对象加入processing set，并且从dirty set中删除。

2：当对象处理完毕，这个对象会从processing set中删除，如果这个对象在dirty set中，该对象会插入队列队尾

### 并发保证
1：如果这个时候有处理中对象到来，这个对象只会加入到dirty set中，不会插入队列中，保证只有一个reconciling loop在处理这个对象。

### 队列的处理
1：第一次插入使用Add方法

2：在reconcile loop中，先获取Get，对象处理（很多场景这里仅仅是个key），处理成功：调用Forget让计数清零，处理失败：调用AddRateLimited重新插入队列，
   最后都要调用Done。
