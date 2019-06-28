<?php
echo $message = "test email";

require_once('ses.php');

$ses = new SimpleEmailServices('access_key', 'secret_key');
$m = new SimpleEmailServiceMessage();

$m->addTo('sample@mail.com');
$m->setFrom('verifiedmain_fro_amazon@mail.com');
$m->setSubject('test subject');
$m->setMessageFromString($message);

var_dump($ses->sendEmail($m));
?>