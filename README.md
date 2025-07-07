# Storagebox

Hetzner storagebox utility to upload and download files.

## Usage

Uploading files

```bash
cat cat.jpg | storagebox upload remote/location/cat.jpg.encrypted
```

Uploading folders

```bash
tar -cvf - ./folder | storagebox upload remote/location/folder.tar.encrypted
```

Downloading files

```bash
storagebox download remote/location/cat.jpg.encrypted > cat.jpg
```
