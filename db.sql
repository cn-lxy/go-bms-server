--! @brief      Deprecated.
-- orm 如何批量导入数据

create table admins (
    id integer primary key autoincrement,
    name varchar(255) not null,
    password varchar(255) not null
);

create table users (
    id integer primary key autoincrement,
    name varchar(255) not null,
    password varchar(255) not null
);
