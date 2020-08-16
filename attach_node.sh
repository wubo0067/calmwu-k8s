#!/usr/bin/env bash

################################################
#
# KUBE_VERSION     the expected kubernetes version
# eg.  ./attach_node.sh
#           --docker-version 17.06.2-ce-1 \
#           --token 264db1.30bcc2b89969a4ca \
#           --endpoint 192.168.0.80:6443
#           --cluster-dns 172.19.0.10
################################################

set -e -x

PKG=pkg

export CLOUD_TYPE=public
export nvidiaDriverInstallUrl=http://aliacs-k8s-cn-hangzhou.oss-cn-hangzhou.aliyuncs.com/public/pkg/run/attach/1.9/nvidia-gpu.sh
export openapi=http://cs-anony.aliyuncs.com

public::common::log(){
    echo $(date +"[%Y%m%d %H:%M:%S]: ") $1
}

public::common::region()
{
    if [ "$CLOUD_TYPE" == "public" ];then
        region=$(curl --retry 5  -sSL http://100.100.100.200/latest/meta-data/region-id)
        if [ "" == "$region" ];then
            kube::common::log "can not get regionid and instanceid! \
                curl --retry 5 -sSL http://100.100.100.200/latest/meta-data/region-id" && exit 256
        fi
        if [ "$region" == "cn-beijing" \
            -o "$region" == "cn-shanghai" \
            -o "$region" == "cn-shenzhen" \
            -o "$region" == "cn-qingdao" \
            -o "$region" == "cn-zhangjiakou" ];then

            region=cn-hangzhou
        fi
        if [ "$region" == "cn-hongkong" ];then
            export KUBE_REPO_PREFIX=registry.cn-hangzhou.aliyuncs.com/acs
        fi
        export REGION=$region
    else
        public::common::log "Do nothing for aglity"
    fi
}

public::common::os_env()
{
    ubu=$(cat /etc/issue|grep "Ubuntu 16.04"|wc -l)
    cet=$(cat /etc/centos-release|grep "CentOS"|wc -l)
    suse=$(cat /etc/os-release |  grep "SUSE"|wc -l)
    redhat=$(cat /etc/redhat-release|grep "Red Hat"|wc -l)
    alios=$(cat /etc/redhat-release|grep "Alibaba"|wc -l)
    if [ "$ubu" == "1" ];then
        export OS="Ubuntu"
    elif [ "$cet" == "1" ];then
        export OS="CentOS"
    elif [ "$suse" == "1" ];then
        export OS="SUSE"
    elif [ "$redhat" == "1" ];then
        export OS="RedHat"
    elif [ "$alios" == "1" ];then
        export OS="AliOS"
    else
        public::common::log "unkown os...   exit"
        exit 1
    fi
}

public::common::prepare_package(){
    PKG_TYPE=$1
    PKG_VERSION=$2
    if [ ! -f ${PKG_TYPE}-${PKG_VERSION}.tar.gz ];then
        if [ -z $PKG_FILE_SERVER ] ;then
            public::common::log "local file ${PKG_TYPE}-${PKG_VERSION}.tar.gz does not exist, And PKG_FILE_SERVER is not config"
            public::common::log "installer does not known where to download installer binary package without PKG_FILE_SERVER env been set. Error: exit"
            exit 1
        fi
        public::common::log "local file ${PKG_TYPE}-${PKG_VERSION}.tar.gz does not exist, trying to download from [$PKG_FILE_SERVER]"
        curl --retry 4 $PKG_FILE_SERVER/$CLOUD_TYPE/pkg/$PKG_TYPE/${PKG_TYPE}-${PKG_VERSION}.tar.gz \
                > ${PKG_TYPE}-${PKG_VERSION}.tar.gz || (public::common::log "download failed with 4 retry,exit 1" && exit 1)
    fi
    tar -xvf ${PKG_TYPE}-${PKG_VERSION}.tar.gz || (public::common::log "untar ${PKG_VERSION}.tar.gz failed!, exit" && exit 1)
}

public::common::nodeid()
{
    if [ "$CLOUD_TYPE" == "public" ];then
        region=$(curl --retry 5  -sSL http://100.100.100.200/latest/meta-data/region-id)
        insid=$(curl --retry 5  -sSL http://100.100.100.200/latest/meta-data/instance-id)
        if [ "" == "$region" -o "" == "$insid" ];then
            kube::common::log "can not get regionid and instanceid! \
            curl --retry 5 -sSL http://100.100.100.200/latest/meta-data/region-id" && exit 256
        fi
        export NODE_ID=$region.$insid
    else
        public::common::log "Do nothing for aglity"
    fi
}

function version_gt() { test "$(echo "$@" | tr " " "\n" | sort -V | head -n 1)" != "$1"; }

function retry()
{
        local n=0;local try=$1
        local cmd="${@: 2}"
        [[ $# -le 1 ]] && {
            echo "Usage $0 <retry_number> <Command>";
        }
        set +e
        until [[ $n -ge $try ]]
        do
                $cmd && break || {
                        echo "Command Fail.."
                        ((n++))
                        echo "retry $n :: [$cmd]"
                        sleep 2;
                        }
        done
        set -e
}

public::common::install_package(){

    public::docker::install

    public::common::prepare_package "kubernetes" $KUBE_VERSION

    sed -i "s#registry.cn-hangzhou.aliyuncs.com#registry-vpc.$REGION.aliyuncs.com#g" pkg/kubernetes/$KUBE_VERSION/module/*
    if [ -f /proc/sys/fs/may_detach_mounts ];then
        sed -i "/fs.may_detach_mounts/ d" /etc/sysctl.conf
        echo "fs.may_detach_mounts=1" >> /etc/sysctl.conf
        sysctl -p|| true
    fi

    if [ "$OS" == "CentOS" ] || [ "$OS" == "RedHat" ] || [ "$OS" == "AliOS" ];then
        dir=pkg/kubernetes/$KUBE_VERSION/rpm

        # Install nfs client and socat
        yum install -y socat fuse fuse-libs nfs-utils nfs-utils-lib pciutils

        yum localinstall -y `ls $dir | xargs -I '{}' echo -n "$dir/{} "`

        sed -i '/net.bridge.bridge-nf-call-iptables/d' /usr/lib/sysctl.d/00-system.conf
        sed -i '$a net.bridge.bridge-nf-call-iptables = 1' /usr/lib/sysctl.d/00-system.conf
        echo 1 > /proc/sys/net/bridge/bridge-nf-call-iptables
    elif [ "$OS" == "Ubuntu" ];then
        dir=pkg/kubernetes/$KUBE_VERSION/debain
        dpkg -i `ls $dir | xargs -I '{}' echo -n "$dir/{} "`
    elif [ "$OS" == "SUSE" ];then
        dir=pkg/kubernetes/$KUBE_VERSION/rpm
        zypper --no-gpg-checks install -y `ls $dir | xargs -I '{}' echo -n "$dir/{} "`
    fi

    sed -i "s#--cluster-dns=10.96.0.10 --cluster-domain=cluster.local#--cluster-dns=$CLUSTER_DNS \
    --pod-infra-container-image=$KUBE_REPO_PREFIX/pause-amd64:3.0 \
    --enable-controller-attach-detach=false \
    --cluster-domain=cluster.local --cloud-provider=external \
    --hostname-override=$NODE_ID --provider-id=$NODE_ID#g" \
    /etc/systemd/system/kubelet.service.d/10-kubeadm.conf
    if docker info|grep 'Cgroup Driver'|grep cgroupfs; then
        sed -i -e 's/cgroup-driver=systemd/cgroup-driver=cgroupfs/g' \
        /etc/systemd/system/kubelet.service.d/10-kubeadm.conf
    fi

    # set node resource reserve when available memory > 3.5G
    if [[ $(cat /proc/meminfo | awk '/MemTotal/ {print $2}') -gt $((1024*512*7)) ]]; then
		sed -i -e "s/--cgroup-driver=/--system-reserved=memory=300Mi \
		--kube-reserved=memory=400Mi \
		--eviction-hard=imagefs.available<15%,memory.available<300Mi,nodefs.available<10%,nodefs.inodesFree<5% --cgroup-driver=/g" \
		/etc/systemd/system/kubelet.service.d/10-kubeadm.conf
	fi

    systemctl daemon-reload ; systemctl enable kubelet.service ; systemctl start kubelet.service
}

public::docker::install()
{
    set +e
    docker version > /dev/null 2>&1
    i=$?
    set -e
    v=$(docker version|grep Version|awk '{gsub(/-/, ".");print $2}'|uniq)
    if [ $i -eq 0 ]; then
        if [[ "$DOCKER_VERSION" == "$v" ]];then
            public::common::log "docker has been installed , return. $DOCKER_VERSION"
            return
        fi
    fi
    public::common::prepare_package "docker" $DOCKER_VERSION
    if [ "$OS" == "CentOS" ] || [ "$OS" == "RedHat" ] || [ "$OS" == "AliOS" ];then
        if [ "$(rpm -qa docker-engine-selinux|wc -l)" == "1" ];then
            yum erase -y docker-engine-selinux
        fi
        if [ "$(rpm -qa docker-engine|wc -l)" == "1" ];then
            yum erase -y docker-engine
        fi
        if [ "$(rpm -qa docker-ce|wc -l)" == "1" ];then
            yum erase -y docker-ce
        fi
        if [ "$(rpm -qa container-selinux|wc -l)" == "1" ];then
            yum erase -y container-selinux
        fi

        if [ "$(rpm -qa docker-ee|wc -l)" == "1" ];then
            yum erase -y docker-ee
        fi

        local pkg=pkg/docker/$DOCKER_VERSION/rpm/
        yum localinstall -y `ls $pkg |xargs -I '{}' echo -n "$pkg{} "`
    elif [ "$OS" == "Ubuntu" ];then
        if [ "$need_reinstall" == "true" ];then
            if [ "$(echo $v|grep ee|wc -l)" == "1" ];then
                apt purge -y docker-ee docker-ee-selinux
            elif [ "$(echo $v|grep ce|wc -l)" == "1" ];then
                apt purge -y docker-ce docker-ce-selinux container-selinux
            else
                apt purge -y docker-engine
            fi
        fi
        dir=pkg/docker/$DOCKER_VERSION/debain
        dpkg -i `ls $dir | xargs -I '{}' echo -n "$dir/{} "`
    elif [ "$OS" == "SUSE" ];then
        if [ "$(rpm -qa docker-engine-selinux|wc -l)" == "1" ];then
            zypper rm -y docker-engine-selinux
        fi
        if [ "$(rpm -qa docker-engine|wc -l)" == "1" ];then
            zypper rm -y docker-engine
        fi
        if [ "$(rpm -qa docker-ce|wc -l)" == "1" ];then
            zypper rm -y docker-ce
        fi
        if [ "$(rpm -qa container-selinux|wc -l)" == "1" ];then
            zypper rm -y container-selinux
        fi

        if [ "$(rpm -qa docker-ee|wc -l)" == "1" ];then
            zypper rm -y docker-ee
        fi
        local pkg=pkg/docker/$KUBE_VERSION/rpm/
        zypper  --no-gpg-checks install -y `ls $pkg |xargs -I '{}' echo -n "$pkg{} "`
    else
        public::common::log "install docker with [unsupported OS version] error!"
        exit 1
    fi
    public::docker::config
}

public::docker::config()
{
    iptables -P FORWARD ACCEPT
    if [ "$OS" == "CentOS" ] || [ "$OS" == "RedHat" ] || [ "$OS" == "AliOS" ];then
        #setenforce 0
        sed -i -e 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config
        sed -i 's#LimitNOFILE=infinity#LimitNOFILE=1048576#g' /lib/systemd/system/docker.service
    fi

    if public::common::region4mirror "$REGION" ; then
       if [ "$DOCKER_VERSION" == "17.06.2.ce" ];then
           sed -i "s#ExecStart=/usr/bin/dockerd#ExecStart=/usr/bin/dockerd -s overlay2 \
           --storage-opt overlay2.override_kernel_check=true \
           --registry-mirror=https://pqbap4ya.mirror.aliyuncs.com --log-driver=json-file \
           --log-opt max-size=100m --log-opt max-file=10#g" /lib/systemd/system/docker.service
       else
           sed -i "s#ExecStart=/usr/bin/dockerd#ExecStart=/usr/bin/dockerd -s overlay2 \
           --storage-opt overlay2.override_kernel_check=true \
           --exec-opt native.cgroupdriver=systemd \
           --registry-mirror=https://pqbap4ya.mirror.aliyuncs.com --log-driver=json-file \
           --log-opt max-size=100m --log-opt max-file=10#g" /lib/systemd/system/docker.service
       fi
    else
       if [ "$DOCKER_VERSION" == "17.06.2.ce" ];then
           sed -i "s#ExecStart=/usr/bin/dockerd#ExecStart=/usr/bin/dockerd -s overlay2 \
           --storage-opt overlay2.override_kernel_check=true \
           --log-driver=json-file \
           --log-opt max-size=100m --log-opt max-file=10#g" /lib/systemd/system/docker.service
       else
           sed -i "s#ExecStart=/usr/bin/dockerd#ExecStart=/usr/bin/dockerd -s overlay2 \
           --storage-opt overlay2.override_kernel_check=true \
           --exec-opt native.cgroupdriver=systemd \
           --log-driver=json-file \
           --log-opt max-size=100m --log-opt max-file=10#g" /lib/systemd/system/docker.service
       fi
    fi
    sed -i "/ExecStart=/a\ExecStartPost=/usr/sbin/iptables -P FORWARD ACCEPT" /lib/systemd/system/docker.service

    systemctl daemon-reload ; systemctl enable  docker.service; systemctl restart docker.service
}


public::common::region4mirror(){
    local region=$1

    if [ "$region" == "cn-hangzhou" \
         -o "$region" == "cn-beijing" \
         -o "$region" == "cn-shanghai" \
         -o "$region" == "cn-qingdao" \
         -o "$region" == "cn-zhangjiakou" \
         -o "$region" == "cn-shenzhen" \
         -o "$region" == "cn-huhehaote" ];then

        true
    else
        false
    fi
}

public::main::node_up()
{
	public::common::nodeid

    public::common::install_package

    if [ $AUTO_FDISK != "" ]; then
      public::main::mount_disk
    fi

    if [[ $KUBE_VERSION = *"1.8"* ]]; then
        echo "Not support GPU yet for $KUBE_VERSION!"
    else
        public::main::nvidia_gpu
    fi

    if [ $FROM_ESS != "" ]; then
      public::main::attach_label
    fi

    discover="--discovery-token-unsafe-skip-ca-verification"

    kubeadm join --node-name "$NODE_ID" --token $TOKEN $APISERVER_LB $discover

    retry 30 grep cluster /etc/kubernetes/kubelet.conf
    cp /etc/kubernetes/kubelet.conf /etc/kubernetes/kubelet.old

    if [[ "$APISERVER_LB" != "1.1.1.1" ]];then
        if [[ -f /etc/kubernetes/kubelet.conf ]];then
            sed -i "s#server: https://.*\$#server: https://$APISERVER_LB#g" /etc/kubernetes/kubelet.conf
        fi
    fi

    systemctl restart kubelet
}

public::main::attach_label()
{
  if [ $LABELS != "" ]; then
    if [ "1" == "$GPU_FOUNDED" ];then
      sed -i "s/^Environment=\"KUBELET_EXTRA_ARGS=--node-labels=/Environment=\"KUBELET_EXTRA_ARGS=--node-labels=$LABELS,/" /etc/systemd/system/kubelet.service.d/10-kubeadm.conf
    else
      KUBELET_EXTRA_ARGS="--node-labels=$LABELS"
      #sed -i '/^ExecStart=$/iEnvironment="KUBELET_EXTRA_ARGS=--feature-gates=DevicePlugins=true"' /etc/systemd/system/kubelet.service.d/10-kubeadm.conf
      sed -i "/^ExecStart=$/iEnvironment=\"KUBELET_EXTRA_ARGS=$KUBELET_EXTRA_ARGS\"" /etc/systemd/system/kubelet.service.d/10-kubeadm.conf
    fi
    systemctl daemon-reload
    systemctl restart kubelet
  else
    echo "No extra labels for this node"
  fi
}

public::main::mount_disk() {
    #config /etc/fstab and mount device
    DISK_ATTACH_POINT="/dev/xvdb"
    if [ -b "/dev/vdb" ]; then
        DISK_ATTACH_POINT="/dev/vdb"
    fi
    set +e
    cat /etc/fstab | grep $DISK_ATTACH_POINT >/dev/null 2>&1
    i=$?
    set -e
    if [ $i -eq 0 ]; then
        echo "disk ${DISK_ATTACH_POINT} is already mounted"
        return
    fi

    if [ -b "$DISK_ATTACH_POINT" ]; then
        echo "start to format disk ${DISK_ATTACH_POINT}"

        local n=0
        local max_retry=6
        # Clean all partition info
        while [[ "$n" -le $max_retry  ]]
        do
            set +e
            # Check if the disk has partition info and delete it
            fdisk -S 56 "${DISK_ATTACH_POINT}" <<-EOF | grep "${DISK_ATTACH_POINT}" | grep "Linux"
p
wq
EOF
            i=$?
            if [ $i -eq 0 ]; then
                echo "disk ${DISK_ATTACH_POINT} is already partitioned, deleting the partition info"
                fdisk -S 56 "${DISK_ATTACH_POINT}" <<-EOF
d

wq
EOF
            else
                break
            fi
            n=`expr $n + 1`
            set -e
            sleep 3
        done

        fdisk -S 56 "${DISK_ATTACH_POINT}" <<-EOF
n
p
1


wq
EOF

        sleep 5
        mkfs.ext4 -i 8192 "${DISK_ATTACH_POINT}1"

        docker_flag=0
        if [ -d "/var/lib/docker" ]; then
            docker_flag=1
            service docker stop
            rm -fr /var/lib/docker
        fi
        mkdir /var/lib/docker
        kubelet_flag=0
        if [ -d "/var/lib/kubelet" ]; then
            kubelet_flag=1
            service kubelet stop
            rm -fr /var/lib/kubelet
        fi
        mkdir /var/lib/kubelet

        if [ ! -d /var/lib/container ]; then
            mkdir /var/lib/container/
        fi
        mount -t ext4 "${DISK_ATTACH_POINT}1" /var/lib/container/
        mkdir /var/lib/container/kubelet /var/lib/container/docker
        echo "${DISK_ATTACH_POINT}1    /var/lib/container/     ext4    defaults        0 0" >>/etc/fstab
        echo "/var/lib/container/kubelet /var/lib/kubelet none defaults,bind 0 0" >>/etc/fstab
        echo "/var/lib/container/docker /var/lib/docker none defaults,bind 0 0" >>/etc/fstab
        mount -a

        if [ $docker_flag == 1 ]; then
            service docker start
        fi
        if [ $kubelet_flag == 1 ]; then
            service kubelet start
        fi
        df -h
    else
        echo "no need to mount disk."
    fi
}

public::main::nvidia_gpu()
{
    current_region=$(curl --retry 5  -sSL http://100.100.100.200/latest/meta-data/region-id)
    if [ "" == "$current_region" ];then
        echo "can not get regionid and instanceid! \
            curl --retry 5 -sSL http://100.100.100.200/latest/meta-data/region-id" && exit 256
    fi

    if [ $current_region == "ap-southeast-3" ] || \
        [ $current_region == "ap-northeast-1" ] || \
        [ $current_region == "ap-southeast-1" ] || \
        [ $current_region == "ap-southeast-2" ] || \
        [ $current_region == "eu-central-1" ] || \
        [ $current_region == "us-east-1" ] || \
        [ $current_region == "cn-hongkong" ] || \
        [ $current_region == "us-west-1" ]; then
        export nvidiaDriverInstallUrl=http://aliacs-k8s-ap-southeast-1.oss-ap-southeast-1.aliyuncs.com/public/pkg/run/attach/1.9/nvidia-gpu.sh
    fi
    curl -L ${nvidiaDriverInstallUrl} -o nvidia-gpu.sh
    chmod u+x nvidia-gpu.sh
    source nvidia-gpu.sh
    public::nvidia::enable_gpu_capability
}

public::common::get_node_info()
{
    if [ "$CLOUD_TYPE" == "public" ];then
        insid=$(curl --retry 5  -sSL http://100.100.100.200/latest/meta-data/instance-id)
        info=$(curl --retry 5 -H "Date:`date -R`" -sfSL "$openapi/token/${OPENAPI_TOKEN}/instance/${insid}/node_info"|grep '\w')
        eval "$info"
        export DOCKER_VERSION=$docker_version
        export CLUSTER_DNS=$cluster_dns
        export TOKEN=$token
        export APISERVER_LB=$endpoint
    fi
}

public::common::callback()
{
    if [ "$CLOUD_TYPE" == "public" ];then
        curl -H "Date:`date -R`" -X POST -sfSL "${callback_url}"
        echo "======================================="
        echo "                 SUCCESS               "
        echo "======================================="
    fi
}

main()
{
    while [[ $# -gt 0 ]]
    do
    key="$1"

    case $key in
        --docker-version)
            export DOCKER_VERSION=$2
            shift
        ;;
        --cluster-dns)
            export CLUSTER_DNS=$2
            shift
        ;;
         --token)
            export TOKEN=$2
            shift
        ;;
        --endpoint)
            export APISERVER_LB=$2
            shift
        ;;
        --openapi-token)
            export OPENAPI_TOKEN=$2
            shift
        ;;
        --ess)
            export FROM_ESS=$2
            export AUTO_FDISK="true"
            shift
        ;;
        --auto-fdisk)
            export AUTO_FDISK="true"
        ;;
        --labels)
            export LABELS=$2
            shift
        ;;
        *)
            public::common::log "unknown option [$key]"
            exit 1
        ;;
    esac
    shift
    done

    if [ "$OPENAPI_TOKEN" != "" ] ;
    then
        # Using default cidr.
        public::common::get_node_info
    fi

    if [ "$DOCKER_VERSION" == "" ] ;
    then
        # Using default cidr.
        public::common::log "DOCKER_VERSION $DOCKER_VERSION is not set."
        exit 1
    fi


    KUBE_VERSION=$(curl -k https://$APISERVER_LB/version|grep gitVersion |awk '{print $2}'|cut -f2 -d \")
    export KUBE_VERSION=${KUBE_VERSION:1}
    if [ "$KUBE_VERSION" == "" ] ;
    then
        # Using default cidr.
        public::common::log "KUBE_VERSION $KUBE_VERSION is failed to set."
        exit 1
    fi

    if [ "$TOKEN" == "" ] ;
    then
        # Using default cidr.
        public::common::log "TOKEN $TOKEN is not set."
        exit 1
    fi

    if [ "$CLUSTER_DNS" == "" ] ;
    then
        # Using default cidr.
        public::common::log "CLUSTER_DNS $CLUSTER_DNS is not set."
        exit 1
    fi

    # KUBE_REPO_PREFIX
    public::common::nodeid
    public::common::os_env

	public::common::region

	# 首先从本地读取相应版本的tar包。当所需要的安装包不存在的时候
	# 如果设置了参数PKG_FILE_SERVER，就从该Server上下载。
    # 如果是在公有云上执行，可以使用内网oss地址
	if [ "$PKG_FILE_SERVER" == "" ];then
	    export PKG_FILE_SERVER=http://aliacs-k8s-$REGION.oss-$REGION.aliyuncs.com
        if [ "$region" == "cn-shanghai-finance-1" ]; then
            export PKG_FILE_SERVER=http://aliacs-k8s-$REGION.oss-$REGION-internal.aliyuncs.com
        fi
	fi

	# 安装Kubernetes时候会启动一些AddOn插件的镜像。
	# 改插件设置镜像仓库的前缀。
	if [ "$KUBE_REPO_PREFIX" == "" ];then
	    export KUBE_REPO_PREFIX=registry.$REGION.aliyuncs.com/acs
	fi

	public::main::node_up

    if [ "$OPENAPI_TOKEN" != "" ] ;
    then
        public::common::callback
    fi

}


main "$@"
