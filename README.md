Some monitoring stuff for icinga:

A Nagios compatible plugin to check radius auth, uses https://github.com/wichert/pyrad
This is cobbled together from the sample code for an auth request from the pyrad docs and the nagios API documentation http://nagios.sourceforge.net/docs/3_0/pluginapi.html

A service notification for Telegram that accepts all icinga vars as arguments. This is mainly to workaround a limitation in Icinga Director which is unable to supply environment vars  to scripts. This is picking up on the ideas of https://github.com/Icinga/icinga2/pull/5170/commits, just my interpretation of this in golang.

A check script that checks performance counters from Nextclouds status api.
