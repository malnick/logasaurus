# Logit
Logit is a command line utility that queries elasticsearch in realtime, so you can tail logs just like you used to. 

## Usage

```logit define $service_name```

You'll be prompted to enter a valid elasticsearch query that 'defines' a service that is sending logs to ES.

```logit tail $service_name```

Logit will present the results from the search as a stream to stdin. Since the query is over standard http sockets, it'll return the query every 1s by default.

## Manpage

#### NAME
    logit -- query ES logs on the CL

#### SYNOPSIS
    logit [-d | --define string] [-i | --intervel time-in-seconds] [-v | --verbose] [-e | --elasticsearch-uri string-uri] [-p | --port elasticsearch-port] [-in | --es-index elasticsearch-index] [-c | --configuration path-to-config]

#### DESCRIPTION
    The logit utility queries elasticsearch for logs based on a valid elasticsearch query. All requests are made to elasticsearch's REST endpoint over HTTP (HTTPS will be an option down the road). 

    Logit maintains a YAML configuration file where you can pre-set service definitions. You can leverage a one-time temporary service definition by using the ```define``` directive on the CLI.  

    Many configurations in the config file can be overridden on the CLI as well. 


