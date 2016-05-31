<?php
require 'config/env.php';
require 'config/Constants.php';
$timestamp = time();

require 'lib/CSSmin.php';
$generationComment = '<!-- //Generated using the build script. DO NOT MODIFY -->'."\n";

$cssFile = 'css/jools.css';
$fh = fopen($cssFile, 'r');
$cssData = fread($fh, filesize($cssFile));
fclose($fh);
$cssData = str_replace('STATIC_URL', STATIC_URL, $cssData);

$cssMin = new CSSmin();
$cssResultFile = 'css/generated/jools-'.$timestamp.'.min.css';
$fh = fopen($cssResultFile, 'w') or die("can't open file");
if(DEBUG == 0) {
	$stringData = $cssMin->run($cssData);
} else {
	$stringData = $cssData;
}
fwrite($fh, $stringData);
fclose($fh);

$cssIncludeFile = 'include/generated/internal_css.php';
$cssTemplateFile = 'view/generated/internal_css.html';
$stringData = $generationComment . '<link rel="stylesheet" type="text/css" href="'. STATIC_URL .'/'.$cssResultFile.'"></link>';

$fh = fopen($cssIncludeFile, 'w') or die("can't open file");
fwrite($fh, $stringData);
fclose($fh);

$fh = fopen($cssTemplateFile, 'w') or die("can't open file");
fwrite($fh, $stringData);
fclose($fh);

$cssFile = 'css/jools-mobile.css';
$fh = fopen($cssFile, 'r');
$cssData = fread($fh, filesize($cssFile));
fclose($fh);
$cssData = str_replace('STATIC_URL', STATIC_URL, $cssData);

$cssMin = new CSSmin();
$cssResultFile = 'css/generated/jools-mobile-'.$timestamp.'.min.css';
$fh = fopen($cssResultFile, 'w') or die("can't open file");
if(DEBUG == 0) {
	$stringData = $cssMin->run($cssData);
} else {
	$stringData = $cssData;
}
fwrite($fh, $stringData);
fclose($fh);

$cssTemplateFile = 'view/generated/internal_css_mobile.html';
$stringData = $generationComment . '<link rel="stylesheet" type="text/css" href="'. STATIC_URL .'/'.$cssResultFile.'"></link>';

$fh = fopen($cssTemplateFile, 'w') or die("can't open file");
fwrite($fh, $stringData);
fclose($fh);

$jsFile = 'js/jools.js';
$fh = fopen($jsFile, 'r');
$jsData = fread($fh, filesize($jsFile));
fclose($fh);

require 'lib/JSmin.php';
$jsResultFile = 'js/generated/jools-'.$timestamp.'.min.js';
$fh = fopen($jsResultFile, 'w') or die("can't open file");
$jsMin = new JSmin($jsData);
if(DEBUG == 0) {
    $stringData = $jsMin->min();
} else {
    $stringData = $jsData;
}
fwrite($fh, $stringData);
fclose($fh);

$jsIncludeFile = 'include/generated/internal_js.php';
$jsTemplateFile = 'view/generated/internal_js.html';
$stringData = $generationComment . '<script type="text/javascript" src="'. STATIC_URL .'/'. $jsResultFile.'"></script>';
$fh = fopen($jsIncludeFile, 'w') or die("can't open file");
fwrite($fh, $stringData);
fclose($fh);
$fh = fopen($jsTemplateFile, 'w') or die("can't open file");
fwrite($fh, $stringData);
fclose($fh);

//Cleanup old concatenated css and js files
$filesToRemove = '';
$jsFiles = glob('js/generated/jools-[0-9]*.js');
$cssFiles = glob('css/generated/jools-[0-9]*.css');
$cssMobileFiles = glob('css/generated/jools-mobile-*.css');
$jsFileCount = count($jsFiles);
$cssFileCount = count($cssFiles);
$cssMobileFileCount = count($cssMobileFiles);
//Keeping current and one previous copy
for($jsIndex = 0; $jsIndex < $jsFileCount - 2; $jsIndex++) {
	$filesToRemove .= $jsFiles[$jsIndex] . ' ';
}
for($cssIndex = 0; $cssIndex < $cssFileCount - 2; $cssIndex++) {
	$filesToRemove .= $cssFiles[$cssIndex] . ' ';
}
for($cssMobileIndex = 0; $cssMobileIndex < $cssMobileFileCount - 2; $cssMobileIndex++) {
	$filesToRemove .= $cssMobileFiles[$cssMobileIndex] . ' ';
}
system("rm $filesToRemove");
echo "Build successfully completed\nFiles removed - $filesToRemove\n";
