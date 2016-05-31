<?php

/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */

/**
 * Description of CallbackController
 *
 * @author snalin
 */
class CallbackController extends Controller {
    public function google() {
        require_once 'lib/google_contacts.php';
    }
    
    public function yahoo() {
        require_once 'lib/yahoo_contacts.php';
    }
}
