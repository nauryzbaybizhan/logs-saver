-- URL
create
or replace function insert_visits_url_if_not_exist(input_url text,
                                                    api_key_id int,
                                                    visits_account_id integer,
                                                    input_api_key text)
    returns int
    language plpgsql
as
$$
declare
ua_id integer;
begin
select id
from visits_url
where url = input_url into ua_id;
if
ua_id IS NULL then
        insert into visits_url (url, visits_api_keys_id) values (input_url, api_key_id) returning id into ua_id;
end if;
return ua_id;
end;
$$;
-- USER_AGENT
create
or replace function insert_visits_ua_if_not_exist(input_ua text, api_key_id int)
    returns int
    language plpgsql
as
$$
declare
ua_id integer;
begin
select id
from visits_ua
where ua = input_ua into ua_id;
if
ua_id IS NULL then
        insert into visits_ua (ua, api_key) values (input_ua, api_key_id) returning id into ua_id;
end if;
return ua_id;

end;
$$;

-- IP
create
or replace function insert_visits_ip_if_not_exist(input_address text,
                                                         input_bot boolean,
                                                         input_data_center boolean,
                                                         input_tor boolean,
                                                         input_proxy boolean,
                                                         input_vpn boolean,
                                                         input_country text,
                                                         input_domain_list text[],
                                                         api_key_id integer,
                                                         visits_account_id integer,
                                                         input_api_key text
)
    returns int
    language plpgsql
as
$$
declare
ip_id integer;
begin
select id
from visits_ip
where address = input_address::inet
into ip_id;
if
ip_id IS NULL then
        insert into visits_ip (address, bot, data_center, tor, proxy, vpn, country, domain_list, visits_api_keys_id, checked, visits_account_id, api_key)
        values (input_address::inet, input_bot, input_data_center, input_tor, input_proxy, input_vpn, input_country,
                input_domain_list, api_key_id, true, visits_account_id, input_api_key)
        returning id into ip_id;
end if;
return ip_id;


end;
$$;

-- API_KEY

create
or replace function insert_visits_api_keys_if_not_exist(input_api_key text)
    returns int
    language plpgsql
as
$$
declare
api_key_id integer;
begin
select id
from visits_api_keys
where api_key = input_api_key into api_key_id;
if
api_key_id IS NULL then
        insert into visits_api_keys (api_key, checked) values (input_api_key, true) returning id into api_key_id;
end if;
return api_key_id;


end;
$$;

-- ACCOUNT
create
or replace function insert_visits_account_if_not_exist(input_user_id text, api_key_id int, input_api_key text)
    returns int
    language plpgsql
as
$$
declare
acc_id integer;
begin
select id
from visits_accounts
where user_id = input_user_id into acc_id;
if
acc_id IS NULL then
        insert into visits_accounts (user_id, visits_api_keys_id, checked, api_key) values (input_user_id, api_key_id, true, input_api_key) returning id into acc_id;
end if;
return acc_id;

end;
$$;

-- DEVICE

create
or replace function insert_visits_device_if_not_exist(input_user text,
                                                             input_ua text,
                                                             input_type smallint,
                                                             input_time timestamp,
                                                             api_key_id int,
                                                             input_api_key text) returns integer
    language plpgsql
as
$$
declare
device_id  integer;
    account_id
integer;
    ua_id
integer;

begin

select *
from insert_visits_ua_if_not_exist(input_ua) into ua_id;
select *
from insert_visits_account_if_not_exist(input_user, api_key_id) into account_id;

select id
from visits_devices
where ua = ua_id
  and api_key = api_key_id
  and type = input_type into device_id;
if
device_id IS NULL then
        insert into visits_devices(account_id, type, visits_api_keys_id, ua, created, checked, api_key)
        values (account_id, input_type, api_key_id, ua_id, input_time, true, input_api_key)
        returning id into device_id;
end if;
return device_id;

end;
$$;