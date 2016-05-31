<?php

/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */

/**
 * Description of DB
 *
 * @author snalin
 */
class DB {
    private static $db = NULL;
    private static $sessionDB = NULL;
    
    public static function getConnection() {
        if (self::$db == NULL) {
            try {
                self::$db = new PDO("mysql:host=" . DB_HOST . ";dbname=jools", 'jools', '@ceVentur@P#one!x');
                self::$db->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);
            } catch (PDOException $e) {
                echo $e->getMessage();
            }
        }
        return self::$db;
    }
    
    public static function getSessionConnection() {
        if (self::$sessionDB == NULL) {
            try {
                self::$sessionDB = new PDO("mysql:host=" . DB_HOST . ";dbname=SecureSessions", 'joolsSessions', '$ec%re$2ec$om&th!ng*');
            } catch (PDOException $e) {
                echo $e->getMessage();
            }
        }
        return self::$sessionDB;
    }
}

