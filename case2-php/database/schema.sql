-- Assessment-provided schema + seed data
-- Run this before migrations to set up base tables

CREATE TABLE IF NOT EXISTS `Merchants` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_id` int(40) NOT NULL,
  `merchant_name` varchar(40) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` bigint(20) NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

INSERT INTO Merchants VALUES
  (1, 1, 'merchant 1', now(), 1, now(), 1),
  (2, 2, 'Merchant 2', now(), 2, now(), 2);

CREATE TABLE IF NOT EXISTS `Outlets` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint(20) NOT NULL,
  `outlet_name` varchar(40) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` bigint(20) NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

INSERT INTO Outlets VALUES
  (1, 1, 'Outlet 1', now(), 1, now(), 1),
  (2, 2, 'Outlet 1', now(), 2, now(), 2),
  (3, 1, 'Outlet 2', now(), 1, now(), 1);

CREATE TABLE IF NOT EXISTS `Transactions` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `merchant_id` bigint(20) NOT NULL,
  `outlet_id` bigint(20) NOT NULL,
  `bill_total` double NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` bigint(20) NOT NULL,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_by` bigint(20) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;

INSERT INTO Transactions VALUES
  (1,  1, 1, 2000, '2026-08-01 12:30:04', 1, '2026-08-01 12:30:04', 1),
  (2,  1, 1, 2500, '2026-08-01 17:20:14', 1, '2026-08-01 17:20:14', 1),
  (3,  1, 1, 4000, '2026-08-02 12:30:04', 1, '2026-08-02 12:30:04', 1),
  (4,  1, 1, 1000, '2026-08-04 12:30:04', 1, '2026-08-04 12:30:04', 1),
  (5,  1, 1, 7000, '2026-08-05 16:59:30', 1, '2026-08-05 16:59:30', 1),
  (6,  1, 3, 2000, '2026-08-02 18:30:04', 1, '2026-08-02 18:30:04', 1),
  (7,  1, 3, 2500, '2026-08-03 17:20:14', 1, '2026-08-03 17:20:14', 1),
  (8,  1, 3, 4000, '2026-08-04 12:30:04', 1, '2026-08-04 12:30:04', 1),
  (9,  1, 3, 1000, '2026-08-04 12:31:04', 1, '2026-08-04 12:31:04', 1),
  (10, 1, 3, 7000, '2026-08-05 16:59:30', 1, '2026-08-05 16:59:30', 1),
  (11, 2, 2, 2000, '2026-08-01 18:30:04', 2, '2026-08-01 18:30:04', 2),
  (12, 2, 2, 2500, '2026-08-02 17:20:14', 2, '2026-08-02 17:20:14', 2),
  (13, 2, 2, 4000, '2026-08-03 12:30:04', 2, '2026-08-03 12:30:04', 2),
  (14, 2, 2, 1000, '2026-08-04 12:31:04', 2, '2026-08-04 12:31:04', 2),
  (15, 2, 2, 7000, '2026-08-05 16:59:30', 2, '2026-08-05 16:59:30', 2),
  (16, 2, 2, 2000, '2026-08-05 18:30:04', 2, '2026-08-05 18:30:04', 2),
  (17, 2, 2, 2500, '2026-08-06 17:20:14', 2, '2026-08-06 17:20:14', 2),
  (18, 2, 2, 4000, '2026-08-07 12:30:04', 2, '2026-08-07 12:30:04', 2),
  (19, 2, 2, 1000, '2026-08-08 12:31:04', 2, '2026-08-08 12:31:04', 2),
  (20, 2, 2, 7000, '2026-08-09 16:59:30', 2, '2026-08-09 16:59:30', 2),
  (21, 2, 2, 1000, '2026-08-10 12:31:04', 2, '2026-08-10 12:31:04', 2),
  (22, 2, 2, 7000, '2026-08-11 16:59:30', 2, '2026-08-11 16:59:30', 2);
