# Combined Login & Register Implementation

## Summary
Successfully combined login and register forms into a single page (`login.html`) without JavaScript, using pure CSS.

## Changes Made

### 1. Frontend Changes

#### `web/templates/login.html`
- Combined both login and register forms in one file
- Added CSS-only toggle system using hidden checkbox
- Added toggle buttons to switch between forms
- Register form now shows when `ShowRegister` flag is true
- Both forms maintain all original functionality:
  - Password visibility toggles
  - File uploads for profile images
  - Form validation
  - Error message display

### 2. Backend Changes

#### `Internal/api/register.go`
- Updated `RegisterPageData` struct to include `ShowRegister` flag
- Modified `Get_RegisterHandler` to use `login.html` template with `ShowRegister: true`
- Updated `renderRegister` function to use `login.html` and set `ShowRegister: true`
- Now when user visits `/register`, they see the register form automatically

#### `Internal/api/login.go`
- Updated `LoginPageData` struct to include `ShowRegister` flag (for consistency)
- Login still uses same `login.html` template but shows login form by default

#### `Internal/api/server.go`
- No changes needed - routes remain the same
- `/login` shows login form
- `/register` shows register form (same page, different state)

### 3. File Management
- Original `register.html` backed up as `register-old-backup.html`

## How It Works

1. **URL-based navigation:**
   - `/login` → Shows login form (checkbox unchecked)
   - `/register` → Shows register form (checkbox checked via `ShowRegister` flag)

2. **Toggle mechanism:**
   - Uses hidden checkbox `#toggle-register`
   - CSS selectors show/hide forms based on checkbox state
   - Toggle buttons at top allow switching between forms

3. **Backend integration:**
   - Both forms POST to their respective endpoints (`/login` and `/register`)
   - Server determines which form to show via `ShowRegister` flag
   - Error handling works for both forms

## User Experience

- Clicking "Login" in sidebar → Opens login form
- Clicking "Register" in sidebar → Opens register form
- User can toggle between forms using buttons
- No JavaScript required
- All animations and transitions via CSS
- Maintains form state on validation errors

## Testing Checklist

- [ ] `/login` shows login form by default
- [ ] `/register` shows register form by default
- [ ] Toggle buttons switch between forms
- [ ] Login form submits correctly
- [ ] Register form submits correctly (with file upload)
- [ ] Password visibility toggles work
- [ ] Error messages display correctly
- [ ] Form validation works
- [ ] Sidebar buttons navigate correctly
