# Profile Images Directory

This directory contains all user profile images for the forum application.

## Contents

- `default.png` - Default profile image used when users don't upload their own photo
- `user_*.png/jpg` - Profile images uploaded during registration (format: `user_<username>.<ext>`)
- `profile_*.*` - Profile images uploaded via the edit profile page (format: `profile_<userid>_<timestamp>.<ext>`)

## Directory Structure

```
database/profile-img/
├── default.png              # Default profile image
├── user_*.{png,jpg,gif}     # Registration profile images
└── profile_*.*              # Updated profile images
```

## File Naming Conventions

1. **Default Image**: `default.png`
2. **Registration Images**: `user_<username>.<ext>`
   - Example: `user_john_doe.png`
3. **Updated Profile Images**: `profile_<userid>_<timestamp>.<ext>`
   - Example: `profile_123_1760538517.jpg`

## Access

Profile images are served via the `/profile-img/` route in the web application.

Example: `http://localhost:8080/profile-img/default.png`

## Migration Notes

All profile images were migrated from `web/static/img/` to this directory on October 21, 2025.

### Changes Made:
1. Moved all `user_*` and `profile_*` images from `web/static/img/` to `database/profile-img/`
2. Updated upload paths in:
   - `Internal/api/editProfile.go` (photo upload handler)
   - `Internal/api/register.go` (registration photo handler)
3. Added file server route in `Internal/api/server.go` to serve `/profile-img/`
4. Updated templates to use new path:
   - `web/templates/root.html`
   - `web/templates/profile.html`
   - `web/templates/edit-profile.html`

## Image Requirements

- **Supported formats**: PNG, JPG, JPEG, GIF
- **Maximum file size**: 5MB
- **Recommended dimensions**: 200x200 pixels or higher (square)
