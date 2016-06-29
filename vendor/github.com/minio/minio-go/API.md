## API Documentation

### Minio client object creation
Minio client object is created using minio-go:
```go
package main

import (
    "fmt"

    "github.com/minio/minio-go"
)

func main() {
    secure := true // Make HTTPS requests by default.
    s3Client, err := minio.New("s3.amazonaws.com", "YOUR-ACCESSKEYID", "YOUR-SECRETACCESSKEY", secure)
    if err != nil {
        fmt.Println(err)
        return
    }
}
```

s3Client can be used to perform operations on S3 storage. APIs are described below.

### Bucket operations

* [`MakeBucket`](#MakeBucket)
* [`ListBuckets`](#ListBuckets)
* [`BucketExists`](#BucketExists)
* [`RemoveBucket`](#RemoveBucket)
* [`ListObjects`](#ListObjects)
* [`ListIncompleteUploads`](#ListIncompleteUploads)

### Object operations

* [`GetObject`](#GetObject)
* [`PutObject`](#PutObject)
* [`CopyObject`](#CopyObject)
* [`StatObject`](#StatObject)
* [`RemoveObject`](#RemoveObject)
* [`RemoveIncompleteUpload`](#RemoveIncompleteUpload)

### File operations.

* [`FPutObject`](#FPutObject)
* [`FGetObject`](#FPutObject)

### Bucket policy operations.

* [`SetBucketPolicy`](#SetBucketPolicy)
* [`GetBucketPolicy`](#GetBucketPolicy)
* [`RemoveBucketPolicy`](#RemoveBucketPolicy)

### Presigned operations

* [`PresignedGetObject`](#PresignedGetObject)
* [`PresignedPutObject`](#PresignedPutObject)
* [`PresignedPostPolicy`](#PresignedPostPolicy)

### Bucket operations
---------------------------------------
<a name="MakeBucket">
#### MakeBucket(bucketName string, location string) error
Create a new bucket.

__Parameters__
* `bucketName` _string_ - Name of the bucket.
* `location` _string_ - region valid values are _us-west-1_, _us-west-2_,  _eu-west-1_, _eu-central-1_, _ap-southeast-1_, _ap-northeast-1_, _ap-southeast-2_, _sa-east-1_

__Return Values__
* `err` _error_

__Example__
```go
err := s3Client.MakeBucket("mybucket", "us-west-1")
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println("Successfully created mybucket.")
```
---------------------------------------
<a name="ListBuckets">
#### ListBuckets() ([]BucketInfo, error)
Lists all buckets.

__Return Values__
* `buckets` _[]BucketInfo_ : array of `BucketInfo` objects that contain:
  * `Name` _string_ : name of the bucket
  * `CreationDate` _time.Time_ : date the bucket was created
* `err` _error_

__Example__
```go
buckets, err := s3Client.ListBuckets()
if err != nil {
    fmt.Println(err)
    return
}
for _, bucket := range buckets {
    fmt.Println(bucket)
}
```
---------------------------------------
<a name="BucketExists">
#### BucketExists(bucketName string) error
Check if bucket exists.

__Parameters__
* `bucketName` _string_ : name of the bucket

__Return Values__
* `err` _error_

__Example__
```go
err := s3Client.BucketExists("mybucket")
if err != nil {
    fmt.Println(err)
    return
}
```
---------------------------------------
<a name="RemoveBucket">
#### RemoveBucket(bucketName string) error
Remove a bucket.

__Parameters__
* `bucketName` _string_ : name of the bucket

__Return Values__
* `err` _error_

__Example__
```go
err := s3Client.RemoveBucket("mybucket")
if err != nil {
    fmt.Println(err)
    return
}
```
---------------------------------------
<a name="GetBucketPolicy">
#### GetBucketPolicy(bucketName string, objectPrefix string) (BucketPolicy, error)
Get access permissions on a bucket or a prefix.

__Parameters__
* `bucketName` _string_ : name of the bucket
* `objectPrefix` _string_ : name of the object prefix

__Return Values__
* `bucketPolicy` _BucketPolicy_ : string that contains: 'none', 'readonly', 'readwrite', or 'writeonly'
* `err` _error_

__Example__
```go
bucketPolicy, err := s3Client.GetBucketPolicy("mybucket", "")
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println("Access permissions for mybucket is", bucketPolicy)
```
---------------------------------------
<a name="SetBucketPolicy">
#### SetBucketPolicy(bucketname string, objectPrefix string, policy BucketPolicy) error
Set access permissions on bucket or an object prefix.

__Parameters__
* `bucketName` _string_: name of the bucket
* `objectPrefix` _string_ : name of the object prefix
* `policy` _BucketPolicy_: policy can be _BucketPolicyNone_, _BucketPolicyReadOnly_, _BucketPolicyReadWrite_, _BucketPolicyWriteOnly_

__Return Values__
* `err` _error_

__Example__
```go
err := s3Client.SetBucketPolicy("mybucket", "myprefix", BucketPolicyReadWrite)
if err != nil {
    fmt.Println(err)
    return
}
```
---------------------------------------
<a name="RemoveBucketPolicy">
#### RemoveBucketPolicy(bucketname string, objectPrefix string) error
Remove existing permissions on bucket or an object prefix.

__Parameters__
* `bucketName` _string_: name of the bucket
* `objectPrefix` _string_ : name of the object prefix

__Return Values__
* `err` _error_

__Example__
```go
err := s3Client.RemoveBucketPolicy("mybucket", "myprefix")
if err != nil {
    fmt.Println(err)
    return
}
```

---------------------------------------
<a name="ListObjects">
#### ListObjects(bucketName string, prefix string, recursive bool, doneCh chan struct{}) <-chan ObjectInfo
List objects in a bucket.

__Parameters__
* `bucketName` _string_: name of the bucket
* `objectPrefix` _string_: the prefix of the objects that should be listed
* `recursive` _bool_: `true` indicates recursive style listing and `false` indicates directory style listing delimited by '/'
* `doneCh`   chan struct{} : channel for pro-actively closing the internal go routine

__Return Values__
* `<-chan ObjectInfo` _chan ObjectInfo_: Read channel for all the objects in the bucket, the object is of the format:
  * `objectInfo.Key` _string_: name of the object
  * `objectInfo.Size` _int64_: size of the object
  * `objectInfo.ETag` _string_: etag of the object
  * `objectInfo.LastModified` _time.Time_: modified time stamp

__Example__
```go
// Create a done channel to control 'ListObjects' go routine.
doneCh := make(chan struct{})

// Indicate to our routine to exit cleanly upon return.
defer close(doneCh)

isRecursive := true
objectCh := s3Client.ListObjects("mybucket", "myprefix", isRecursive, doneCh)
for object := range objectCh {
    if object.Err != nil {
        fmt.Println(object.Err)
        return
    }
    fmt.Println(object)
}

```

---------------------------------------
<a name="ListIncompleteUploads">
#### ListIncompleteUploads(bucketName string, prefix string, recursive bool, doneCh chan struct{}) <-chan ObjectMultipartInfo
List partially uploaded objects in a bucket.

__Parameters__
* `bucketname` _string_: name of the bucket
* `prefix` _string_: prefix of the object names that are partially uploaded
* `recursive` bool: directory style listing when false, recursive listing when true
* `doneCh`   chan struct{} : channel for pro-actively closing the internal go routine

__Return Values__
* `<-chan ObjectMultipartInfo` _chan ObjectMultipartInfo_ : emits multipart objects of the format:
  * `multiPartObjInfo.Key` _string_: name of the incomplete object
  * `multiPartObjInfo.UploadID` _string_: upload ID of the incomplete object
  * `multiPartObjInfo.Size` _int64_: size of the incompletely uploaded object

__Example__
```go
// Create a done channel to control 'ListObjects' go routine.
doneCh := make(chan struct{})

// Indicate to our routine to exit cleanly upon return.
defer close(doneCh)

isRecursive := true
multiPartObjectCh := s3Client.ListIncompleteUploads("mybucket", "myprefix", isRecursive, doneCh)
for multiPartObject := range multiPartObjectCh {
    if multiPartObject.Err != nil {
        fmt.Println(multiPartObject.Err)
        return
    }
    fmt.Println(multiPartObject)
}
```

---------------------------------------
### Object operations
<a name="GetObject">
#### GetObject(bucketName string, objectName string) *Object
Download an object.

__Parameters__
* `bucketName` _string_: name of the bucket
* `objectName` _string_: name of the object

__Return Values__
* `object` _*Object_ : represents object reader.

__Example__
```go
object, err := s3Client.GetObject("mybucket", "photo.jpg")
if err != nil {
    fmt.Println(err)
    return
}
localFile _ := os.Open("/tmp/local-file")
if _, err := io.Copy(localFile, object); err != nil {
    fmt.Println(err)
    return
}
```
---------------------------------------
---------------------------------------
<a name="FGetObject">
#### FGetObject(bucketName string, objectName string, filePath string) error
Callback is called with `error` in case of error or `null` in case of success

__Parameters__
* `bucketName` _string_: name of the bucket
* `objectName` _string_: name of the object
* `filePath` _string_: path to which the object data will be written to

__Return Values__
* `err` _error_

__Example__
```go
err := s3Client.FGetObject("mybucket", "photo.jpg", "/tmp/photo.jpg")
if err != nil {
    fmt.Println(err)
    return
}
```
---------------------------------------
<a name="PutObject">
#### PutObject(bucketName string, objectName string, reader io.Reader, contentType string) (int64, error)
Upload contents from `io.Reader` to objectName.

__Parameters__
* `bucketName` _string_: name of the bucket
* `objectName` _string_: name of the object
* `reader` _io.Reader_: Any golang object implementing io.Reader
* `contentType` _string_: content type of the object.

__Return Values__
* `totalUploadedSize` _int64_: Size in bytes of uploaded file
* `err` _error_

__Example__
```go
file, err := os.Open("my-testfile")
if err != nil {
	fmt.Println(err)
    return
}
defer file.Close()

n, err := s3Client.PutObject("my-bucketname", "my-objectname", object, "application/octet-stream")
if err != nil {
    fmt.Println(err)
    return
}
```

---------------------------------------
<a name="CopyObject">
#### CopyObject(bucketName string, objectName string, objectSource string, conditions CopyConditions) error
Copy a source object into a new object with the provided name in the provided bucket.

__Parameters__
* `bucketName` _string_: name of the bucket
* `objectName` _string_: name of the object
* `objectSource` _string_: name of the object source.
* `conditions` _CopyConditions_: Collection of supported CopyObject conditions. ['x-amz-copy-source', 'x-amz-copy-source-if-match', 'x-amz-copy-source-if-none-match', 'x-amz-copy-source-if-unmodified-since', 'x-amz-copy-source-if-modified-since']

__Return Values__
* `err` _error_

__Example__
```go
// All following conditions are allowed and can be combined together.

// Set copy conditions.
var copyConds = minio.NewCopyConditions()
// Set modified condition, copy object modified since 2014 April.
copyConds.SetModified(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

// Set unmodified condition, copy object unmodified since 2014 April.
// copyConds.SetUnmodified(time.Date(2014, time.April, 0, 0, 0, 0, 0, time.UTC))

// Set matching ETag condition, copy object which matches the following ETag.
// copyConds.SetMatchETag("31624deb84149d2f8ef9c385918b653a")

// Set matching ETag except condition, copy object which does not match the following ETag.
// copyConds.SetMatchETagExcept("31624deb84149d2f8ef9c385918b653a")

err := s3Client.CopyObject("my-bucketname", "my-objectname", "/my-sourcebucketname/my-sourceobjectname", copyConds)
if err != nil {
    fmt.Println(err)
    return
}
```

---------------------------------------
<a name="FPutObject">
#### FPutObject(bucketName string, objectName string, filePath string, contentType string) (int64, error)
Uploads the object using contents from a file

__Parameters__
* `bucketName` _string_: name of the bucket
* `objectName` _string_: name of the object
* `filePath` _string_: file path of the file to be uploaded
* `contentType` _string_: content type of the object

__Return Values__
* `totalUploadedSize` _int64_: Size in bytes of uploaded file
* `err`: _error_

__Example__
```go
n, err := s3Client.FPutObject("my-bucketname", "my-objectname", "/tmp/my-filename.csv", "application/csv")
if err != nil {
    fmt.Println(err)
    return
}
```
---------------------------------------
<a name="StatObject">
#### StatObject(bucketName string, objectName string) (ObjectInfo, error)
Get metadata of an object.

__Parameters__
* `bucketName` _string_: name of the bucket
* `objectName` _string_: name of the object

__Return Values__
* `objInfo`   _ObjectInfo_ : object stat info for following format:
  * `objInfo.Size` _int64_: size of the object
  * `objInfo.ETag` _string_: etag of the object
  * `objInfo.ContentType` _string_: Content-Type of the object
  * `objInfo.LastModified` _time.Time_: modified time stamp
* `err` _error_

__Example__
```go
objInfo, err := s3Client.StatObject("mybucket", "photo.jpg")
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(objInfo)
```
---------------------------------------
<a name="RemoveObject">
#### RemoveObject(bucketName string, objectName string) error
Remove an object.

__Parameters__
* `bucketName` _string_: name of the bucket
* `objectName` _string_: name of the object

__Return Values__
* `err` _error_

__Example__
```go
err := s3Client.RemoveObject("mybucket", "photo.jpg")
if err != nil {
    fmt.Println(err)
    return
}
```
---------------------------------------
<a name="RemoveIncompleteUpload">
#### RemoveIncompleteUpload(bucketName string, objectName string) error
Remove an partially uploaded object.

__Parameters__
* `bucketName` _string_: name of the bucket
* `objectName` _string_: name of the object

__Return Values__
* `err` _error_

__Example__
```go
err := s3Client.RemoveIncompleteUpload("mybucket", "photo.jpg")
if err != nil {
    fmt.Println(err)
    return
}
```

### Presigned operations
---------------------------------------
<a name="PresignedGetObject">
#### PresignedGetObject(bucketName, objectName string, expiry time.Duration, reqParams url.Values) (*url.URL, error)
Generate a presigned URL for GET.

__Parameters__
* `bucketName` _string_: name of the bucket.
* `objectName` _string_: name of the object.
* `expiry` _time.Duration_: expiry in seconds.
* `reqParams` _url.Values_ : additional response header overrides supports _response-expires_, _response-content-type_, _response-cache-control_, _response-content-disposition_

__Return Values__
* `url` _*url.URL_ : Presigned URL for GET on an object
* `err` _error_

__Example__
```go
// Set request parameters for content-disposition.
reqParams := make(url.Values)
reqParams.Set("response-content-disposition", "attachment; filename=\"your-filename.txt\"")

// Generates a presigned url which expires in a day.
presignedURL, err := s3Client.PresignedGetObject("mybucket", "photo.jpg", time.Second * 24 * 60 * 60, reqParams)
if err != nil {
    fmt.Println(err)
    return
}
```

---------------------------------------
<a name="PresignedPutObject">
#### PresignedPutObject(bucketName string, objectName string, expiry time.Duration) (*url.URL, error)
Generate a presigned URL for PUT.
<blockquote>
NOTE: you can upload to S3 only with specified object name.
</blockquote>

__Parameters__
* `bucketName` _string_: name of the bucket
* `objectName` _string_: name of the object
* `expiry` _time.Duration_: expiry in seconds

__Return Values__
* `url` _*url.URL_ : Presigned URL for PUT on an object
* `err` _error

__Example__
```go
// Generates a url which expires in a day.
presignedURL, err := s3Client.PresignedPutObject("mybucket", "photo.jpg", time.Second * 24 * 60 * 60)
if err != nil {
    fmt.Println(err)
    return
}
```

---------------------------------------
<a name="PresignedPostPolicy">
#### PresignedPostPolicy(policy PostPolicy) (*url.URL, map[string]string, error)
PresignedPostPolicy we can provide policies specifying conditions restricting
what you want to allow in a POST request, such as bucket name where objects can be
uploaded, key name prefixes that you want to allow for the object being created and more.

We need to create our policy first:
```go
policy := minio.NewPostPolicy()
```
Apply upload policy restrictions:
```go
policy.SetBucket("my-bucketname")
policy.SetKey("my-objectname")
policy.SetExpires(time.Now().UTC().AddDate(0, 0, 10)) // expires in 10 days

// Only allow 'png' images.
policy.SetContentType("image/png")

// Only allow content size in range 1KB to 1MB.
policy.SetContentLengthRange(1024, 1024*1024)
```
Get the POST form key/value object:
```go
url, formData, err := s3Client.PresignedPostPolicy(policy)
if err != nil {
    fmt.Println(err)
    return
}
```

POST your content from the command line using `curl`:
```go
fmt.Printf("curl ")
for k, v := range m {
    fmt.Printf("-F %s=%s ", k, v)
}
fmt.Printf("-F file=@/etc/bash.bashrc ")
fmt.Printf("%s\n", url)
```
