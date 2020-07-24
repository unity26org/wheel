package help

import ()

var Content = `
Usage:
  wheel new APP_PATH [options]             # Creates new app

  wheel generate SUBJECT NAME ATTRIBUTES   # Adds new CRUD to an existing app. 
                                           # SUBJECT: scaffold/model/entity/handler. 
                                           # NAME: name of the model, entity or handler
                                           # ATTRIBUTES: when not a handler, is a pair of column name
                                           # and column type separated by ":" i.e. description:string
                                           # Available types are: 
                                           # string/text/integer/decimal/datetime/bool/reference.
                                           # When a handler "attributes" are functions inside handler.
                                           
Options:
  -d, [--database]                         # Preconfigure for selected database (options: mysql/postgresql)
  -G, [--skip-git]                         # Skip .gitignore file

More:
  -h, [--help]                             # Show this help message and quit
  -v, [--version]                          # Show Wheel version number and quit
`
