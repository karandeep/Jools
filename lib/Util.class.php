<?php

/*
 * To change this license header, choose License Headers in Project Properties.
 * To change this template file, choose Tools | Templates
 * and open the template in the editor.
 */

/**
 * Description of Util
 *
 * @author snalin
 */
class Util {
    const STRING = 1;
    const INT = 2;
    
    public static function localize($str) {
        //TODO: add localization
        return $str;
    }
    
    public static function getCurrencySymbol() {
        return '<span class="WebRupee">Rs.</span>';
    }
    
    public static function formatCurrency($amount) {
        return money_format('%!i', $amount);
    }
    
    public static function sanitizeName($name) {
        //Make sure there are no semicolons in name - which we're using as separator
        $name = str_replace(';', '-', $name);
        if(empty($name)) {
            $name = '~';
        }
        return $name;
    }
    
    public static function sanitizeString($str) {
        $temp = strtolower($str);
        $result = str_replace(" ", "-", $temp);
        return $result;
    }
    
    public static function sanitizeStringForUrl($str) {
        $temp = strtolower($str);
        $result = str_replace(" ", "-", $temp);
        return $result;
    }
    
    public static function getParam($name, $type = self::STRING) {
        //TODO: Make things safer
        if(isset($_REQUEST[$name]) ) {
            $param = htmlentities($_REQUEST[$name]);
            if($type == self::STRING) {
                return $param;
            } else if($type == self::INT) {
                return intval($param);
            }
        }
        return PARAM_FETCH_FAILED;
    }
    
    public static function restorePlus($str) {
        return str_replace(array('%20', ' '), '+', $str);
    }
    
    public static function encrypt($text) {
        $encryptedText = trim(base64_encode(mcrypt_encrypt(MCRYPT_RIJNDAEL_256, 
                SALT, $text, MCRYPT_MODE_ECB, 
                mcrypt_create_iv(mcrypt_get_iv_size(MCRYPT_RIJNDAEL_256, MCRYPT_MODE_ECB), MCRYPT_RAND))));
        return $encryptedText;
    }
    
    public static function decrypt($text) {
        return trim(mcrypt_decrypt(MCRYPT_RIJNDAEL_256, 
                SALT, base64_decode($text), MCRYPT_MODE_ECB, 
                mcrypt_create_iv(mcrypt_get_iv_size(MCRYPT_RIJNDAEL_256, MCRYPT_MODE_ECB), MCRYPT_RAND)));
    }
    
    public static function isValidEmail($email) {
        $regex = '/^[_a-z0-9-]+(\.[_a-z0-9-]+)*@[a-z0-9-]+(\.[a-z0-9-]+)*(\.[a-z]{2,3})$/';
        // Run the preg_match() function on regex against the email address
        if (!preg_match($regex, $email)) {
            return false;
        }
        return true;
    }
    
    public static function urlencode_rfc3986($input) {
        if (is_array($input)) {
            return array_map(array('Util', 'urlencode_rfc3986'), $input);
        } else if (is_scalar($input)) {
            return str_replace(
                    '+', ' ', str_replace('%7E', '~', rawurlencode($input))
            );
        } else {
            return '';
        }
    }
    
    public static function generateHash($input) {
        return hash('sha512', SALT . $input . SALT);
    }
}
