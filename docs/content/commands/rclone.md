---
date: 2017-09-30T14:21:35+01:00
title: "rclone"
slug: rclone
url: /commands/rclone/
---
## rclone

Sync files and directories to and from local and remote object stores - v1.38

### Synopsis



Rclone is a command line program to sync files and directories to and
from various cloud storage systems and using file transfer services, such as:

  * Amazon Drive
  * Amazon S3
  * Backblaze B2
  * Box
  * Dropbox
  * FTP
  * Google Cloud Storage
  * Google Drive
  * HTTP
  * Hubic
  * Microsoft Azure Blob Storage
  * Microsoft OneDrive
  * Openstack Swift / Rackspace cloud files / Memset Memstore
  * QingStor
  * SFTP
  * Yandex Disk
  * The local filesystem

Features

  * MD5/SHA1 hashes checked at all times for file integrity
  * Timestamps preserved on files
  * Partial syncs supported on a whole file basis
  * Copy mode to just copy new/changed files
  * Sync (one way) mode to make a directory identical
  * Check mode to check for file hash equality
  * Can sync to and from network, eg two different cloud accounts

See the home page for installation, usage, documentation, changelog
and configuration walkthroughs.

  * https://rclone.org/


```
rclone [flags]
```

### Options

```
      --acd-templink-threshold int        Files >= this size will be downloaded via their tempLink. (default 9G)
      --acd-upload-wait-per-gb duration   Additional time per GB to wait after a failed complete upload to see if it appears. (default 3m0s)
      --ask-password                      Allow prompt for password for encrypted configuration. (default true)
      --azureblob-chunk-size int          Upload chunk size. Must fit in memory. (default 4M)
      --azureblob-upload-cutoff int       Cutoff for switching to chunked upload (default 256M)
      --b2-chunk-size int                 Upload chunk size. Must fit in memory. (default 96M)
      --b2-hard-delete                    Permanently delete files on remote removal, otherwise hide files.
      --b2-test-mode string               A flag string for X-Bz-Test-Mode header.
      --b2-upload-cutoff int              Cutoff for switching to chunked upload (default 190.735M)
      --b2-versions                       Include old versions in directory listings.
      --backup-dir string                 Make backups into hierarchy based in DIR.
      --bind string                       Local address to bind to for outgoing connections, IPv4, IPv6 or name.
      --box-upload-cutoff int             Cutoff for switching to multipart upload (default 50M)
      --buffer-size int                   Buffer size when copying files. (default 16M)
      --bwlimit BwTimetable               Bandwidth limit in kBytes/s, or use suffix b|k|M|G or a full timetable.
      --checkers int                      Number of checkers to run in parallel. (default 8)
  -c, --checksum                          Skip based on checksum & size, not mod-time & size
      --config string                     Config file. (default "/home/ncw/.rclone.conf")
      --contimeout duration               Connect timeout (default 1m0s)
  -L, --copy-links                        Follow symlinks and copy the pointed to item.
      --cpuprofile string                 Write cpu profile to file
      --crypt-show-mapping                For all files listed show how the names encrypt.
      --delete-after                      When synchronizing, delete files on destination after transfering
      --delete-before                     When synchronizing, delete files on destination before transfering
      --delete-during                     When synchronizing, delete files during transfer (default)
      --delete-excluded                   Delete files on dest excluded from sync
      --disable string                    Disable a comma separated list of features.  Use help to see a list.
      --drive-auth-owner-only             Only consider files owned by the authenticated user.
      --drive-chunk-size int              Upload chunk size. Must a power of 2 >= 256k. (default 8M)
      --drive-formats string              Comma separated list of preferred formats for downloading Google docs. (default "docx,xlsx,pptx,svg")
      --drive-list-chunk int              Size of listing chunk 100-1000. 0 to disable. (default 1000)
      --drive-shared-with-me              Only show files that are shared with me
      --drive-skip-gdocs                  Skip google documents in all listings.
      --drive-trashed-only                Only show files that are in the trash
      --drive-upload-cutoff int           Cutoff for switching to chunked upload (default 8M)
      --drive-use-trash                   Send files to the trash instead of deleting permanently. (default true)
      --dropbox-chunk-size int            Upload chunk size. Max 150M. (default 128M)
  -n, --dry-run                           Do a trial run with no permanent changes
      --dump-auth                         Dump HTTP headers with auth info
      --dump-bodies                       Dump HTTP headers and bodies - may contain sensitive info
      --dump-filters                      Dump the filters to the output
      --dump-headers                      Dump HTTP headers - may contain sensitive info
      --exclude stringArray               Exclude files matching pattern
      --exclude-from stringArray          Read exclude patterns from file
      --fast-list                         Use recursive list if available. Uses more memory but fewer transactions.
      --files-from stringArray            Read list of source-file names from file
  -f, --filter stringArray                Add a file-filtering rule
      --filter-from stringArray           Read filtering patterns from a file
      --gcs-location string               Default location for buckets (us|eu|asia|us-central1|us-east1|us-east4|us-west1|asia-east1|asia-noetheast1|asia-southeast1|australia-southeast1|europe-west1|europe-west2).
      --gcs-storage-class string          Default storage class for buckets (MULTI_REGIONAL|REGIONAL|STANDARD|NEARLINE|COLDLINE|DURABLE_REDUCED_AVAILABILITY).
  -h, --help                              help for rclone
      --ignore-checksum                   Skip post copy check of checksums.
      --ignore-existing                   Skip all files that exist on destination
      --ignore-size                       Ignore size when skipping use mod-time or checksum.
  -I, --ignore-times                      Don't skip files that match size and time - transfer all files
      --immutable                         Do not modify files. Fail if existing files have been modified.
      --include stringArray               Include files matching pattern
      --include-from stringArray          Read include patterns from file
      --local-no-unicode-normalization    Don't apply unicode normalization to paths and filenames
      --log-file string                   Log everything to this file
      --log-level string                  Log level DEBUG|INFO|NOTICE|ERROR (default "NOTICE")
      --low-level-retries int             Number of low level retries to do. (default 10)
      --max-age string                    Don't transfer any file older than this in s or suffix ms|s|m|h|d|w|M|y
      --max-depth int                     If set limits the recursion depth to this. (default -1)
      --max-size int                      Don't transfer any file larger than this in k or suffix b|k|M|G (default off)
      --memprofile string                 Write memory profile to file
      --min-age string                    Don't transfer any file younger than this in s or suffix ms|s|m|h|d|w|M|y
      --min-size int                      Don't transfer any file smaller than this in k or suffix b|k|M|G (default off)
      --modify-window duration            Max time diff to be considered the same (default 1ns)
      --no-check-certificate              Do not verify the server SSL certificate. Insecure.
      --no-gzip-encoding                  Don't set Accept-Encoding: gzip.
      --no-traverse                       Don't traverse destination file system on copy.
      --no-update-modtime                 Don't update destination mod-time if files identical.
      --old-sync-method                   Deprecated - use --fast-list instead
  -x, --one-file-system                   Don't cross filesystem boundaries.
      --onedrive-chunk-size int           Above this size files will be chunked - must be multiple of 320k. (default 10M)
      --onedrive-upload-cutoff int        Cutoff for switching to chunked upload - must be <= 100MB (default 10M)
  -q, --quiet                             Print as little stuff as possible
      --retries int                       Retry operations this many times if they fail (default 3)
      --s3-acl string                     Canned ACL used when creating buckets and/or storing objects in S3
      --s3-storage-class string           Storage class to use when uploading S3 objects (STANDARD|REDUCED_REDUNDANCY|STANDARD_IA)
      --size-only                         Skip based on size only, not mod-time or checksum
      --skip-links                        Don't warn about skipped symlinks.
      --stats duration                    Interval between printing stats, e.g 500ms, 60s, 5m. (0 to disable) (default 1m0s)
      --stats-log-level string            Log level to show --stats output DEBUG|INFO|NOTICE|ERROR (default "INFO")
      --stats-unit string                 Show data rate in stats as either 'bits' or 'bytes'/s (default "bytes")
      --streaming-upload-cutoff int       Cutoff for switching to chunked upload if file size is unknown. Upload starts after reaching cutoff or when file ends. (default 100k)
      --suffix string                     Suffix for use with --backup-dir.
      --swift-chunk-size int              Above this size files will be chunked into a _segments container. (default 5G)
      --syslog                            Use Syslog for logging
      --syslog-facility string            Facility for syslog, eg KERN,USER,... (default "DAEMON")
      --timeout duration                  IO idle timeout (default 5m0s)
      --tpslimit float                    Limit HTTP transactions per second to this.
      --tpslimit-burst int                Max burst of transactions for --tpslimit. (default 1)
      --track-renames                     When synchronizing, track file renames and do a server side move if possible
      --transfers int                     Number of file transfers to run in parallel. (default 4)
  -u, --update                            Skip files that are newer on the destination.
      --user-agent string                 Set the user-agent to a specified string. The default is rclone/ version (default "rclone/v1.38")
  -v, --verbose count[=-1]                Print lots more stuff (repeat for more)
  -V, --version                           Print the version number
```

### SEE ALSO
* [rclone authorize](/commands/rclone_authorize/)	 - Remote authorization.
* [rclone cat](/commands/rclone_cat/)	 - Concatenates any files and sends them to stdout.
* [rclone check](/commands/rclone_check/)	 - Checks the files in the source and destination match.
* [rclone cleanup](/commands/rclone_cleanup/)	 - Clean up the remote if possible
* [rclone config](/commands/rclone_config/)	 - Enter an interactive configuration session.
* [rclone copy](/commands/rclone_copy/)	 - Copy files from source to dest, skipping already copied
* [rclone copyto](/commands/rclone_copyto/)	 - Copy files from source to dest, skipping already copied
* [rclone cryptcheck](/commands/rclone_cryptcheck/)	 - Cryptcheck checks the integrity of a crypted remote.
* [rclone cryptdecode](/commands/rclone_cryptdecode/)	 - Cryptdecode returns unencrypted file names.
* [rclone dbhashsum](/commands/rclone_dbhashsum/)	 - Produces a Dropbbox hash file for all the objects in the path.
* [rclone dedupe](/commands/rclone_dedupe/)	 - Interactively find duplicate files delete/rename them.
* [rclone delete](/commands/rclone_delete/)	 - Remove the contents of path.
* [rclone genautocomplete](/commands/rclone_genautocomplete/)	 - Output completion script for a given shell.
* [rclone gendocs](/commands/rclone_gendocs/)	 - Output markdown docs for rclone to the directory supplied.
* [rclone listremotes](/commands/rclone_listremotes/)	 - List all the remotes in the config file.
* [rclone ls](/commands/rclone_ls/)	 - List all the objects in the path with size and path.
* [rclone lsd](/commands/rclone_lsd/)	 - List all directories/containers/buckets in the path.
* [rclone lsjson](/commands/rclone_lsjson/)	 - List directories and objects in the path in JSON format.
* [rclone lsl](/commands/rclone_lsl/)	 - List all the objects path with modification time, size and path.
* [rclone md5sum](/commands/rclone_md5sum/)	 - Produces an md5sum file for all the objects in the path.
* [rclone mkdir](/commands/rclone_mkdir/)	 - Make the path if it doesn't already exist.
* [rclone mount](/commands/rclone_mount/)	 - Mount the remote as a mountpoint. **EXPERIMENTAL**
* [rclone move](/commands/rclone_move/)	 - Move files from source to dest.
* [rclone moveto](/commands/rclone_moveto/)	 - Move file or directory from source to dest.
* [rclone ncdu](/commands/rclone_ncdu/)	 - Explore a remote with a text based user interface.
* [rclone obscure](/commands/rclone_obscure/)	 - Obscure password for use in the rclone.conf
* [rclone purge](/commands/rclone_purge/)	 - Remove the path and all of its contents.
* [rclone rcat](/commands/rclone_rcat/)	 - Copies standard input to file on remote.
* [rclone rmdir](/commands/rclone_rmdir/)	 - Remove the path if empty.
* [rclone rmdirs](/commands/rclone_rmdirs/)	 - Remove empty directories under the path.
* [rclone sha1sum](/commands/rclone_sha1sum/)	 - Produces an sha1sum file for all the objects in the path.
* [rclone size](/commands/rclone_size/)	 - Prints the total size and number of objects in remote:path.
* [rclone sync](/commands/rclone_sync/)	 - Make source and dest identical, modifying destination only.
* [rclone tree](/commands/rclone_tree/)	 - List the contents of the remote in a tree like fashion.
* [rclone version](/commands/rclone_version/)	 - Show the version number.

###### Auto generated by spf13/cobra on 30-Sep-2017
