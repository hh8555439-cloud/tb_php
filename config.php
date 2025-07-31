<?php
$db = new PDO('mysql:host=localhost;dbname=tb_demo;charset=utf8mb4', 'root', 'root');
$db->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
$db->setAttribute(PDO::ATTR_EMULATE_PREPARES, false); // 强制使用真正的预处理
?>