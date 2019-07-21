drop table posts;
drop table threads;
drop table sessions;
drop table users;
drop table demands;


create table users (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  first_name varchar(255),
  last_name  varchar(255),
  username   varchar(255),
  role       varchar(255),
  email      varchar(255) not null unique,
  password   varchar(255) not null,
  created_at timestamp not null,
  last_sign_in_at timestamp
);
/*,
confirmed_at timestamp,
invited_at timestamp,
aud        varchar(255),
confirmation_token varchar(255),
confirmation_sent_at timestamp,

recovery_token varchar(255),
recovery_sent_at timestamp,

email_change_token varchar(255),
email_change varchar(255),
email_change_sent_at timestamp,


raw_app_meta_data text,
raw_user_meta_data text,

is_super_admin boolean,

updated_at timestamp
*/


create table providers (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  user_id    integer references users(id),
  mobile_phone  varchar(255),
  company_name  varchar(255),
  company_addr  varchar(255),
  company_city  varchar(255),
  company_zip   varchar(64),
  company_country varchar(255),

	equipment     text,
  eligible_items text,
	operating_countries text
);
/*,
licence bytea*/


CREATE TABLE images (
  imgname text,
  img bytea
);


create table sessions (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  email      varchar(255),
  user_id    integer references users(id),
  created_at timestamp not null
);

create table threads (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  topic      text,
  user_id    integer references users(id),
  created_at timestamp not null
);

create table posts (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  body       text,
  user_id    integer references users(id),
  thread_id  integer references threads(id),
  created_at timestamp not null
);

create table demands (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  object     text,
  collection text,
  delivery   text,
  timeframe  text,
  user_id    integer references users(id),
  created_at timestamp not null,
  status     integer
);

create table bids (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  starting_bid numeric CHECK (starting_bid > 0),
  currency text,
  provider_id    integer references providers(id),
  demand_id  integer references demands(id),
  pick_up_time text,
  drop_off_time text,
  quote_expiration text,
  note text,
  payment text,
  created_at timestamp not null,
  status     integer
);

create table messages (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  body       text,
  user_id    integer references users(id),
  demand_id  integer references demands(id),
  created_at timestamp not null
);
