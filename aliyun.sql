create table admins
(
    id       int unsigned auto_increment
        primary key,
    account  varchar(20) not null,
    password varchar(20) not null
)
    charset = utf8;

create table book_type
(
    id   int unsigned auto_increment
        primary key,
    type varchar(20) not null
)
    charset = utf8;

create table books
(
    isbn        varchar(20)                        not null
        primary key,
    name        varchar(20)                        not null,
    type_id     int unsigned                       not null,
    author      varchar(20)                        not null,
    public      varchar(30)                        not null,
    public_date varchar(10)                        not null,
    register    datetime default CURRENT_TIMESTAMP null,
    stock       int unsigned                       not null,
    constraint fk_book_type_id
        foreign key (type_id) references book_type (id)
)
    charset = utf8;

create table users
(
    id       int unsigned auto_increment
        primary key,
    name     varchar(20)                           not null,
    account  varchar(20)                           not null,
    password varchar(20) default '123456'          null,
    sex      varchar(1)                            not null,
    college  varchar(30)                           not null,
    birthday date                                  null,
    register datetime    default CURRENT_TIMESTAMP null
)
    charset = utf8;

create table borrow
(
    id          int unsigned auto_increment
        primary key,
    uid         int unsigned                           not null,
    bid         varchar(20)                            not null,
    days        int unsigned default 30                null,
    borrow_date datetime     default CURRENT_TIMESTAMP null,
    back_date   datetime                               null,
    constraint fk_bid
        foreign key (bid) references books (isbn),
    constraint fk_uid
        foreign key (uid) references users (id)
)
    charset = utf8;

create definer = lxy@`%` trigger borrow_insert_before_trigger
    before insert
    on borrow
    for each row
begin
    declare _stock int;
    declare msg varchar(20);
    select stock into _stock from books where isbn = NEW.bid;
    if _stock = 0 then
        set msg = '库存不足';
        signal sqlstate 'HY000' set message_text = msg;
    else
        update `books` set stock=stock - 1 where isbn = new.bid;
    end if;
end;

create definer = lxy@`%` trigger borrow_update_after_trigger
    after update
    on borrow
    for each row
begin
    update `books` set stock=stock + 1 where isbn = NEW.bid;
end;

