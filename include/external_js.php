<script type="text/javascript" src="<?= STATIC_URL ?>/js/libs/jquery/jquery-2.0.3.min.js"></script>
<script type="text/javascript" src="<?= STATIC_URL ?>/js/libs/jqueryui/jquery-ui-1.10.3.custom.min.js"></script>
<?php if (!empty($response->data['includeBlocksIt'])) { ?>
    <script type="text/javascript" src="<?= STATIC_URL ?>/js/libs/jquery/jquery.imagesloaded.js"></script>
    <script type="text/javascript" src="<?= STATIC_URL ?>/js/libs/blocksit/blocksit.min.js"></script>
    <?php
}