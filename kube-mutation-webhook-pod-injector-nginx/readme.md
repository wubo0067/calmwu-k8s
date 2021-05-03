#### 安装cfssl工具

1. 下载cfssl工具

   ```
   wget https://pkg.cfssl.org/R1.2/cfssl_linux-amd6
   
   wget https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
   ```

2. 拷贝到命令目录

   ```
   mv cfssljson_linux-amd64 /usr/local/bin/cfssljson
   
   mv cfssl-certinfo_linux-amd64 /usr/local/bin/cfssl-certinfo
   
   mv cfssl_linux-amd64 /usr/local/bin/cfssl
   ```

3. 创建ca和server证书，在ca目录

   ```
   ./create-ca.sh
   ```

   

#### 获取cabundle

- cat ca.pem |base64|tr -d "\n"

- 将输出填写入替换掉文件mutatingwebhook.yaml中${CA_BUNDLE}

  ```
    clientConfig:
      service:
        name: nginx-injector-pod-webhook-svc
        namespace: nginx-injector-pod-webhook
        path: "/mutate"
      caBundle: ${CA_BUNDLE}
  ```

```

```



#### 部署

- 在apiserver中启用admission

  ```shell
  --enable-admission-plugins=NodeRestriction,MutatingAdmissionWebhook,ValidatingAdmissionWebhook
  ```

- 创建namespace以及server证书的secret。

  ```
  ./create-nginx-injector-pod-webhook-certs-secret.sh
  ```

- 创建sidecar的configmap。

  ```
  kubectl apply -f nginx-sidecar-configmap.yaml
  ```

- 创建sidecar容器niginx的配置文件，configmap，

  ```
  kubectl apply -f nginx-configmap.yaml
  ```

- 创建MutatingWebhookConfiguration资源。

  ```
  kubectl apply -f mutatingwebhook.yaml
  ```

- 查看命名空间下所有资源。

  ```
  kubectl api-resources --verbs=list --namespaced -o name | xargs -n 1 kubectl get --show-kind --ignore-not-found -n <namespace>
  ```

- 创建service

  ```
  kubectl apply -f nginx-injector-pod-webhook-svc.yaml
  ```



#### 制作镜像

- 执行make build，make image

  ![1619768552873](C:\Users\wubo0\AppData\Roaming\Typora\typora-user-images\1619768552873.png)

- 拷贝镜像

  ```
  docker save littlebull/nginx-injector-pod-webhook-server:v0.0.1 -o nginjector-wb.tar
  
  scp nginjector-wb.tar root@192.168.6.132:/home/calmwu/Dev/Downloads
  ```

- 在集群节点导入镜像

  ```
  ctr -n=k8s.io images import nginjector-wb.tar
  ```



#### 部署、测试

- 部署

  ```
  kubectl apply -f nginx-injector-pod-webhook-deployment.yaml
  ```

- 测试ping接口，返回pong说明服务ok

  ```
  curl --insecure https://10.96.200.99:8443/ping
  ```


- 创建测试用namespace，打注入标签

  ```
  kubectl create namespace ns-test-injector-nginx
  
  kubectl label namespace ns-test-injector-nginx nginx-injection=enabled
  ```

- 验证注入的deployment，重点**nginx-injector-pod-webhook/inject: "true"**

  ```
  apiVersion: extensions/v1beta1
  kind: Deployment
  metadata:
   name: sleep
   namespace: ns-test-injector-nginx
  spec:
   replicas: 1
   template:
    metadata:
     annotations:
       nginx-injector-pod-webhook/inject: "true"
    spec:
     containers:
     - name: sleep
       image: busybox
       command: ["/bin/sleep","infinity"]
       imagePullPolicy: IfNotPresent
  ```

  

