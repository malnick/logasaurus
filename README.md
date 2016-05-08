Build Status: [![Circle CI](https://circleci.com/gh/malnick/logasaurus/tree/master.svg?style=svg)](https://circleci.com/gh/malnick/logasaurus/tree/master)
<img style="float: center;" src="https://dl.dropboxusercontent.com/u/77193293/logasaurus.png">
# Logasaurus 
Logasurous (loga command) is a command line utility that queries elasticsearch in realtime, so you can tail logs just like you used to. 

## Build & Configure & Run

```make build```

Will result in a binary in the form ```loga_$VERSION-$REVISION``` being created in this directory. If you checkout a tag, the appendage will be only the tagged version.

If you're on OSx I recommend placing the build binary in `/usr/local/bin` or somewhere else in your $PATH. 

Upon building, you can run loga for version information:

```
./loga_$VERSION -version
```

### Execute Once to Build loga.yaml Locally

```./loga_$VERSION```

Creates the default `loga.yaml` in the current directory:

```yaml
defined_queries:
  example: foo AND bar
sync_interval: 5
sync_depth: 10
elasticsearch_url: localhost
elasticsearch_port: "9200"
elasticsearch_index: ""
highlight_query: true
start_time: 0
count: 500
log_verbose: false
```

Update this file, replacing at minimum the default defined query (or you can override this on the CLI with a one-time query using `-d`), and the elasticsearch URL.

Once done, execute your first lookup (example with defined query in config):

```
./loga_$VERSION -s example
```

or do a one-time query on the CLI:

```
 ./loga_$VERSION -d foo AND bar
```

### Path to loga.yaml
You can override the config location is `-c` - don't use ~ or other shell expansion, provide the fully qualified path if you use this option.

## Usage

#### Defined query on the CLI:

```loga -d "some_query AND another_query"```

Will return matched messages from the last 10 minutes (see -sd override below) and resync backwards 5 seconds every 5 seconds (see -si override below).

#### Defined service in loga.yaml:

```loga -s my_service_name```

Will return the query lookup from 'my_service_name' which should be in the 'define' section of the loga.yaml.

loga will present the results from the search as a stream to stdin. Since the query is over standard http sockets, it'll return the query every 1s by default.

## Manpage

**NAME**
  
loga -- query ES logs on the CL

**SYNOPSIS**

loga [-d | --define string] [-i | --intervel time-in-seconds] [-v | --verbose] [-e | --elasticsearch-uri string-uri] [-p | --port elasticsearch-port] [-in | --es-index elasticsearch-index] [-c | --configuration path-to-config]

**DESCRIPTION**

The logasaurous (loga command) utility queries elasticsearch for logs based on a valid elasticsearch query. All requests are made to elasticsearch's REST endpoint over HTTP (HTTPS will be an option down the road). 

Logasaurous maintains a YAML configuration file where you can pre-set service definitions. You can leverage a one-time temporary service definition by using the ```define``` directive on the CLI.  

Many configurations in the config file can be overridden on the CLI as well. 

#### -c | Config 
  Override the default configuration path. Default is ~/.loga.yaml on osx and /etc/loga.yaml on *nix distros. 
  Ex: loga -c /fully/qualified/path/loga.yaml

#### -co | Count 
  Override the default count of queries to return. Default is 500.
  Ex: loga -d "some_query" -co 10 # Returns 10 queries from most recent. 

#### -d | Define 
  A temporary service definition. Must be a valid elasticsearch query. Can not be used with -s.
  Ex: loga -d "some_value AND \"a-long-string\""

#### -e | Elasticsearch URL
  Override for `elasticsearch_uri` in config file. Default is localhost.
  Ex: loga -d "some_query" -e my.elastic.com

#### -h | Enable Host Output
  Outputs the hostname for the log message before the message in cyan
  Ex: loga -d "some_query" -h

#### -hl | Highlight Query
  Highlights the string in the message that contains a match to your query. Outputs in yellow.
  Ex: logo -d "some_query" -hl

#### -s | Service Abstraction
  A defined service in the loga.yaml. Can not be used with -d.
  Ex: loga -s my_defined_service_in_loga.yaml

#### -si | Sync Interval 
  Time in seconds between elasticsearch queries. Default is 5s.
  Ex: loga -d "some_query" -si 10

#### -sd | Sync Depth
  Time in minutes to sync backwards - only affects first sync. Start time is always time.Meow() but this might change. 
  Ex: loga -d "some_query" -sd 120

#### -st | Start Time
  Time in past in minutes to start the search.
  Ex: loga -d "some_query" -st 20 # Starts the search 20 minutes in the past to the sync depth, so a window 30-20 minutes ago if used with defualt sync depth of 10 minutes. It will update itself every 5 seconds by default.

#### -p | Port
  Override for `elasticsearch_port` in config file. Default is 9300.       
  Ex: loga -d "some_query" -p 4500

#### -v | Verbose 
  Verbose output.
  Ex: Figure it out. 

## Tested
### Test Me

`make test`

### Tested On
- Elasticsearch: 1.4.4 - 2.3.2
- Logstash: 1.5 - 2.3.2
