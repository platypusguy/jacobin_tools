# jacobin_tools
Various tools related to the Jacobin JVM project (see [jacobin.org](http://www.jacobin.org))

These include:
* **gob2csv** - Converts the .gob file creates of the classes in the JDK 17 distribution into a CSV file format on standard output.
* **cleansrc** - Cleans source files of non-displayable characters with options for replacement and deletion.

The jacobin tools are installed by executing ```install.sh``` (Linux, MacOS, and Unix) or ```install.bat``` (Windows). After compilation, they reside in ```$GOPATH/bin```. It is recommended that ```$GOPATH/bin``` be an element of the user's executable PATH. 
