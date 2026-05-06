-- ============================================================
-- Relief Atlas  --  Sample / seed data
CONNECT disaster_user/DisasterApp123!@XEPDB1
-- NOTE: Admin and User records must be created via the API
--       (POST /api/admins, POST /api/users) so passwords are
--       properly bcrypt-hashed by the application.
-- ============================================================

-- Disasters
INSERT INTO Disaster (disaster_name, disaster_type, start_date, status)
VALUES ('Hurricane Delta', 'Hurricane', TO_DATE('2024-09-15', 'YYYY-MM-DD'), 'Active');

INSERT INTO Disaster (disaster_name, disaster_type, start_date, status)
VALUES ('Wildfire Omega', 'Wildfire', TO_DATE('2024-08-01', 'YYYY-MM-DD'), 'Contained');

INSERT INTO Disaster (disaster_name, disaster_type, start_date, status)
VALUES ('Flood Surge Gamma', 'Flood', TO_DATE('2024-10-03', 'YYYY-MM-DD'), 'Active');

COMMIT;

-- Affected Areas
INSERT INTO Affected_Areas (area_name, severity, population, disaster_id)
VALUES ('Coastal Region A', 'Critical',
        75000,
        (SELECT disaster_id FROM Disaster WHERE disaster_name = 'Hurricane Delta'));

INSERT INTO Affected_Areas (area_name, severity, population, disaster_id)
VALUES ('Downtown Zone B', 'High',
        42000,
        (SELECT disaster_id FROM Disaster WHERE disaster_name = 'Hurricane Delta'));

INSERT INTO Affected_Areas (area_name, severity, population, disaster_id)
VALUES ('North Hills District', 'Moderate',
        18000,
        (SELECT disaster_id FROM Disaster WHERE disaster_name = 'Wildfire Omega'));

INSERT INTO Affected_Areas (area_name, severity, population, disaster_id)
VALUES ('River Valley East', 'Critical',
        61000,
        (SELECT disaster_id FROM Disaster WHERE disaster_name = 'Flood Surge Gamma'));

COMMIT;

-- Shelters
INSERT INTO Shelter (shelter_name, capacity, location, occupied_number, contact_number, area_id)
VALUES ('City Hall Emergency Shelter', 500, '1 Municipal Plaza', 320,
        '+1-555-0101',
        (SELECT area_id FROM Affected_Areas WHERE area_name = 'Coastal Region A'));

INSERT INTO Shelter (shelter_name, capacity, location, occupied_number, contact_number, area_id)
VALUES ('Central Sports Arena', 1200, '88 Stadium Road', 750,
        '+1-555-0202',
        (SELECT area_id FROM Affected_Areas WHERE area_name = 'Downtown Zone B'));

INSERT INTO Shelter (shelter_name, capacity, location, occupied_number, contact_number, area_id)
VALUES ('Community Recreation Centre', 300, '45 Park Avenue', 110,
        '+1-555-0303',
        (SELECT area_id FROM Affected_Areas WHERE area_name = 'North Hills District'));

INSERT INTO Shelter (shelter_name, capacity, location, occupied_number, contact_number, area_id)
VALUES ('Valley Primary School', 450, '12 School Lane', 390,
        '+1-555-0404',
        (SELECT area_id FROM Affected_Areas WHERE area_name = 'River Valley East'));

COMMIT;

-- Distributions (admin_id NULL — assign via API after creating an admin)
INSERT INTO Distribution (material_name, quantity, distribution_date, area_id, admin_id)
VALUES ('Rice Bags (50 kg)', 200,
        TO_DATE('2024-09-18', 'YYYY-MM-DD'),
        (SELECT area_id FROM Affected_Areas WHERE area_name = 'Coastal Region A'),
        NULL);

INSERT INTO Distribution (material_name, quantity, distribution_date, area_id, admin_id)
VALUES ('Drinking Water (5 L)', 1500,
        TO_DATE('2024-09-19', 'YYYY-MM-DD'),
        (SELECT area_id FROM Affected_Areas WHERE area_name = 'Downtown Zone B'),
        NULL);

INSERT INTO Distribution (material_name, quantity, distribution_date, area_id, admin_id)
VALUES ('First Aid Kits', 80,
        TO_DATE('2024-08-05', 'YYYY-MM-DD'),
        (SELECT area_id FROM Affected_Areas WHERE area_name = 'North Hills District'),
        NULL);

INSERT INTO Distribution (material_name, quantity, distribution_date, area_id, admin_id)
VALUES ('Blankets', 600,
        TO_DATE('2024-10-05', 'YYYY-MM-DD'),
        (SELECT area_id FROM Affected_Areas WHERE area_name = 'River Valley East'),
        NULL);

COMMIT;

-- Donations (user_id NULL — assign via API after creating users)
INSERT INTO Donations (amount, donation_date, user_id)
VALUES (5000.00, TO_DATE('2024-09-16', 'YYYY-MM-DD'), NULL);

INSERT INTO Donations (amount, donation_date, user_id)
VALUES (1200.50, TO_DATE('2024-09-20', 'YYYY-MM-DD'), NULL);

INSERT INTO Donations (amount, donation_date, user_id)
VALUES (750.00, TO_DATE('2024-08-03', 'YYYY-MM-DD'), NULL);

INSERT INTO Donations (amount, donation_date, user_id)
VALUES (3000.00, TO_DATE('2024-10-04', 'YYYY-MM-DD'), NULL);

COMMIT;
