CREATE SCHEMA IF NOT EXISTS app;

CREATE TABLE app.c_tools (
tool_id BIGINT PRIMARY KEY,
tool_id_father BIGINT,
name VARCHAR(50) NOT NULL,
description VARCHAR(200),
order_by INTEGER DEFAULT 1 NOT NULL,
created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
created_by VARCHAR(50),
modified_on TIMESTAMP,
modified_by VARCHAR(50),
configuration TEXT NOT NULL,
path VARCHAR(50) NOT NULL,
is_active BOOLEAN DEFAULT TRUE,


CONSTRAINT fk_c_tools_father
    FOREIGN KEY (tool_id_father)
    REFERENCES app.c_tools(tool_id)

);
CREATE INDEX idx_c_tools_father ON app.c_tools(tool_id_father);
CREATE INDEX idx_c_tools_active ON app.c_tools(is_active);

CREATE TABLE app.c_profiles (
profile_id BIGINT PRIMARY KEY,
name VARCHAR(50) NOT NULL,
description VARCHAR(200),
order_by INTEGER DEFAULT 1 NOT NULL,
created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
created_by VARCHAR(50),
modified_on TIMESTAMP,
modified_by VARCHAR(50),
is_active BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_c_profiles_active ON app.c_profiles(is_active);


CREATE TABLE app.c_groups (
group_id BIGINT PRIMARY KEY,
name VARCHAR(50) NOT NULL,
description VARCHAR(200),
order_by INTEGER DEFAULT 1 NOT NULL,
created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
created_by VARCHAR(50),
modified_on TIMESTAMP,
modified_by VARCHAR(50),
is_active BOOLEAN DEFAULT TRUE
);

CREATE INDEX idx_c_groups_active ON app.c_groups(is_active);

CREATE TABLE app.r_groups_profiles (
group_id BIGINT NOT NULL,
profile_id BIGINT NOT NULL,
created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
created_by VARCHAR(50),
modified_on TIMESTAMP,
modified_by VARCHAR(50),
is_active BOOLEAN DEFAULT TRUE,


PRIMARY KEY (group_id, profile_id),

CONSTRAINT fk_grp_prf_group
    FOREIGN KEY (group_id)
    REFERENCES app.c_groups(group_id),

CONSTRAINT fk_grp_prf_profile
    FOREIGN KEY (profile_id)
    REFERENCES app.c_profiles(profile_id)


);
CREATE INDEX idx_r_gp_group ON app.r_groups_profiles(group_id);
CREATE INDEX idx_r_gp_profile ON app.r_groups_profiles(profile_id);

CREATE TABLE app.r_profiles_tools (
profile_id BIGINT NOT NULL,
tool_id BIGINT NOT NULL,
created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
created_by VARCHAR(50),
modified_on TIMESTAMP,
modified_by VARCHAR(50),
operation VARCHAR(250) DEFAULT '*' NOT NULL,
is_active BOOLEAN DEFAULT TRUE,


PRIMARY KEY (profile_id, tool_id),

CONSTRAINT fk_prf_tls_profile
    FOREIGN KEY (profile_id)
    REFERENCES app.c_profiles(profile_id),

CONSTRAINT fk_prf_tls_tool
    FOREIGN KEY (tool_id)
    REFERENCES app.c_tools(tool_id)


);

CREATE INDEX idx_r_pt_profile ON app.r_profiles_tools(profile_id);
CREATE INDEX idx_r_pt_tool ON app.r_profiles_tools(tool_id);

CREATE TABLE app.r_users_profiles (
user_id BIGINT NOT NULL,
profile_id BIGINT NOT NULL,
created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
created_by VARCHAR(50),
modified_on TIMESTAMP,
modified_by VARCHAR(50),
is_active BOOLEAN DEFAULT TRUE,


PRIMARY KEY (user_id, profile_id)

);

CREATE INDEX idx_r_up_user ON app.r_users_profiles(user_id);
CREATE INDEX idx_r_up_profile ON app.r_users_profiles(profile_id);


CREATE TABLE t_users (
    id              BIGSERIAL PRIMARY KEY,
    name          VARCHAR(100) NOT NULL,
    last_name      VARCHAR(100) NOT NULL,
    second_last_name      VARCHAR(100),
    email           VARCHAR(255) NOT NULL UNIQUE,
    hash_password   TEXT NOT NULL,
    created_on      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    modified_on     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active       BOOLEAN DEFAULT TRUE
);

CREATE USER auth_service WITH PASSWORD 'chag3_m3!!';
GRANT CONNECT ON DATABASE "app-db" TO auth_service;
GRANT USAGE ON SCHEMA app TO auth_service;
GRANT SELECT, INSERT, UPDATE, DELETE 
ON ALL TABLES IN SCHEMA app 
TO auth_service;
GRANT USAGE, SELECT 
ON ALL SEQUENCES IN SCHEMA app 
TO auth_service;
ALTER DEFAULT PRIVILEGES IN SCHEMA app
GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO auth_service;
ALTER DEFAULT PRIVILEGES IN SCHEMA app
GRANT USAGE, SELECT ON SEQUENCES TO auth_service;



CREATE TABLE t_refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token_hash TEXT NOT NULL,
    device TEXT,
    ip_address INET,
    user_agent TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
 	session_started_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP not null,
    revoked BOOLEAN DEFAULT FALSE,

    CONSTRAINT fk_refresh_user
        FOREIGN KEY (user_id)
        REFERENCES t_users(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_refresh_user_id 
ON t_refresh_tokens(user_id);

CREATE INDEX idx_refresh_active 
ON t_refresh_tokens(user_id, revoked, expires_at);