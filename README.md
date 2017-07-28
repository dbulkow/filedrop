# File Drop
I wanted a way to share files between groups that use different operating systems. Exported file shares isn't working so people tend to email, sometimes large, files. Filedrop provides a web-based interface that everyone can use. Files are uploaded with an explicit expiration so they are cleaned up automatically.
# Build
```
git clone https://github.com/dbulkow/filedrop.git
cd filedrop
go generate
go build
```
# Deploy in docker
The HTML files are built into the binary for easier deployment. For the golang build image this means we need to generate them before building the filedrop image.
```
go generate
docker build -t filedrop .
```
Use environment variables to configure filedrop.

| Env Variable     | Description |
| ------------     | ----------- |
| FILEDROP_ROOT    | The directory in which the downloads subdir will be created|
| FILEDROP_SERVER_URL     | Advertised URL. Used after uploading a file to tell the user what URL to send to recipients |
| FILEDROP_ADDRESS | Listen address - defaults to "127.0.0.1". Set to "0.0.0.0" to listen on all interfaces |
| FILEDROP_PORT    | Port number on which to listen |

Be sure to set the advertised URL so recipients know how to reach the server.
```
docker run -d -p 8080:8080 -e FILEDROP_SERVER_URL="http://my-machine.corp.com" --name filedrop filedrop
```

Since file storage is intended to be ephemeral I've just been letting files live in the container. For a site with more traffic it might be desirable to save files elsewhere. Add a volume to the docker run command line:
```
-v <some downloads directory>:/downloads
```

# License

MIT License

Copyright (c) 2017 David Bulkow

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
