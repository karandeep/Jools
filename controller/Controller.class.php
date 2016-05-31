<?php

/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */

/**
 * Description of Controller
 *
 * @author snalin
 */
class Controller {
    private $connClosed;
    
    public function __construct() {
        $this->connClosed = false;
        $this->response = new Response();
    }
    
    public function __destruct() {
        $this->close();
    }
    
    public function close() {
        if($this->connClosed) {
            return;
        }
        $this->connClosed = true;
        exit;
    }
    
    public function sendJsonResponse(Response $response) {
        header('Content-type: application/json');
        $data = $response->getData();
        echo json_encode($data);
        $this->close();
    }
}