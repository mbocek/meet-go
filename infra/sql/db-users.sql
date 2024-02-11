CREATE USER meet_users WITH PASSWORD 'meet_users';
ALTER USER meet_users CREATEDB;
CREATE DATABASE meet_users WITH OWNER = meet_users ENCODING = 'UTF8';