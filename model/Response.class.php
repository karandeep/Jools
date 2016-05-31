<?php

/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */

/**
 * Description of Response
 *
 * @author snalin
 */
class Response {
    public $success;
    public $message;
    public $error_code;
    public $data;

    public function __construct() {
        $this->data = array();
    }
    
    public function getData() {
        $data = array(
            'success' => $this->success,
            'message' => $this->message,
            'error_code' => $this->error_code,
            'data' => $this->data,
        );
        return $data;
    }
}
