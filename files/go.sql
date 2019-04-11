select t.id, title, content, created_date, priority, 
case when c.name is null then 'NA' else c.name end from
    task t, status s, user u left outer join  category c
    on c.id=t.cat_id where u.username=? and s.id=t.task_status_id 
    and u.id=t.user_id  and s.status='PENDING' and t.hide!=1 order by t.created_date asc
