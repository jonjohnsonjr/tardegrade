# tardegrade

This is a tool for selectively keeping certain files in a tar.gz stream without breaking the overall structure of the gzip stream.

This is useful for manipulating e.g. APKs that rely on a specific concatenation of gzip members in their data format.

## Example usage

Let's say I have an APK that contains the `crane` binary and an SBOM.

```
curl -L https://packages.wolfi.dev/os/aarch64/crane-0.18.0-r0.apk | tar -tvz
-rwxrwxrwx 0/0             512 2024-01-17 18:37 .SIGN.RSA.wolfi-signing.rsa.pub
-rw-r--r-- 0/0             410 2024-01-17 18:21 .PKGINFO
drwxr-xr-x 0/0               0 2024-01-17 18:21 usr
drwxr-xr-x 0/0               0 2024-01-17 18:21 usr/bin
-rwxr-xr-x 0/0         9580416 2024-01-17 18:21 usr/bin/crane
drwxr-xr-x 0/0               0 2024-01-17 18:21 var
drwxr-xr-x 0/0               0 2024-01-17 18:21 var/lib
drwxr-xr-x 0/0               0 2024-01-17 18:21 var/lib/db
drwxr-xr-x 0/0               0 2024-01-17 18:21 var/lib/db/sbom
-rw-r--r-- 0/0            2133 2024-01-17 18:21 var/lib/db/sbom/crane-0.18.0-r0.spdx.json
```

If I want to drop that SBOM:

```
curl -s -L https://packages.wolfi.dev/os/aarch64/crane-0.18.0-r0.apk | tardegrade $(curl -s -L https://packages.wolfi.dev/os/aarch64/crane-0.18.0-r0.apk | tar -t | head -n 5) | tar -tv
-rwxrwxrwx  0 root   root      512 Jan 17 10:37 .SIGN.RSA.wolfi-signing.rsa.pub
-rw-r--r--  0 root   root      410 Jan 17 10:21 .PKGINFO
drwxr-xr-x  0 root   root        0 Jan 17 10:21 usr
drwxr-xr-x  0 root   root        0 Jan 17 10:21 usr/bin
-rwxr-xr-x  0 root   root  9580416 Jan 17 10:21 usr/bin/crane
```

Notably, this tool only knows about tar and gzip.
It will not fix up the `.PKGINFO` file to update the metadata (datahash, size, etc.).

## Why is this useful?

I'm just using it to generate testdata for [melange](https://github.com/chainguard-dev/melange).
