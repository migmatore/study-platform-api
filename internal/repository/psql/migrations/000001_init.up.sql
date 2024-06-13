CREATE TABLE institutions
(
    id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name        VARCHAR(200) NOT NULL UNIQUE,
    description VARCHAR(1000)
);

CREATE TABLE roles
(
    id   INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE users
(
    id             INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    full_name      VARCHAR(100) NOT NULL,
    phone          VARCHAR(20),
    email          VARCHAR(50)  NOT NULL UNIQUE,
    password_hash  VARCHAR(100) NOT NULL,
    role_id        INT          NOT NULL REFERENCES roles (id),
    institution_id INT REFERENCES institutions (id) ON DELETE CASCADE
);

CREATE TABLE classrooms
(
    id           INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title        VARCHAR(100) NOT NULL,
    description  VARCHAR(1000),
    teacher_id   INT          NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    max_students INT          NOT NULL DEFAULT 1
);

CREATE TABLE classroom_students
(
    id           INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    classroom_id INT NOT NULL REFERENCES classrooms(id) ON DELETE CASCADE,
    student_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE lessons
(
    id           INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    title        VARCHAR(100) NOT NULL,
    classroom_id INT          NOT NULL REFERENCES classrooms (id) ON DELETE CASCADE,
    content      JSONB,
    active       BOOLEAN DEFAULT FALSE
);

INSERT INTO roles(name)
VALUES ('admin'),
       ('student'),
       ('teacher');
