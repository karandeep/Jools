<?php

/*
 * Copyright 2012 Google Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

require_once 'lib/google-api-php-client/src/Google_Client.php';

$client = new Google_Client();
$client->setApplicationName('Google Contacts');
$client->setScopes("http://www.google.com/m8/feeds/");

if (isset($_GET['code'])) {
    $client->authenticate();
    $_SESSION['token'] = $client->getAccessToken();
}

if (isset($_SESSION['token'])) {
    $client->setAccessToken($_SESSION['token']);
}

if ($client->getAccessToken()) {
    $req = new Google_HttpRequest("https://www.google.com/m8/feeds/contacts/default/full?max-results=50000");
    $val = $client->getIo()->authenticatedRequest($req);

    $response = $val->getResponseBody();
    $doc = new DOMDocument;
    $doc->recover = true;
    $doc->loadXML($response);

    $xpath = new DOMXPath($doc);
    $xpath->registerNamespace('gd', 'http://schemas.google.com/g/2005');
    $result = $xpath->query('//gd:email');

    $allEmails = '';
    $allNames = '';
    $emailData = array();
    $nameData = array();
    $count = 0;
    foreach ($result as $contact) {
        $email = $contact->getAttribute('address');
        $name = $contact->parentNode->getElementsByTagName('title')->item(0)->textContent;
        $name = Util::sanitizeName($name);
        $emailData[$count] = $email;
        $nameData[$count] = $name;
        $count++;
    }
    asort($nameData);
    foreach($nameData as $index => $name) {
        $allEmails .= $emailData[$index]. ';';
        $allNames .= "$name;";
    }
	//error_log("Params:" . $_SESSION['id'] . "Network:" . GOOGLE . "All Emails: $allEmails, AllNames: $allNames");
    User::storeContactEmails($_SESSION['id'], GOOGLE, $allEmails, $allNames);
    // The access token may have been updated lazily.
    $_SESSION['token'] = $client->getAccessToken();

    print '<script type="text/javascript">'
            //. 'window.opener.popupClosed();'
            . 'window.opener.popupClosed('.GOOGLE.');'
            . 'window.close();'
            . '</script>';
} else {
    $auth = $client->createAuthUrl();
}

if (isset($auth)) {
    print '<script type="text/javascript">window.location="' . $auth . '"</script>';
}
