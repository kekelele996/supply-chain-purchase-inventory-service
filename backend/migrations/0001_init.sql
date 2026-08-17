-- 餐饮供应链 API 初始表结构（与 GORM AutoMigrate 保持一致，供参考/手工建表）
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(200) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'operator',
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    INDEX idx_users_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS suppliers (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    contact_person VARCHAR(50) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(100) NULL,
    address VARCHAR(200) NOT NULL,
    categories JSON NOT NULL,
    rating DOUBLE NOT NULL DEFAULT 3.0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    deleted_at DATETIME(3) NULL,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    INDEX idx_suppliers_status (status),
    INDEX idx_suppliers_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS inventory_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(20) NOT NULL,
    supplier_id BIGINT UNSIGNED NOT NULL,
    batch_no VARCHAR(50) NOT NULL UNIQUE,
    quantity DOUBLE NOT NULL DEFAULT 0,
    unit VARCHAR(20) NOT NULL,
    min_threshold DOUBLE NOT NULL DEFAULT 10.0,
    expiry_date DATE NOT NULL,
    storage_location VARCHAR(100) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'normal',
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    INDEX idx_inventory_supplier (supplier_id),
    INDEX idx_inventory_status (status),
    INDEX idx_inventory_category (category)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS purchase_orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(30) NOT NULL UNIQUE,
    supplier_id BIGINT UNSIGNED NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'draft',
    total_amount DECIMAL(12,2) NOT NULL DEFAULT 0,
    creator_id BIGINT UNSIGNED NOT NULL,
    approver_id BIGINT UNSIGNED NULL,
    approved_at DATETIME(3) NULL,
    completed_at DATETIME(3) NULL,
    notes TEXT NULL,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    INDEX idx_orders_supplier (supplier_id),
    INDEX idx_orders_status (status),
    INDEX idx_orders_creator (creator_id),
    INDEX idx_orders_approver (approver_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS purchase_order_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL,
    inventory_item_id BIGINT UNSIGNED NOT NULL,
    quantity DOUBLE NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(12,2) NOT NULL DEFAULT 0,
    INDEX idx_items_order (order_id),
    INDEX idx_items_inventory (inventory_item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    action VARCHAR(20) NOT NULL,
    target_type VARCHAR(50) NOT NULL,
    target_id BIGINT UNSIGNED NOT NULL,
    detail JSON NULL,
    created_at DATETIME(3) NULL,
    INDEX idx_logs_user (user_id),
    INDEX idx_logs_action (action),
    INDEX idx_logs_target (target_type, target_id),
    INDEX idx_logs_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
