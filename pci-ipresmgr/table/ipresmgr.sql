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
    k8sresource_type int NOT NULL,                 -- 资源类型，Deployment和StatefulSet proto.K8SApiResourceKindType
    ip VARCHAR(16) NOT NULL,                       -- 分配的ip
    mac VARCHAR(16) NOT NULL,                      -- mac地址
    netregional_id VARCHAR(128) NOT NULL,          -- 用到的网络域id  
    subnet_id VARCHAR(36) NOT NULL,                -- 用到的子网id  
    port_id VARCHAR(48) NOT NULL,                  -- PortID
    subnetgatewayaddr VARCHAR(16) NOT NULL,        -- 子网网关地址
    alloc_time TIMESTAMP NOT NULL,                 -- ip从nsp分配的时间   
    is_bind TINYINT NOT NULL,                      -- ip是否绑定，0：没有绑定，1：绑定
    bind_podid VARCHAR(36) NULL,                   -- 绑定的podid，解绑后StatefuSet这个podid不能清除
    bind_time TIMESTAMP NULL DEFAULT '0000-00-00 00:00:00',                      -- 绑定的时间 
    PRIMARY KEY(k8sresource_id, port_id),
    INDEX(bind_podid),
    INDEX(k8sresource_type),
    INDEX(is_bind)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS tbl_K8SResourceIPRecycle (
    srv_instance_name VARCHAR(32) NOT NULL,        -- 服务实例名字，这个资源由该服务实例负责回收
    k8sresource_id VARCHAR(128) NOT NULL,          -- k8sclusterid-namespace-resource_name
    k8sresource_type int NOT NULL,                 -- 资源类型，Deployment和StatefulSet proto.K8SApiResourceKindType
    replicas INT NOT NULL,                         -- pod数量
    create_time TIMESTAMP NOT NULL,                -- 释放资源插入时间
    nspresource_release_time TIMESTAMP NOT NULL,   -- ip资源归还给nsp的时间，租期到期时间
    recycle_object_id VARCHAR(64) NOT NULL,        -- 回收对象id，每次都不同
    -- unbind_count INT NOT NULL DEFAULT 0,            取消绑定的数量    
    -- netregional_id VARCHAR(128) NOT NULL,           释放用到的网络域id
    -- subnet_id VARCHAR(36) NOT NULL,                 释放用到的子网id
    -- port_id VARCHAR(48) NOT NULL,                   PortID，释放只需要这个参数  
    -- subnetgatewayaddr VARCHAR(16) NOT NULL,         子网网关地址    
    -- nsp_resources BLOB,                             释放的ip列表，[{ip,portid},....]  
    PRIMARY KEY(k8sresource_id),
    INDEX(srv_instance_name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS tbl_K8SResourceIPRecycleHistroy (
    id INT UNSIGNED AUTO_INCREMENT, 
    k8sresource_id VARCHAR(128) NOT NULL,          -- k8sclusterid-namespace-resource_name
    k8sresource_type int NOT NULL,                 -- 资源类型，Deployment和StatefulSet proto.K8SApiResourceKindType
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

CREATE TABLE IF NOT EXISTS tbl_K8SJobNetInfo (
    k8sresource_id VARCHAR(128) NOT NULL,          -- k8sclusterid-namespace-resource_name
    k8sresource_type int NOT NULL,                 -- 资源类型，Job CronJob proto.K8SApiResourceKindType
    netregional_id VARCHAR(128) NOT NULL,          -- 用到的网络域id  
    subnet_id VARCHAR(36) NOT NULL,                -- 用到的子网id
    subnetgatewayaddr VARCHAR(16) NOT NULL,        -- 子网网关地址   
    PRIMARY KEY(k8sresource_id)         
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS tbl_K8SJobIPBind (
    k8sresource_id VARCHAR(128) NOT NULL,          -- k8sclusterid-namespace-resource_name    
    ip VARCHAR(16) NOT NULL,                       -- 分配的ip
    bind_podid VARCHAR(36) NULL,                   -- 绑定的podid，解绑后StatefuSet这个podid不能清除  
    port_id VARCHAR(48) NOT NULL,                  -- PortID  
    PRIMARY KEY(k8sresource_id, port_id),
    INDEX(bind_podid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS tbl_K8SScaleDownMark (
    recycle_mark_id VARCHAR(64) NOT NULL,          -- 回收标记id，每次都不同
    k8sresource_id VARCHAR(128) NOT NULL,          -- k8sclusterid-namespace-resource_name  
    PRIMARY KEY(recycle_mark_id, k8sresource_id)  
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS tbl_Test (
    id INT UNSIGNED AUTO_INCREMENT, 
    k8sresource_id VARCHAR(128) NOT NULL,   -- k8sclusterid-namespace-resource_name
    nspresource_release_time TIMESTAMP NOT NULL,   -- ip归还给nsp的时间，租期到期时间
    subnet_id VARCHAR(36),                         -- 释放用到的子网id  
    create_time TIMESTAMP NULL DEFAULT '1970-01-02 00:00:00',                -- 插入时间
    use_flag INT NOT NULL, -- 测试悲观锁 0: 没有使用，1：使用
    nsp_resources BLOB,                            -- 释放的ip列表，{ip,mac}   
    PRIMARY KEY(id),
    UNIQUE KEY k8sresid (k8sresource_id),
    INDEX(subnet_id)    
) ENGINE=InnoDB DEFAULT CHARSET=utf8;