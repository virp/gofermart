create table if not exists users
(
    id            uuid primary key,
    login         text           not null unique,
    password_hash text           not null,
    balance       decimal(10, 2) not null default 0,
    withdrawn     decimal(10, 2) not null default 0
);

create table if not exists orders
(
    id          bigint primary key,
    user_id     uuid        not null,
    status      text        not null,
    uploaded_at timestamptz not null,
    accrual     decimal(10, 2),
    foreign key (user_id) references users (id)
);

create table if not exists withdrawals
(
    id           bigint primary key,
    user_id      uuid           not null,
    sum          decimal(10, 2) not null,
    processed_at timestamptz    not null,
    foreign key (user_id) references users (id)
);