USE db_ipresmgr;

ALTER TABLE tbl_IPResMgrSrvRegister ADD srv_pid int NOT NULL AFTER srv_addr;
ALTER TABLE tbl_IPResMgrSrvRegister DROP PRIMARY KEY, ADD PRIMARY KEY(srv_instance_name, srv_pid);