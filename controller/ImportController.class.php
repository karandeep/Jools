<?php

class ImportController extends Controller {
	public function google() {
		$userId = Util::getParam("userId");
		$_SESSION['id'] = $userId;
		require_once 'lib/google_contacts.php';
	}

	public function yahoo() {
		$userId = Util::getParam("userId");
		$_SESSION['id'] = $userId;
		require_once 'lib/yahoo_contacts.php';
	}
}
