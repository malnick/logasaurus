<img style="float: center;" src="https://dl.dropboxusercontent.com/u/77193293/masked_logitexample.png">
# Logit
Logit is a command line utility that queries elasticsearch in realtime, so you can tail logs just like you used to. 

## Build & Configure & Run

1. ```go build logit.go```

If you're on OSx I recommend placing the build binary in `/usr/local/bin` or somewhere else in your $PATH. 

#### config.yaml
The config.yaml is located in ```~/.logit/config.yaml```. YOU NEED TO MOVE THE config file to this location, making the dot dir along the way.

Make sure to update the config.yaml with your elasticsearch URI and port assignements.

You can override the config location is `-c` - don't use ~ or other shell expansion, provide the fully qualified path if you use this option.

## Usage

#### Defined query on the CLI:

```logit -d "some_query AND another_query"```

Will return matched messages from the last 10 minutes (see -sd override below) and resync backwards 5 seconds every 5 seconds (see -si override below).

#### Defined service in config.yaml:

```logit -s my_service_name```

Will return the query lookup from 'my_service_name' which should be in the 'define' section of the config.yaml.

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

    -c | config string
      Override the default configuration path. Default is ~/.logit.yaml on osx and /etc/logit.yaml on *nix distros. 

    -d | define string
      A temporary service definition. Must be a valid elasticsearch query. Can not be used with -s.

    -e | elasticsearch-uri string
      Override for `elasticsearch_uri` in config file. Default is localhost.

    -s | service abstraction
      A defined service in the config.yaml. Can not be used with -d.

    -si | sync interval number
      Time in seconds between elasticsearch queries. Default is 1s.

    -sd | sync depth number
      Time in minutes to sync backwards - only affects first sync

    -in | --index string
      Override for `logstash_index`. Default is logstash-*.

    -p | --port number
      Override for `elasticsearch_port` in config file. Default is 9300.       

    -v | --verbose 
      Verbose output.

