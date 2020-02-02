USE db_ipresmgr;

ALTER TABLE tbl_K8SScaleDownMark ADD current_replicas INT NOT NULL AFTER k8sresource_type;