-- ============================================================
-- Relief Atlas  --  PL/SQL procedures, functions, triggers, views
-- ============================================================
CONNECT disaster_user/DisasterApp123!@XEPDB1

-- ------------------------------------------------------------
-- TRIGGER: prevent shelter occupancy from exceeding capacity
-- ------------------------------------------------------------
CREATE OR REPLACE TRIGGER trg_shelter_capacity
    BEFORE INSERT OR UPDATE OF occupied_number, capacity ON Shelter
    FOR EACH ROW
BEGIN
    IF :NEW.occupied_number > :NEW.capacity THEN
        RAISE_APPLICATION_ERROR(
            -20001,
            'Shelter capacity exceeded. Capacity: ' || :NEW.capacity ||
            ', Attempted occupancy: ' || :NEW.occupied_number
        );
    END IF;
END trg_shelter_capacity;
/

-- ------------------------------------------------------------
-- TRIGGER: set distribution_date to SYSDATE if not provided
-- ------------------------------------------------------------
CREATE OR REPLACE TRIGGER trg_distribution_date
    BEFORE INSERT ON Distribution
    FOR EACH ROW
    WHEN (NEW.distribution_date IS NULL)
BEGIN
    :NEW.distribution_date := SYSDATE;
END trg_distribution_date;
/

-- ------------------------------------------------------------
-- FUNCTION: total donations across all records
-- ------------------------------------------------------------
CREATE OR REPLACE FUNCTION fn_total_donations
RETURN NUMBER
AS
    v_total NUMBER;
BEGIN
    SELECT NVL(SUM(amount), 0) INTO v_total FROM Donations;
    RETURN v_total;
END fn_total_donations;
/

-- ------------------------------------------------------------
-- FUNCTION: total donations for a specific user
-- ------------------------------------------------------------
CREATE OR REPLACE FUNCTION fn_user_donations(p_user_id IN NUMBER)
RETURN NUMBER
AS
    v_total NUMBER;
BEGIN
    SELECT NVL(SUM(amount), 0) INTO v_total
    FROM Donations
    WHERE user_id = p_user_id;
    RETURN v_total;
END fn_user_donations;
/

-- ------------------------------------------------------------
-- FUNCTION: available spots in a shelter
-- ------------------------------------------------------------
CREATE OR REPLACE FUNCTION fn_shelter_available(p_shelter_id IN NUMBER)
RETURN NUMBER
AS
    v_cap     NUMBER;
    v_occ     NUMBER;
BEGIN
    SELECT capacity, occupied_number
    INTO v_cap, v_occ
    FROM Shelter
    WHERE shelter_id = p_shelter_id;
    RETURN v_cap - v_occ;
EXCEPTION
    WHEN NO_DATA_FOUND THEN
        RAISE_APPLICATION_ERROR(-20002, 'Shelter not found: ' || p_shelter_id);
END fn_shelter_available;
/

-- ------------------------------------------------------------
-- PROCEDURE: register a disaster + its first affected area
-- ------------------------------------------------------------
CREATE OR REPLACE PROCEDURE sp_register_disaster(
    p_disaster_name  IN  VARCHAR2,
    p_disaster_type  IN  VARCHAR2,
    p_start_date     IN  DATE,
    p_status         IN  VARCHAR2,
    p_area_name      IN  VARCHAR2,
    p_severity       IN  VARCHAR2,
    p_population     IN  NUMBER,
    p_disaster_id    OUT NUMBER,
    p_area_id        OUT NUMBER
)
AS
BEGIN
    INSERT INTO Disaster (disaster_name, disaster_type, start_date, status)
    VALUES (p_disaster_name, p_disaster_type, p_start_date, p_status)
    RETURNING disaster_id INTO p_disaster_id;

    INSERT INTO Affected_Areas (area_name, severity, population, disaster_id)
    VALUES (p_area_name, p_severity, p_population, p_disaster_id)
    RETURNING area_id INTO p_area_id;

    COMMIT;
EXCEPTION
    WHEN OTHERS THEN
        ROLLBACK;
        RAISE;
END sp_register_disaster;
/

-- ------------------------------------------------------------
-- PROCEDURE: distribute aid to an affected area
-- ------------------------------------------------------------
CREATE OR REPLACE PROCEDURE sp_distribute_aid(
    p_material_name   IN  VARCHAR2,
    p_quantity        IN  NUMBER,
    p_area_id         IN  NUMBER,
    p_admin_id        IN  NUMBER,
    p_dist_id         OUT NUMBER
)
AS
    v_count NUMBER;
BEGIN
    SELECT COUNT(*) INTO v_count
    FROM Affected_Areas WHERE area_id = p_area_id;

    IF v_count = 0 THEN
        RAISE_APPLICATION_ERROR(-20003, 'Affected area not found: ' || p_area_id);
    END IF;

    INSERT INTO Distribution (material_name, quantity, distribution_date, area_id, admin_id)
    VALUES (p_material_name, p_quantity, SYSDATE, p_area_id, p_admin_id)
    RETURNING distribution_id INTO p_dist_id;

    COMMIT;
EXCEPTION
    WHEN OTHERS THEN
        ROLLBACK;
        RAISE;
END sp_distribute_aid;
/

-- ------------------------------------------------------------
-- PROCEDURE: update shelter occupancy safely
-- ------------------------------------------------------------
CREATE OR REPLACE PROCEDURE sp_update_occupancy(
    p_shelter_id    IN NUMBER,
    p_new_occupancy IN NUMBER
)
AS
    v_cap NUMBER;
BEGIN
    SELECT capacity INTO v_cap
    FROM Shelter WHERE shelter_id = p_shelter_id
    FOR UPDATE;

    IF p_new_occupancy > v_cap THEN
        RAISE_APPLICATION_ERROR(
            -20001,
            'Cannot set occupancy ' || p_new_occupancy ||
            ': exceeds capacity ' || v_cap
        );
    END IF;

    UPDATE Shelter
    SET occupied_number = p_new_occupancy
    WHERE shelter_id = p_shelter_id;

    COMMIT;
EXCEPTION
    WHEN NO_DATA_FOUND THEN
        ROLLBACK;
        RAISE_APPLICATION_ERROR(-20002, 'Shelter not found: ' || p_shelter_id);
    WHEN OTHERS THEN
        ROLLBACK;
        RAISE;
END sp_update_occupancy;
/

-- ------------------------------------------------------------
-- VIEW: disaster summary with area counts and population
-- ------------------------------------------------------------
CREATE OR REPLACE VIEW v_disaster_summary AS
SELECT
    d.disaster_id,
    d.disaster_name,
    d.disaster_type,
    d.start_date,
    d.status,
    COUNT(a.area_id)             AS affected_area_count,
    NVL(SUM(a.population), 0)   AS total_population_affected
FROM Disaster d
LEFT JOIN Affected_Areas a ON d.disaster_id = a.disaster_id
GROUP BY d.disaster_id, d.disaster_name, d.disaster_type, d.start_date, d.status;

-- ------------------------------------------------------------
-- VIEW: shelter occupancy status with area context
-- ------------------------------------------------------------
CREATE OR REPLACE VIEW v_shelter_status AS
SELECT
    s.shelter_id,
    s.shelter_name,
    s.capacity,
    s.occupied_number,
    s.capacity - s.occupied_number                              AS available_spots,
    ROUND((s.occupied_number / NULLIF(s.capacity, 0)) * 100, 1) AS occupancy_pct,
    s.location,
    s.contact_number,
    a.area_name,
    a.severity,
    a.disaster_id
FROM Shelter s
JOIN Affected_Areas a ON s.area_id = a.area_id;

-- ------------------------------------------------------------
-- VIEW: distribution log with area and admin names
-- ------------------------------------------------------------
CREATE OR REPLACE VIEW v_distribution_log AS
SELECT
    dist.distribution_id,
    dist.material_name,
    dist.quantity,
    dist.distribution_date,
    a.area_name,
    a.severity,
    adm.admin_name,
    adm.email AS admin_email
FROM Distribution dist
JOIN Affected_Areas a   ON dist.area_id  = a.area_id
LEFT JOIN Admins    adm ON dist.admin_id = adm.admin_id;
