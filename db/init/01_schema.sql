-- ============================================================
-- Relief Atlas  --  Oracle XE 21c schema
-- ============================================================
-- Switch from SYSDBA to the application user so all objects
-- are owned by disaster_user, not SYS.
CONNECT disaster_user/DisasterApp123!@XEPDB1

-- Disaster events
CREATE TABLE Disaster (
    disaster_id     NUMBER GENERATED ALWAYS AS IDENTITY,
    disaster_name   VARCHAR2(100)  NOT NULL,
    disaster_type   VARCHAR2(50),
    start_date      DATE,
    status          VARCHAR2(50)   DEFAULT 'Active',
    CONSTRAINT pk_disaster PRIMARY KEY (disaster_id)
);

-- Areas affected by a disaster
CREATE TABLE Affected_Areas (
    area_id         NUMBER GENERATED ALWAYS AS IDENTITY,
    area_name       VARCHAR2(100)  NOT NULL,
    severity        VARCHAR2(50),
    population      NUMBER(12),
    disaster_id     NUMBER,
    CONSTRAINT pk_area PRIMARY KEY (area_id),
    CONSTRAINT fk_area_disaster FOREIGN KEY (disaster_id)
        REFERENCES Disaster(disaster_id) ON DELETE CASCADE
);

-- Shelters inside affected areas
CREATE TABLE Shelter (
    shelter_id      NUMBER GENERATED ALWAYS AS IDENTITY,
    shelter_name    VARCHAR2(100)  NOT NULL,
    capacity        NUMBER(10)     DEFAULT 0,
    location        VARCHAR2(200),
    occupied_number NUMBER(10)     DEFAULT 0,
    contact_number  VARCHAR2(20),
    area_id         NUMBER,
    CONSTRAINT pk_shelter        PRIMARY KEY (shelter_id),
    CONSTRAINT fk_shelter_area   FOREIGN KEY (area_id)
        REFERENCES Affected_Areas(area_id) ON DELETE CASCADE,
    CONSTRAINT chk_shelter_cap   CHECK (occupied_number >= 0 AND capacity >= 0)
);

-- Administrator accounts
CREATE TABLE Admins (
    admin_id        NUMBER GENERATED ALWAYS AS IDENTITY,
    admin_name      VARCHAR2(100)  NOT NULL,
    email           VARCHAR2(200)  NOT NULL,
    password        VARCHAR2(255)  NOT NULL,
    CONSTRAINT pk_admin       PRIMARY KEY (admin_id),
    CONSTRAINT uq_admin_email UNIQUE (email)
);

-- Aid distribution records
CREATE TABLE Distribution (
    distribution_id   NUMBER GENERATED ALWAYS AS IDENTITY,
    material_name     VARCHAR2(100) NOT NULL,
    quantity          NUMBER(12)    NOT NULL,
    distribution_date DATE          DEFAULT SYSDATE,
    area_id           NUMBER,
    admin_id          NUMBER,
    CONSTRAINT pk_distribution  PRIMARY KEY (distribution_id),
    CONSTRAINT fk_dist_area     FOREIGN KEY (area_id)
        REFERENCES Affected_Areas(area_id) ON DELETE CASCADE,
    CONSTRAINT fk_dist_admin    FOREIGN KEY (admin_id)
        REFERENCES Admins(admin_id) ON DELETE SET NULL,
    CONSTRAINT chk_dist_qty     CHECK (quantity > 0)
);

-- Registered donors / public users
CREATE TABLE Users (
    user_id         NUMBER GENERATED ALWAYS AS IDENTITY,
    user_name       VARCHAR2(100)  NOT NULL,
    user_email      VARCHAR2(200)  NOT NULL,
    user_phoneno    VARCHAR2(20),
    password        VARCHAR2(255)  NOT NULL,
    CONSTRAINT pk_users       PRIMARY KEY (user_id),
    CONSTRAINT uq_user_email  UNIQUE (user_email)
);

-- Donation records
CREATE TABLE Donations (
    donation_id     NUMBER GENERATED ALWAYS AS IDENTITY,
    amount          NUMBER(12,2)   NOT NULL,
    donation_date   DATE           DEFAULT SYSDATE,
    user_id         NUMBER,
    CONSTRAINT pk_donation      PRIMARY KEY (donation_id),
    CONSTRAINT fk_donation_user FOREIGN KEY (user_id)
        REFERENCES Users(user_id) ON DELETE SET NULL,
    CONSTRAINT chk_donation_amt CHECK (amount > 0)
);
