CREATE DATABASE IF NOT EXISTS db_ipresmgr;

USE db_ipresmgr;

CREATE TABLE IF NOT EXISTS tbl_K8SResourceIPBind (
    k8sresource_id VARCHAR(128) NOT NULL,          -- k8sclusterid-namespace-resource_name
    k8sresource_type VARCHAR(32) NOT NULL,         -- 资源类型名，Deployment和StatefulSet
    ip VARCHAR(16) NOT NULL,                    -- 分配的ip
    is_bind TINYINT NOT NULL,                   -- ip是否绑定，0：没有绑定，1：绑定
    bind_podid VARCHAR(36),                     -- 绑定的podid，解绑后StatefuSet这个podid不能清除
    bind_time TIMESTAMP,                        -- 绑定的时间
    create_time TIMESTAMP,                      -- ip从nsp分配的时间    
    PRIMARY KEY(k8sresource_id, ip),
    INDEX(bind_podid),
    INDEX(k8sresource_type)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS tbl_K8SResourceIPRecycle (
    k8sresource_id VARCHAR(128) NOT NULL,          -- k8sclusterid-namespace-resource_name
    k8sresource_release_time TIMESTAMP NOT NULL,   -- k8s资源释放时间
    ipresource_release_time TIMESTAMP NOT NULL,    -- ip资源归还给nsp的时间，租期到期时间
    netregional_id VARCHAR(128),                   -- 释放用到的网络域id
    subnet_id VARCHAR(36),                         -- 释放用到的子网id
    ips BLOB,                                      -- 释放的ip列表，ip1:ip2:ip3
    PRIMARY KEY(k8sresource_id),
    INDEX(netregional_id),
    INDEX(subnet_id)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS tbl_K8SResourceIPRecycleHistroy (
    k8sresource_id VARCHAR(128) NOT NULL,          -- k8sclusterid-namespace-resource_name
    ip_release_time TIMESTAMP NOT NULL,            -- ip归还给nsp的时间，租期到期时间
    netregional_id VARCHAR(128),                   -- 释放用到的网络域id
    subnet_id VARCHAR(36),                         -- 释放用到的子网id
    ips BLOB,                                      -- 释放的ip列表，ip1:ip2:ip3   
    PRIMARY KEY(k8sresource_id),
    INDEX(netregional_id),
    INDEX(subnet_id)    
) ENGINE=InnoDB DEFAULT CHARSET=latin1;