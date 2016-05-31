<html>
<head>
<?php
require_once __DIR__ .'/../config/env.php';
require_once __DIR__ .'/../config/Constants.php';
require_once __DIR__ .'/../include/generated/internal_css.php';
?>
</head>
<body>
    <h1>
        Page not found
    </h1>
    <style type="text/css">
        #goog-wm { } 
        #goog-wm h3.closest-match { } 
        #goog-wm h3.closest-match a { } 
        #goog-wm h3.other-things { } 
        #goog-wm ul li { } 
        #goog-wm li.search-goog { display: list-item; }
    </style>
    <script type="text/javascript">
    var GOOG_FIXURL_LANG = 'en';
    var GOOG_FIXURL_SITE = '<?= BASE_URL ?>';
    </script>
    <script type="text/javascript" src="https://linkhelp.clients.google.com/tbproxy/lh/wm/fixurl.js"></script>
    <script type="text/javascript" src="<?= STATIC_URL ?>/js/libs/jquery/jquery-1.10.2.min.js"></script>
    <script type="text/javascript" src="<?= STATIC_URL ?>/js/libs/jqueryui/jquery-ui-1.10.3.custom.min.js"></script>
<?php
require_once __DIR__ .'/../include/generated/internal_js.php';
?>
    <script type="text/javascript">
        var configObj = null;
        var trackObj = null;
        var userObj = null;
        var utilObj = null;
        $(document).ready(function() {
            configObj = new Config('<?= STATIC_URL ?>', '<?= BASE_URL ?>', '', '404');
            utilObj = new Util();
            userObj = new User(false, localStorage.getItem('userData'), -1);
            trackObj = new Track();
            trackObj.count("visit", "404", window.location.pathname);
        });
    </script>
    </body>
</html>
