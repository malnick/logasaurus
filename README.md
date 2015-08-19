# Logit
Logit is a command line utility that queries elasticsearch in realtime, so you can tail logs just like you used to. 

## Usage

```logit define $service_name```

You'll be prompted to enter a valid elasticsearch query that 'defines' a service that is sending logs to ES.

```logit tail $service_name```

Logit will present the results from the search as a stream to stdin. Since the query is over standard http sockets, it'll return the query every 1s by default.

## Manpage

**NAME**
  
logit -- query ES logs on the CL

**SYNOPSIS**

logit [-d | --define string] [-i | --intervel time-in-seconds] [-v | --verbose] [-e | --elasticsearch-uri string-uri] [-p | --port elasticsearch-port] [-in | --es-index elasticsearch-index] [-c | --configuration path-to-config]

**DESCRIPTION**

The logit utility queries elasticsearch for logs based on a valid elasticsearch query. All requests are made to elasticsearch's REST endpoint over HTTP (HTTPS will be an option down the road). 

Logit maintains a YAML configuration file where you can pre-set service definitions. You can leverage a one-time temporary service definition by using the ```define``` directive on the CLI.  

Many configurations in the config file can be overridden on the CLI as well. 

  Options:

    -c | --config string
      Override the default configuration path. Default is ~/.logit.yaml on osx and /etc/logit.yaml on *nix distros. 

    -d | --define string
      A temporary service definition. Must be a valid elasticsearch query.

    -e | --elasticsearch-uri string
      Override for `elasticsearch_uri` in config file. Default is localhost.

    -i | --interval number
      Time in seconds between elasticsearch queries. Default is 1s.

    -in | --index string
      Override for `logstash_index`. Default is logstash-*.

    -p | --port number
      Override for `elasticsearch_port` in config file. Default is 9300.       

    -v | --verbose 
      Verbose output.

**LICENSE**

 WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
                      Version 2, December 2004

   Copyright (C) 2004 Sam Hocevar <sam@hocevar.net>

   Everyone is permitted to copy and distribute verbatim or modified
   copies of this license document, and changing it is allowed as long
   as the name is changed.

              DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
     TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION

    0. You just DO WHAT THE FUCK YOU WANT TO.

