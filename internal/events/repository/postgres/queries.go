package events_postgres_repository

const PreInsertQueries = `
alter table visits_accounts drop constraint visits_accounts_api_key_fk;

alter table visits
    drop constraint visits_api_key__fk,
    drop constraint visits_account_fk,
    drop constraint visits_ip_fk,
    drop constraint visits_url_fk,
    drop constraint visits_ua_fk;`

const PostInsertQueries = `

	
	alter table visits_accounts
	    add constraint visits_accounts_api_key_fk foreign key (api_key) references visits_api_keys;
	
	
	alter table visits
	    add constraint visits_api_key__fk foreign key (api_key) references visits_api_keys,
	    add constraint visits_account_fk foreign key (account) references visits_accounts,
	    add constraint visits_ip_fk foreign key (ip) references visits_ip,
	    add constraint visits_url_fk foreign key (url) references visits_url,
	    add constraint visits_ua_fk foreign key (ua) references visits_ua;
	
	
	analyze visits_accounts;
	analyze visits_api_keys;
	analyze visits_ip;
	analyze visits_ua;
	analyze visits_url;
	analyze visits;`
