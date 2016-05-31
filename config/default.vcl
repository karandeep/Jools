backend default {
    .host = "127.0.0.1";
    .port = "81";
    .connect_timeout = 1s;
}

sub vcl_recv {
  if (req.url ~ "^/jewelry") {
#    || req.url ~ "^/$") {
    unset req.http.cookie;
  }
}