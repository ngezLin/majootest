CREATE TABLE IF NOT EXISTS `Merchants` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) NOT NULL,
  `merchant_name` varchar(40) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` bigint(20) NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `Outlets` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint(20) NOT NULL,
  `outlet_name` varchar(40) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` bigint(20) NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `merchant_id` (`merchant_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `Transactions` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint(20) NOT NULL,
  `outlet_id` bigint(20) NOT NULL,
  `bill_total` double NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` bigint(20) NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `updated_by` bigint(20) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `merchant_id` (`merchant_id`),
  KEY `outlet_id` (`outlet_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- Seed initial data
INSERT INTO `Merchants` (`id`, `user_id`, `merchant_name`, `created_at`, `created_by`, `updated_at`, `updated_by`) VALUES
(1, 1, 'merchant 1', '2026-08-01 07:00:00', 1, '2026-08-01 07:00:00', 1),
(2, 2, 'Merchant 2', '2026-08-01 07:00:00', 1, '2026-08-01 07:00:00', 1);

INSERT INTO `Outlets` (`id`, `merchant_id`, `outlet_name`, `created_at`, `created_by`, `updated_at`, `updated_by`) VALUES
(1, 1, 'outlet 1', '2026-08-01 07:00:00', 1, '2026-08-01 07:00:00', 1),
(2, 2, 'Outlet 2', '2026-08-01 07:00:00', 1, '2026-08-01 07:00:00', 1);

INSERT INTO `Transactions` (`id`, `merchant_id`, `outlet_id`, `bill_total`, `created_at`, `created_by`, `updated_at`, `updated_by`) VALUES
(1, 1, 1, 2000, '2026-08-01 10:00:00', 1, '2026-08-01 10:00:00', 1),
(2, 1, 1, 2500, '2026-08-01 11:00:00', 1, '2026-08-01 11:00:00', 1),
(3, 1, 1, 4000, '2026-08-02 12:00:00', 1, '2026-08-02 12:00:00', 1),
(4, 1, 1, 2000, '2026-08-02 13:00:00', 1, '2026-08-02 13:00:00', 1),
(5, 2, 2, 2500, '2026-08-03 14:00:00', 1, '2026-08-03 14:00:00', 1),
(6, 1, 1, 1000, '2026-08-04 15:00:00', 1, '2026-08-04 15:00:00', 1),
(7, 2, 2, 4000, '2026-08-04 16:00:00', 1, '2026-08-04 16:00:00', 1),
(8, 2, 2, 1000, '2026-08-04 17:00:00', 1, '2026-08-04 17:00:00', 1),
(9, 1, 1, 7000, '2026-08-05 18:00:00', 1, '2026-08-05 18:00:00', 1),
(10, 2, 2, 7000, '2026-08-05 19:00:00', 1, '2026-08-05 19:00:00', 1);
