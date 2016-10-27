<?php
	header("Cache-Control: no-store, no-cache, must-revalidate");
	header("Pragma: no-cache");
	header("Content-Type: application/json");

	$module = escapeshellcmd($_GET['module']);
	$shell_file = dirname(__FILE__) . '/modules/' . $module . '.sh';
	echo shell_exec($shell_file);
