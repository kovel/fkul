DROP TABLE IF EXISTS fkul_world_championship_bombardier;

CREATE TABLE fkul_world_championship_bombardier
(
    footballer VARCHAR(255) NOT NULL,
    goals INT NOT NULL DEFAULT 0,
    year INT NOT NULL
);