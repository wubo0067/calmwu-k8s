CREATE DATABASE IF NOT EXISTS db_ipresmgr;

USE db_ipresmgr;

CREATE TABLE IF NOT EXISTS tbl_IPResMgrSrvRegister (
    srv_instance_name VARCHAR(32) NOT NULL,          -- 服务实例名字
    srv_addr VARCHAR(32) NOT NULL,                   -- 服务监听的地址x.x.x.x:port
    register_time TIMESTAMP NOT NULL,                -- 服务注册时间
    PRIMARY KEY(srv_instance_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS tbl_K8SResourceIPBind (
    k8sresource_id VARCHAR(128) NOT NULL,          -- k8sclusterid-namespace-resource_name
    k8sresource_type VARCHAR(32) NOT NULL,         -- 资源类型名，Deployment和StatefulSet
    ip VARCHAR(16) NOT NULL,                       -- 分配的ip
    mac VARCHAR(16) NOT NULL,                      -- mac地址
    netregional_id VARCHAR(128) NOT NULL,          -- 用到的网络域id  
    subnet_id VARCHAR(36) NOT NULL,                -- 用到的子网id  
    port_id VARCHAR(48) NOT NULL,                  -- PortID
    subnetgatewayaddr VARCHAR(16) NOT NULL,        -- 子网网关地址
    alloc_time TIMESTAMP NOT NULL DEFAULT '0000-00-00 00:00:00',                -- ip从nsp分配的时间   
    is_bind TINYINT NOT NULL,                      -- ip是否绑定，0：没有绑定，1：绑定
    bind_podid VARCHAR(36) NULL,                   -- 绑定的podid，解绑后StatefuSet这个podid不能清除
    bind_time TIMESTAMP NULL DEFAULT '0000-00-00 00:00:00',                      -- 绑定的时间 
    PRIMARY KEY(k8sresource_id, ip),
    INDEX(bind_podid),
    INDEX(k8sresource_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS tbl_K8SResourceIPRecycle (
    srv_instance_name VARCHAR(32) NOT NULL,        -- 服务实例名字，这个资源由该服务实例负责回收
    k8sresource_id VARCHAR(128) NOT NULL,          -- k8sclusterid-namespace-resource_name
    replicas INT NOT NULL,                         -- pod数量
    unbind_count INT NOT NULL DEFAULT 0,           -- 取消绑定的数量
    create_time TIMESTAMP NOT NULL,                -- 释放资源插入时间
    nspresource_release_time TIMESTAMP NOT NULL,   -- ip资源归还给nsp的时间，租期到期时间
    netregional_id VARCHAR(128) NOT NULL,          -- 释放用到的网络域id
    subnet_id VARCHAR(36) NOT NULL,                -- 释放用到的子网id
    port_id VARCHAR(48) NOT NULL,                  -- PortID，释放只需要这个参数  
    subnetgatewayaddr VARCHAR(16) NOT NULL,        -- 子网网关地址    
    nsp_resources BLOB,                            -- 释放的ip列表，{ip,mac}  
    PRIMARY KEY(k8sresource_id),
    INDEX(port_id),
    INDEX(srv_instance_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS tbl_K8SResourceIPRecycleHistroy (
    id INT UNSIGNED AUTO_INCREMENT, 
    k8sresource_id VARCHAR(128) NOT NULL,          -- k8sclusterid-namespace-resource_name
    replicas INT NOT NULL,                         -- pod数量    
    nspresource_release_time TIMESTAMP NOT NULL,   -- ip归还给nsp的时间，租期到期时间
    netregional_id VARCHAR(128) NOT NULL,          -- 释放用到的网络域id
    subnet_id VARCHAR(36) NOT NULL,                -- 释放用到的子网id
    port_id VARCHAR(48) NOT NULL,                  -- PortID，释放只需要这个参数    
    create_time TIMESTAMP NOT NULL,                -- 插入时间
    nsp_resources BLOB,                            -- 释放的ip列表，{ip,mac}   
    PRIMARY KEY(id),
    INDEX(k8sresource_id),
    INDEX(port_id),
    INDEX(subnet_id)    
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS tbl_Test (
    id INT UNSIGNED AUTO_INCREMENT, 
    k8sresource_id VARCHAR(128) NOT NULL,          -- k8sclusterid-namespace-resource_name
    nspresource_release_time TIMESTAMP NOT NULL,   -- ip归还给nsp的时间，租期到期时间
    subnet_id VARCHAR(36),                         -- 释放用到的子网id  
    create_time TIMESTAMP NULL DEFAULT '0000-00-00 00:00:00',                -- 插入时间
    nsp_resources BLOB,                            -- 释放的ip列表，{ip,mac}   
    PRIMARY KEY(id),
    INDEX(k8sresource_id),
    INDEX(subnet_id)    
) ENGINE=InnoDB DEFAULT CHARSET=utf8;