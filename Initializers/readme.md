### 初始化

通过该配置创建资源对象 InitializerConfiguration 之后，就会在每个 Pod 的 metadata.initializers.pending 字段中添加 custom-pod-initializer 字段。该初始化控制器会定期扫描新的 Pod，一旦在 Pod 的 pending 字段中检测到自己的名称，就会执行其逻辑，执行完逻辑之后就会将 pending 字段下的自己的名称删除。

只有在 pending 字段下的列表中的第一个Initializers可以对资源进行操作，当所有的Initializers执行完成，并且 pending字段为空时，该对象就会被认为初始化成功。

你可能会注意到一个问题：如果 kube-apiserver 不能显示这些资源，那么用户级控制器是如何处理资源的呢？

为了解决这个问题，kube-apiserver 暴露了一个 ?includeUninitialized 查询参数，它会返回所有的资源对象（包括未初始化的）。