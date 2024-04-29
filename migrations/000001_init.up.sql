-- Создание таблицы пользователей
CREATE TABLE ee_users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash CHAR(64) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы курсов
CREATE TABLE ee_courses (
    course_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    meta TEXT,
    instructor_id INTEGER REFERENCES ee_users(user_id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы владельцев таблиц
CREATE TABLE ee_course_owners (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES ee_users(user_id),
    course_id INTEGER REFERENCES ee_courses(course_id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы разделов курса
CREATE TABLE ee_course_sections (
    section_id SERIAL PRIMARY KEY,
    course_id INTEGER REFERENCES ee_courses(course_id),
    title VARCHAR(255) NOT NULL,
    "order" INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы уроков
CREATE TABLE ee_lessons (
    lesson_id SERIAL PRIMARY KEY,
    section_id INTEGER REFERENCES ee_course_sections(section_id),
    title VARCHAR(255) NOT NULL,
    content_text TEXT,
    video_id INTEGER, -- Связь с видео будет установлена отдельно
    "order" INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы видео
CREATE TABLE ee_videos (
    video_id SERIAL PRIMARY KEY,
    lesson_id INTEGER UNIQUE REFERENCES ee_lessons(lesson_id),
    url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы изображений
CREATE TABLE ee_images (
    image_id SERIAL PRIMARY KEY,
    lesson_id INTEGER REFERENCES ee_lessons(lesson_id),
    url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы тестов
CREATE TABLE ee_quizzes (
    quiz_id SERIAL PRIMARY KEY,
    lesson_id INTEGER REFERENCES ee_lessons(lesson_id),
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы вопросов к тестам
CREATE TABLE ee_quiz_questions (
    question_id SERIAL PRIMARY KEY,
    quiz_id INTEGER REFERENCES ee_quizzes(quiz_id),
    question_text TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы ответов на вопросы
CREATE TABLE ee_quiz_answers (
    answer_id SERIAL PRIMARY KEY,
    question_id INTEGER REFERENCES ee_quiz_questions(question_id),
    user_id INTEGER REFERENCES ee_users(user_id),
    answer_text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(question_id, user_id)
);

-- Комментарии для таблицы Users
COMMENT ON TABLE ee_users IS 'Хранит информацию о пользователях платформы';
COMMENT ON COLUMN ee_users.user_id IS 'Уникальный идентификатор пользователя';
COMMENT ON COLUMN ee_users.username IS 'Имя пользователя, используемое для входа в систему';
COMMENT ON COLUMN ee_users.email IS 'Электронная почта пользователя';
COMMENT ON COLUMN ee_users.password_hash IS 'Хэш пароля пользователя';
COMMENT ON COLUMN ee_users.created_at IS 'Дата и время создания записи пользователя';
COMMENT ON COLUMN ee_users.updated_at IS 'Дата и время последнего обновления записи пользователя';

-- Комментарии для таблицы Courses
COMMENT ON TABLE ee_courses IS 'Содержит информацию о курсах, доступных на платформе';
COMMENT ON COLUMN ee_courses.course_id IS 'Уникальный идентификатор курса';
COMMENT ON COLUMN ee_courses.title IS 'Название курса';
COMMENT ON COLUMN ee_courses.description IS 'Описание курса';
COMMENT ON COLUMN ee_courses.meta IS 'Мета информация курса';
COMMENT ON COLUMN ee_courses.instructor_id IS 'Идентификатор пользователя, создавшего курс';
COMMENT ON COLUMN ee_courses.created_at IS 'Дата и время создания курса';
COMMENT ON COLUMN ee_courses.updated_at IS 'Дата и время последнего обновления курса';

-- Комментарии для таблицы CourseSections
COMMENT ON TABLE ee_course_sections IS 'Описывает разделы в курсах для структурирования материала';
COMMENT ON COLUMN ee_course_sections.section_id IS 'Уникальный идентификатор раздела курса';
COMMENT ON COLUMN ee_course_sections.course_id IS 'Идентификатор курса, к которому относится раздел';
COMMENT ON COLUMN ee_course_sections.title IS 'Название раздела курса';
COMMENT ON COLUMN ee_course_sections.order IS 'Порядковый номер раздела в курсе для упорядочивания';
COMMENT ON COLUMN ee_course_sections.created_at IS 'Дата и время создания раздела';
COMMENT ON COLUMN ee_course_sections.updated_at IS 'Дата и время последнего обновления раздела';

-- Комментарии для таблицы Lessons
COMMENT ON TABLE ee_lessons IS 'Хранит информацию об уроках, включая видео, тексты и задания';
COMMENT ON COLUMN ee_lessons.lesson_id IS 'Уникальный идентификатор урока';
COMMENT ON COLUMN ee_lessons.section_id IS 'Идентификатор раздела, к которому относится урок';
COMMENT ON COLUMN ee_lessons.title IS 'Название урока';
COMMENT ON COLUMN ee_lessons.content_text IS 'Текстовое содержимое урока, если применимо';
COMMENT ON COLUMN ee_lessons.order IS 'Порядковый номер урока в разделе для упорядочивания';
COMMENT ON COLUMN ee_lessons.created_at IS 'Дата и время создания урока';
COMMENT ON COLUMN ee_lessons.updated_at IS 'Дата и время последнего обновления урока';

-- Комментарии для таблицы Videos
COMMENT ON TABLE ee_videos IS 'Содержит информацию о видео, прикрепленных к урокам';
COMMENT ON COLUMN ee_videos.video_id IS 'Уникальный идентификатор видео';
COMMENT ON COLUMN ee_videos.lesson_id IS 'Идентификатор урока, к которому прикреплено видео';
COMMENT ON COLUMN ee_videos.url IS 'URL видеофайла';
COMMENT ON COLUMN ee_videos.created_at IS 'Дата и время загрузки видео';
COMMENT ON COLUMN ee_videos.updated_at IS 'Дата и время последнего обновления информации о видео';

-- Комментарии для таблицы Images
COMMENT ON TABLE ee_images IS 'Содержит информацию об изображениях, используемых в уроках';
COMMENT ON COLUMN ee_images.image_id IS 'Уникальный идентификатор изображения';
COMMENT ON COLUMN ee_images.lesson_id IS 'Идентификатор урока, к которому относится изображение';
COMMENT ON COLUMN ee_images.url IS 'URL изображения';
COMMENT ON COLUMN ee_images.created_at IS 'Дата и время загрузки изображения';
COMMENT ON COLUMN ee_images.updated_at IS 'Дата и время последнего обновления информации об изображении';

-- Комментарии для таблицы Quizzes
COMMENT ON TABLE ee_quizzes IS 'Описывает тесты или квизы, связанные с уроками';
COMMENT ON COLUMN ee_quizzes.quiz_id IS 'Уникальный идентификатор теста';
COMMENT ON COLUMN ee_quizzes.lesson_id IS 'Идентификатор урока, к которому относится тест';
COMMENT ON COLUMN ee_quizzes.title IS 'Название теста';
COMMENT ON COLUMN ee_quizzes.created_at IS 'Дата и время создания теста';
COMMENT ON COLUMN ee_quizzes.updated_at IS 'Дата и время последнего обновления теста';

-- Комментарии для таблицы QuizQuestions
COMMENT ON TABLE ee_quiz_questions IS 'Содержит вопросы для тестов';
COMMENT ON COLUMN ee_quiz_questions.question_id IS 'Уникальный идентификатор вопроса';
COMMENT ON COLUMN ee_quiz_questions.quiz_id IS 'Идентификатор теста, к которому относится вопрос';
COMMENT ON COLUMN ee_quiz_questions.question_text IS 'Текст вопроса';
COMMENT ON COLUMN ee_quiz_questions.correct_answer IS 'Правильный ответ на вопрос';
COMMENT ON COLUMN ee_quiz_questions.created_at IS 'Дата и время создания вопроса';
COMMENT ON COLUMN ee_quiz_questions.updated_at IS 'Дата и время последнего обновления вопроса';

-- Комментарии для таблицы QuizAnswers
COMMENT ON TABLE ee_quiz_answers IS 'Хранит ответы пользователей на вопросы тестов';
COMMENT ON COLUMN ee_quiz_answers.answer_id IS 'Уникальный идентификатор ответа';
COMMENT ON COLUMN ee_quiz_answers.question_id IS 'Идентификатор вопроса, на который дан ответ';
COMMENT ON COLUMN ee_quiz_answers.user_id IS 'Идентификатор пользователя, оставившего ответ';
COMMENT ON COLUMN ee_quiz_answers.answer_text IS 'Текст ответа пользователя';
COMMENT ON COLUMN ee_quiz_answers.is_correct IS 'Показывает, является ли ответ правильным';
COMMENT ON COLUMN ee_quiz_answers.created_at IS 'Дата и время создания ответа';
COMMENT ON COLUMN ee_quiz_answers.updated_at IS 'Дата и время последнего обновления ответа';
