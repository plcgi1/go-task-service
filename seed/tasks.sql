INSERT INTO task(status)
SELECT
        'NEW'
FROM generate_series(1, 1000) AS s(i);