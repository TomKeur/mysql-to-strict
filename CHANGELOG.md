# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2018-01-17

### Added

- A Changelog
- Added .sql files to .gitignore.
- Using the (databasename+timestamp+output.sql) as outputfile, when no file flag is given.

### Changed

- build_docker_image.sh, using VERSION from version file tag image with the correct version.

### Fixed

- When adding the UNKNOWN to ENUM's, there was no check if there already was a column with the name unknown. When this is the case, the empy value will be removed instead of trying to add another UNKNOWN column.

## [1.0.0] - 2018-01-16

- Initial release.