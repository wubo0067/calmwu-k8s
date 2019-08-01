1. 创建secret

    kubectl create -f kata-secret.yaml
    
    挂载到pod中，使用volume
    ```
    apiVersion: v1
    kind: Pod
    metadata:
    labels:
        name: db
    name: db
    spec:
    volumes:
    - name: secrets
        secret:
        secretName: mysecret
    containers:
    - image: gcr.io/my_project_id/pg:v1
        name: db
        volumeMounts:
        - name: secrets
        mountPath: "/etc/secrets"
        readOnly: true
        ports:
        - name: cp
        containerPort: 5432
        hostPort: 5432
    ```


2. 创建configmap

    kubectl create -f kata-cm.yaml
