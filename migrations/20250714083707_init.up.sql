CREATE TABLE task (
      id SERIAL PRIMARY KEY,
      status VARCHAR(20) DEFAULT 'NEW',
      error_message VARCHAR(255),
      count_of_tryings INT,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON COLUMN "public"."task"."status" IS 'Статус задачи';
COMMENT ON COLUMN "public"."task"."error_message" IS 'Сообщение об ошибке';
COMMENT ON COLUMN "public"."task"."count_of_tryings" IS 'Количество запусков одной и той же задачи';
