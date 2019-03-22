在Pod级别设置Volume

查看pod中某个容器的log：kubectl logs volume-pod -c busybox
登录到日志产生容器查看：kubectl exec -it volume-pod -c tomcat -- ls /usr/local/tomcat/logs
                      kubectl exec -it volume-pod -c tomcat -- tail /usr/local/tomcat/logs/catalina.*log