#!/usr/bin/env bash
./cockroach sql --insecure --host="cockroach" --execute=" CREATE DATABASE IF NOT EXISTS app_database; \
                                        CREATE USER IF NOT EXISTS  app_user; \
                                        GRANT ALL ON DATABASE app_database TO app_user; \
                                        use app_database; \
                                        CREATE TABLE IF NOT EXISTS  users (
                                        id UUID NOT NULL DEFAULT gen_random_uuid(),
                                        name STRING NULL,
                                        lastname STRING NULL,
                                        email STRING NOT NULL,
                                        password STRING NOT NULL,
                                        userrole string NOT NULL,
                                        createdon TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
                                        updatedon TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT '1970-01-01',
                                        lastlogin TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT '1970-01-01',
                                        CONSTRAINT \"primary\" PRIMARY KEY (id ASC)); \
                                        CREATE TABLE IF NOT EXISTS  movies (
                                        id UUID NOT NULL DEFAULT gen_random_uuid(),
                                        title STRING NOT NULL,
                                        year STRING NOT NULL,
                                        createdby STRING NOT NULL,
                                        createdon TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now(),
                                        updatedby STRING NOT NULL DEFAULT '-',
                                        updatedon TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT '1970-01-01',
                                        CONSTRAINT \"primary\" PRIMARY KEY (id ASC)); \
                                        CREATE TABLE IF NOT EXISTS movie_directors (
                                        movie_id UUID NOT NULL,
                                        director STRING NOT NULL,
                                        createdby STRING NOT NULL,
                                        createdon TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT now());"