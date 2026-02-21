-- +goose Up

CREATE TABLE departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL CHECK (trim(name) <> ''),
    parent_id INTEGER NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_department_parent
        FOREIGN KEY (parent_id)
        REFERENCES departments(id)
        ON UPDATE CASCADE
        ON DELETE SET NULL
);

CREATE INDEX idx_departments_parent_id ON departments(parent_id);



CREATE TABLE employees (
    id SERIAL PRIMARY KEY,
    department_id INTEGER NOT NULL,
    full_name VARCHAR(200) NOT NULL CHECK (trim(full_name) <> ''),
    position VARCHAR(200) NOT NULL CHECK (trim(position) <> ''),
    hired_at DATE NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_employee_department
        FOREIGN KEY (department_id)
        REFERENCES departments(id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT
);

CREATE INDEX idx_employees_department_id ON employees(department_id);



-- +goose Down

DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS departments;