# Byte Pieper
People depend on various data store, some hosted by 3rd parties and
some running on premise and they all need backups.

Since all of them have similar requirements towards security and
availability, it makes sense to provide a common set of tools.

## Pipeline
Each pipeline consists of one input, any number of filters and one
output.

### Inputs
The pipeline supports multiple collectors for various data sources.

#### file
Reads a file.

#### tar
Reads a directory and streams it to the next filter in the pipeline.

### Filters
Filters can be chained.

#### rot13
Demo filter for testing.

### Outputs
Outputs write the data at the end of the pipeline to some location.

#### file
Wrtie to a local file.

## Configuration
There is a json based configuration which defines the pipelines.

## Examples
See [examples](examples/)

## Implementation
Byte piper is implemented in Go. The inputs implement io.Reader and
read data from various locations. Filters implement io.Reader, reading
from the previous element in the pipeline. Outputs implement io.Writer
. The last step is just a io.Copy to the output from the last filter in
the chain.
