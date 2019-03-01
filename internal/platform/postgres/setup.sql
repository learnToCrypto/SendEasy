drop table posts;
drop table threads;
drop table sessions;
drop table users;
drop table demands;


create table users (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  name       varchar(255),
  first_name varchar(255),
  last_name  varchar(255),
  aud        varchar(255),
  role       varchar(255),

  email      varchar(255) not null unique,
  password   varchar(255) not null,

  confirmed_at timestamp,
  invited_at timestamp,

  confirmation_token varchar(255),
  confirmation_sent_at timestamp,

  recovery_token varchar(255),
  recovery_sent_at timestamp,

  email_change_token varchar(255),
  email_change varchar(255),
  email_change_sent_at timestamp,
  last_sign_in_at timestamp,

  raw_app_meta_data text,
  raw_user_meta_data text,

  is_super_admin boolean,

  created_at timestamp not null,
  updated_at timestamp
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

create table messages (
  id         serial primary key,
  uuid       varchar(64) not null unique,
  body       text,
  user_id    integer references users(id),
  demand_id  integer references demands(id),
  created_at timestamp not null
);
