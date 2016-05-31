<?php

/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */

/**
 * Description of User
 *
 * @author snalin
 */
class User {
    public static function storeContactEmails($encUserId, $network, $emails, $names) {
        $conn = DB::getConnection();
        $sql = 'INSERT IGNORE INTO Emails (encUserId, network, emails, names) VALUES (:encUserId, :network, :emails, :names);';
        $st = $conn->prepare($sql);
        $st->execute(array('encUserId' => $encUserId,
            'network' => $network,
            'emails' => $emails,
            'names' => $names,
		));

		if($network == GOOGLE) {
			$sql = "UPDATE User SET gmailImport = 1 WHERE encId = :encId LIMIT 1";
		} else if($network == YAHOO) {
			$sql = "UPDATE User SET yahooImport = 1 WHERE encId = :encId LIMIT 1";
		}
		$st = $conn->prepare($sql);
		$st->execute( array(
			'encId' => $encUserId,
		));
    }
}
