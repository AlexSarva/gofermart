package storagepg

const ddl = `
CREATE TABLE if not exists public.users (
    id uuid primary key ,
    username text unique ,
    passwd text,
    cookie text,
    cookie_expires timestamp,
    created timestamptz default now()
);
CREATE TABLE if not exists public.orders (
     user_id uuid references public.users(id),
     order_num text primary key,
     accrual float4,
     status text,
     created timestamptz default now()
);
delete from public.orders where user_id in (select user_id from public.users where username like 'test%');
delete from public.users where username like 'test%';

create table if not exists public.withdraw (
   user_id uuid references public.users(id),
   order_num text primary key,
   withdraw float4,
   created timestamptz default now()
);
delete from public.withdraw where user_id in (select user_id from public.users where username like 'test%');

create or replace view public.balance as (
 with total as (
     select user_id, sum(accrual) total from public.orders group by user_id
 ),
      withdraw as (
          select user_id, sum(withdraw) withdraw from public.withdraw group by user_id
      )
 select users.id user_id,
        coalesce(total.total,0)::numeric(10,2) total,
        coalesce(withdraw.withdraw,0)::numeric(10,2) withdraw,
        (coalesce(total.total,0) - coalesce(withdraw.withdraw,0))::numeric(10,2) current
 from public.users
          left join total on total.user_id = users.id
          left join withdraw on withdraw.user_id = users.id);
`
