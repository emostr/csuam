CREATE TABLE users (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    full_name TEXT NOT NULL,
    role ENUM('teacher', 'librarian', 'head_teacher') NOT NULL DEFAULT 'teacher',
    created_at DATETIME NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE materials (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    category ENUM('awards', 'photos', 'videos', 'library', 'documents') NOT NULL,
    status ENUM('pending', 'approved', 'rejected') NOT NULL DEFAULT 'approved',
    file_key TEXT,
    file_name TEXT,
    file_mime TEXT,
    file_size BIGINT,
    content MEDIUMTEXT,
    content_format TEXT,
    `condition` TEXT NOT NULL DEFAULT '',
    location TEXT NOT NULL DEFAULT '',
    origin_date DATE,
    created_by BIGINT,
    created_at DATETIME NOT NULL DEFAULT current_timestamp(),
    updated_at DATETIME NOT NULL DEFAULT current_timestamp(),
    CONSTRAINT fk_materials_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE INDEX idx_materials_category ON materials(category);
CREATE INDEX idx_materials_status ON materials(status);
CREATE INDEX idx_materials_created_at ON materials(created_at);
CREATE INDEX idx_materials_origin_date ON materials(origin_date);

CREATE TABLE material_permissions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    material_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    granted_by BIGINT,
    created_at DATETIME NOT NULL DEFAULT current_timestamp(),
    UNIQUE KEY uq_material_permissions (material_id, user_id),
    CONSTRAINT fk_permissions_material FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE,
    CONSTRAINT fk_permissions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_permissions_granted_by FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE loans (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    material_id BIGINT NOT NULL,
    borrower_name TEXT NOT NULL,
    note TEXT NOT NULL DEFAULT '',
    taken_at DATETIME NOT NULL DEFAULT current_timestamp(),
    due_date DATE NOT NULL,
    returned_at DATETIME,
    CONSTRAINT fk_loans_material FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE INDEX idx_loans_returned ON loans(returned_at);
