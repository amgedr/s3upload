s3upload
========
s3upload is a command-line tool for uploading files and folders from multiple locations to Amazon S3.

## Configuration
Since s3upload currently has no command-line parameters, the S3 authentication keys and list of source folders and destination buckets are setup in the configuration file, .s3upload.conf. Which is created automatically in the current user's home folder when the program is executed for the first time. After setting up the configuration file just run the program manually or via a cron job:

```
s3upload
```

Here is the sample configurations file created if it doesn't exist:

```
[Auth]
AccessKey="your_amazon_accesskey"  # Get from Amazon
SecretKey="your_amazon_secretkey"  # Get from Amazon

[Setup]
#Overwrite the files if they exist in the buccket. Please note that setting
#this option to true may incur additional cost
overwrite=false

[Locations]
# local directory to upload. Use \\ on Windows.
source="/backup/"

# S3 bucket and path
destination="https://example-bucket.s3.amazonaws.com/"
```

To upload the contents of more than one source folder, repeat the last two lines. For example:

```
source="/backup1/"
destination="https://example-bucket1.s3.amazonaws.com/"

source="/backup2/"
destination="https://example-bucket2.s3.amazonaws.com/"

source="/backup3/"
destination="https://example-bucket3.s3.amazonaws.com/"
```



