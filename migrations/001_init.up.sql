CREATE TABLE IF NOT EXISTS user_logins (
    id          BIGSERIAL PRIMARY KEY,
    user_id     UUID NOT NULL,
    login_time  TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    UNIQUE (user_id, login_time)
);

CREATE TABLE IF NOT EXISTS daily_unique_users (
    day     DATE NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY (day, user_id)
);
CREATE INDEX IF NOT EXISTS idx_daily_day ON daily_unique_users (day);

CREATE TABLE IF NOT EXISTS monthly_unique_users (
    month   DATE NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY (month, user_id)
    );
CREATE INDEX IF NOT EXISTS idx_monthly_month ON monthly_unique_users (month);