CREATE DATABASE disaster_management;
CREATE TABLE Disaster (
    disaster_id INT AUTO_INCREMENT,
    disaster_name VARCHAR(100),
    disaster_type VARCHAR(50),
    start_date DATE,
    status VARCHAR(50),
    
    CONSTRAINT pk_disaster PRIMARY KEY (disaster_id)
);
CREATE TABLE Affected_Areas (
    area_id INT AUTO_INCREMENT,
    area_name VARCHAR(100),
    severity VARCHAR(50),
    population INT,
    disaster_id INT,

    CONSTRAINT pk_area PRIMARY KEY (area_id),
    CONSTRAINT fk_area_disaster FOREIGN KEY (disaster_id)
        REFERENCES Disaster(disaster_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
CREATE TABLE Shelter (
    shelter_id INT AUTO_INCREMENT,
    shelter_name VARCHAR(100),
    capacity INT,
    location VARCHAR(100),
    occupied_number INT,
    contact_number VARCHAR(15),
    area_id INT,

    CONSTRAINT pk_shelter PRIMARY KEY (shelter_id),
    CONSTRAINT fk_shelter_area FOREIGN KEY (area_id)
        REFERENCES Affected_Areas(area_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
CREATE TABLE Admins (
    admin_id INT AUTO_INCREMENT,
    admin_name VARCHAR(100),
    email VARCHAR(100),
    password VARCHAR(100),

    CONSTRAINT pk_admin PRIMARY KEY (admin_id),
    CONSTRAINT uq_admin_email UNIQUE (email)
);
CREATE TABLE Distribution (
    distribution_id INT AUTO_INCREMENT,
    material_name VARCHAR(100),
    quantity INT,
    distribution_date DATE,
    area_id INT,
    admin_id INT,

    CONSTRAINT pk_distribution PRIMARY KEY (distribution_id),

    CONSTRAINT fk_distribution_area FOREIGN KEY (area_id)
        REFERENCES Affected_Areas(area_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    CONSTRAINT fk_distribution_admin FOREIGN KEY (admin_id)
        REFERENCES Admins(admin_id)
        ON DELETE SET NULL
        ON UPDATE CASCADE
);
CREATE TABLE Users (
    user_id INT AUTO_INCREMENT,
    user_name VARCHAR(100),
    user_email VARCHAR(100),
    user_phoneno VARCHAR(15),
    password VARCHAR(100),

    CONSTRAINT pk_users PRIMARY KEY (user_id),
    CONSTRAINT uq_user_email UNIQUE (user_email)
);
CREATE TABLE Donations (
    donation_id INT AUTO_INCREMENT,
    amount DECIMAL(10,2),
    donation_date DATE,
    user_id INT,

    CONSTRAINT pk_donation PRIMARY KEY (donation_id),

    CONSTRAINT fk_donation_user FOREIGN KEY (user_id)
        REFERENCES Users(user_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
