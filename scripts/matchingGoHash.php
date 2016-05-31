<?php
//http://www.ietf.org/rfc/rfc4648.txt

$toEnc = "Get hash for this string";
$nakedHexHash = hash('sha512', $toEnc);
$nakedHash = array_shift( unpack('H*', $nakedHexHash) );
$nakedPackedHash = pack('H*', $nakedHexHash);
$genHash = str_replace( array('+', '/'), array('-','_'), base64_encode($nakedPackedHash));
$correctHash = "A-QOvs6ugc-3uycAqOm3hooFAH6kPC9PnyL-esa_5dyf_1bvucYeoxzPN2Xq1fIPvN4kkxfjMHh9lIb9xpm0IA==";

if($genHash != $correctHash) {
    echo "Hashes don't match yet\n\r";
} else {
    echo "Hashes match!!!!!!!!!!\n\r";
}
echo "Generated Hash: $genHash";
echo "\n\r";
echo "Naked hex hash: $nakedHexHash";
echo "\n\r";
echo "Naked unpacked hash: $nakedHash";
echo "\n\r";
echo "Naked packed hash: $nakedPackedHash";
echo "\n\r";
