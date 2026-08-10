CREATE TYPE user_role AS ENUM ('teacher', 'librarian', 'head_teacher');
CREATE TYPE material_category AS ENUM ('awards', 'photos', 'videos', 'library', 'documents');
CREATE TYPE material_status AS ENUM ('pending', 'approved', 'rejected');

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    full_name TEXT NOT NULL,
    role user_role NOT NULL DEFAULT 'teacher',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE materials (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    category material_category NOT NULL,
    status material_status NOT NULL DEFAULT 'approved',
    file_key TEXT,
    file_name TEXT,
    file_mime TEXT,
    file_size BIGINT,
    content TEXT,
    content_format TEXT,
    condition TEXT NOT NULL DEFAULT '',
    location TEXT NOT NULL DEFAULT '',
    origin_date DATE,
    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_materials_category ON materials(category);
CREATE INDEX idx_materials_status ON materials(status);
CREATE INDEX idx_materials_created_at ON materials(created_at);
CREATE INDEX idx_materials_origin_date ON materials(origin_date);

CREATE TABLE material_permissions (
    id BIGSERIAL PRIMARY KEY,
    material_id BIGINT NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    granted_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (material_id, user_id)
);

CREATE TABLE loans (
    id BIGSERIAL PRIMARY KEY,
    material_id BIGINT NOT NULL REFERENCES materials(id) ON DELETE CASCADE,
    borrower_name TEXT NOT NULL,
    note TEXT NOT NULL DEFAULT '',
    taken_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    due_date DATE NOT NULL,
    returned_at TIMESTAMPTZ
);

CREATE INDEX idx_loans_material ON loans(material_id);
CREATE INDEX idx_loans_returned ON loans(returned_at);
