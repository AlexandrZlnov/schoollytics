### Schoollytics

**Schoollytics** - сервис позволяющий анализировать успеваемость школьника на основании полученных оценок.
Сервис взаимодействует с личным кабинетом школьника в электронном дневнике Московской Электронной Школы (МЭШ) - https://school.mos.ru/

После аутентификации сервис получает данные об успеваемости школьника в формате JSON и парсит их в базу данных.

В дальнейшем, после доработки, Schoollytics позволит получать сводные отчеты и срезы данных по следующим параметрам и не только: 
- расчет необходимого количества оценок определенного номинала для получения желаемого среднего балла
- текущий средний балл по предметам в разных временных рамках:
    - по триместрам
    - весь текущий период обучения
    - последний месяц
    - любой выбранный диапазон времени    
- текущий общий средний балл по всем предметам (основной показатель успеваемости) с возможностью выбора временных рамок, аналогично описанным выше.
- сводный отчет об успеваемость (на основе среднего балла) по предмету на разных временных этапах:
    - по неделям или месяцам с начала года
    - по полугодиям
- отчет об успеваемости по предмету за последний месяц в сравнении с предшествующим периодом (позволит оценить динамику среднего балла за прошедшие 30 дней)

Для прохождения аутентификации используются возможности по автоматизации управления браузером Chromium билбиотеки playwright-go <a href="https://github.com/playwright-community/playwright-go">playwright-go</a>

Для хранения данных используется PostgreSQL.

* * * *
#### Версии
**Текущая версия - v0.0.1**
Как работает, на текущий момент:
- с помощью [playwright-go](github.com/playwright-community/playwright-go) создается экземпляр браузера Chromium и страница авторизации - https://school.mos.ru
- средствами playwright имитируем клик кнопки войти
- ждем редиректа на страницу авторизации через портал МОС.РУ по адресу - https://login.mos.ru/sps/login/methods/**
- в соответствующие поля передаем значение login и password
- login и password вводятся в окне терминала по запросу программы (в дальнейшем будет исправлено)
- клик кнопки войти по `button id="bind"`
- ожидаем переход на страницу дневника школьника
- извлекаем куки с именем "aupd_token" в ней храниться токен сессии
- отправляем POST запрос по адресу `https://school.mos.ru/api/ej/acl/v1/sessions` в тело запроса включаем токен. Обязательно включаем в заголовки `x-mes-subsystem", "familyweb`.
В ответе получаем JSON который сохраняем в файл `student_info.json` и десериализукм в структуру `domain.StudentInfo`. Из которой получаем `StudentID`
- отправляем GET запрос по адресу `https://school.mos.ru/api/family/web/v1/subject_marks?student_id=%d` где `%d - studentId`
В ответе получаем JSON с оценкам школьника, которые сохраняем в файл `response.json`
- дальнейшие действия производятся с файлами данных школьника `response.json` и `student_info.json`, которые десериализуются в структуры `StudentPerformance`
и `Students` и передаются в базу данных.
- база данных под управлением PostgreSQL работает только на локальной машине (будет исправлено)

#### Структура
- `cmd/`                
&nbsp;&nbsp;&nbsp;&nbsp; `main.go`  
- `internal/`  
&nbsp;&nbsp;&nbsp;&nbsp; `domain/`            - структуры и конфигурационные файлы  
&nbsp;&nbsp;&nbsp;&nbsp; `repository/`        - функции работы с базой данных  
&nbsp;&nbsp;&nbsp;&nbsp; `service/`           - основная логика приложения  
- `.env`                  - переменные окружения  
- `response.json`         - json с оценками школьника  
- `student_info.json`     - json с информацие о школьнике  

#### Запуск
go run cmd/main.go 

* * * *
#### DB
Схема БД - https://dbdesigner.page.link/JsCY2VXuUtUWk9Ky6

<details>
<summary>Дамп БД, только схема</summary>
    
        -- grades
        CREATE TABLE public.grades (
            id integer NOT NULL,
            student_id integer NOT NULL,
            subject_id integer NOT NULL,
            external_id bigint,
            value numeric(5,2) NOT NULL,
            weight integer DEFAULT 1 NOT NULL,
            control_form_name character varying(100),
            date date NOT NULL,
            original_grade_system_type character varying(50),
            period_id integer,
            CONSTRAINT chk_value CHECK ((value >= (0)::numeric)),
            CONSTRAINT chk_weight CHECK ((weight > 0))
        );
        ALTER TABLE public.grades OWNER TO postgres;
        CREATE SEQUENCE public.grades_id_seq
            AS integer
            START WITH 1
            INCREMENT BY 1
            NO MINVALUE
            NO MAXVALUE
            CACHE 1;
        ALTER SEQUENCE public.grades_id_seq OWNER TO postgres;
        ALTER SEQUENCE public.grades_id_seq OWNED BY public.grades.id;
        ALTER TABLE ONLY public.grades ALTER COLUMN id SET DEFAULT nextval('public.grades_id_seq'::regclass);
        ALTER TABLE ONLY public.grades
            ADD CONSTRAINT grades_pkey PRIMARY KEY (id);
        ALTER TABLE ONLY public.grades
            ADD CONSTRAINT uniq_external_id UNIQUE (external_id);
        ALTER TABLE ONLY public.grades
            ADD CONSTRAINT fk_periods_id FOREIGN KEY (period_id) REFERENCES public.periods(id);
        ALTER TABLE ONLY public.grades
            ADD CONSTRAINT fk_student FOREIGN KEY (student_id) REFERENCES public.students(id) ON DELETE CASCADE;
        ALTER TABLE ONLY public.grades
            ADD CONSTRAINT fk_subject FOREIGN KEY (subject_id) REFERENCES public.subjects(id) ON DELETE CASCADE;

        --periods
        CREATE TABLE public.periods (
            id integer NOT NULL,
            start_date character varying(10) NOT NULL,
            end_date character varying(10) NOT NULL,
            title character varying(10) NOT NULL,
            dynamic character varying(10) NOT NULL,
            value character varying(10),
            count integer,
            target jsonb,
            fixed_value character varying(10),
            start_iso character varying(15),
            end_iso character varying(15)
        );
        ALTER TABLE public.periods OWNER TO postgres;
        CREATE SEQUENCE public.periods_id_seq
            AS integer
            START WITH 1
            INCREMENT BY 1
            NO MINVALUE
            NO MAXVALUE
            CACHE 1;
        ALTER SEQUENCE public.periods_id_seq OWNER TO postgres;
        ALTER SEQUENCE public.periods_id_seq OWNED BY public.periods.id;
        ALTER TABLE ONLY public.periods ALTER COLUMN id SET DEFAULT nextval('public.periods_id_seq'::regclass);
        ALTER TABLE ONLY public.periods
            ADD CONSTRAINT periods_pkey PRIMARY KEY (id);

        -- schools
        CREATE TABLE public.schools (
            id integer NOT NULL,
            school_id integer NOT NULL,
            name text NOT NULL,
            shortname text NOT NULL,
            organization_id text NOT NULL
        );
        ALTER TABLE public.schools OWNER TO postgres;
        CREATE SEQUENCE public.schools_id_seq
            START WITH 1
            INCREMENT BY 1
            NO MINVALUE
            NO MAXVALUE
            CACHE 1;
        ALTER SEQUENCE public.schools_id_seq OWNER TO postgres;
        ALTER TABLE public.schools ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
            SEQUENCE NAME public.schools_id_seq1
            START WITH 1
            INCREMENT BY 1
            NO MINVALUE
            NO MAXVALUE
            CACHE 1
        );
        ALTER TABLE ONLY public.schools
            ADD CONSTRAINT schools_pkey PRIMARY KEY (id);
        ALTER TABLE ONLY public.schools
            ADD CONSTRAINT unq_schools_organization_id UNIQUE (organization_id);

        -- students
        CREATE TABLE public.students (
            id integer NOT NULL,
            user_id integer NOT NULL,
            profile_id integer NOT NULL,
            guid text,
            first_name character varying(100) NOT NULL,
            last_name character varying(100) NOT NULL,
            middle_name character varying(100),
            phone_number character varying(20),
            authentication_token text,
            person_id character varying(255),
            pswrd_change_required boolean DEFAULT false,
            regional_auth character varying(50),
            date_of_birth date,
            sex character varying(10),
            school_id integer,
            CONSTRAINT students_sex_check CHECK (((sex)::text = ANY ((ARRAY['male'::character varying, 'female'::character varying])::text[])))
        );
        ALTER TABLE public.students OWNER TO postgres;
        ALTER TABLE public.students ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
            SEQUENCE NAME public.students_id_seq
            START WITH 1
            INCREMENT BY 1
            NO MINVALUE
            NO MAXVALUE
            CACHE 1
        );
        ALTER TABLE ONLY public.students
            ADD CONSTRAINT students_pkey PRIMARY KEY (id);
        ALTER TABLE ONLY public.students
            ADD CONSTRAINT unq_students_user_id UNIQUE (user_id);	
        ALTER TABLE ONLY public.students
            ADD CONSTRAINT fk_school_id FOREIGN KEY (school_id) REFERENCES public.schools(id) ON DELETE CASCADE;

        -- subjects
        CREATE TABLE public.subjects (
            id integer NOT NULL,
            name character varying(100) NOT NULL,
            external_id integer
        );
        ALTER TABLE public.subjects OWNER TO postgres;
        ALTER TABLE public.subjects ALTER COLUMN id ADD GENERATED ALWAYS AS IDENTITY (
            SEQUENCE NAME public.subjects_id_seq
            START WITH 1
            INCREMENT BY 1
            NO MINVALUE
            NO MAXVALUE
            CACHE 1
        );
        ALTER TABLE ONLY public.subjects
            ADD CONSTRAINT subjects_pkey PRIMARY KEY (id);
        ALTER TABLE ONLY public.subjects
            ADD CONSTRAINT unk_subjects_external_id UNIQUE (external_id);
    </details>

* * * *
#### Планы на доработку
- исключить дублирование структур хранящих данные о школьнике StudentInfo - Students 
- добавить проверку наличия файла базы данных 
- добавить создание базы данных и таблицами
- добавить индексы в таблицах БД
- добавить проверку наличия оценок у школьника (сейчас они будут дублироваться в БД)
- добавить загрузку данных для подключения к БД из ENV
- добавить валидацию входного токена
- убрать дублирование токена (оставить только в заголовке)
- рассмотреть альтернативные методы получения токена (например, через API)
- добавить возможность headless-режима при аутентификации
- расширить функционал в части аналитики данных школьника



