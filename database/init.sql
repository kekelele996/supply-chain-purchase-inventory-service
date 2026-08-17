-- 餐饮供应链 API 数据库初始化脚本（容器首次启动时由 MySQL 自动执行）
-- 说明：表结构由 CREATE TABLE IF NOT EXISTS 保证；GORM AutoMigrate 会在应用启动时做增量同步。
-- 种子数据：3 个默认用户（密码 admin123 / manager123 / operator123，bcrypt 哈希）、4 个供应商、4 条库存记录。

CREATE DATABASE IF NOT EXISTS supplychain_db DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE supplychain_db;

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

-- 种子用户（密码见注释）
INSERT INTO users (username, password_hash, role, created_at, updated_at) VALUES
('admin', '$2a$10$q.XZyGM4QV1ZpRkja5jY5elATIMOY7vVpObVDVsEJwnIf3yW.lor6', 'admin', NOW(3), NOW(3)),
('manager', '$2a$10$LsySS.iuD0S2ogPJ.a2YrO5MSIPZOJ1dtxiu8WuknIB7KutciTGBK', 'manager', NOW(3), NOW(3)),
('operator', '$2a$10$mVw/OjDE5QXTKX0G.bSA6OPuqejRP5F54bSVzzkEL.VW//3p00.LW', 'operator', NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE username = VALUES(username);

-- 种子供应商
INSERT INTO suppliers (name, contact_person, phone, email, address, categories, rating, status, created_at, updated_at) VALUES
('绿源蔬菜基地', '王强', '13800000001', 'wang@greenfarm.cn', '北京市朝阳区蔬菜基地1号', JSON_ARRAY('蔬菜','水果'), 4.5, 'active', NOW(3), NOW(3)),
('鲜丰肉类供应', '李娜', '13800000002', 'lina@meat.cn', '上海市浦东新区肉联厂路2号', JSON_ARRAY('肉类'), 4.0, 'active', NOW(3), NOW(3)),
('海味源水产', '赵海', '13800000003', 'zhaohai@seafood.cn', '广州市海珠区水产市场3号', JSON_ARRAY('海鲜'), 3.5, 'active', NOW(3), NOW(3)),
('百味调料行', '陈香', '13800000004', 'chenxiang@seasoning.cn', '成都市锦江区调料市场4号', JSON_ARRAY('调料','干货'), 3.0, 'suspended', NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE name = VALUES(name);

-- 种子库存
INSERT INTO inventory_items (name, category, supplier_id, batch_no, quantity, unit, min_threshold, expiry_date, storage_location, status, created_at, updated_at) VALUES
('有机番茄', 'vegetable', 1, 'VEG-20260801', 120, 'kg', 30, DATE_ADD(CURDATE(), INTERVAL 5 DAY), 'A区-01架', 'normal', NOW(3), NOW(3)),
('猪五花肉', 'meat', 2, 'MEAT-20260802', 8, 'kg', 20, DATE_ADD(CURDATE(), INTERVAL 3 DAY), 'B区-冷冻2号', 'low', NOW(3), NOW(3)),
('东海带鱼', 'seafood', 3, 'SEA-20260720', 50, '箱', 10, DATE_SUB(CURDATE(), INTERVAL 1 DAY), 'C区-冷藏3号', 'expired', NOW(3), NOW(3)),
('郫县豆瓣酱', 'seasoning', 4, 'SEA-20260701', 200, '瓶', 50, DATE_ADD(CURDATE(), INTERVAL 60 DAY), 'D区-调料架', 'normal', NOW(3), NOW(3))
ON DUPLICATE KEY UPDATE batch_no = VALUES(batch_no);
