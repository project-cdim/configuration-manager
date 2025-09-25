# Changelog 0.1

## [0.1.1](https://github.com/project-cdim/configuration-manager/compare/v0.1.0...v0.1.1) - 2025-09-24

The changes from v0.1.0 are as follows:

### Features

- Add "Unit" Vertex to database data
- Add publish functionality when Synchronize Hardware Information API completes
- Added "Unit" Vertex to database data
- Added publish functionality when Synchronize Hardware Information API completes
- The specification of the 'Update Resource Annotation for 1 Resource' API has been changed. When multiple resources exist within the same device unit and annotation is updated for any one of these resources, the annotation of other resources within the same device unit will also be updated accordingly.

### Breaking Changes

- Updated required Go version from 1.23+ to 1.24+.

### Bug Fixes

- Fixed issues with escape characters in path parameters, query parameters, and request body for Configuration Information Management API
