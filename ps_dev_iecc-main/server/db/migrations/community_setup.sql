-- Community Dashboard Database Setup
-- Run this script to create the required tables and seed resource entries

-- -------------------------------------------------------
-- 1. Core community table
-- -------------------------------------------------------
CREATE TABLE IF NOT EXISTS `community` (
  `id`               INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `name`             VARCHAR(255)  NOT NULL,
  `icon`             VARCHAR(100)  NOT NULL DEFAULT 'bx-group',
  `established_date` DATE          NOT NULL,
  `total_points`     INT           NOT NULL DEFAULT 0,
  `member_count`     INT           NOT NULL DEFAULT 0,
  `rating`           DECIMAL(3,1)  NOT NULL DEFAULT 4.5,
  `reliability`      TINYINT       NOT NULL DEFAULT 85 COMMENT 'Percentage 0-100',
  `quality`          TINYINT       NOT NULL DEFAULT 78 COMMENT 'Percentage 0-100',
  `frequency`        TINYINT       NOT NULL DEFAULT 92 COMMENT 'Percentage 0-100',
  `status`           CHAR(1)       NOT NULL DEFAULT '1' COMMENT '1=active,0=inactive',
  `created_at`       DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at`       DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -------------------------------------------------------
-- 2. Community members
-- -------------------------------------------------------
CREATE TABLE IF NOT EXISTS `community_members` (
  `id`           INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `community_id` INT UNSIGNED NOT NULL,
  `user_id`      VARCHAR(50)  NOT NULL,
  `points`       INT          NOT NULL DEFAULT 0,
  `joined_at`    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_community_user` (`community_id`, `user_id`),
  KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -------------------------------------------------------
-- 3. Community activities / feed
-- -------------------------------------------------------
CREATE TABLE IF NOT EXISTS `community_activities` (
  `id`           INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `community_id` INT UNSIGNED NOT NULL,
  `title`        VARCHAR(255) NOT NULL,
  `description`  TEXT,
  `icon`         VARCHAR(100) NOT NULL DEFAULT 'bx-bell',
  `created_at`   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_community_id` (`community_id`),
  KEY `idx_created_at`   (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -------------------------------------------------------
-- 4. Community events & announcements
-- -------------------------------------------------------
CREATE TABLE IF NOT EXISTS `community_events` (
  `id`           INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `community_id` INT UNSIGNED NOT NULL,
  `title`        VARCHAR(255) NOT NULL,
  `description`  TEXT,
  `event_date`   DATE         NOT NULL,
  `type`         ENUM('event','announcement') NOT NULL DEFAULT 'event',
  `created_at`   DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_community_id` (`community_id`),
  KEY `idx_event_date`   (`event_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -------------------------------------------------------
-- 5. Community targets
-- -------------------------------------------------------
CREATE TABLE IF NOT EXISTS `community_targets` (
  `id`             INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `community_id`   INT UNSIGNED NOT NULL,
  `weekly_target`  INT          NOT NULL DEFAULT 100,
  `weekly_current` INT          NOT NULL DEFAULT 0,
  `updated_at`     DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uq_community_target` (`community_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -------------------------------------------------------
-- 6. Community mandates (task checklist)
-- -------------------------------------------------------
CREATE TABLE IF NOT EXISTS `community_mandates` (
  `id`           INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `community_id` INT UNSIGNED NOT NULL,
  `title`        VARCHAR(255) NOT NULL,
  `completed`    TINYINT(1)   NOT NULL DEFAULT 0,
  `sort_order`   INT          NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  KEY `idx_community_id` (`community_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- -------------------------------------------------------
-- 7. Register Community Dashboard resource entries
--    Adjust res_group and sort_by to match your environment.
-- -------------------------------------------------------
INSERT INTO `master_resource_v2`
  (`path`, `icon`, `menu`, `name`, `element`, `status`, `api_for`, `activity`, `res_group`, `sort_by`, `group`)
VALUES
  -- Frontend app route (shows in sidebar)
  ('/community', 'bx-group', 1, 'Community', 'AppCommunityDashboard', '1', 'app', 0, 1, 90, 'Community'),

  -- Backend API routes (protected via ScopeMiddleware)
  ('/community/details',    '', 0, '', '', '1', 'api', 0, 1, 91, 'Community'),
  ('/community/members',    '', 0, '', '', '1', 'api', 0, 1, 92, 'Community'),
  ('/community/activities', '', 0, '', '', '1', 'api', 0, 1, 93, 'Community'),
  ('/community/events',     '', 0, '', '', '1', 'api', 0, 1, 94, 'Community'),
  ('/community/targets',    '', 0, '', '', '1', 'api', 0, 1, 95, 'Community');
