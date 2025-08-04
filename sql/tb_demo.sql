-- phpMyAdmin SQL Dump
-- version 4.8.5
-- https://www.phpmyadmin.net/
--
-- 主机： localhost
-- 生成日期： 2025-08-04 10:30:05
-- 服务器版本： 8.0.12
-- PHP 版本： 7.3.4

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- 数据库： `tb_demo`
--

-- --------------------------------------------------------

--
-- 表的结构 `comments`
--

CREATE TABLE `comments` (
  `id` int(11) NOT NULL,
  `content` text NOT NULL,
  `user_id` int(11) NOT NULL,
  `goods_id` int(11) NOT NULL,
  `to_user_id` int(11) DEFAULT NULL,
  `root_id` int(11) DEFAULT NULL,
  `type` enum('root','answer') NOT NULL DEFAULT 'root',
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `to_answer_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- 转存表中的数据 `comments`
--

INSERT INTO `comments` (`id`, `content`, `user_id`, `goods_id`, `to_user_id`, `root_id`, `type`, `create_time`, `to_answer_id`) VALUES
(14, '+3再水', 2, 1, NULL, NULL, 'root', '2025-07-29 03:41:03', NULL),
(16, '+1', 2, 1, NULL, NULL, 'root', '2025-07-29 03:57:22', NULL),
(17, '我要回复你了', 3, 1, NULL, NULL, 'root', '2025-07-29 03:57:47', NULL),
(18, '踩踩踩', 3, 1, 2, 17, 'answer', '2025-07-29 04:57:21', 17),
(19, '你真是天才', 2, 1, 3, 17, 'answer', '2025-07-29 04:57:36', 18),
(21, '哈哈哈', 2, 1, 3, 17, 'answer', '2025-07-29 04:58:55', 17),
(22, '你才是', 3, 1, 2, 17, 'answer', '2025-07-29 08:05:27', 19),
(23, '可以吗', 3, 5, 3, 0, 'root', '2025-07-29 09:13:19', -1),
(26, 'gg', 3, 1, 2, 17, 'answer', '2025-07-29 10:27:28', 21),
(27, 'shid', 3, 1, 2, 17, 'answer', '2025-07-29 10:28:47', 21),
(28, '5', 3, 1, 3, 17, 'answer', '2025-07-29 10:34:24', 27),
(29, '来看玫琳凯了', 7, 6, 7, 0, 'root', '2025-07-30 02:37:14', -1),
(31, '；lllll', 7, 5, 3, 23, 'answer', '2025-07-30 02:37:42', 24),
(33, '斤斤计较', 7, 5, 3, 0, 'root', '2025-07-30 02:38:30', -1),
(34, '8484984498', 7, 5, 3, 0, 'root', '2025-07-30 02:38:36', -1),
(35, 'lllll', 3, 5, 7, 30, 'answer', '2025-07-30 02:39:54', 32),
(37, '<a>123</a>', 3, 7, 3, 0, 'root', '2025-07-31 06:44:50', -1),
(38, 'undefined', 3, 7, 3, 0, 'root', '2025-08-01 02:08:44', -1),
(39, 'undefined', 3, 35, 3, 0, 'root', '2025-08-01 02:09:09', -1),
(40, 'undefined', 3, 35, 3, 39, 'answer', '2025-08-01 02:10:16', 39),
(41, '好', 3, 35, 3, 0, 'root', '2025-08-01 02:10:23', -1);

-- --------------------------------------------------------

--
-- 表的结构 `messages`
--

CREATE TABLE `messages` (
  `id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `content` text COLLATE utf8_unicode_ci NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- 转存表中的数据 `messages`
--

INSERT INTO `messages` (`id`, `user_id`, `content`, `created_at`) VALUES
(1, 2, '第一条留言', '2025-07-28 03:21:35'),
(5, 3, '第二条留言', '2025-07-29 08:12:01'),
(6, 7, '中文klkkkk', '2025-07-30 02:36:53'),
(7, 3, '<a>123</a>', '2025-07-30 02:41:10'),
(8, 3, '<a>123g</a>', '2025-07-31 06:44:42'),
(35, 3, '可以的', '2025-07-31 11:35:42');

-- --------------------------------------------------------

--
-- 表的结构 `users`
--

CREATE TABLE `users` (
  `id` int(11) NOT NULL,
  `username` varchar(50) COLLATE utf8_unicode_ci NOT NULL,
  `password` varchar(255) COLLATE utf8_unicode_ci NOT NULL,
  `email` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `role` enum('admin','user') COLLATE utf8_unicode_ci DEFAULT 'user',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- 转存表中的数据 `users`
--

INSERT INTO `users` (`id`, `username`, `password`, `email`, `role`, `created_at`) VALUES
(1, 'test', '$2y$10$RxSueXc/B57mYhx6X9Gte.k79nZVe/PNi7bs1FzFTKTe7k9tgembC', 'test@qq.com', 'user', '2025-07-25 09:21:27'),
(2, 'user', 'user', '2@qq.com', 'user', '2025-07-28 02:45:38'),
(3, 'root', 'root', '1@qq.com', 'admin', '2025-07-28 03:28:43'),
(4, 'SA24225162', '11', '3@qq.com', 'user', '2025-07-28 03:30:36'),
(5, 'XX', 'XX', 'test1@qq.com', 'user', '2025-07-28 03:33:06'),
(6, 'root1111111', '11', 'te@qq.com', 'user', '2025-07-30 02:34:34'),
(7, 'hankinTest', 'hankinTest', 'hankinTest@qq.com', 'user', '2025-07-30 02:35:24');

--
-- 转储表的索引
--

--
-- 表的索引 `comments`
--
ALTER TABLE `comments`
  ADD PRIMARY KEY (`id`),
  ADD KEY `goods_id` (`goods_id`),
  ADD KEY `root_id` (`root_id`);

--
-- 表的索引 `messages`
--
ALTER TABLE `messages`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_id` (`user_id`);

--
-- 表的索引 `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`);

--
-- 在导出的表使用AUTO_INCREMENT
--

--
-- 使用表AUTO_INCREMENT `comments`
--
ALTER TABLE `comments`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=42;

--
-- 使用表AUTO_INCREMENT `messages`
--
ALTER TABLE `messages`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=36;

--
-- 使用表AUTO_INCREMENT `users`
--
ALTER TABLE `users`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
