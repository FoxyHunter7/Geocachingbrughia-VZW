# Geocaching Brughia Data Directory

This directory contains the SQLite database and is mounted into the Docker container.

**Important:** This directory should be backed up regularly!

## Backup

To backup the database:
```bash
cp data/geocaching.db data/geocaching.db.backup-$(date +%Y%m%d)
```

## Permissions

The directory needs to be writable by UID 65532 (nonroot user in distroless):
```bash
sudo chown -R 65532:65532 data/
```
