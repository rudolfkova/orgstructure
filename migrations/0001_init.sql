-- +goose Up

CREATE TABLE departments (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    parent_id BIGINT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT departments_parent_fk
        FOREIGN KEY (parent_id)
        REFERENCES departments(id)
        ON DELETE CASCADE
);

CREATE UNIQUE INDEX departments_root_name_unique
ON departments (name)
WHERE parent_id IS NULL;

CREATE UNIQUE INDEX departments_parent_name_unique
ON departments (parent_id, name)
WHERE parent_id IS NOT NULL;



CREATE TABLE employees (
    id BIGSERIAL PRIMARY KEY,
    department_id BIGINT NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    position VARCHAR(200) NOT NULL,
    hired_at DATE NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT employees_department_fk
        FOREIGN KEY (department_id)
        REFERENCES departments(id)
        ON DELETE CASCADE
);

CREATE INDEX employees_department_idx
ON employees (department_id);

CREATE INDEX employees_created_at_idx
ON employees (created_at);



-- +goose Down

DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS departments;