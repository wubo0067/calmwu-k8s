start transaction

select * from `tbl_K8SResourceIPBind` where k8sresource_id='cluster-1:default:kata-nginx-deployment' AND k8sresource_type=0 AND is_bind=0 LIMIT 1 FOR UPDATE;

update tbl_K8SResourceIPBind set is_bind=1, bind_poduniquename='calmwu', bind_time='2019-10-17 14:17:41' where port_id='port-1SJtdl2m6tzi4uaukh1AjRB2bJx' and is_bind=0 and k8sresource_id='cluster-1:default:kata-nginx-deployment' and k8sresource_type=0;

update tbl_K8SResourceIPBind set is_bind=0, bind_poduniquename='calmwu', bind_time='2019-10-17 14:17:41' where port_id='port-1SJtdl2m6tzi4uaukh1AjRB2bJx' and is_bind=1 and k8sresource_id='cluster-1:default:kata-nginx-deployment' and k8sresource_type=0;