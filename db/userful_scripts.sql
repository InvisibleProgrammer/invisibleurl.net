-- get audit log events

select al.recorded_at, u.email_address, lt.name, al.description
from audit_log al 
  inner join lt_audit_log_events lt on lt.event_id = al.event_id
  inner join users u on u.user_id = al.user_id;
