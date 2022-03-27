-- API KEYS

create table visits_api_keys
(
    id       serial
        constraint visits_api_keys_pk
            primary key,
    api_key  text                  not null,
    checked  boolean default false,
    "quote " integer default 50000 not null
);

comment on table visits_api_keys is 'List of api keys with quotas';

-- IP

create table visits_ip
(
    id                 serial
        constraint visits_ip_pk
            primary key,
    address            inet    not null,
    bot                boolean not null,
    data_center        boolean not null,
    tor                boolean not null,
    proxy              boolean not null,
    vpn                boolean not null,
    country            text    not null,
    domain_list        text[]  default '{}'::text[],
    visits_api_keys_id integer
        constraint visits_ip_visits_api_keys_id_fk
            references visits_api_keys,
    checked            boolean default false,
    visits_account_id  integer
        constraint visits_ip_visits_account_id_fk
            references visits_accounts,
    api_key            text
);

comment on table visits_ip is 'Contains ip information';



create unique index visits_ip_address_uindex
    on visits_ip (address);

create unique index visits_ip_id_uindex
    on visits_ip (id);

-- URL

create table visits_url
(
    id                 serial
        constraint visits_url_pk
            primary key,
    url                text    not null,
    visits_api_keys_id integer not null
        constraint visits_url_visits_api_keys_id_fk
            references visits_api_keys,
    checked            boolean default false,
    visits_account_id  integer
        constraint visits_url_visits_account_id_fk
            references visits_accounts,
    api_key            text
);

comment on table visits_url is 'Visited urls';



create unique index visits_url_url_uindex
    on visits_url (url);

create unique index visits_url_id_uindex
    on visits_url (id);

-- USER_AGENT

create table visits_ua
(
    id                 serial
        constraint visits_ua_pk
            primary key,
    ua                 text    not null,
    visits_api_keys_id integer not null
        constraint visits_ua_visits_api_keys_id_fk
            references visits_api_keys,
    checked            boolean default false,
    visits_account_id  integer
        constraint visits_ua_visits_account_id_fk
            references visits_accounts,
    api_key            text
);

comment on table visits_ua is 'Contains user agents';


create unique index visits_ua_ua_uindex
    on visits_ua (ua);

create unique index visits_ua_id_uindex
    on visits_ua (id);


create unique index visits_api_keys_id_uindex
    on visits_api_keys (id);

create unique index visits_api_keys_key_uindex
    on visits_api_keys (api_key);

-- ACCOUNTS

create table visits_accounts
(
    id                 serial
        constraint visits_accounts_pk
            primary key,
    user_id            text                                not null,
    ips                integer[] default '{}'::integer[]   not null,
    countries          integer[] default '{}'::integer[]   not null,
    total_visits       integer   default 0                 not null,
    visits_api_keys_id integer                             not null
        constraint visits_accounts_visits_api_keys_id_fk
            references visits_api_keys,
    last_ip            inet,
    created            timestamp default CURRENT_TIMESTAMP not null,
    last_updated       timestamp default CURRENT_TIMESTAMP not null,
    checked            boolean   default false,
    api_key            text
);


create unique index visits_accounts_id_uindex
    on visits_accounts (id);

create unique index visits_accounts_user_id_uindex
    on visits_accounts (user_id);

-- VISITS
create table visits
(
    id                 bigserial
        constraint visits_pk
            primary key,
    visits_api_keys_id integer   not null
        constraint visits_visits_api_keys_id_fk
            references visits_api_keys,
    account            integer   not null
        constraint visits_account_fk
            references visits_accounts,
    ip                 integer   not null
        constraint visits_ip_fk
            references visits_ip,
    url                integer   not null
        constraint visits_url_fk
            references visits_url,
    device             bigint    not null
        constraint visits_ua_fk
            references visits_ua,
    time               timestamp not null,
    checked            boolean default false,
    api_key            text
);


-- DEVICE

create table visits_devices
(
    id                 bigserial
        constraint visits_devices_pk
            primary key,
    account_id         integer                                                not null
        constraint visits_devices_account__fk
            references visits_accounts,
    type               smallint  default 1                                    not null,
    visits_api_keys_id integer                                                not null
        constraint visits_devices_visits_api_keys_id_fk
            references visits_api_keys,
    ua                 integer                                                not null
        constraint visits_devices_ua_fk
            references visits_ua,
    created            timestamp default (now())::timestamp without time zone not null,
    checked            boolean   default false,
    api_key            text
);



create unique index visits_devices_id_uindex
    on visits_devices (id);

create unique index visits_devices_uniq_index
    on visits_devices (account_id, visits_api_keys_id, ua);
