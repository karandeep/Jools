<?php
require_once 'include/all.php';
//Check in case json response is being requested
if(!isset($_SERVER['HTTP_ACCEPT']) || strpos($_SERVER['HTTP_ACCEPT'], 'javascript') === FALSE) {
    require_once 'include/doctype.php';
}

$controller = Util::getParam('controller');
if($controller == -1) {
    $controller = 'Home';
}
$action = Util::getParam('action');
if($action == -1) {
    $action = 'display';
}

$class = $controller .'Controller';
if(!class_exists($class)) {
    error_log("Invalid value for controller $controller");
    exit;
}

$controllerObj = new $class();
if(!method_exists($controllerObj, $action)) {
    error_log("Invalid value for action $action");
    exit;
}

session_start();
$controllerObj->$action();